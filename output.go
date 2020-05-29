package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/mrombout/cf-audit-cli-plugin/client"
	"os"
	"strings"
	"text/tabwriter"
)

var Reset = "\033[0m"
var Bold = "\033[1m"
var White = "\033[37m"
var LightCyan = "\033[96m"
var Red = "\033[31m"

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

func ErrorColor(input string) string {
	return Bold + Red + input + Reset
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
			builder.WriteString(fmt.Sprintf("space %s", EntityNameColor(strings.Join(request.SpaceGuids, ", ")) + " "))
		}
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

func printFailed() {
	fmt.Println(ErrorColor("FAILED"))
}