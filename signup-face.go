package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	face "github.com/Kagami/go-face"
)

// Signup has /signup/face api logic
func Signup(w http.ResponseWriter, r *http.Request) {
	//리퀘스트 온거 파싱
	imgFile, mdata, err := GetData(r)
	if err != nil {
		fmt.Println("GetData error, Error: ", err.Error())
		return
	}
	//파싱한거 보내서 feature vector 얻어오기
	var feature face.Descriptor
	feature, err = GetFeature(imgFile)
	if err != nil {
		fmt.Println("GetFeature error, Error: ", err.Error())
		return
	}

	// mdata -> json go value
	md := new(Meta)
	rs := strings.NewReader(mdata)
	dec := json.NewDecoder(rs)
	err = dec.Decode(md)
	if err != nil {
		fmt.Println("error in parsing the meta data, Error: ", err.Error())
		return
	}
	//DB에 저장 - 지금은 JSON파일로 저장
	fmt.Println(mdata)
	fmt.Println(md.Name)
	fmt.Println(feature)
}
