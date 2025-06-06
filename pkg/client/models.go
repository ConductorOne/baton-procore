package client

import "time"

type Company struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type Project struct {
	Id                      int64          `json:"id"`
	AccountingProjectNumber *string        `json:"accounting_project_number"`
	Active                  bool           `json:"active"`
	Address                 *string        `json:"address"`
	City                    *string        `json:"city"`
	Company                 Company        `json:"company"`
	CompletionDate          string         `json:"completion_date"`
	CountryCode             string         `json:"country_code"`
	County                  *string        `json:"county"`
	CreatedAt               time.Time      `json:"created_at"`
	CreatedBy               CreatedBy      `json:"created_by"`
	CustomFields            map[string]any `json:"custom_fields"`
	DeliveryMethod          *string        `json:"delivery_method"`
	DesignatedMarketArea    *string        `json:"designated_market_area"`
	DisplayName             string         `json:"display_name"`
	EstimatedValue          string         `json:"estimated_value"`
	IsDemo                  bool           `json:"is_demo"`
	Latitude                *float64       `json:"latitude"`
	Longitude               *float64       `json:"longitude"`
	Name                    string         `json:"name"`
	OriginCode              *string        `json:"origin_code"`
	OriginData              *string        `json:"origin_data"`
	OriginId                *string        `json:"origin_id"`
	OwnersProjectId         *string        `json:"owners_project_id"`
	ParentJob               *string        `json:"parent_job"`
	ParentJobId             *string        `json:"parent_job_id"`
	Phone                   *string        `json:"phone"`
	PhotoId                 *string        `json:"photo_id"`
	ProjectBidTypeId        *string        `json:"project_bid_type_id"`
	ProjectNumber           *string        `json:"project_number"`
	ProjectOwnerTypeId      *string        `json:"project_owner_type_id"`
	ProjectRegionId         *string        `json:"project_region_id"`
	ProjectStage            *string        `json:"project_stage"`
	ProjectedFinishDate     string         `json:"projected_finish_date"`
	Sector                  string         `json:"sector"`
	Stage                   string         `json:"stage"`
	StartDate               string         `json:"start_date"`
	StateCode               *string        `json:"state_code"`
	StoreNumber             *string        `json:"store_number"`
	TimeZone                string         `json:"time_zone"`
	TotalValue              string         `json:"total_value"`
	UpdatedAt               time.Time      `json:"updated_at"`
	WorkScope               string         `json:"work_scope"`
	Zip                     *string        `json:"zip"`
}

type CreatedBy struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type Contact struct {
	Id                   int64       `json:"id"`
	ContactInfo          ContactInfo `json:"contact"`
	EmployeeId           string      `json:"employee_id"`
	FirstName            string      `json:"first_name"`
	IsEmployee           bool        `json:"is_employee"`
	LastName             string      `json:"last_name"`
	UserId               *int64      `json:"user_id"`
	UserUUID             *string     `json:"user_uuid"`
	WorkClassificationId int64       `json:"work_classification_id"`
	OriginId             string      `json:"origin_id"`
}

type ContactInfo struct {
	IsActive bool   `json:"is_active"`
	Email    string `json:"email"`
}

type User struct {
	Address                   string             `json:"address"`
	Avatar                    string             `json:"avatar"`
	BusinessId                string             `json:"business_id"`
	BusinessPhone             string             `json:"business_phone"`
	BusinessPhoneExtension    int                `json:"business_phone_extension"`
	City                      string             `json:"city"`
	ContactId                 int                `json:"contact_id"`
	CountryCode               string             `json:"country_code"`
	CreatedAt                 time.Time          `json:"created_at"`
	EmailAddress              string             `json:"email_address"`
	EmailSignature            string             `json:"email_signature"`
	EmployeeId                string             `json:"employee_id"`
	ERPIntegratedAccountant   bool               `json:"erp_integrated_accountant"`
	FaxNumber                 string             `json:"fax_number"`
	FirstName                 string             `json:"first_name"`
	Id                        int                `json:"id"`
	Initials                  string             `json:"initials"`
	IsActive                  bool               `json:"is_active"`
	IsEmployee                bool               `json:"is_employee"`
	JobTitle                  string             `json:"job_title"`
	LastActivatedAt           time.Time          `json:"last_activated_at"`
	LastLoginAt               time.Time          `json:"last_login_at"`
	LastName                  string             `json:"last_name"`
	MobilePhone               string             `json:"mobile_phone"`
	Name                      string             `json:"name"`
	Notes                     string             `json:"notes"`
	OriginId                  string             `json:"origin_id"`
	OriginData                string             `json:"origin_data"`
	StateCode                 string             `json:"state_code"`
	UpdatedAt                 time.Time          `json:"updated_at"`
	WelcomeEmailSentAt        time.Time          `json:"welcome_email_sent_at"`
	Zip                       string             `json:"zip"`
	WorkClassificationId      int                `json:"work_classification_id"`
	PermissionTemplate        PermissionTemplate `json:"permission_template"`
	CompanyPermissionTemplate PermissionTemplate `json:"company_permission_template"`
	Vendor                    Vendor             `json:"vendor"`
	Role                      string             `json:"role"`
	VerifiedEmployee          bool               `json:"verified_employee"`
}

type PermissionTemplate struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	ProjectSpecific bool   `json:"project_specific"`
	Type            string `json:"type"`
}

type Vendor struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
