package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	face "github.com/Kagami/go-face"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	//DB.json 불러오기
	file, err := os.OpenFile(DB_name,
		os.O_RDONLY,
		os.FileMode(0644))
	if err != nil {
		RespJSON(w, false, err)
		return
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, file)
	if err != nil {
		RespJSON(w, false, err)
		return
	}
	bb := buf.Bytes()

	jm := make(map[string]face.Descriptor)
	if !isEmpty(bb) {
		err = json.Unmarshal(bb, &jm)
		if err != nil {
			RespJSON(w, false, err)
			return
		}
	} else {
		RespJSON(w, false, errors.New("You need to sign up"))
		return
	}

	// request에서 img, meta-data 파싱하기
	imgFile, mdata, err := GetData(r)
	if err != nil {
		RespJSON(w, false, err)
		return
	}
	nf, err := GetFeature(imgFile)
	if err != nil {
		RespJSON(w, false, err)
		return
	}

	md := new(Meta)
	rs := strings.NewReader(mdata)
	dec := json.NewDecoder(rs)
	dec.Decode(md)
	fmt.Println(md) //tmp log

	//이름 찾기
	var of face.Descriptor
	for k, v := range jm {
		if k == md.Name {
			of = v
		}
	}
	var failF face.Descriptor
	if of == failF {
		RespJSON(w, false, errors.New("There no your name. You need to sign up"))
		return
	}

	//얼굴 매칭하기
	rec, err := face.NewRecognizer("models")
	if err != nil {
		RespJSON(w, false, err)
		return
	}
	var s []face.Descriptor
	var c []int32
	s = append(s, of)
	c = append(c, 0)
	//l := []string{"Its your face!!"}
	rec.SetSamples(s, c)

	var emsg error
	catID := rec.ClassifyTreshold(nf, 0.4) //낮을수록 같기 힘듦. 잘 조정하도록
	if catID < 0 {
		emsg = errors.New("it's not your face. You need to sign up")
		RespJSON(w, false, emsg)
		return
	}
	// reponse
	RespJSON(w, true, emsg)
}
