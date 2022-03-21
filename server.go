package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gocw/helpers"
	"io/ioutil"
	"net/http"
)

var (
	s []byte
	b bool
	i int
)

func startServer() {
	finish := make(chan bool)

	var err error

	routerAlpha := mux.NewRouter().StrictSlash(true)
	routerAlpha.HandleFunc("/alpha", aioHandler(0))

	routerSTT := mux.NewRouter().StrictSlash(true)
	routerSTT.HandleFunc("/stt", aioHandler(2))

	routerTTS := mux.NewRouter().StrictSlash(true)
	routerTTS.HandleFunc("/tts", aioHandler(1))

	routerAlexa := mux.NewRouter().StrictSlash(true)
	routerAlexa.HandleFunc("/alexa", aioHandler(3))

	go func() {
		err = http.ListenAndServe(":3000", routerAlexa)
		checkError(err)
	}()

	go func() {
		err = http.ListenAndServe(":3001", routerAlpha)
		checkError(err)
	}()

	go func() {
		err = http.ListenAndServe(":3002", routerSTT)
		checkError(err)
	}()

	go func() {
		err = http.ListenAndServe(":3003", routerTTS)
		checkError(err)
	}()

	fmt.Println("Microservices running. Listening to ports 3000 -> 3003.")

	<-finish
}

func aioHandler(sv int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(500)
			fmt.Fprintf(w, "Internal server error processing your request.")
			return
		}

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

func checkError(e error) {
	if e != nil {
		fmt.Println("An error occurred loading the microservices. Please retry.")
	}
}
