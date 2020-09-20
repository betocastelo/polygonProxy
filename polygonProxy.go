package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/betocastelo/polygonProxy/dataModel"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var apiKey = flag.String("a", "", "Polygon API key")
var apiKeyString string
var ticker = flag.String("t", "", "Ticker")

var callWaitTime = 1500 * time.Millisecond

func main() {
	flag.Parse()
	apiKeyString = "?apiKey=" + *apiKey
	today := time.Now()
	startDate := today.Add(-2 * 24 * 365 * time.Hour)

	for date := startDate; date.Before(today); date = date.Add(24*time.Hour) {
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			continue
		}
		data := getOpenClose(getDateString(date))
		printData(data)
		time.Sleep(callWaitTime)
	}
}

func printData(data dataModel.OpenClose) {
	fmt.Printf("%s,%d,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f\n", data.Date, data.Volume, data.Open, data.Close, data.High, data.Low, data.PreMarket, data.AfterHours)
}

func getDateString(date time.Time) string {
	year, month, day := date.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func getOpenClose(dateString string) dataModel.OpenClose {
	endPoint := "https://api.polygon.io/v1/open-close"
	url := endPoint + "/" + *ticker + "/" + dateString + apiKeyString

	response, err := http.Get(url)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		if response.StatusCode == 429 {
			log.Printf("Too many requests. Extra wait for %s", dateString)
			time.Sleep(5*callWaitTime)
		}
		time.Sleep(callWaitTime)
		log.Printf("Retrying %s\n", dateString)
		return getOpenClose(dateString)
	}

	respBody, _ := ioutil.ReadAll(response.Body)

	var data dataModel.OpenClose
	_ = json.Unmarshal(respBody, &data)

	if data.Status != "OK" {
		log.Printf("Status for %s to be 'OK', but recieved '%s'", dateString, data.Status)
		data.Date = dateString
	}
	return data
}
