package runner

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var callbackURL string

func notifyFlowStart(flowUUID string) {
	formData := url.Values{
		"flow_uuid": {flowUUID},
	}

	resp, err := http.PostForm(callbackURL+"/start", formData)
	if err != nil {
		log.Printf("Flow %s failed to send start notification: %s", flowUUID, err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Flow %s failed to send start notification: %s", flowUUID, err.Error())
	}
	log.Printf("Flow %s sent started notification: %s", flowUUID, body)
}
func notifyFlowSuccess(flowUUID string, downloadURL string) {
	formData := url.Values{
		"flow_uuid":   {flowUUID},
		"storage_url": {downloadURL},
	}

	resp, err := http.PostForm(callbackURL+"/success", formData)
	if err != nil {
		log.Printf("Flow %s failed to send success notification: %s", flowUUID, err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Flow %s failed to send success notification: %s", flowUUID, err.Error())
	}
	log.Printf("Flow %s sent success notification: %s", flowUUID, body)
}
func notifyFlowFail(flowUUID string, message string) {
	formData := url.Values{
		"flow_uuid": {flowUUID},
		"message":   {message},
	}

	resp, err := http.PostForm(callbackURL+"/fail", formData)
	if err != nil {
		log.Printf("Flow %s failed to send fail notification: %s", flowUUID, err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Flow %s failed to send fail notification: %s", flowUUID, err.Error())
	}
	log.Printf("Flow %s sent fail notification: %s", flowUUID, body)
}
