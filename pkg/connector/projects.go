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

type projectBuilder struct {
	client *client.Client
}

func (o *projectBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return projectResourceType
}

func projectResource(project client.Project) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"address":   project.Address,
		"city":      project.City,
		"company":   project.Company.Name,
		"active":    project.Active,
		"phone":     project.Phone,
		"createdAt": project.CreatedAt,
		"updatedAt": project.UpdatedAt,
	}
	return resourceSdk.NewGroupResource(
		project.Name,
		companyResourceType,
		project.Id,
		[]resourceSdk.GroupTraitOption{
			resourceSdk.WithGroupProfile(profile),
		},
	)
}

// List returns all the projects from the database as resource objects.
// Projects include a ProjectTrait because they are the 'shape' of a standard project.
func (o *projectBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	projects, err := o.client.GetProjects(ctx, parentResourceID.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-procore: error getting companies: %w", err)
	}

	rv := make([]*v2.Resource, 0, len(projects))
	for _, project := range projects {
		resource, err := projectResource(project)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-procore: error converting project to resource: %w", err)
		}
		rv = append(rv, resource)
	}

	return rv, "", nil, nil
}

func (o *projectBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *projectBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newProjectBuilder(client *client.Client) *projectBuilder {
	return &projectBuilder{
		client: client,
	}
}
