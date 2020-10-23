package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	face "github.com/Kagami/go-face"
	"github.com/gorilla/mux"
)

type Request struct {
	Name  string `json:"name"`
	Image []byte `json:"image"`
}

type Resp struct {
	Feature face.Descriptor `json:"feature"`
}

const (
	SuccesMsg string = "signup success"
	ErrorMsg  string = "you got some errors"
)

/*
// signup/face
func Signup(w http.ResponseWriter, r *http.Request) {

}
*/

//main
func main() {
	fmt.Println("Facial-Auth-server started")
	router := mux.NewRouter()

	URL := "http://116.89.189.52:8080/get/feature"

	// signup/face
	router.HandleFunc("/signup/face", func(w http.ResponseWriter, r *http.Request) {
		//리퀘스트 온거 파싱
		imgFile, name, err := getData(r)
		if err != nil {
			fmt.Println("getData error, Error: ", err.Error())
			return
		}
		//파싱한거 보내서 feature vector 얻어오기
		var feature face.Descriptor
		feature, err = getFeature(imgFile, URL)
		if err != nil {
			fmt.Println("getFeature error, Error: ", err.Error())
			return
		}
		//DB에 저장 - 지금은 JSON파일로 저장
		fmt.Println(name)
		fmt.Println(feature)
		//Signup(w, r)
	}).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8081", router))
}

func getData(r *http.Request) (multipart.File, string, error) {
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
		return f, "", fmt.Errorf("error in FormValue(\"meta-data\"): %s\n", fmt.Errorf("no such key"))
	}

	return imgFile, mdata, err
}

func getFeature(r io.Reader, u string) (face.Descriptor, error) {
	// make multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	var f face.Descriptor //nil descriptor
	// make part
	part, err := writer.CreateFormField("face-img")
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
	req, err := http.NewRequest("POST", u, body)
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
	return rs.Feature, nil
}
