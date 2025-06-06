package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-procore/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	client *client.Client
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

func contactResource(user client.Contact) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"email":      user.ContactInfo.Email,
		"isEmployee": user.IsEmployee,
		// contact, previously known as reference person, is an individual without a procore account
		// https://support.procore.com/faq/what-is-a-contact-in-procore-and-which-project-tools-support-the-concept
		"contact": true,
	}

	return resourceSdk.NewUserResource(
		fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		userResourceType,
		user.Id,
		[]resourceSdk.UserTraitOption{
			resourceSdk.WithUserProfile(profile),
			resourceSdk.WithEmail(user.ContactInfo.Email, true),
			resourceSdk.WithEmployeeID(user.EmployeeId),
		},
	)
}

func userResource(user client.User) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"email":      user.EmailAddress,
		"isEmployee": user.IsEmployee,
		// contact, previously known as reference person, is an individual without a procore account
		// https://support.procore.com/faq/what-is-a-contact-in-procore-and-which-project-tools-support-the-concept
		"contact": false,
	}
	// __AUTO_GENERATED_PRINT_VAR_START__
	fmt.Println(fmt.Sprintf("userResource user: %+v", user)) // __AUTO_GENERATED_PRINT_VAR_END__

	status := v2.UserTrait_Status_STATUS_ENABLED
	if !user.IsActive {
		status = v2.UserTrait_Status_STATUS_DISABLED
	}

	_type := v2.UserTrait_ACCOUNT_TYPE_HUMAN
	if !user.IsEmployee {
		// managed service account, for installed apps
		_type = v2.UserTrait_ACCOUNT_TYPE_SERVICE
	}

	return resourceSdk.NewUserResource(
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
			resourceSdk.WithAccountType(_type),
		},
	)
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	contacts, err := o.client.GetContacts(ctx, parentResourceID.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-procore: error getting users: %w", err)
	}

	rv := make([]*v2.Resource, 0, len(contacts))
	for _, contact := range contacts {
		resource, err := contactResource(contact)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-procore: error converting user to resource: %w", err)
		}
		rv = append(rv, resource)
	}

	users, err := o.client.GetCompanyUsers(ctx, parentResourceID.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-procore: error getting users: %w", err)
	}

	for _, user := range users {
		resource, err := userResource(user)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-procore: error converting user to resource: %w", err)
		}
		rv = append(rv, resource)
	}
	return rv, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client *client.Client) *userBuilder {
	return &userBuilder{
		client: client,
	}
}
