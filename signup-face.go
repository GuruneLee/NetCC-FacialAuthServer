package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	face "github.com/Kagami/go-face"
)

type SignUp_Meta struct {
	Name string `json:name`
}

// Signup has /signup/face api logic
func Signup(w http.ResponseWriter, r *http.Request) {
	//리퀘스트 온거 파싱
	imgFile, mdata, err := GetData(r)
	if err != nil {
		fmt.Println("GetData error, Error: ", err.Error())
		RespJson(w, false, err)
		return
	}
	//파싱한거 보내서 feature vector 얻어오기
	var feature face.Descriptor
	feature, err = GetFeature(imgFile)
	if err != nil {
		fmt.Println("GetFeature error, Error: ", err.Error())
		RespJson(w, false, err)
		return
	}

	// mdata -> json go value
	md := new(SignUp_Meta)
	rs := strings.NewReader(mdata)
	dec := json.NewDecoder(rs)
	err = dec.Decode(md)
	if err != nil {
		fmt.Println("error in parsing the meta data, Error: ", err.Error())
		RespJson(w, false, err)
		return
	}
	//DB에 저장 - 지금은 JSON파일로 저장
	err = makeFile(md, feature, DB_name)
	if err != nil {
		fmt.Println("error in making File, Error: ", err.Error())
		RespJson(w, false, err)
		return
	}

	// 성공한 response
	RespJson(w, true, nil)

}

func makeFile(m *SignUp_Meta, f face.Descriptor, fn string) error {
	// open file
	file, err := os.OpenFile(
		fn, //file name
		os.O_CREATE|os.O_RDWR,
		os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("openFile error - " + err.Error())
	}
	defer file.Close()

	// 파일 json으로 읽어오기
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, file)
	if err != nil {
		return fmt.Errorf("io.Copy error - " + err.Error())
	}
	bb := buf.Bytes()

	jm := make(map[string]interface{})
	if !isEmpty(bb) {
		err = json.Unmarshal(bb, &jm)
		if err != nil {
			return fmt.Errorf("Unmarshal error - " + err.Error())
		}
	}

	// json에 meta-data와 face-decriptor추가
	n := m.Name
	if jm[n] != nil {
		return fmt.Errorf("There is same named account...우린 동명이인은 고려안해요")
	} else {
		jm[n] = f
	}

	// 파일에 다시 쓰기
	jsonBytes, err := json.MarshalIndent(jm, "", "  ")
	if err != nil {
		return fmt.Errorf("Marshal error - " + err.Error())
	}
	file.WriteAt(jsonBytes, 0)

	return nil

}

func isEmpty(b []byte) bool {
	if len(b) == 0 {
		return true
	}
	return false
}
