package main

import (
	"fmt"
	"gocw/helpers"
	"io/ioutil"
	"net/http"
)

var (
	s []byte
	b bool
	i int
)

func aioHandler(sv int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q, _ := ioutil.ReadAll(r.Body)
		sq := string(q)
		switch sv {
		case 0:
			s, b, i = helpers.GetAnswer(sq)
		case 1:
			s, b, i = helpers.TextToSpeech(sq)
		case 2:
			s, b, i = helpers.SpeechToText(sq)
		case 3:
			s, b, i = helpers.AnswerQuestion(sq)
		}

		if b != false {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
		w.WriteHeader(i)
		fmt.Fprintf(w, string(s))
	}
}
