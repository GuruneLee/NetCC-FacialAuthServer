package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	face "github.com/Kagami/go-face"
)

//GetData get 'meta-data' and 'multipart.File'
func GetData(r *http.Request) (multipart.File, string, error) {
	r.ParseMultipartForm(32 << 20)
	var f multipart.File //nil file

	imgFile, _, err := r.FormFile("user-face")
	if err != nil {
		return f, "", fmt.Errorf("error in FormFile(\"user-face\"): " + err.Error())
	}
	defer imgFile.Close()

	/*
			imgByte, err := ioutil.ReadAll(file)
			if err != nil {
			return f, "", err
		}
	*/

	//mdata is just 'user name' now
	mdata := r.PostFormValue("meta-data")
	if mdata == "" {
		return f, "", fmt.Errorf("error in FormValue(\"meta-data\"): %v", fmt.Errorf("no such key"))
	}

	return imgFile, mdata, err
}

// GetFeature request the feature-vec to 'face-ai-server'
func GetFeature(r io.Reader) (face.Descriptor, error) {
	// make multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	var f face.Descriptor //nil descriptor
	// make part
	part, err := writer.CreateFormFile("face-img", "face-img")
	if err != nil {
		return f, err
	}
	_, err = io.Copy(part, r)
	if err != nil {
		return f, err
	}
	err = writer.Close()
	if err != nil {
		return f, err
	}

	// req 생성
	req, err := http.NewRequest("POST", URL, body)
	if err != nil {
		return f, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// rep받아서 json으로 바꾸기
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return f, err
	}
	defer resp.Body.Close()

	if resp.Body == nil {
		return f, errors.New("empty response body")
	}
	rs := new(Resp)
	json.NewDecoder(resp.Body).Decode(rs)

	if resp.StatusCode == http.StatusUnauthorized {
		return f, fmt.Errorf("responsed error msg: " + rs.Error)
	}
	return rs.Feature, nil
}
