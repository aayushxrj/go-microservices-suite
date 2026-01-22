package main

import (
	"broker/cmd/event"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
	case "log":
		log.Println("HandleSubmission: log called")
		// app.logItem(w, requestPayload.Log)
		app.logEventViaRabbit(w, requestPayload.Log)
	case "mail":
		log.Println("HandleSubmission: mail called")
		app.sendMail(w, requestPayload.Mail)
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

func (app *Config) logItem(w http.ResponseWriter, l LogPayload) {
	log.Println("logItem called")

	// create some json we will send to the logger microservice
	jsonData, _ := json.MarshalIndent(l, "", "\t")

	// call the logger microservice
	req, err := http.NewRequest("POST", "http://logger-service:8083/log", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, fmt.Errorf("error calling logger service"))
		return
	}

	var jsonFromService jsonResponse

	jsonFromService.Error = false
	jsonFromService.Message = "logged"

	_ = app.writeJSON(w, http.StatusAccepted, jsonFromService)

}

func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {
	log.Println("sendMail called")

	// create some json we will send to the mail microservice
	jsonData, _ := json.MarshalIndent(m, "", "\t")

	// call the mail microservice
	mailServiceURL := "http://mail-service:8084/send"

	req, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, fmt.Errorf("error calling mail service"))
		return
	}

	var jsonFromService jsonResponse

	jsonFromService.Error = false
	jsonFromService.Message = "Message sent to " + m.To
	_ = app.writeJSON(w, http.StatusAccepted, jsonFromService)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	err = app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		log.Println(err)
	}

}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		log.Println(err)
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "    ")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
