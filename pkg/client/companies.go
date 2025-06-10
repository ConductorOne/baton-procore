package client

import (
	"context"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

const perPage = 100

func (c *Client) GetCompanies(ctx context.Context, page int) ([]Company, *http.Response, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, CompaniesURL, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	values := req.URL.Query()
	values.Set("page", fmt.Sprintf("%d", page))
	values.Set("per_page", fmt.Sprintf("%d", perPage))
	req.URL.RawQuery = values.Encode()

	var target []Company
	var rateLimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&rateLimitData),
	)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("baton-procore: error getting companies: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("unexpected status code: %d, expected: %d", res.StatusCode, http.StatusOK)
	}

	return target, res, &rateLimitData, nil
}
