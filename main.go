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
		fmt.Println("AFTER PARSING:\n", responseStruct)
		fmt.Println("TESTING ATTRIBUTES:\n", responseStruct.Position)
	}	
}

type ApiResponse struct {
	Time int `json:"timestamp"`
	Message string `json:"message"`
	Position IssPosition `json:"iss_position"`
}

type IssPosition struct {
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func getPosition (body []byte) (*ApiResponse, error) {
	var response = new(ApiResponse)
	err := json.Unmarshal(body, &response)
	if err != nil {
		handleError (err)
	}

	return response, err
}

func handleError (err error) () {
	fmt.Println("GET request returned an error:\n", err)
}