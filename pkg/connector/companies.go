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

type companyBuilder struct {
	client *client.Client
}

func (o *companyBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return companyResourceType
}

func companyResource(company client.Company) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"isActive": company.IsActive,
	}
	return resourceSdk.NewGroupResource(
		company.Name,
		companyResourceType,
		company.Id,
		[]resourceSdk.GroupTraitOption{
			resourceSdk.WithGroupProfile(profile),
		},
		resourceSdk.WithAnnotation(
			&v2.ChildResourceType{ResourceTypeId: projectResourceType.Id},
		),
	)
}

// List returns all the companys from the database as resource objects.
// Companys include a CompanyTrait because they are the 'shape' of a standard company.
func (o *companyBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	companies, err := o.client.GetCompanies(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-procore: error getting companies: %w", err)
	}

	rv := make([]*v2.Resource, 0, len(companies))
	for _, company := range companies {
		resource, err := companyResource(company)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-procore: error converting company to resource: %w", err)
		}
		rv = append(rv, resource)
	}
	return rv, "", nil, nil
}

func (o *companyBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *companyBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newCompanyBuilder(client *client.Client) *companyBuilder {
	return &companyBuilder{
		client: client,
	}
}
