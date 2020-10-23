package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type loginInfo struct {
	ID   string `json:"id"`
	PW   string `json:"pw"`
	Name string `json:"name"`
}

func postIMGreq(url string, filepath string, loginData loginInfo) (*http.Request, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	//multipart start
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//part1 : image
	part1, err := writer.CreateFormFile("file", "send_img")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part1, file)

	//part2 : json
	part2, err := writer.CreateFormFile("user_data", "login_info")
	if err != nil {
		return nil, err
	}
	jsonbytes, err := json.Marshal(loginData)
	if err != nil {
		return nil, err
	}
	jsonReader := bytes.NewReader(jsonbytes)
	_, err = io.Copy(part2, jsonReader)

	//multipart end
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err

}

func main() {
	// 1. make request
	url := "http://116.89.189.52:8080/get/feature"
	filepath := "face-img.JPG"
	loginData := loginInfo{"chlee", "chlee", "ChangHa"}
	req, err := postIMGreq(url, filepath, loginData)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Do(request)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	fmt.Printf("name: %v\n", body)

}
