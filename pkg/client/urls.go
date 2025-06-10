package client

const (
	BaseURL      = "https://api.procore.com/rest/v1.0"
	CompaniesURL = BaseURL + "/companies"
	ProjectsURL  = BaseURL + "/projects"

	// https://developers.procore.com/reference/rest/company-people?version=latest
	// internal and external users are referred to as "people" in Procore.
	CompanyPeopleURL = CompaniesURL + "/%s/people"

	ProjectPeopleURL = ProjectsURL + "/%s/people"

	// https://developers.procore.com/reference/rest/company-users?version=latest
	// internal only, this endpoint brings back way more data than the People endpoint.
	CompanyUsersURL = CompaniesURL + "/%s/users"

	ProjectUsersURL = ProjectsURL + "/%s/users"
)
