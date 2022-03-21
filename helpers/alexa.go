/**
The helpers package contains alexa.go, alpha.go, stt.go and tts.go.

alexa.go sends requests to each microservice in order to get the
answer to a question.
The order is as follows:
SpeechToText -> GetAnswer -> TextToSpeech -> Response to initial request.
**/

package helpers

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

/**
Function that goes through the whole process of answering a spoken question
via speech.

Input:
s string : Question to be answered. Input is in JSON format.

Returns:
[]byte : byte array of converted speech or error message.
bool   : True if error.
int    : http status code.

The converted speech is returned as a string of a base64 wav.
**/
func AnswerQuestion(s string) ([]byte, bool, int) {
	var body []byte // request body

	c := &http.Client{}                                                                          // instantiation of http client
	req, err := http.NewRequest("POST", "http://127.0.0.1:3002/stt", bytes.NewReader([]byte(s))) // prepare and send req to endpoint

	a, b, d := checkReq(err) // check if there was an error preparing the request
	if b != false {
		return a, b, d //returns 500 error
	}

	r, err := c.Do(req) // send request to endpoint

	x, y, z := checkResp(err) // check for invalid response
	if y != false {
		return x, y, z // returns 500 error
	}

	defer r.Body.Close()

	body, err = ioutil.ReadAll(r.Body) // read response body

	if err != nil {
		return []byte("Error parsing response."), true, http.StatusInternalServerError // returns 500 error if response can not be parsed
	}

	// The following chunk of code does the same as above for different endpoints and using responses from each microservice.
	req, err = http.NewRequest("POST", "http://127.0.0.1:3001/alpha", bytes.NewReader(body))

	a, b, d = checkReq(err)
	if b != false {
		return a, b, d
	}

	r, err = c.Do(req)

	x, y, z = checkResp(err)
	if y != false {
		return x, y, z
	}

	defer r.Body.Close()

	body, err = ioutil.ReadAll(r.Body)

	if err != nil {
		return []byte("Error parsing response."), true, http.StatusInternalServerError
	}

	req, err = http.NewRequest("POST", "http://127.0.0.1:3003/tts", bytes.NewReader(body))

	a, b, d = checkReq(err)
	if b != false {
		return a, b, d
	}

	r, err = c.Do(req)

	x, y, z = checkResp(err)
	if y != false {
		return x, y, z
	}

	defer r.Body.Close()

	body, err = ioutil.ReadAll(r.Body)

	if err != nil {
		return []byte("Error parsing response."), true, http.StatusInternalServerError
	}

	return body, false, http.StatusOK // Return answer in JSON format "{"speech" : "ANSWER_HERE"}
}

/**
Function that checks for errors preparing requests.

Input:
e error : error variable.

Returns:
[]byte : byte array of error message.
bool   : True if error.
int    : http status code.

Returns nil, false, 200 if e is nil.
**/
func checkReq(e error) ([]byte, bool, int) {
	if e != nil {
		return []byte("Error preparing request."), true, http.StatusInternalServerError // return 500 error if e != nil
	}

	return nil, false, http.StatusOK // return ok signal if error == nil.
}

/**
Function that checks for errors in request response.

Input:
e error : error variable.

Returns:
[]byte : byte array of error message.
bool   : True if error.
int    : http status code.

Returns nil, false, 200 if e is nil.
**/
func checkResp(e error) ([]byte, bool, int) {
	if e != nil {
		return []byte("Error contacting external API."), true, http.StatusInternalServerError // return 500 error if e != nil
	}

	return nil, false, http.StatusOK // return ok signal if error == nil.
}
