package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"data-explorer-convertor/models"
)

var order = 1

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Please provide input, output filenames and title as arguments.")
		os.Exit(1)
	}

	file1 := os.Args[1]
	file2 := os.Args[2]
	title := os.Args[3]

	err := checkFileExists(file1)
	if err != nil {
		fmt.Println("Error: Input file does not exist.")
		os.Exit(1)
	}

	convertToDataExplorer(file1, file2, title)
}

func checkFileExists(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

func convertToDataExplorer(inputFile string, outputFile string, dashboardTitle string) {
	inputData, _ := os.ReadFile(inputFile)
	var grafanaDashboard models.GrafanaDashboard
	json.Unmarshal(inputData, &grafanaDashboard)

	outputJson := make(map[string]interface{})
	header := make(map[string]interface{})
	outputJson["header"] = header
	header["dateTimeRange"] = true
	header["dropdowns"] = []interface{}{}

	tabs := []models.Tab{}

	tab := models.Tab{}
	tab.Key = dashboardTitle
	tab.Order = 1
	tab.Title = dashboardTitle
	tab.Type = "metrics"

	queryList := []models.QueryItem{}

	for _, panel := range grafanaDashboard.Panels {
		getExprsFromPanel(panel, &queryList)
	}

	tab.QueriesList = queryList

	tabs = append(tabs, tab)
	outputJson["tabs"] = tabs

	byteData, _ := json.MarshalIndent(outputJson, "", "    ")

	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileData := strings.ReplaceAll(string(byteData), "\\u0026", "&")
	_, err = f.WriteString(fileData)
	if err != nil {
		fmt.Println(err)
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getExprsFromPanel(panel models.Panel, queryList *[]models.QueryItem) {
	for _, target := range panel.Targets {
		title := ""
		if target.QueryExpr != nil {
			if len(panel.Targets) > 1 {
				if target.Legend != nil {
					// If "legendFormat" exists and is a string, split it
					legendFormat := *target.Legend
					legendSplit := strings.Split(legendFormat, " - ")
					title = panel.Title + " - " + legendSplit[len(legendSplit)-1] + " query"
				}
			} else {
				title = panel.Title + " query"
			}

			query := models.QueryItem{}
			query.Name = title
			if panel.IsLines != nil && *panel.IsLines {
				*query.ChartType = "line"
			}
			// query["data_source_name"] = "Apica Monitoring" //

			options := models.QueryOptions{}
			options.Description = title
			options.Order = order
			order += 1
			options.Plot = models.NewQueryPlot()
			options.UpperLimit = ""
			query.Options = options

			queryExpr := *target.QueryExpr
			queryExpr = strings.ReplaceAll(queryExpr, "$namespace", "{{namespace}}")
			if queryExpr != strings.ReplaceAll(queryExpr, ",service=~\"$service\"", "") {
				queryExpr = strings.ReplaceAll(queryExpr, ",service=~\"$service\"", "")
				queryExpr = strings.ReplaceAll(queryExpr, ", quantile=", " quantile=")
				queryExpr = strings.ReplaceAll(queryExpr, ",quantile=", "quantile=")
			} else if queryExpr != strings.ReplaceAll(queryExpr, ", service=~\"$service\"", "") {
				queryExpr = strings.ReplaceAll(queryExpr, ", service=~\"$service\"", "")
				queryExpr = strings.ReplaceAll(queryExpr, ", quantile=", " quantile=")
				queryExpr = strings.ReplaceAll(queryExpr, ",quantile=", "quantile=")
			}
			queryExpr += "&duration=1h&step=5m"
			query.Query = queryExpr
			query.Schema = ""

			*queryList = append(*queryList, query)
		}
	}

	if panel.Panels != nil {
		for _, p := range *panel.Panels {
			getExprsFromPanel(p, queryList)
		}
	}
}
