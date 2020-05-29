package main

import (
	"code.cloudfoundry.org/cli/cf/flags"
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/mrombout/cf-audit-cli-plugin/client"
)

const CmdAuditEvents = "audit-events"
const CmdServiceEvents = "service-events"
const CmdServiceBindingEvents = "service-binding-events"

const FlagOrganization = "organization"
const FlagSpace = "space"

var auditEventFlags flags.FlagContext
var serviceEventsFlags flags.FlagContext
var serviceBindingEventsFlags flags.FlagContext

func main() {
	auditEventFlags = flags.New()
	auditEventFlags.NewStringFlag(FlagOrganization, "o", "Organization")
	auditEventFlags.NewStringFlag(FlagSpace, "s", "Space")
	// TODO: flag to change types
	// TODO: flag to change limit (per_page)

	serviceEventsFlags = flags.New()
	serviceEventsFlags.NewStringFlag(FlagOrganization, "o", "Organization")
	serviceEventsFlags.NewStringFlag(FlagSpace, "s", "Space")
	// TODO: flag to change limit (per_page)

	serviceBindingEventsFlags = flags.New()
	serviceBindingEventsFlags.NewStringFlag(FlagOrganization, "o", "Organization")
	serviceBindingEventsFlags.NewStringFlag(FlagSpace, "s", "Space")
	// TODO: flag to change limit (per_page)

	plugin.Start(new(AuditPlugin))
}

type AuditPlugin struct{}

func (c *AuditPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	cfClient := client.CloudFoundryClient{
		CliConnection: cliConnection,
	}

	command := args[0]
	flagsAndParameters := args[1:]

	switch command {
	case CmdAuditEvents:
		c.runAuditEvents(cliConnection, flagsAndParameters, cfClient)
	case CmdServiceEvents:
		c.runAuditService(cliConnection, flagsAndParameters, cfClient)
	case CmdServiceBindingEvents:
		c.runAuditServiceBinding(cliConnection, flagsAndParameters, cfClient)
	}
}

func (c *AuditPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "AuditPlugin",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     CmdAuditEvents,
				HelpText: "Show recent audit events",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s [-o ORG] [-s SPACE]", CmdAuditEvents),
					Options: map[string]string {
						"-o": "Organization",
						"-s": "Space",
					},
				},
			},
			{
				Name:     CmdServiceEvents,
				HelpText: "Show recent service audit events",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s [-o ORG] [-s SPACE] SERVICE_NAME", CmdServiceEvents),
					Options: map[string]string {
						"-o": "Organization",
						"-s": "Space",
					},
				},
			},
			{
				Name:     CmdServiceBindingEvents,
				HelpText: "Show recent service binding audit events",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s [-o ORG] [-s SPACE] SERVICE_NAME", CmdServiceBindingEvents),
					Options: map[string]string {
						"-o": "Organization",
						"-s": "Space",
					},
				},
			},
		},
	}
}
