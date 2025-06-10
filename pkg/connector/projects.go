package connector

import (
	"context"
	"fmt"
	"strconv"
	"sync"

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
	client     *client.Client
	m          sync.Mutex
	usersCache *map[string]int
}

func (o *projectBuilder) fillCache(ctx context.Context, companyId string) error {
	var page = 1
	users, res, err := o.client.GetCompanyUsers(ctx, companyId, page)
	if err != nil {
		return fmt.Errorf("baton-procore: error getting users for company %s: %w", companyId, err)
	}
	for _, u := range users {
		(*o.usersCache)[u.EmailAddress] = u.Id
	}

	for client.HasNextPage(res) {
		page++
		users, res, err = o.client.GetCompanyUsers(ctx, companyId, page)
		if err != nil {
			return fmt.Errorf("baton-procore: error getting users for company %s: %w", companyId, err)
		}
		for _, u := range users {
			(*o.usersCache)[u.EmailAddress] = u.Id
		}
	}
	return nil
}

func (o *projectBuilder) GetUserId(ctx context.Context, companyId, email string) (int, error) {
	o.m.Lock()
	defer o.m.Unlock()
	userId, ok := (*o.usersCache)[email]
	if ok {
		return userId, nil
	}

	err := o.fillCache(ctx, companyId)
	if err != nil {
		return 0, fmt.Errorf("baton-procore: error filling user cache for company %s: %w", companyId, err)
	}

	userId, ok = (*o.usersCache)[email]
	if !ok {
		return 0, fmt.Errorf("baton-procore: user with email %s not found in company %s", email, companyId)
	}
	return userId, nil
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

	var page = 1
	var err error
	if pToken.Token != "" {
		page, err = strconv.Atoi(pToken.Token)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-terraform-cloud: failed to parse page token: %w", err)
		}
	}

	projects, res, err := o.client.GetProjects(ctx, parentResourceID.Resource, page)
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

	var nextPage string
	if client.HasNextPage(res) {
		nextPage = strconv.Itoa(page + 1)
	}

	return rv, nextPage, nil, nil
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
	// get company id from resource groupTrait
	groupTrait, err := resourceSdk.GetGroupTrait(resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-procore: error getting group traits: %w", err)
	}
	traits := groupTrait.GetProfile().AsMap()
	companyId, ok := traits["company_id"].(string)
	if !ok {
		return nil, "", nil, fmt.Errorf("baton-procore: company_id not found in project resource profile")
	}

	users, err := o.client.GetProjectUsers(ctx, companyId, resource.Id.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-procore: error getting users: %w", err)
	}

	rv := make([]*v2.Grant, 0, len(users))
	for _, user := range users {
		// using company user id because project users have a different id, even if they are the same user.
		companyUserId, err := o.GetUserId(ctx, companyId, user.EmailAddress)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-procore: error getting user id for email %s: %w", user.EmailAddress, err)
		}

		principalID, err := resourceSdk.NewResourceID(userResourceType, companyUserId)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-procore: failed to create user resource ID: %w", err)
		}
		rv = append(rv, grant.NewGrant(
			resource,
			companyMembership,
			principalID,
		))
	}
	return rv, "", nil, nil
}

func newProjectBuilder(client *client.Client, userCache *map[string]int) *projectBuilder {
	return &projectBuilder{
		client:     client,
		m:          sync.Mutex{},
		usersCache: userCache,
	}
}
