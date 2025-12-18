package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	log.Println("/ called")

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	err := app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		log.Println(err)
	}
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	log.Println("/handle called")

	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		log.Println("HandleSubmission: auth called")
		app.authenticate(w, requestPayload.Auth)
	default:
		log.Println("HandleSubmission: default called")
		app.errorJSON(w, fmt.Errorf("unknown action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	log.Println("authenticate called")

	// create some json we will send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the auth microservice
	req, err := http.NewRequest("POST", "http://authentication-service:8082/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	// handle response from auth microservice
	if resp.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, fmt.Errorf("invalid credentials"))
		return
	} else if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, fmt.Errorf("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(resp.Body).Decode(&jsonFromService)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse 
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	err = app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		log.Println(err)
	}
}