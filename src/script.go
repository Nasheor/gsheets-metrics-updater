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

var column []string

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
	params.Add("include", "taskListNames")
	params.Add("includeCompletedTasks", "true")
	params.Add("tagIds", "3871")
	params.Add("matchAllTags", "true")
	params.Add("matchAllExcludedTags", "false")
	params.Add("createdAfterDate", fromDate)
	params.Add("createdBeforeDate", toDate)

	//Adding the query params to the base url
	endpoint.RawQuery = params.Encode()
	return endpoint.String()
}

//GET request to the projects API for information about tassk
func getDataFromEndpoint(endpoint string) Data {
	// Setting the type of request and authorization data
	req, err := http.NewRequest("GET", endpoint, nil)
	token, _ := ioutil.ReadFile("accesstoken.txt")
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
	taskListNames := []int{1137238, 1138697}
	count := 0
	validTagPresent := false
	for i := range data.Tasks {
		validTagPresent = false
		task := data.Tasks[i]
		for j := range task.Tags {
			tag := task.Tags[j]
			if tag.Name == team {
				count = count + 1
				validTagPresent = true
			}
		}
		if !validTagPresent {
			if task.TaskListID == taskListNames[0] && team == "frontend" {
				count = count + 1
			} else if task.TaskListID == taskListNames[1] && team == "backend" {
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
		spreadsheetID = "Frontend Sheet Identifier"
		measurable = "Number of bugs generated in the last 7 days"
	} else {
		spreadsheetID = "Backend Sheet Identifier"
		measurable = "Number of bugs generated in the last 7 days"
	}
	return spreadsheetID, measurable
}

//Check if the measurable is present in the scorecard
func checkMeasurablePresent(data [][]interface{}, measurable string) bool {
	isMeasurablePresent := false
	for _, row := range data {
		for _, cell := range row {
			if cell == measurable {
				isMeasurablePresent = true
				break
			}
		}
	}
	return isMeasurablePresent
}

//Create missing data in the scorecard i.e., current week or measurable or both
func createData(client *http.Client, dataToWrite string, identifier string, data [][]interface{}, spreadsheetID string, sheetID int) {
	srv, _ := sheets.New(client)
	startIndex := -1
	endIndex := -1
	var majorDimension string

	// Deciding where to insert the empty row or column based on the row or column identifier
	for rowIndex, row := range data {
		for columnIndex, item := range row {
			if identifier == "Total healthscore" && item == identifier {
				startIndex = rowIndex
				endIndex = rowIndex + 1
				majorDimension = "ROWS"
				break
			} else if identifier == "Result" && item == identifier {
				startIndex = columnIndex + 1
				endIndex = columnIndex + 2
				majorDimension = "COLUMNS"
				break
			}
		}
		if startIndex != -1 {
			break
		}
	}

	// Inserting the row or column
	var idr sheets.InsertDimensionRequest
	var r sheets.Request
	var d sheets.DimensionRange

	d.SheetId = int64(sheetID)
	d.Dimension = majorDimension
	d.EndIndex = int64(endIndex)
	d.StartIndex = int64(startIndex)

	idr.Range = &d

	r.InsertDimension = &idr

	requests := []*sheets.Request{&r}

	rb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}
	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, rb).Do()

	if err != nil {
		fmt.Println(err)
	}

	var writeRange string
	var tmp []interface{}
	if majorDimension == "ROWS" {
		writeRange = "A" + strconv.Itoa(startIndex+1) + ":Z" + strconv.Itoa(endIndex)
		tmp = []interface{}{" ", dataToWrite, "0%", "0", "0.00%", " ", " ", " ", " ", " "}
	} else {
		writeRange = column[startIndex] + "1:" + column[startIndex] + "2"
		tmp = []interface{}{dataToWrite}
	}
	metrics := [][]interface{}{tmp}
	var vr sheets.ValueRange
	vr.MajorDimension = majorDimension
	vr.Range = writeRange
	vr.Values = metrics

	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, &vr).ValueInputOption("RAW").ResponseValueRenderOption("FORMATTED_VALUE").Do()
}

//Get the week number from the start of the year
func getCurrentWeekNumber() int {
	now := time.Now().UTC()
	_, week := now.ISOWeek()
	return week
}

// Flags true if the current week number is already present in the scorecard
func isCurrentWeekPresent(firstRow []interface{}) bool {
	isWeekPresent := false
	weekNumber := getCurrentWeekNumber()
	for _, element := range firstRow {
		parsedElement, _ := strconv.Atoi(fmt.Sprintf("%v", element))
		if parsedElement == weekNumber {
			isWeekPresent = true
			break
		}
	}
	return isWeekPresent
}

// Updates the scorecard with the bugs count
func updateScoreCard(client *http.Client, bugs int, data [][]interface{}, spreadsheetID string) {
	srv, _ := sheets.New(client)
	var c string
	for _, row := range data {
		for columnIndex, item := range row {
			if item == strconv.Itoa(getCurrentWeekNumber()) {
				c = column[columnIndex]
				break
			}
		}
		if len(c) > 1 {
			break
		}
	}

	index := len(data) - 1
	writeRange := c + strconv.Itoa(index)
	tmp := []interface{}{bugs}
	weeklyBugs := [][]interface{}{tmp}
	var vr sheets.ValueRange
	vr.MajorDimension = "COLUMNS"
	vr.Range = writeRange
	vr.Values = weeklyBugs

	_, err := srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, &vr).ValueInputOption("RAW").ResponseValueRenderOption("FORMATTED_VALUE").Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Scorecard updated at ", writeRange)

}

//Writes to the Frontend or Backend scorecard
func writeToScoreCard(team string, bugs int, sheetID int) {
	b, err := ioutil.ReadFile("credentials.json")
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
	data := resp.Values
	column = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	isMeasurablePresent := checkMeasurablePresent(data, measurable)
	var identifier string

	if !isMeasurablePresent {
		fmt.Println("Measurable: '", measurable, "' not present. Creating measurable now...")
		identifier = "Total healthscore"
		createData(client, measurable, identifier, data, spreadsheetID, sheetID)
		fmt.Println("Measurable Created")
	} else {
		fmt.Println("Measurable already present")
	}

	currentWeek := strconv.Itoa(getCurrentWeekNumber())
	if !isCurrentWeekPresent(data[0]) {
		identifier = "Result"
		fmt.Println("Creating Current week: ", currentWeek)
		createData(client, currentWeek, identifier, data, spreadsheetID, sheetID)
		fmt.Println("Current week created")
	} else {
		fmt.Println("Current Week: ", currentWeek, " already present")
	}
	fmt.Println("Updating Scorecard...")
	resp, err = srv.Spreadsheets.Values.Get(spreadsheetID, "A1:Z100").Do()
	data = resp.Values

	updateScoreCard(client, bugs, data, spreadsheetID)
}

func main() {

	// Getting the current date and seperating the day of month, month and year
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

	fmt.Println("Total Bugs in the last 7 days: ", totalBugs)
	fmt.Println("Number of Frontend bugs in the last 7 days: ", frontEndBugs)
	fmt.Println("Number of Backend bugs in the last 7 days: ", backEndBugs)
	fmt.Println("Writing to Frontend Scorecard")
	writeToScoreCard("frontend", frontEndBugs, 11111111)
	fmt.Println("Frontend scorecard updated")
	fmt.Println("Writing to Backend Scorecard")
	writeToScoreCard("backend", backEndBugs, 111111111)
}
