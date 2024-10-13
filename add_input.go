package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/akamensky/argparse"
)

type AddCmd struct {
	cmd *argparse.Command

	Company  *string
	Role     *string
	Location *string

	HybridFlag *bool
	RemoteFlag *bool

	CompanyArg         *string
	RoleArg            *string
	LocationArg        *string
	WorkTypeArg        *string
	ApplicationDateArg *string
	ResponseDateArg    *string
	ResponseArg        *string

	CommentArg *string
}

func SetupAddCmd(parser *argparse.Parser) *AddCmd {

	var fields = &AddCmd{}

	cmd := parser.NewCommand("add", "Will add a job application.")

	fields.cmd = cmd
	fields.Company = cmd.StringPositional(&argparse.Options{Help: "Name of the company."})
	fields.Role = cmd.StringPositional(&argparse.Options{Help: "Position or a role at the company."})
	fields.Location = cmd.StringPositional(&argparse.Options{Help: `Work location,specify the city and then is it remote or hybrid with -r and -b flags.`})

	fields.RemoteFlag = cmd.Flag("r", "remote", &argparse.Options{Help: "Will flag this as the remote work model."})
	fields.HybridFlag = cmd.Flag("b", "hybrid", &argparse.Options{Help: "Will flag this as a hybrid work model with the main location specified with the location parameter."})

	fields.CompanyArg = cmd.String("", "company", &argparse.Options{Help: "Will have higher priority compared to the positional argument."})
	fields.RoleArg = cmd.String("", "role", &argparse.Options{Help: "Will have higher priority compared to the positional argument."})
	fields.LocationArg = cmd.String("", "location", &argparse.Options{Help: "Will have higher priority compared to the positional argument."})
	fields.WorkTypeArg = cmd.String("t", "type", &argparse.Options{Help: "Work model: hybrid or remote."})
	fields.ApplicationDateArg = cmd.String("", "sdate", &argparse.Options{Help: "If omitted will have a value of current date (ignoring the time)."})
	fields.ResponseDateArg = cmd.String("", "rdate", &argparse.Options{Help: "The date response is received, in case an old (resloved) application is being added."})
	fields.ResponseArg = cmd.String("", "response", &argparse.Options{Help: "Company's response (accepted, denied, hr interview, tech interview ...)."})

	fields.CommentArg = cmd.String("", "comment", &argparse.Options{Help: "General purpose comment."})

	return fields
}

func (cmd *AddCmd) Happened() bool {
	return cmd.cmd.Happened()
}

func (cmd *AddCmd) GetCompany() (string, error) {
	if *cmd.CompanyArg != "" {
		return *cmd.CompanyArg, nil
	} else if *cmd.Company != "" {
		return *cmd.Company, nil
	} else {
		return "", errors.New("no company provided")
	}
}

func (cmd *AddCmd) GetRole() (string, error) {
	if *cmd.RoleArg != "" {
		return *cmd.RoleArg, nil
	} else if *cmd.Role != "" {
		return *cmd.Role, nil
	} else {
		return "", errors.New("no role provided")
	}
}

// location does not have to be provided, job can be fully remote
func (cmd *AddCmd) GetLocation() string {
	if *cmd.LocationArg != "" {
		return *cmd.LocationArg
	} else if *cmd.Location != "" {
		return *cmd.Location
	} else {
		return ""
	}
}

func (cmd *AddCmd) IsHybrid() bool {
	return (*cmd.HybridFlag || *cmd.WorkTypeArg == "hybrid")
}

func (cmd *AddCmd) IsRemote() bool {
	return (*cmd.RemoteFlag || *cmd.WorkTypeArg == "remote")
}

// func (cmd *AddCmd) GetWorkType() string {
// 	EMPTY_VAL := "      " // 6 spaces -> len(remote|hybrid|onsite) == 6 )

// 	var remoteValue = EMPTY_VAL
// 	if cmd.IsRemote() {
// 		remoteValue = "remote"
// 	}

// 	var hybridValue = EMPTY_VAL
// 	if cmd.IsHybrid() {
// 		hybridValue = "hybrid"
// 	}

// 	var onSiteValue = EMPTY_VAL
// 	if cmd.IsHybrid() || (!cmd.IsHybrid() && !cmd.IsRemote()) {
// 		onSiteValue = "onsite"
// 	}

// 	return fmt.Sprintf("%v/%v/%v", onSiteValue, hybridValue, remoteValue)
// }

func (cmd *AddCmd) GetApplicationDate(format string) (time.Time, error) {
	if *cmd.ApplicationDateArg != "" {
		dt, err := time.Parse(format, *cmd.ApplicationDateArg)
		if err != nil {
			errMsg := "application date in invalid format: " + *cmd.ApplicationDateArg
			return time.Time{}, errors.New(errMsg)
		}

		return dt, nil
	} else {
		return time.Now(), nil
	}
}

func (cmd *AddCmd) GetResponseDate(format string) (time.Time, error) {
	if *cmd.ResponseDateArg != "" {
		dt, err := time.Parse(*cmd.ResponseArg, format)
		if err != nil {
			errMsg := "response date in invalid format: " + *cmd.ResponseDateArg
			return time.Time{}, errors.New(errMsg)
		}

		return dt, nil
	} else {
		return time.Time{}, nil
	}
}

func (cmd *AddCmd) GetResponse() string {
	if *cmd.ResponseArg != "" {
		return *cmd.ResponseArg
	} else {
		return "No response"
	}
}

func (cmd *AddCmd) GetComment() string {
	if *cmd.CommentArg == "" {
		return " "
	} else {
		return *cmd.CommentArg
	}
}

func formatWorkType(location string, isRemote, isHybrid bool) string {
	EMPTY_VAL := "      " // 6 spaces -> len(remote|hybrid|onsite) == 6 )

	var remoteValue = EMPTY_VAL
	if isRemote {
		remoteValue = "remote"
	}

	var hybridValue = EMPTY_VAL
	if isHybrid {
		hybridValue = "hybrid"
	}

	var onSiteValue = EMPTY_VAL
	if isHybrid || (!isHybrid && !isRemote) {
		if location != "" {
			onSiteValue = location
		}
	}

	return fmt.Sprintf("%v/%v/%v", onSiteValue, hybridValue, remoteValue)
}

func (cmd *AddCmd) AsApplication(dateFormat string) *ApplicationData {
	appData := ApplicationData{}

	cmp, _ := cmd.GetCompany()
	appData.Company = cmp

	role, _ := cmd.GetRole()
	appData.Role = role

	appData.Location = cmd.GetLocation()

	appData.WorkType = formatWorkType(cmd.GetLocation(), cmd.IsRemote(), cmd.IsHybrid())

	appData.IsHybrid = cmd.IsHybrid()
	appData.IsRemote = cmd.IsRemote()

	appDate, err := cmd.GetApplicationDate(dateFormat)
	if err != nil {
		fmt.Println("Failed ot parse input: ", err)
	}
	appData.ApplicationDate = appDate.Local().Format(dateFormat)

	respDate, err := cmd.GetResponseDate(dateFormat)
	if err != nil {
		fmt.Println("Failed to parse input: ", err)
	}
	appData.ResponseDate = respDate.Local().Format(dateFormat)

	appData.Response = cmd.GetResponse()

	appData.Comment = cmd.GetComment()

	return &appData
}
