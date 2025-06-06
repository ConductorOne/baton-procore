package client

import (
	"context"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func (c *Client) GetProjects(ctx context.Context, companyId string) ([]Project, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ProjectsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Procore-Company-Id", companyId)
	values := req.URL.Query()
	values.Set("company_id", companyId)
	req.URL.RawQuery = values.Encode()

	var target []Project
	var rateLimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&rateLimitData),
	)
	// __AUTO_GENERATED_PRINT_VAR_START__
	fmt.Println(fmt.Sprintf("GetProjects target: %+v", target)) // __AUTO_GENERATED_PRINT_VAR_END__

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		logBody(ctx, res.Body)
		return nil, fmt.Errorf("unexpected status code: %d, expected: %d", res.StatusCode, http.StatusOK)
	}

	return target, nil
}
