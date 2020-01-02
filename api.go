package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type dupeMessage struct {
	Message string `json:"message"`
}

func DupeCheck(sha256 []byte) (bool, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:8081/dupe/%v", sha256), nil)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	var dm dupeMessage
	err = json.Unmarshal(body, &dm)
	if err != nil {
		return false, err
	}

	if dm.Message == "IS_DUPE" {
		return true, nil
	} else {
		return false, nil
	}

}
