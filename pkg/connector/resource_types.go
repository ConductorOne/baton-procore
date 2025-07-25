package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

// The user resource type is for all user objects from the database.
var userResourceType = &v2.ResourceType{
	Id:          "user",
	DisplayName: "User",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_USER},
}

var companyResourceType = &v2.ResourceType{
	Id:          "company",
	DisplayName: "Company",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_GROUP},
}

var projectResourceType = &v2.ResourceType{
	Id:          "project",
	DisplayName: "Project",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_GROUP},
}
