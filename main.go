package main

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

func main() {
	resp, err := http.Get("http://api.open-notify.org/iss-now.json")

	if err != nil {
		handleError (err)
		os.Exit(1)

	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			handleError (err)
			os.Exit(1)
		}

		fmt.Println(string(contents))
		responseStruct, err := getPosition([]byte(contents))
    fmt.Println("ISS POSITION")
		fmt.Println("Latitude and Longitude:", responseStruct.Position)

    lat := responseStruct.Position.Latitude
    lon := responseStruct.Position.Longitude
    resp2, err2 := http.Get("http://geoservices.tamu.edu/Services/ReverseGeocoding/WebService/v04_01/HTTP/default.aspx?apiKey=183d0aec4c0a4a8e8856e73adce9227d&version=4.10&lat=" + lat + "&lon=" + lon + "&format=json")

    if err2 != nil {
      handleError (err2)
      os.Exit(1)

    } else {
      defer resp.Body.Close()
      contents2, err2 := ioutil.ReadAll(resp2.Body)

      if err2 != nil {
        handleError (err2)
        os.Exit(1)
      }

      fmt.Println("CITY DATA")
      responseStruct2, err2 := getAddress([]byte(contents2))

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

func getPosition (body []byte) (*ApiResponse, error) {
	var response = new(ApiResponse)
	err := json.Unmarshal(body, &response)
	if err != nil {
		handleError (err)
	}

	return response, err
}

func getAddress (body []byte) (*ApiResponse2, error) {
  var response = new (ApiResponse2)
  err := json.Unmarshal(body, &response)
  if err != nil {
    handleError (err)
  }

  return response, err
}

//TODO: better error handling function
func handleError (err error) () {
	fmt.Println("GET request returned an error:\n", err)
}
