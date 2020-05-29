package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/mrombout/cf-audit-cli-plugin/client"
	"os"
)

func (c *AuditPlugin) runAuditService(cliConnection plugin.CliConnection, args []string, cfClient client.CloudFoundryClient) {
	globalConfig, err := parseFlags(serviceEventsFlags, cliConnection, args)
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		_, _ = cliConnection.CliCommand("help", CmdServiceEvents)
		os.Exit(1)
	}

	if len(serviceEventsFlags.Args()) == 0 {
		fmt.Print("Incorrect Usage: the required argument `SERVICE_NAME` was not provided\n\n")
		_, _ = cliConnection.CliCommand("help", CmdServiceEvents)
		os.Exit(1)
	}
	serviceName := serviceEventsFlags.Args()[0]

	service, err := cliConnection.GetService(serviceName)
	if err != nil {
		printFailed()
		fmt.Println(err)
		os.Exit(1)
	}

	request := client.ListAuditEventsRequest{
		OrganizationGuids: []string{globalConfig.organizationGuid},
		SpaceGuids:        []string{globalConfig.spaceGuid},
		TargetGuids:       []string{service.Guid},
	}
	auditEvents, err := cfClient.ListAuditEvents(request)
	if err != nil {
		printFailed()
		fmt.Println(err)
		os.Exit(1)
	}

	printActionString(cliConnection, request)
	printAuditEventTable(auditEvents)
}
