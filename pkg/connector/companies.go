package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-procore/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const companyMembership = "member"

type companyBuilder struct {
	client *client.Client
}

func (o *companyBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return companyResourceType
}

func companyResource(company client.Company) (*v2.Resource, error) {
	profile := map[string]any{
		"isActive": company.IsActive,
	}

	ret, err := resourceSdk.NewGroupResource(
		company.Name,
		companyResourceType,
		company.Id,
		nil,
		resourceSdk.WithResourceProfile(profile),
		resourceSdk.WithAnnotation(
			&v2.ChildResourceType{ResourceTypeId: projectResourceType.Id},
			&v2.ChildResourceType{ResourceTypeId: userResourceType.Id},
		),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *companyBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	page := 1
	var err error
	if pToken.Token != "" {
		page, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, fmt.Errorf("invalid page token format for companies: %w", err)
		}
	}

	var annotations annotations.Annotations
	companies, res, rateLimitDesc, err := o.client.GetCompanies(ctx, page)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get companies: %w", err)
	}
	annotations = *annotations.WithRateLimiting(rateLimitDesc)

	rv := make([]*v2.Resource, 0, len(companies))
	for _, company := range companies {
		resource, err := companyResource(company)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error converting company to resource: %w", err)
		}
		rv = append(rv, resource)
	}

	var nextPage string
	if client.HasNextPage(res) {
		nextPage = strconv.Itoa(page + 1)
	}
	return rv, nextPage, annotations, nil
}

func (o *companyBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			resource,
			companyMembership,
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Member of %s company", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Member of %s company", resource.DisplayName)),
		),
	}, "", nil, nil
}

func (o *companyBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	page := 1
	var err error
	if pToken.Token != "" {
		page, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, fmt.Errorf("invalid page token format for company grants: %w", err)
		}
	}

	var annotations annotations.Annotations
	users, res, rateLimitDesc, err := o.client.GetCompanyUsers(ctx, resource.Id.Resource, page)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get users: %w", err)
	}
	annotations = *annotations.WithRateLimiting(rateLimitDesc)

	rv := make([]*v2.Grant, 0, len(users))
	for _, user := range users {
		principalID, err := resourceSdk.NewResourceID(userResourceType, user.Id)
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to create user resource ID: %w", err)
		}
		rv = append(rv, grant.NewGrant(
			resource,
			companyMembership,
			principalID,
		))
	}

	var nextPage string
	if client.HasNextPage(res) {
		nextPage = strconv.Itoa(page + 1)
	}
	return rv, nextPage, annotations, nil
}

func newCompanyBuilder(client *client.Client) *companyBuilder {
	return &companyBuilder{
		client: client,
	}
}
