/**
The helpers package contains alexa.go, alpha.go, stt.go and tts.go.

tts.go uses the azure speech API to convert text to speech.
**/
package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	ttsURI = "https://uksouth.tts.speech.microsoft.com/cognitiveservices/v1" // The Azure tts endpoint
)

type speechJSON struct {
	Text string `json:"speech"` // Struct for response format
}

/**
Function to convert text to speech using Azure speech API.

Input:
text string : text to be converted. Input is in JSON format.

Returns:
[]byte : byte array of converted speech or error message.
bool   : True if error.
int    : http sstatus code.

The converted speech is returned as a string of a base64 wav. The Azure response is an unencoded wav file.
**/
func TextToSpeech(text string) ([]byte, bool, int) {
	err := godotenv.Load(".env") // load environment variables from .env file.

	// returns internal server if .env cannot be opened/read.
	if err != nil {
		return []byte("Error loading required files."), true, http.StatusInternalServerError
	}

	// Parse text input
	body := pb(text)

	// returns 400 error if body is in invalid format (does not contain speech field for example).
	if body == "" {
		return []byte("Invalid JSON input."), true, http.StatusBadRequest
	}

	c := &http.Client{}                                                        // instantiation of http client.
	req, err := http.NewRequest("POST", ttsURI, bytes.NewReader([]byte(body))) // set request body and method.

	// returns 500 error if there was an error preparing request.
	if err != nil {
		return []byte("Error preparing request."), true, http.StatusInternalServerError
	}

	// set required headers.
	req.Header.Set("X-Microsoft-OutputFormat", "riff-24khz-16bit-mono-pcm")
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("AZURE_KEY"))
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("User-Agent", "gocw/1.0")

	r, err := c.Do(req) // send request and store response.

	// returns 500 error if there was an error sending the request.
	if err != nil {
		return []byte("Error contacting external API."), true, http.StatusInternalServerError
	}

	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(r.Body) // read response body.
		if err != nil {
			return []byte("Error parsing API response."), true, http.StatusBadGateway // returns 502 if can't read response
		} else {
			return toJSON(body), false, http.StatusOK // returns JSON with a "speech" field and string of a base64 encoded wav.
		}
	} else {
		return []byte("Error retrieving converted text."), true, http.StatusBadGateway // Returns 502 error if response status code is not 200
	}
}

/**
Function to parse text field of input JSON.

Input:
s string : the JSON to be parsed.

Returns:
string : parsed JSON value.

Returns an XML to be used in request body if there are no errors.
**/
func pb(t string) string {
	var q TextJSON                       // textJSON struct.
	err := json.Unmarshal([]byte(t), &q) // map text field to q.

	// Return q.Text (empty string) if there is an error in parsing t.
	if err != nil {
		return q.Text
	}

	// Return query in XML format
	return "<speak version='1.0' xml:lang='en-US'><voice xml:lang='en-US' " +
		"xml:gender='Male' name='en-US-ChristopherNeural'>" + q.Text + "</voice></speak>"
}

/**
Function to convert byte array to JSON.

Input:
b []byte : the input to be converted.

Returns:
[]byte : byte array of converted JSON.

The speech field is mapped to a speechJSON struct and encoded to base64.
**/
func toJSON(b []byte) []byte {
	s := speechJSON{base64.StdEncoding.EncodeToString(b)} // encode b and map it to speech field of struct
	o, _ := json.Marshal(s)                               // convert to JSON
	return o                                              // return JSON of encoded b
}
