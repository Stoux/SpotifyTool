package handlers

import (
	json2 "encoding/json"
	"log"
	"net/http"
)

func Json(handler func(writer http.ResponseWriter, request *http.Request) (result interface{}, shouldOutput bool)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var err error

		// Get the response from the server
		if result, output := handler(writer, request); output {
			// Our response is going to be JSON
			writer.Header().Set("Content-Type", "application/json")

			// Marshall it & write as response
			if json, err := json2.Marshal(result); err == nil {
				_, err = writer.Write(json)
			}
		}

		// Check if successful
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Println("[Error] Failed to marshall / write JSON response: ", err)
		}
	}
}

func JsonWithOutput(handler func(writer http.ResponseWriter, request *http.Request) (result interface{})) func(writer http.ResponseWriter, request *http.Request) {
	return Json(func(writer http.ResponseWriter, request *http.Request) (result interface{}, shouldOutput bool) {
		return handler(writer, request), true
	})
}
