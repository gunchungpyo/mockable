package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Mock object for store mock information
type Mock struct {
	ContentType string `json:"content_type"`
	Source      string `json:"source"`
	Endpoint    string `json:"endpoint"`
	StatusCode  int    `json:"status_code"`
}

// ServerConfig object to store server configuration
type ServerConfig struct {
	Port string `json:"port"`
}

// ErrorResponse object for error handling
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var mocks []Mock

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path
	sourcePath := "example/mock/"
	var mock Mock
	for i := range mocks {
		if mocks[i].Endpoint == endpoint {
			mock = mocks[i]
			break
		}
	}
	sourcePath = sourcePath + mock.Source
	file, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		errorResponse := ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()}
		errorJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorJSON)
		log.Println(errorResponse)
		return
	}
	w.Header().Set("Content-Type", mock.ContentType)
	w.WriteHeader(mock.StatusCode)
	w.Write(file)
	log.Println("endpoint called:", endpoint)
}

func main() {
	file, _ := ioutil.ReadFile("config/server.json")
	config := ServerConfig{}
	_ = json.Unmarshal([]byte(file), &config)

	// set up router
	file, _ = ioutil.ReadFile("config/endpoints.json")
	_ = json.Unmarshal([]byte(file), &mocks)
	for i := 0; i < len(mocks); i++ {
		fmt.Println("Register route ", mocks[i].Endpoint)
		http.HandleFunc(mocks[i].Endpoint, jsonHandler)
	}

	fmt.Println("Server started using port", config.Port)
	err := http.ListenAndServe(config.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
