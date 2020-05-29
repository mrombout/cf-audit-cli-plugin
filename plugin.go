// +build !V7

package main

import (
	"code.cloudfoundry.org/cli/plugin"
	plugin_models "code.cloudfoundry.org/cli/plugin/models"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mrombout/cf-audit-cli-plugin/client"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
)

type AuditPlugin struct{}

const CmdAuditEvents = "audit-events"
const CmdServiceEvents = "service-events"
const CmdServiceBindingEvents = "service-binding-events"

var Reset = "\033[0m"
var Bold = "\033[1m"
var White = "\033[37m"
var LightCyan = "\033[96m"

func EntityNameColor(input string) string {
	return Bold + LightCyan + input + Reset
}

func TableHeaderColor(input string) string {
	return Bold + White + input + Reset
}

func TableContentHeaderColor(input string) string {
	return Bold + LightCyan + input + Reset
}

// tabWriter doesn't recognize colors, so in order for the table header to align with the rows this will add empty null-
// chars to compensate for the color characters
func TableColorHack(input string) string {
	return strings.Repeat("\000", len(Bold+LightCyan)) + input + strings.Repeat("\000", len(Reset))
}

func (c *AuditPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	cfClient, err := createCloudFoundryClient(cliConnection)
	if err != nil {
		panic(err)
	}

	if args[0] == CmdAuditEvents {
		c.runAudit(cliConnection, args, cfClient)
	} else if args[0] == CmdServiceEvents {
		c.runAuditService(cliConnection, args, cfClient)
	} else if args[0] == CmdServiceBindingEvents {
		c.runAuditServiceBinding(cliConnection, args, cfClient)
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
					Usage: fmt.Sprintf("cf %s", CmdAuditEvents),
				},
			},
			{
				Name:     CmdServiceEvents,
				HelpText: "Show recent service audit events",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s SERVICE_NAME", CmdServiceEvents),
				},
			},
			{
				Name:     CmdServiceBindingEvents,
				HelpText: "Show recent service binding audit events",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s SERVICE_NAME", CmdServiceBindingEvents),
				},
			},
		},
	}
}

func createCloudFoundryClient(cliConnection plugin.CliConnection) (client.CloudFoundryClient, error) {
	apiEndpoint, err := cliConnection.ApiEndpoint()
	if err != nil {
		return client.CloudFoundryClient{}, nil
	}

	listAuditEventEndpointUrl, err := url.Parse(apiEndpoint + "/v3/audit_events")
	if err != nil {
		return client.CloudFoundryClient{}, nil
	}

	accessToken, err := cliConnection.AccessToken()
	if err != nil {
		return client.CloudFoundryClient{}, nil
	}

	cfClient := client.CloudFoundryClient{
		HttpClient:              http.Client{},
		CliConnection:           cliConnection,
		ListAuditEventsEndpoint: listAuditEventEndpointUrl,
		AccessToken:             accessToken,
	}
	return cfClient, nil
}

func parseFlags(command string, cliConnection plugin.CliConnection, args []string) (flags *flag.FlagSet, orgGuid string, spaceGuid string, err error) {
	organizationName := ""
	spaceName := ""

	flags = flag.NewFlagSet(command, flag.ExitOnError)
	flags.StringVar(&organizationName, "o", "", "Organization")
	flags.StringVar(&spaceName, "s", "", "Space")
	// TODO: flag to change types
	// TODO: flag to changes page
	// TODO: flag to change per_page
	// TODO: flag to change order-by

	err = flags.Parse(args[1:])

	if organizationName != "" {
		var org plugin_models.GetOrg_Model
		org, err = cliConnection.GetOrg(organizationName)
		if err != nil {
			return nil, "", "", err
		}
		orgGuid = org.Guid
	} else {
		currentOrg, err := cliConnection.GetCurrentOrg()
		if err != nil {
			return nil, "", "", err
		}
		orgGuid = currentOrg.Guid
	}

	if spaceName != "" {
		var space plugin_models.GetSpace_Model
		space, err = cliConnection.GetSpace(spaceName)
		if err != nil {
			return nil, "", "", err
		}
		spaceGuid = space.Guid
	} else {
		currentSpace, err := cliConnection.GetCurrentSpace()
		if err != nil {
			return nil, "", "", err
		}
		spaceGuid = currentSpace.Guid
	}

	return flags, orgGuid, spaceGuid, err
}

