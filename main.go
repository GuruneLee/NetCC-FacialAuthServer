package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

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

	URL := "116.89.189.52:8080/get/feature"

	// signup/face
	router.HandleFunc("/signup/face", func(w http.ResponseWriter, r *http.Request) {
		//리퀘스트 온거 파싱
		imgFile, name, err := getData(r)
		//파싱한거 보내서 feature vector 얻어오기
		var feature face.Descriptor
		feature, err = getFeature(imgFile, URL)
		if err != nil {
			fmt.Errorf("getFeature error, Error: ", err.Error())
			return
		}
		//DB에 저장 - 지금은 JSON파일로 저장
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
		return f, "", err
	}
	defer imgFile.Close()

	/*
			imgByte, err := ioutil.ReadAll(file)
			if err != nil {
			return f, "", err
		}
	*/

	//mdata is just 'user name' now
	mdata, _, err := r.FormFile("meta-data")
	if err != nil {
		return f, "", err
	}
	defer mdata.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, mdata)
	if err != nil {
		return f, "", err
	}
	name := buf.String()

	return imgFile, name, err
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

	rs := new(Resp)
	json.NewDecoder(resp.Body).Decode(rs)

	return rs.Feature, nil
}