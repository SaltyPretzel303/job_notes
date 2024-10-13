package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akamensky/argparse"
)

const NAME = "Job Application Manager"
const DESCRIPTION = "This utils allows easier management of job applications."
const DATE_FORMAT = time.DateTime

type ApplicationData struct {
	Company         string
	Role            string
	Location        string
	WorkType        string // on-site/hybrid/remote
	IsHybrid        bool
	IsRemote        bool
	ApplicationDate string
	ResponseDate    string
	Response        string // ACCEPTED_RESPONSE || DENIED_RESPONSE
	Comment         string
}

type SearchData struct {
	Company       string
	Role          string
	Location      string
	WorkType      []string
	IsHybrid      bool
	IsRemote      bool
	OlderThanDate string
}

type Db interface {
	Save(ap *ApplicationData) *ApplicationData
	Search(data *SearchData) []ApplicationData
	SearchCompany(comp string) []ApplicationData
	SearchRole(role string) []ApplicationData
	SearchLocation(location string) []ApplicationData
	SearchOlderThanDays(daysOld int) []ApplicationData
	SearchDenied() []ApplicationData
	SearchAccepted() []ApplicationData
	Close()
}

func (ap *ApplicationData) AsJson() string {
	bytes, _ := json.MarshalIndent(ap, "", "	")

	return string(bytes)
}

func strip(s string, toStrip string) (string, bool) {
	value := strings.Replace(s, toStrip, "", -1)

	return value, (value != s)
}

func parseShortDate(s string, format string) string {

	if len(s) == 0 {
		return ""
	}

	tokens := strings.Split(s, ".")
	day, _ := strconv.Atoi(tokens[0])
	month, _ := strconv.Atoi(tokens[1])
	year := 2024
	if len(tokens) > 2 {
		year, _ = strconv.Atoi(tokens[2])
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).
		Local().Format(format)

}

// TODO return err as well
func ParseAppData(tokens []string) *ApplicationData {
	company := tokens[0]
	role := tokens[1]
	llocation := tokens[2]
	noRemote, isRemote := strip(llocation, "remote")
	location, isHybrid := strip(noRemote, "hybrid")
	apDate := parseShortDate(tokens[3], DATE_FORMAT)
	resDate := parseShortDate(tokens[4], DATE_FORMAT)
	resp := tokens[5]
	comment := tokens[6]

	return &ApplicationData{
		Company:         company,
		Role:            role,
		Location:        location,
		WorkType:        tokens[2],
		IsHybrid:        isHybrid,
		IsRemote:        isRemote,
		ApplicationDate: apDate,
		ResponseDate:    resDate,
		Response:        resp,
		Comment:         comment,
	}
}

func NowAsString() string {
	return time.Now().Format(DATE_FORMAT)
}

func main() {

	var parser = argparse.NewParser(NAME, DESCRIPTION)

	var addFields = SetupAddCmd(parser)
	var getFields = SetupGetCmd(parser)
	var importFields = SetupImportCmd(parser)

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Println("Failed to parser arguments: ", err)
		os.Exit((1))
	}

	db := NewLocalDb()
	defer db.Close()

	if addFields.Happened() {
		newApp := addFields.AsApplication(DATE_FORMAT)

		fmt.Printf("Will save the next item: \n %v\n", newApp.AsJson())
		fmt.Print("Enter Y to proceed, anything else to cancel: ")

		var input string
		_, err := fmt.Scanln(&input)

		if err != nil || strings.ToUpper(input) != "Y" {
			fmt.Println("Aborting.")
		} else {
			fmt.Println("Will save the application.")
			fmt.Println(newApp.AsJson())
			db.Save(newApp)
		}
	} else if getFields.Happened() {
		aps := db.Search(getFields.AsSearchData(DATE_FORMAT))
		fmt.Println("RESULTS")
		fmt.Println("--------------")

		for _, ap := range aps {
			fmt.Println(ap.AsJson())
		}

		fmt.Println("--------------")

	} else if importFields.Happened() {

		file := importFields.GetFile()

		f, err := os.Open(file)

		if err != nil {
			fmt.Println("Failed to open file: ", file, " err: ", err)
			return
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)

		scanner.Scan() // will step over the first/header row
		data := []ApplicationData{}

		for scanner.Scan() {
			line := scanner.Text()

			if line == "end,,,,,,,," {
				fmt.Print("Found the end.")
				break
			}

			dataLine := ParseAppData(strings.Split(line, ","))
			data = append(data, *dataLine)

			fmt.Printf("%+v \n", data)
		}

		if len(data) > 0 {
			fmt.Println("Saving parsed records.")
			for _, item := range data {
				db.Save(&item)
			}
		}
	}

	fmt.Println("Leaving.")

}
