package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func uploadFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileName := filepath.Base(path)
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("mediafile", fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	written, err := io.Copy(part, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Wrote %v kb\n", written>>10)

	part, err = writer.CreateFormFile("metadata", "test.json")
	io.Copy(part, strings.NewReader(`{"jsonKey": "this is some json"}`))

	// don't forget to close the writer before you make the request dingus
	writer.Close()

	req, _ := http.NewRequest("POST", "http://localhost:8081/file", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("Error uploading")
		fmt.Println(err)
	}
}
func uploadMetadata(mf *MediaFile) {
	json, _ := json.Marshal(*mf)
	req, err := http.NewRequest("POST", "http://localhost:8081/mediafile", bytes.NewBuffer(json))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
