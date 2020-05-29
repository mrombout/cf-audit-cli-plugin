package main

import (
	"code.cloudfoundry.org/cli/cf/flags"
	"code.cloudfoundry.org/cli/plugin"
)

type GlobalConfig struct {
	organizationGuid string
	spaceGuid string
}

func parseFlags(fc flags.FlagContext, cliConnection plugin.CliConnection, args []string) (GlobalConfig, error) {
	globalConfig := GlobalConfig{}

	err := fc.Parse(args...)
	if err != nil {
		return GlobalConfig{}, err
	}

	if fc.IsSet(FlagOrganization) {
		org, err := cliConnection.GetOrg(fc.String(FlagOrganization))
		if err != nil {
			panic(err)
		}
		globalConfig.organizationGuid = org.Guid
	} else {
		currentOrg, err := cliConnection.GetCurrentOrg()
		if err != nil {
			panic(err)
		}
		globalConfig.organizationGuid = currentOrg.Guid
	}

	if fc.IsSet(FlagSpace) {
		space, err := cliConnection.GetSpace(FlagSpace)
		if err != nil {
			panic(err)
		}
		globalConfig.spaceGuid = space.Guid
	} else {
		currentSpace, err := cliConnection.GetCurrentSpace()
		if err != nil {
			panic(err)
		}
		globalConfig.spaceGuid = currentSpace.Guid
	}

	return globalConfig, nil
}
