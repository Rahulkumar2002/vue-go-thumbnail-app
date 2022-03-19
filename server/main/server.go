package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type screenshotAPIRequest struct {
	Token          string `json:"token"`
	Url            string `json:"url"`
	Output         string `json:"output"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	ThumbnailWidth int    `json:"thumbnail_width"`
}

type thumbnailRequest struct {
	Url string `json:"url"`
}

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	var decoded thumbnailRequest

	err := json.NewDecoder(r.Body).Decode(&decoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Set env variable token.
	os.Setenv("token", "SCREEN-SHOT-API-KEY")

	// Create a struct to call the ScreenShot api with the following parameters.
	apiRequest := screenshotAPIRequest{
		Token:          os.Getenv("token"),
		Url:            decoded.Url,
		Output:         "json",
		Width:          1920,
		Height:         1080,
		ThumbnailWidth: 500,
	}

	// Convert the struct to json string.
	jsonString, err := json.Marshal(apiRequest)
	checkError(err)

	// Creating a Http Request.

	req, err := http.NewRequest("POST", "https://shot.screenshotapi.net/screenshot", bytes.NewBuffer(jsonString))
	checkError(err)
	req.Header.Set("Content-Type", "application/json")

	// Execute the Http Request

	client := &http.Client{}
	response, err := client.Do(req)
	checkError(err)

	// Tell go to close the resposnse body
	defer response.Body.Close()

	// Read the raw response into a Go struct.
	type screenshotAPIResponse struct {
		Screenshot string `json:"screenshot"`
	}

	var apiResponse screenshotAPIResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	checkError(err)

	// Passing data back to the frontend.

	_, err = fmt.Fprintf(w, `{"screenshot": "%v"}`, apiResponse.Screenshot)
	checkError(err)
}

func main() {
	// Serving static files from ../../frontend/dist directory.
	fs := http.FileServer(http.Dir("../../frontend/dist"))
	http.Handle("/", fs)
	// Using the thumbnailHandler function.
	http.HandleFunc("/api/thumbnail", thumbnailHandler)
	// Starting server
	fmt.Println("Server is listing at Port 3000")
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}
