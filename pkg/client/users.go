package client

import (
	"context"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

// People refers to the users and contacts in Procore.
// contacts are also known as reference users in Procore.
// they are individuals without a Procore account, external to the organization.
//
//   - https://support.procore.com/faq/what-is-a-contact-in-procore-and-which-project-tools-support-the-concept
//   - https://developers.procore.com/reference/rest/company-people?version=latest
func (c *Client) GetCompanyPeople(ctx context.Context, companyId string) ([]Contact, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(CompanyPeopleURL, companyId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Procore-Company-Id", companyId)

	values := req.URL.Query()
	values.Set("filters[reference_users_only]", "true")
	req.URL.RawQuery = values.Encode()

	var target []Contact
	var rateLimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&rateLimitData),
	)

	if err != nil {
		return nil, fmt.Errorf("error getting company users from Procore API: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		logBody(ctx, res.Body)
		return nil, fmt.Errorf("unexpected status code: %d, expected: %d", res.StatusCode, http.StatusOK)
	}

	return target, nil
}

// People refers to the users and contacts in Procore.
// contacts are also known as reference users in Procore.
// they are individuals without a Procore account, external to the organization.
//
//   - https://support.procore.com/faq/what-is-a-contact-in-procore-and-which-project-tools-support-the-concept
//   - https://developers.procore.com/reference/rest/project-people?version=latest
func (c *Client) GetProjectPeople(ctx context.Context, companyId, projectId string, page int) ([]Contact, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(ProjectPeopleURL, projectId), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Procore-Company-Id", companyId)

	values := req.URL.Query()
	values.Set("filters[reference_users_only]", "true")
	values.Set("page", fmt.Sprintf("%d", page))
	values.Set("per_page", fmt.Sprintf("%d", perPage))
	req.URL.RawQuery = values.Encode()

	var target []Contact
	var rateLimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&rateLimitData),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("error getting company users from Procore API: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		logBody(ctx, res.Body)
		return nil, nil, fmt.Errorf("unexpected status code: %d, expected: %d", res.StatusCode, http.StatusOK)
	}

	return target, res, nil
}

func (c *Client) GetCompanyUsers(ctx context.Context, companyId string, page int) ([]User, *http.Response, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(CompanyPeopleURL, companyId), nil)
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