func (c *AuditPlugin) runAudit(cliConnection plugin.CliConnection, args []string, cfClient client.CloudFoundryClient) {
	_, orgGuid, spaceGuid, err := parseFlags(CmdAuditEvents, cliConnection, args)
	if err != nil {
		panic(err)
	}

	request := client.ListAuditEventsRequest{
		SpaceGuids:        []string{spaceGuid},
		OrganizationGuids: []string{orgGuid},
	}
	auditEvents, err := cfClient.ListAuditEvents(request)
	if err != nil {
		panic(err)
	}

	printActionString(cliConnection, request)
	printAuditEventTable(auditEvents)
}

func (c *AuditPlugin) runAuditService(cliConnection plugin.CliConnection, args []string, cfClient client.CloudFoundryClient) {
	flags, orgGuid, spaceGuid, err := parseFlags(CmdServiceEvents, cliConnection, args)
	if err != nil {
		panic(err)
	}

	serviceName := flags.Arg(0)
	if serviceName == "" {
		panic("no service name given, print usage")
	}

	service, err := cliConnection.GetService(serviceName)
	if err != nil {
		panic(err)
	}

	request := client.ListAuditEventsRequest{
		OrganizationGuids: []string{orgGuid},
		SpaceGuids:        []string{spaceGuid},
		TargetGuids:       []string{service.Guid},
	}
	auditEvents, err := cfClient.ListAuditEvents(request)
	if err != nil {
		panic(err)
	}

	printActionString(cliConnection, request)
	printAuditEventTable(auditEvents)
}

// TODO: Maybe just merge this with the audit-service command?
func (c *AuditPlugin) runAuditServiceBinding(cliConnection plugin.CliConnection, args []string, cfClient client.CloudFoundryClient) {
	flags, orgGuid, spaceGuid, err := parseFlags(CmdServiceBindingEvents, cliConnection, args)
	if err != nil {
		panic(err)
	}

	serviceName := flags.Arg(0)
	if serviceName == "" {
		panic("no service name given, print usage")
	}

	service, err := cliConnection.GetService(serviceName)
	if err != nil {
		panic(err)
	}

	request := client.ListAuditEventsRequest{
		OrganizationGuids: []string{orgGuid},
		SpaceGuids:        []string{spaceGuid},
		Types: []string{
			"audit.service_binding.create",
			"audit.service_binding.delete",
		},
	}
	auditEvents, err := cfClient.ListAuditEvents(request)
	if err != nil {
		panic(err)
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

func printActionString(cliConnection plugin.CliConnection, request client.ListAuditEventsRequest) {
	builder := strings.Builder{}

	builder.WriteString("Getting ")

	if len(request.Types) > 0 {
		builder.WriteString(EntityNameColor(strings.Join(request.Types, ", ")) + " ")
	} else {
		builder.WriteString("all ")
	}

	builder.WriteString("events ")

	if len(request.TargetGuids) > 0 {
		// TODO: Show names instead of guids
		builder.WriteString(fmt.Sprintf("for service %s ", EntityNameColor(strings.Join(request.TargetGuids, ", "))))
	}

	if len(request.OrganizationGuids) > 0 {
		// TODO: Show names instead of guids
		builder.WriteString(fmt.Sprintf("in org %s ", EntityNameColor(strings.Join(request.OrganizationGuids, ", "))))

		if len(request.SpaceGuids) > 0 {
			builder.WriteString("/ ")
		}
	}

	if len(request.SpaceGuids) > 0 {
		// TODO: Show names instead of guids
		builder.WriteString(EntityNameColor(strings.Join(request.SpaceGuids, ", ")) + " ")
	}

	username, err := cliConnection.Username()
	if err == nil {
		builder.WriteString("as " + EntityNameColor(username))
	}

	builder.WriteString("...\n\n")

	fmt.Print(builder.String())
}

func printAuditEventTable(auditEvents []client.AuditEventModel) {
	tableFormat := "%s\t %s\t %s\t %s\n"

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	// TODO: use same date format as CF CLI 2020-05-26T11:31:28.00+0200 (check CF CLI code, does CF CLI use hard-coded format or based on locale)
	fmt.Fprintf(w, tableFormat, TableHeaderColor("time"), TableHeaderColor("event"), TableHeaderColor("actor"), TableHeaderColor("description"))
	for _, event := range auditEvents {
		// TODO use different description strategy for each possible audit event type
		fmt.Fprintf(w, tableFormat, TableContentHeaderColor(event.CreatedAt), TableColorHack(event.Type), TableColorHack(event.Actor.Name), TableColorHack("TODO"))
	}
	w.Flush()
}

func main() {
	plugin.Start(new(AuditPlugin))
}
