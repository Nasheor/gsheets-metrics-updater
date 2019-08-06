package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Prints the number of frontend bugs, backend bugs and bugs yet to be classified
func printBugsStats(frontEndBugs int, backEndBugs int, totalBugs int) {
	fmt.Println("Number of Frontend bugs: ", frontEndBugs)
	fmt.Println("Number of Backend bugs: ", backEndBugs)
	fmt.Println("Number of bugs yet to be classified: ", totalBugs-(frontEndBugs+backEndBugs))
}

// Returns three strings - year, month and day of the month
func getYearMonthDay(currentDateTime time.Time) (string, string, string) {
	//Splitting the string "2019-06-19 00:00:00" using "-" as delimiter
	formattedDate := strings.Split(currentDateTime.String(), "-")

	//Splitting the string "19 00:00:00" using " " as delimiter and returning the value at index 0
	formattedDate[2] = strings.Split(formattedDate[2], " ")[0]

	// Returned in Year, Month and Day of Month Format
	return formattedDate[0], formattedDate[1], formattedDate[2]
}

// Generates the URL of the endpoint that needs to be hit
func getEndPoint(fromDate, toDate string) string {
	// Setting the Base URL
	endpoint, _ := url.Parse("https://digitalcrew.teamwork.com/projects/api/v2/projects/263073/tasks.json")

	// Setting the query params
	params := url.Values{}

	params.Add("includeBlockedTasks", "true")
	params.Add("ignoreStartDates", "false")
	params.Add("includeCompletedTasks", "true")
	params.Add("tagIds", "3871")
	params.Add("onlyUntaggedTasks", "false")
	params.Add("matchAllTags", "true")
	params.Add("matchAllExcludedTags", "false")
	params.Add("createdAfterDate", fromDate)
	params.Add("createdBeforeDate", toDate)
	params.Add("createdFilter", "custom")

	//Adding the query params to the base url
	endpoint.RawQuery = params.Encode()
	return endpoint.String()
}

//GET request to the projects API for information about tassk
func getDataFromEndpoint(endpoint string) Data {
	// Setting the type of request and authorization data
	req, err := http.NewRequest("GET", endpoint, nil)
	token, _ := ioutil.ReadFile("../authentication/accesstoken.txt")
	encodedToken := base64.StdEncoding.EncodeToString([]byte(string(token)))
	basicAuth := "Basic " + encodedToken
	req.Header.Set("Authorization", basicAuth)

	//Executing the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal("Error in the request.\n[ERRO] -", err)
	}
	defer resp.Body.Close()

	// Getting json data, decoding them from the raw format and storing them at "record's" address
	var record Data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	return record
}

//Counting the number of tags with a particular tag name
func getBugsCount(team string, data Data) int {
	count := 0
	for i := range data.Tasks {
		task := data.Tasks[i]
		for j := range task.Tags {
			tag := task.Tags[j]
			if tag.Name == team {
				count = count + 1
			}
		}
	}
	return count
}

//Get Scorecard info to access
func getScorecardInfo(team string) (string, string) {
	// Setting the id of the spreadsheet
	spreadsheetID := ""
	measurable := ""
	if team == "frontend" {
		spreadsheetID = "1igV-OBcSFOj6MrNMtsxKYht56_7WeHnXOUMp9Acs-Bs"
		measurable = "Number of bugs generated for frontened in the last 7 days"
	} else {
		spreadsheetID = "1igV-OBcSFOj6MrNMtsxKYht56_7WeHnXOUMp9Acs-Bs"
		measurable = "Number of bugs generated for backend in the last 7 days"
	}
	return spreadsheetID, measurable
}

//Get the week number from the start of the year
func getCurrentWeekNumber() int {
	now := time.Now().UTC()
	_, week := now.ISOWeek()
	return week
}

func isCurrentWeekPresent(firstRow []interface{}) bool {

	// Flags true if the current week number is already present in the scorecard
	isWeekPresent := false

	for _, element := range firstRow {
		parsedElement, _ := strconv.Atoi(fmt.Sprintf("%v", element))
		if parsedElement == getCurrentWeekNumber() {
			isWeekPresent = true
			break
		}
	}
	return isWeekPresent
}

