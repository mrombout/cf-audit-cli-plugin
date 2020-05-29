package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"encoding/json"
	"fmt"
	"github.com/mrombout/cf-audit-cli-plugin/client"
	"os"
)

func (c *AuditPlugin) runAuditServiceBinding(cliConnection plugin.CliConnection, args []string, cfClient client.CloudFoundryClient) {
	globalConfig, err := parseFlags(serviceBindingEventsFlags, cliConnection, args)
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		_, _ = cliConnection.CliCommand("help", CmdServiceBindingEvents)
		os.Exit(1)
	}

	if len(serviceBindingEventsFlags.Args()) == 0 {
		fmt.Print("Incorrect Usage: the required argument `SERVICE_NAME` was not provided\n\n")
		_, _ = cliConnection.CliCommand("help", CmdServiceBindingEvents)
		os.Exit(1)
	}
	serviceName := serviceBindingEventsFlags.Args()[0]

	service, err := cliConnection.GetService(serviceName)
	if err != nil {
		printFailed()
		fmt.Println(err)
		os.Exit(1)
	}

	request := client.ListAuditEventsRequest{
		OrganizationGuids: []string{globalConfig.organizationGuid},
		SpaceGuids:        []string{globalConfig.spaceGuid},
		Types: []string{
			"audit.service_binding.create",
			"audit.service_binding.delete",
		},
	}
	auditEvents, err := cfClient.ListAuditEvents(request)
	if err != nil {
		printFailed()
		fmt.Println(err)
		os.Exit(1)
	}

	var filteredAuditEvents []client.AuditEventModel
	for _, event := range auditEvents {
		if event.Type == "audit.service_binding.create" {
			var data client.AuditServiceBindingCreateData
			_ = json.Unmarshal(event.Data, &data)

			if data.Request.Relationships["service_instance"].Data.Guid == service.Guid {
				filteredAuditEvents = append(filteredAuditEvents, event)
			}
		} else if event.Type == "audit.service_binding.delete" {
			var data client.AuditServiceBindingDeleteData
			_ = json.Unmarshal(event.Data, &data)

			if data.Request.ServiceInstanceGuid == service.Guid {
				filteredAuditEvents = append(filteredAuditEvents, event)
			}
		}
	}

	printActionString(cliConnection, request)
	printAuditEventTable(filteredAuditEvents)
}
