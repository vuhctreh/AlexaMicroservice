package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/alpha", aioHandler(0))
	router.HandleFunc("/tts", aioHandler(1))
	router.HandleFunc("/stt", aioHandler(2))
	router.HandleFunc("/alexa", aioHandler(3))
	log.Fatal(http.ListenAndServe(":3001", router))

	//t, _, _ := helpers.GetAnswer(`{"text": "What is the melting point of ice"}`)
	//t, _, s := helpers.TextToSpeech(`{"bruh": "What is the melting point of ice"}`)
	//fmt.Println(s)
	//ioutil.WriteFile("test.txt", t, 0644)

	//output := helpers.AnswerQuestion(string(helpers.TextToSpeech("What is the capital of France")))
	//ioutil.WriteFile("test.txt", []byte(output), 0644)
}
