package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	type requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var payload requestPayload
	err := app.readJSON(w, r, &payload)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(payload.Email)
	if err != nil {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(payload.Password)
	if err != nil || !valid {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log authentication using logger-service
	err = app.logRequest("authentication", fmt.Sprintf("User %s logged in", user.Email))
	if err != nil {
		_ = app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payloadResponse := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in as %s", user.Email),
		Data:    user,
	}
	_ = app.writeJSON(w, http.StatusAccepted, payloadResponse)
}

func (app *Config) logRequest(name, data string) error {
	var payload struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	payload.Name = name
	payload.Data = data

	jsonData, _ := json.MarshalIndent(payload, "", "\t")
	logServiceURL := "http://logger-service:8083/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return errors.New("error calling logger service")
	}
	return nil
}
