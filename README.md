![Baton Logo](./baton-logo.png)

# `baton-procore` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-procore.svg)](https://pkg.go.dev/github.com/conductorone/baton-procore) ![ci](https://github.com/conductorone/baton-procore/actions/workflows/ci.yaml/badge.svg)

`baton-procore` is a connector for built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Prerequisites
No prerequisites were specified for `baton-procore`

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-procore
baton-procore
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_DOMAIN_URL=domain_url -e BATON_API_KEY=apiKey -e BATON_USERNAME=username public.ecr.aws/conductorone/baton-procore:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-procore/cmd/baton-procore@main

baton-procore

baton resources
```

# Data Model

`baton-procore` will pull down information about the following resources:
- Companies
- Projects
- Users

# Requirements

To use the `baton-procore` connector, you need to set up a Procore application with the following steps:

## Setting up Procore Application

1. **Create a Procore App**
   - Go to https://developers.procore.com/ and create a new application
   - Follow the documentation at https://developers.procore.com/documentation/building-apps-create-new for detailed instructions

2. **Configure the App for Data Connector**
   - After creating the app, navigate to the app's configuration builder section
   - Select `Data Connector Components`
   - Check the service account checkbox to enable service account functionality

3. **Get Your Credentials**
   - Copy the app's Client ID and Client Secret ID from the application settings
   - You'll need these values for the `--client-id` and `--client-secret` flags

4. **Enable Project Directory (For Provisioning)**
   - If you plan to use provisioning features, enable project directory in the projects you want to provision
   - Go to each project's admin section, then navigate to tool settings to enable this feature

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-procore` Command Line Usage

```
baton-procore

Usage:
  baton-procore [flags]
  baton-procore [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  config             Get the connector config schema
  help               Help about any command

Flags:
      --client-id string                                 The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string                             The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
      --external-resource-c1z string                     The path to the c1z file to sync external baton resources with ($BATON_EXTERNAL_RESOURCE_C1Z)
      --external-resource-entitlement-id-filter string   The entitlement that external users, groups must have access to sync external baton resources ($BATON_EXTERNAL_RESOURCE_ENTITLEMENT_ID_FILTER)
  -f, --file string                                      The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                                             help for baton-procore
      --log-format string                                The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string                                 The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --otel-collector-endpoint string                   The endpoint of the OpenTelemetry collector to send observability data to (used for both tracing and logging if specific endpoints are not provided) ($BATON_OTEL_COLLECTOR_ENDPOINT)
      --procore-client-id string                         required: The client ID to use for authentication. ($BATON_PROCORE_CLIENT_ID)
      --procore-client-secret string                     required: The client secret to use for authentication. ($BATON_PROCORE_CLIENT_SECRET)
  -p, --provisioning                                     This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --skip-full-sync                                   This must be set to skip a full sync ($BATON_SKIP_FULL_SYNC)
      --sync-resources strings                           The resource IDs to sync ($BATON_SYNC_RESOURCES)
      --ticketing                                        This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                                          version for baton-procore

Use "baton-procore [command] --help" for more information about a command.
```
