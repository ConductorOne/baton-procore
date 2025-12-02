package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-procore/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"google.golang.org/grpc/codes"
)

type userBuilder struct {
	client *client.Client
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

func userResource(user client.User) (*v2.Resource, error) {
	profile := map[string]any{
		"email":      user.EmailAddress,
		"isEmployee": user.IsEmployee,
		// contact, previously known as reference person, is an individual without a procore account
		// https://support.procore.com/faq/what-is-a-contact-in-procore-and-which-project-tools-support-the-concept
		"contact": false,
	}

	status := v2.UserTrait_Status_STATUS_ENABLED
	if !user.IsActive {
		status = v2.UserTrait_Status_STATUS_DISABLED
	}

	ret, err := resourceSdk.NewUserResource(
		user.Name,
		userResourceType,
		user.Id,
		[]resourceSdk.UserTraitOption{
			resourceSdk.WithUserProfile(profile),
			resourceSdk.WithEmail(user.EmailAddress, true),
			resourceSdk.WithEmployeeID(user.EmployeeId),
			resourceSdk.WithCreatedAt(user.CreatedAt),
			resourceSdk.WithLastLogin(user.LastLoginAt),
			resourceSdk.WithStatus(status),
			resourceSdk.WithAccountType(v2.UserTrait_ACCOUNT_TYPE_HUMAN),
		},
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	page := 1
	var err error
	if pToken.Token != "" {
		page, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, fmt.Errorf("invalid page token format for users: %w", err)
		}
	}

	var annotations annotations.Annotations
	users, res, rateLimitDesc, err := o.client.GetCompanyUsers(ctx, parentResourceID.Resource, page)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get users: %w", err)
	}
	annotations = *annotations.WithRateLimiting(rateLimitDesc)

	rv := make([]*v2.Resource, 0, len(users))
	for _, user := range users {
		resource, err := userResource(user)
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to convert user data to resource: %w", err)
		}
		rv = append(rv, resource)
	}

	var nextPage string
	if client.HasNextPage(res) {
		nextPage = strconv.Itoa(page + 1)
	}
	return rv, nextPage, annotations, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *userBuilder) CreateAccountCapabilityDetails(ctx context.Context) (*v2.CredentialDetailsAccountProvisioning, annotations.Annotations, error) {
	return &v2.CredentialDetailsAccountProvisioning{
		SupportedCredentialOptions: []v2.CapabilityDetailCredentialOption{
			v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
		},
		PreferredCredentialOption: v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
	}, nil, nil
}

func (o *userBuilder) CreateAccount(ctx context.Context, accountInfo *v2.AccountInfo, credentialOptions *v2.CredentialOptions) (
	connectorbuilder.CreateAccountResponse,
	[]*v2.PlaintextData,
	annotations.Annotations,
	error,
) {
	var (
		err      error
		userBody client.UserBody
	)
	pMap := accountInfo.Profile.AsMap()
	companyId, ok := pMap["companyId"].(string)
	if !ok {
		return nil, nil, nil, uhttp.WrapErrors(codes.InvalidArgument, "companyId missing from account creation profile")
	}
	email, ok := pMap["email"].(string)
	if !ok {
		return nil, nil, nil, uhttp.WrapErrors(codes.InvalidArgument, "email missing from account creation profile")
	}
	lastName, ok := pMap["lastName"].(string)
	if !ok {
		return nil, nil, nil, uhttp.WrapErrors(codes.InvalidArgument, "lastName missing from account creation profile")
	}
	isEmployee, _ := pMap["isEmployee"].(bool)
	employeeId, ok := pMap["employeeId"].(string)
	if !ok && isEmployee {
		// if the user is an employee, the employeeId is required
		// https://developers.procore.com/reference/rest/company-users?version=latest#create-company-user
		// user.employee_id: string
		// The ID of the Employee of the Company User when user[is_employee] is set to true
		return nil, nil, nil, uhttp.WrapErrors(codes.InvalidArgument, "employeeId required for employee users but not provided")
	}
	firstName, _ := pMap["firstName"].(string)
	city, _ := pMap["city"].(string)
	address, _ := pMap["address"].(string)
	jobTitle, _ := pMap["jobTitle"].(string)
	userBody = client.UserBody{
		EmailAddress: email,
		LastName:     lastName,
		FirstName:    firstName,
		City:         city,
		Address:      address,
		JobTitle:     jobTitle,
		EmployeeId:   employeeId,
		IsEmployee:   isEmployee,
	}
	switch v := pMap["vendorId"].(type) {
	case int:
		userBody.VendorId = v
	case float64:
		// Convert float64 to int, but check for precision loss
		if v != float64(int64(v)) {
			return nil, nil, nil, uhttp.WrapErrors(codes.InvalidArgument, "vendorId has precision loss when converting from float64")
		}
		userBody.VendorId = int(v)
	case nil:
		// vendorId not provided, which is optional
	default:
		return nil, nil, nil, uhttp.WrapErrors(codes.InvalidArgument, fmt.Sprintf("vendorId has unexpected type %T", v))
	}

	user, err := o.client.CreateCompanyUser(ctx, companyId, client.CreateUserBody{
		User: userBody,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create account: %w", err)
	}
	userResource, err := userResource(*user)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to convert created user to resource: %w", err)
	}

	return &v2.CreateAccountResponse_SuccessResult{
		Resource:              userResource,
		IsCreateAccountResult: true,
	}, nil, nil, nil
}

func newUserBuilder(client *client.Client) *userBuilder {
	return &userBuilder{
		client: client,
	}
}
