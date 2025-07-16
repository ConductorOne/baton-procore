package connector

import (
	"context"
	"fmt"
	"io"

	"github.com/conductorone/baton-procore/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
)

type Connector struct {
	client *client.Client
	// cache is needed because project users ids are different from company users ids, even if
	// they are the same user.
	//	email: company_id
	usersCache map[string]int
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newCompanyBuilder(d.client),
		newProjectBuilder(d.client),
		newUserBuilder(d.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Connector) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Baton Connector",
		Description: "This connector allows you to sync data from Procore.",
		AccountCreationSchema: &v2.ConnectorAccountCreationSchema{
			FieldMap: map[string]*v2.ConnectorAccountCreationSchema_Field{
				"companyId": {
					DisplayName: "Company ID",
					Required:    true,
					Description: "The ID of the company to which the user belongs.",
					Field: &v2.ConnectorAccountCreationSchema_Field_StringField{
						StringField: &v2.ConnectorAccountCreationSchema_StringField{},
					},
					Placeholder: "Company ID",
					Order:       1,
				},
				"email": {
					DisplayName: "Email",
					Required:    true,
					Description: "The email address of the user.",
					Field: &v2.ConnectorAccountCreationSchema_Field_StringField{
						StringField: &v2.ConnectorAccountCreationSchema_StringField{},
					},
					Placeholder: "Email",
					Order:       2,
				},
				"lastName": {
					DisplayName: "User's Last Name",
					Required:    true,
					Description: "The last name of the user.",
					Field: &v2.ConnectorAccountCreationSchema_Field_StringField{
						StringField: &v2.ConnectorAccountCreationSchema_StringField{},
					},
					Placeholder: "Last Name",
					Order:       3,
				},
				"firstName": {
					DisplayName: "User's First Name",
					Required:    false,
					Description: "The first name of the user.",
					Field: &v2.ConnectorAccountCreationSchema_Field_StringField{
						StringField: &v2.ConnectorAccountCreationSchema_StringField{},
					},
					Placeholder: "First Name",
					Order:       4,
				},
				"city": {
					DisplayName: "User's City",
					Required:    false,
					Description: "The city where the user resides.",
					Field: &v2.ConnectorAccountCreationSchema_Field_StringField{
						StringField: &v2.ConnectorAccountCreationSchema_StringField{},
					},
					Placeholder: "City",
					Order:       5,
				},
				"address": {
					DisplayName: "User's Address",
					Required:    false,
					Description: "The address of the user.",
					Field: &v2.ConnectorAccountCreationSchema_Field_StringField{
						StringField: &v2.ConnectorAccountCreationSchema_StringField{},
					},
					Placeholder: "Address",
					Order:       6,
				},
				"jobTitle": {
					DisplayName: "User's Job Title",
					Required:    false,
					Description: "The job title of the user.",
					Field: &v2.ConnectorAccountCreationSchema_Field_StringField{
						StringField: &v2.ConnectorAccountCreationSchema_StringField{},
					},
					Placeholder: "Job Title",
					Order:       7,
				},
				"isEmployee": {
					DisplayName: "Is Employee",
					Required:    false,
					Description: "Indicates if the user is an employee.",
					Field: &v2.ConnectorAccountCreationSchema_Field_StringField{
						StringField: &v2.ConnectorAccountCreationSchema_StringField{},
					},
					Placeholder: "Is Employee",
					Order:       8,
				},
				"isActive": {
					DisplayName: "Is Active",
					Required:    false,
					Description: "Indicates if the user is currently active.",
					Field: &v2.ConnectorAccountCreationSchema_Field_StringField{
						StringField: &v2.ConnectorAccountCreationSchema_StringField{},
					},
					Placeholder: "Is Active",
					Order:       9,
				},
			},
		},
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, clientId, clientSecret string) (*Connector, error) {
	client, err := client.New(ctx, clientId, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("error creating Procore client: %w", err)
	}
	return &Connector{
		client:     client,
		usersCache: make(map[string]int),
	}, nil
}
