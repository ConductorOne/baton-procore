package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	// Add the SchemaFields for the Config.
	ClientId = field.StringField(
		"procore-client-id",
		field.WithDescription("The client ID to use for authentication."),
		field.WithRequired(true),
		field.WithDisplayName("Client ID"),
	)

	ClientSecret = field.StringField(
		"procore-client-secret",
		field.WithDescription("The client secret to use for authentication."),
		field.WithRequired(true),
		field.WithDisplayName("Client Secret"),
		field.WithIsSecret(true),
	)

	ConfigurationFields = []field.SchemaField{ClientId, ClientSecret}

	// FieldRelationships defines relationships between the ConfigurationFields that can be automatically validated.
	// For example, a username and password can be required together, or an access token can be
	// marked as mutually exclusive from the username password pair.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run -tags=generate ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("Procore"),
)
