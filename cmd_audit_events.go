package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/mrombout/cf-audit-cli-plugin/client"
	"os"
)

func (c *AuditPlugin) runAuditEvents(cliConnection plugin.CliConnection, args []string, cfClient client.CloudFoundryClient) {
	globalConfig, err := parseFlags(auditEventFlags, cliConnection, args)
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		_, _ = cliConnection.CliCommand("help", CmdServiceEvents)
		os.Exit(1)
	}

	request := client.ListAuditEventsRequest{
		SpaceGuids:        []string{globalConfig.spaceGuid},
		OrganizationGuids: []string{globalConfig.organizationGuid},
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
