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

const projectMembership = "member"

type projectBuilder struct {
	client *client.Client
}

func getCompanyId(resource *v2.Resource) (string, error) {
	groupTrait, err := resourceSdk.GetGroupTrait(resource)
	if err != nil {
		return "", fmt.Errorf("baton-procore: error getting group traits: %w", err)
	}
	traits := groupTrait.GetProfile().AsMap()
	companyId, ok := traits["company_id"].(string)
	if !ok {
		return "", fmt.Errorf("baton-procore: company_id not found in project resource profile")
	}
	return companyId, nil
}

func (o *projectBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return projectResourceType
}

func projectResource(project client.Project) (*v2.Resource, error) {
	profile := map[string]any{
		"company_id":   fmt.Sprintf("%d", project.Company.Id),
		"company_name": project.Company.Name,
		"active":       project.Active,
	}
	return resourceSdk.NewGroupResource(
		project.Name,
		projectResourceType,
		project.Id,
		[]resourceSdk.GroupTraitOption{
			resourceSdk.WithGroupProfile(profile),
		},
	)
}

func (o *projectBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	page := 1
	var err error
	if pToken.Token != "" {
		page, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-procore: failed to parse page token: %w", err)
		}
	}

	var annotations annotations.Annotations
	projects, res, rateLimitDesc, err := o.client.GetProjects(ctx, parentResourceID.Resource, page)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-procore: error getting projects: %w", err)
	}
	annotations = *annotations.WithRateLimiting(rateLimitDesc)

	rv := make([]*v2.Resource, 0, len(projects))
	for _, project := range projects {
		resource, err := projectResource(project)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-procore: error converting project to resource: %w", err)
		}
		rv = append(rv, resource)
	}

	var nextPage string
	if client.HasNextPage(res) {
		nextPage = strconv.Itoa(page + 1)
	}

	return rv, nextPage, annotations, nil
}

func (o *projectBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			resource,
			projectMembership,
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Member of %s project", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Member of %s project", resource.DisplayName)),
		),
	}, "", nil, nil
}

func (o *projectBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	page := 1
	var err error
	if pToken.Token != "" {
		page, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-procore: failed to parse page token: %w", err)
		}
	}

	// get company id from resource groupTrait
	companyId, err := getCompanyId(resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-procore: error getting company id from project resource: %w", err)
	}

	var annotations annotations.Annotations
	users, res, rateLimitDesc, err := o.client.GetProjectUsers(ctx, companyId, resource.Id.Resource, page)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-procore: error getting users: %w", err)
	}
	annotations = *annotations.WithRateLimiting(rateLimitDesc)

	rv := make([]*v2.Grant, 0, len(users))
	for _, user := range users {
		principalID, err := resourceSdk.NewResourceID(userResourceType, user.Id)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-procore: failed to create user resource ID: %w", err)
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

func (o *projectBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	projectId := entitlement.Resource.Id.Resource
	companyId, err := getCompanyId(entitlement.Resource)
	if err != nil {
		return nil, fmt.Errorf("baton-procore: error getting company id from project resource: %w", err)
	}
	userId, err := strconv.Atoi(principal.Id.Resource)
	if err != nil {
		return nil, fmt.Errorf("baton-procore: failed to parse user id from grant principal: %w", err)
	}

	err = o.client.AddUserToProject(ctx, companyId, projectId, userId)
	if err != nil {
		return nil, fmt.Errorf("baton-procore: error adding user to project: %w", err)
	}

	return nil, nil
}

func (o *projectBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	entitlement := grant.Entitlement
	projectId := entitlement.Resource.Id.Resource
	companyId, err := getCompanyId(entitlement.Resource)
	if err != nil {
		return nil, fmt.Errorf("baton-procore: error getting company id from project resource: %w", err)
	}
	userId, err := strconv.Atoi(grant.Principal.Id.Resource)
	if err != nil {
		return nil, fmt.Errorf("baton-procore: failed to parse user id from grant principal: %w", err)
	}

	err = o.client.RemoveUserFromProject(ctx, companyId, projectId, userId)
	if err != nil {
		return nil, fmt.Errorf("baton-procore: error removing user from project: %w", err)
	}

	return nil, nil
}

func newProjectBuilder(client *client.Client) *projectBuilder {
	return &projectBuilder{
		client: client,
	}
}
