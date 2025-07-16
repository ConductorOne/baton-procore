package client

const (
	BaseURL         = "https://api.procore.com/rest"
	GetCompaniesURL = BaseURL + "/v1.0/companies"
	GetProjectsURL  = BaseURL + "/v1.1/projects"

	// https://developers.procore.com/reference/rest/company-users?version=latest
	CompanyUsersURL = BaseURL + "/v1.3/companies/%s/users"

	ProjectUsersURL = BaseURL + "/v1.0/projects/%s/users"

	// https://developers.procore.com/reference/rest/project-users?version=latest#add-company-user-to-project
	AddUserToProjectURL = ProjectUsersURL + "/%d/actions/add"

	// https://developers.procore.com/reference/rest/project-users?version=latest#remove-a-user-from-the-project
	RemoveUserFromProjectURL = ProjectUsersURL + "/%d/actions/remove"
)
