package client

import (
	"context"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

// contacts are also known as reference users in Procore.
// they are individuals without a Procore account, external to the organization.
//
//	https://support.procore.com/faq/what-is-a-contact-in-procore-and-which-project-tools-support-the-concept
func (c *Client) GetContacts(ctx context.Context, companyId string) ([]Contact, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(PeopleURL, companyId), nil)
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

func (c *Client) GetCompanyUsers(ctx context.Context, companyId string) ([]User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(PeopleURL, companyId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var target []User
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
