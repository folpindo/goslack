package goslack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func get(url string) (response []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("API GET '%s' failed with code %d", url, resp.StatusCode)
		return
	}

	response, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

func post(url string, body url.Values) (response []byte, err error) {
	resp, err := http.PostForm(url, body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("API POST '%s' failed with code %d", url, resp.StatusCode)
		return
	}

	response, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return

}

func getWebsocketURL(token string) (wsurl string, id string, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	body, err := get(url)
	if err != nil {
		return
	}

	var respObj RtmStart
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return
	}

	if !respObj.Success() {
		err = fmt.Errorf("Slack error: %s", respObj.GetError())
		return
	}

	wsurl = respObj.URL
	id = respObj.Self.ID
	return
}
