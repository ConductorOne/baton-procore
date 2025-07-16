package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func (c *Client) GetCompanyUsers(ctx context.Context, companyId string, page int) ([]User, *http.Response, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(CompanyUsersURL, companyId), nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create request: %w", err)
	}
	values := req.URL.Query()
	values.Set("page", fmt.Sprintf("%d", page))
	values.Set("per_page", fmt.Sprintf("%d", perPage))
	req.URL.RawQuery = values.Encode()

	var target []User
	var rateLimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&rateLimitData),
	)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("error getting company users from Procore API: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("unexpected status code: %d, expected: %d", res.StatusCode, http.StatusOK)
	}

	return target, res, &rateLimitData, nil
}

func (c *Client) GetProjectUsers(ctx context.Context, companyId, projectId string, page int) ([]User, *http.Response, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(ProjectUsersURL, projectId), nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Procore-Company-Id", companyId)

	values := req.URL.Query()
	values.Set("page", fmt.Sprintf("%d", page))
	values.Set("per_page", fmt.Sprintf("%d", perPage))
	req.URL.RawQuery = values.Encode()

	var target []User
	var rateLimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&rateLimitData),
	)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting project users from Procore API: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("unexpected status code: %d, expected: %d", res.StatusCode, http.StatusOK)
	}

	return target, res, &rateLimitData, nil
}

func (c *Client) CreateCompanyUser(ctx context.Context, companyId string, body CreateUserBody) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(CompanyUsersURL, companyId), bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Procore-Company-Id", companyId)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		logBody(ctx, res.Body)
		return fmt.Errorf("error creating company user in Procore API: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		logBody(ctx, res.Body)
		return fmt.Errorf("unexpected status code: %d, expected: %d", res.StatusCode, http.StatusOK)
	}

	return nil
}
