package client

import (
	"context"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func (c *Client) GetProjects(ctx context.Context, companyId string, page int) ([]Project, *http.Response, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, GetProjectsURL, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Procore-Company-Id", companyId)
	values := req.URL.Query()
	values.Set("company_id", companyId)
	values.Set("page", fmt.Sprintf("%d", page))
	values.Set("per_page", fmt.Sprintf("%d", perPage))
	req.URL.RawQuery = values.Encode()

	var target []Project
	var rateLimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&rateLimitData),
	)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("baton-procore: error getting projects: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("unexpected status code: %d, expected: %d", res.StatusCode, http.StatusOK)
	}
	return target, res, &rateLimitData, nil
}

// https://developers.procore.com/reference/rest/project-users?version=latest#add-company-user-to-project
func (c *Client) AddUserToProject(ctx context.Context, companyId, projectId string, userId int) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(AddUserToProjectURL, projectId, userId), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Procore-Company-Id", companyId)

	res, err := c.Do(req)
	if err != nil {
		logBody(ctx, res.Body)
		return fmt.Errorf("baton-procore: error adding user to project: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		logBody(ctx, res.Body)
		return fmt.Errorf("unexpected status code: %d, expected: %d", res.StatusCode, http.StatusNoContent)
	}
	return nil
}

// https://developers.procore.com/reference/rest/project-users?version=latest#remove-a-user-from-the-project
func (c *Client) RemoveUserFromProject(ctx context.Context, companyId, projectId string, userId int) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(RemoveUserFromProjectURL, projectId, userId), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Procore-Company-Id", companyId)

	res, err := c.Do(req)
	if err != nil {
		logBody(ctx, res.Body)
		return fmt.Errorf("baton-procore: error removing user from project: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		logBody(ctx, res.Body)
		return fmt.Errorf("unexpected status code: %d, expected: %d", res.StatusCode, http.StatusNoContent)
	}
	return nil
}
