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

func userResource(user client.User) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"email":      user.ContactInfo.Email,
		"isEmployee": user.IsEmployee,
		// contact, previously known as reference person, is an individual without a procore account
		// https://support.procore.com/faq/what-is-a-contact-in-procore-and-which-project-tools-support-the-concept
		"contact": user.UserId == nil,
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

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	// TODO: get all users first then get all contact users
	// company users, then company people with reference users filter
	return nil, "", nil, nil
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