// Creates the respective week number (current weeek from the start of the year)
func createWeekInScorecard(data [][]interface{}) [][]interface{} {
	newData := make([][]interface{}, 30)
	breakpoint := -1

	if !isCurrentWeekPresent(data[0]) {
		for rowIndex, row := range data {
			for columnIndex, item := range row {
				if item == "Result" {
					j := columnIndex
					breakpoint = j + 1
				}
				if breakpoint == columnIndex || breakpoint == len(row) {
					newRowData := append(make([]interface{}, 0), data[rowIndex][:breakpoint]...)
					if rowIndex == 0 {
						newRowData = append(newRowData, getCurrentWeekNumber())
					} else {
						newRowData = append(newRowData, 0)
					}
					newRowData = append(newRowData, data[rowIndex][breakpoint:]...)
					newData[rowIndex] = newRowData
					break
				}
			}
		}
		return newData
	}
	return data
}

// If the measruable "Number of bugs" isn't already present, create the measurable
func createMeasurable(data [][]interface{}, measurableData string) [][]interface{} {
	measurable := []interface{}{" ", measurableData, "0%", "0", "0.00%", " ", " ", " ", " ", " "}
	if isCurrentWeekPresent(data[0]) {
		measurable = []interface{}{" ", measurableData, "0%", "0", "0.00%", "0", " ", " ", " ", " ", " "}
	}
	dataWithMeasurable := make([][]interface{}, 50)
	count := 0
	measurablePresent := false
	for _, row := range data {
		if count == (len(data) - 2) {
			dataWithMeasurable[count] = row
			dataWithMeasurable[count+1] = measurable
			measurablePresent = true
		} else {
			if measurablePresent {
				dataWithMeasurable[len(data)] = row
			} else {
				dataWithMeasurable[count] = row
			}
		}
		count++
	}
	return dataWithMeasurable
}

//Writes to the Frontend or Backend scorecard
func writeToScoreCard(team string, bugs int) {
	b, err := ioutil.ReadFile("../authentication/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Setting the id and measurable of the spreadsheet
	spreadsheetID, measurable := getScorecardInfo(team)

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, "A1:Z100").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// Holds the current data present in the spreadsheet
	rawData := resp.Values
	var vr sheets.ValueRange

	isMeasurablePresent := false
	for _, row := range rawData {
		for _, cell := range row {
			if cell == measurable {
				isMeasurablePresent = true
				break
			}
		}
	}
	dataToWrite := rawData
	if !isMeasurablePresent {
		fmt.Println("Measurable: '", measurable, "' not present. Creating measurable now...")
		dataToWrite = createMeasurable(dataToWrite, measurable)
		fmt.Println("Measurable Created")
	} else {
		fmt.Println("Measurable already present")
	}
	if !isCurrentWeekPresent(dataToWrite[0]) {
		fmt.Println("Creating Current week: ", getCurrentWeekNumber())
		dataToWrite = createWeekInScorecard(dataToWrite)
		fmt.Println("Current week created")
	} else {
		fmt.Println("Week already present")
	}

	for rowIndex, row := range dataToWrite {
		for _, item := range row {
			if item == measurable {
				dataToWrite[rowIndex][5] = bugs
				break
			}

		}
	}
	writeRange := "A1"
	vr.Values = dataToWrite[:(len(dataToWrite) - 1)]
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, &vr).ValueInputOption("RAW").ResponseValueRenderOption("UNFORMATTED_VALUE").Do()

	writeRange = "A" + strconv.Itoa(len(dataToWrite))

	vr.Values = dataToWrite[(len(dataToWrite) - 1):]
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, &vr).ValueInputOption("RAW").Do()

}

func main() {

	//Getting the current date and seperating the day of month, month and year
	currentDateTime := time.Now()
	fromDateTime := currentDateTime.AddDate(0, 0, -7)
	currentYear, currentMonth, currentDayOfMonth := getYearMonthDay(currentDateTime)
	fromYear, fromMonth, fromDayOfMonth := getYearMonthDay(fromDateTime)

	//Form the createdBeforeDate and createdAfterDate params yyyymmdd
	fromDate := fromYear + fromMonth + fromDayOfMonth
	toDate := currentYear + currentMonth + currentDayOfMonth

	//Get data from the end point
	endpoint := getEndPoint(fromDate, toDate)
	data := getDataFromEndpoint(endpoint)
	frontEndBugs := getBugsCount("frontend", data)
	backEndBugs := getBugsCount("backend", data)
	totalBugs := len(data.Tasks)

	fmt.Println("Total Bugs: ", totalBugs)
	fmt.Println("Number of Frontend bugs: ", frontEndBugs)
	fmt.Println("Number of Backend bugs: ", backEndBugs)
	fmt.Println("Writing to Frontend Scorecard")
	writeToScoreCard("frontend", frontEndBugs)
	fmt.Println("Frontend scorecard updated")
	fmt.Println("Writing to Backend Scorecard")
	writeToScoreCard("backend", backEndBugs)
}
