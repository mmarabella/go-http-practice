package main

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

const issApiUrl string = "http://api.open-notify.org/iss-now.json"
const geocodeApiUrl string = "http://geoservices.tamu.edu/Services/ReverseGeocoding/WebService/v04_01/HTTP/default.aspx?apiKey=183d0aec4c0a4a8e8856e73adce9227d&version=4.10&lat="

func main() {
	resp, getPositionErr := http.Get(issApiUrl)

	if getPositionErr != nil {
		handleError (getPositionErr)
		os.Exit(1)

	} else {
		defer resp.Body.Close()
		contents, readPositionErr := ioutil.ReadAll(resp.Body)

		if readPositionErr != nil {
			handleError (readPositionErr)
			os.Exit(1)
		}

		fmt.Println(string(contents))

    // TODO: handle error?
		responseStruct, _ := unmarshalPosition([]byte(contents))
    fmt.Println("ISS POSITION")
		fmt.Println("Latitude and Longitude:", responseStruct.Position)

    lat := responseStruct.Position.Latitude
    lon := responseStruct.Position.Longitude
    resp2, getCityErr := http.Get(geocodeApiUrl + lat + "&lon=" + lon + "&format=json")

    if getCityErr != nil {
      handleError (getCityErr)
      os.Exit(1)

    } else {
      defer resp.Body.Close()
      contents2, readCityErr := ioutil.ReadAll(resp2.Body)

      if readCityErr != nil {
        handleError (readCityErr)
        os.Exit(1)
      }

      fmt.Println("CITY DATA")

      // TODO: handle error?
      responseStruct2, _ := unmarshalAddress([]byte(contents2))

      if responseStruct2.QueryStatusCode == "Unknown" {
        fmt.Println("The space station is not currently over a US city")
      } else {
        fmt.Println("The space station is over", responseStruct2.StreetAddresses[0].City, responseStruct2.StreetAddresses[0].State)
      }
    }
  }
}

// structs for json parsing

type ApiResponse struct {
	Time int `json:"timestamp"`
	Message string `json:"message"`
	Position IssPosition `json:"iss_position"`
}

type IssPosition struct {
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type ApiResponse2 struct {
  TransactionId string `json:"TransactionId"`
  Version string `json:"Version"`
  QueryStatusCode string `json:"QueryStatusCode"`
  TimeTaken string `json:"TimeTaken"`
  Exception string `json:"Exception"`
  ErrorMessage string `json:"ErrorMessage"`
  StreetAddresses []StreetAddresses `json:"StreetAddresses"`
}

type StreetAddresses struct {
  TransactionId string `json:"TransactionId"`
  Version string `json:"Version"`
  QueryStatusCode string `json:"QueryStatusCode"`
  TimeTaken string `json:"TimeTaken"`
  Exception string `json:"Exception"`
  ErrorMessage string `json:"ErrorMessage"`
  APN string `json:"APN"`
  StreetAddress string `json:"StreetAddress"`
  City string `json:"City"`
  State string `json:"State"`
  Zip string `json:"Zip"`
  ZipPlus4 string `json:"ZipPlus4"`
}

// functions to unmarshal JSON

func unmarshalPosition (body []byte) (*ApiResponse, error) {
	var response = new(ApiResponse)
	err := json.Unmarshal(body, &response)
	if err != nil {
		handleError (err)
	}

	return response, err
}

func unmarshalAddress (body []byte) (*ApiResponse2, error) {
  var response = new (ApiResponse2)
  err := json.Unmarshal(body, &response)
  if err != nil {
    handleError (err)
  }

  return response, err
}

//TODO: better error handling function
func handleError (err error) () {
	fmt.Println("Error encountered:\n", err)
}
