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

// Resp implements the response to w (ResponseWriter)
type Resp struct {
	Error string `json:error`
}

// GetFeatureResp implempents the response from /get/feature (face-ai-server)
type GetFeatureResp struct {
	Feature face.Descriptor `json:"feature"`
	Error   string          `json:"error"`
}

// Meta implements the form of meta-data part of r (Request)
type Meta struct {
	Name string `json:name`
}

//GetData get 'meta-data' and 'multipart.File' from r (Request)
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

// GetFeature requests the feature-vec to 'face-ai-server'
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
	rs := new(GetFeatureResp)
	json.NewDecoder(resp.Body).Decode(rs)

	if resp.StatusCode == http.StatusUnauthorized {
		return f, fmt.Errorf("responsed error msg: " + rs.Error)
	}
	return rs.Feature, nil
}

// RespJSON makes app/json response to w (ResponseWriter), and sends it
func RespJSON(w http.ResponseWriter, r bool, e error) {
	var code int  //http status code
	var es string //error string
	if r == true {
		code = http.StatusAccepted
		es = ""
	} else {
		code = http.StatusUnauthorized
		es = e.Error()
	}

	re := Resp{es}
	result, err := json.Marshal(re)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(result)
}
