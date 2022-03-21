/**
The main package contains main.go and server.go.

server.go creates channels to start routers on ports
3000 to 3003. Requests to each endpoint are handled
by a single function.
**/

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gocw/helpers"
	"io/ioutil"
	"net/http"
)

var (
	s []byte // response to request
	b bool   // error boolean. true = error.
	i int    // http status code
)

func startServer() {
	finish := make(chan bool) // Create channel

	var err error // error variable for listenAndServe

	// Create routers for each endpoint. Handled by aioHandler.
	routerAlpha := mux.NewRouter().StrictSlash(true)
	routerAlpha.HandleFunc("/alpha", aioHandler(0))

	routerSTT := mux.NewRouter().StrictSlash(true)
	routerSTT.HandleFunc("/stt", aioHandler(2))

	routerTTS := mux.NewRouter().StrictSlash(true)
	routerTTS.HandleFunc("/tts", aioHandler(1))

	routerAlexa := mux.NewRouter().StrictSlash(true)
	routerAlexa.HandleFunc("/alexa", aioHandler(3))

	// Attach router to each port and serve them.
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

	<-finish // Close channel
}

/**
Function to handle requests for each router.

Input:
sv int : an int representing the service requested.

Returns:
func(http.ResponseWriter, *http.Request) : a response function.

The response function is wrapped and returned by aioHandler in
order to pass variables (sv int) to it.
**/
func aioHandler(sv int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := ioutil.ReadAll(r.Body) // Read request body.

		// Set Content-Type to plain text, error code to 500 and body to error message
		// if ioutil.ReadAll returns an error.
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(500)
			fmt.Fprintf(w, "Internal server error processing your request.")
			return
		}

		sq := string(q) // Convert query to string.
		switch sv {     // switch to process request based on sv. Calls appropriate helper function.
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
			w.Header().Set("Content-Type", "text/plain; charset=utf-8") // set content-type if error
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8") // set content-type
		}
		w.WriteHeader(i)          // set response body
		fmt.Fprintf(w, string(s)) // send response
	}
}

/**
Fucntion to check for errors on router instantiation.

Input:
e error : error returned from instantiation function.

Returns:
void

Logs to console if an error occurred whilst setting up routers
**/
func checkError(e error) {
	if e != nil {
		fmt.Println("An error occurred loading the microservices. Please retry.")
	}
}
