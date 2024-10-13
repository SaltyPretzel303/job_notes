package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/akamensky/argparse"
)

type SearchCmd struct {
	cmd *argparse.Command

	Company *string
	Role    *string

	Location   *string
	HybridFlag *bool
	RemoteFlag *bool

	Older *string
}

var ALLOWED_SEARCH_BY = [...]string{"company", "role", "location"}

func SetupGetCmd(parser *argparse.Parser) *SearchCmd {
	var searchCmd = SearchCmd{}

	cmd := parser.NewCommand("get", "Will search the job applications by the first argument.")
	searchCmd.cmd = cmd

	searchCmd.Company = cmd.String("", "company", &argparse.Options{Help: "Name (or part of the name) of the company you applied at."})
	searchCmd.Role = cmd.String("", "role", &argparse.Options{Help: "Role (or part of the) name that you applied for."})

	searchCmd.Location = cmd.String("l", "location", &argparse.Options{Help: "Location (or part of the location) you applied for."})
	searchCmd.RemoteFlag = cmd.Flag("r", "remote", &argparse.Options{Help: "Is the application for remote work model."})
	searchCmd.HybridFlag = cmd.Flag("y", "hybrid", &argparse.Options{Help: "Is the application for the hybrid work model."})

	searchCmd.Older = cmd.String("o", "older", &argparse.Options{Help: "Search for applications older than (in days).", Default: "0"})

	return &searchCmd
}

func (cmd *SearchCmd) Happened() bool {
	return cmd.cmd.Happened()
}

func (cmd *SearchCmd) GetCompany() string {
	return *cmd.Company
}

func (cmd *SearchCmd) GetRole() string {
	return *cmd.Role
}

func (cmd *SearchCmd) IsRemote() bool {
	return *cmd.RemoteFlag
}

func (cmd *SearchCmd) IsHybrid() bool {
	return *cmd.HybridFlag
}

func (cmd *SearchCmd) GetWorkType() []string {
	locs := []string{}

	// if *cmd.OnSiteFlag {
	// 	locs = append(locs, "onsite")
	// }

	if *cmd.HybridFlag {
		locs = append(locs, "hybrid")
	}

	if *cmd.RemoteFlag {
		locs = append(locs, "remote")
	}

	return locs
}

func (cmd *SearchCmd) GetLocation() string {
	return *cmd.Location
}

func (cmd *SearchCmd) GetOlderThanDate(format string) string {
	days, err := strconv.Atoi(*cmd.Older)

	if err != nil {
		fmt.Println("Failed to parse older than data: ", err)
	}

	return time.Now().AddDate(0, 0, -1*days).Local().Format(format)
}

func (cmd *SearchCmd) AsSearchData(dateFormat string) *SearchData {
	return &SearchData{
		Company:       cmd.GetCompany(),
		Role:          cmd.GetRole(),
		Location:      cmd.GetLocation(),
		WorkType:      cmd.GetWorkType(),
		IsHybrid:      cmd.IsHybrid(),
		IsRemote:      cmd.IsRemote(),
		OlderThanDate: cmd.GetOlderThanDate(dateFormat),
	}
}
