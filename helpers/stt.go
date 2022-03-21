/**
The helpers package contains alexa.go, alpha.go, stt.go and tts.go.

stt.go uses the azure speech API to convert speech to text.
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
	sttURI = "https://uksouth.stt.speech.microsoft.com/speech/recognition/conversation/cognitiveservices/v1?language=en-GB" // Azure speech API endpoint
)

type Response struct {
	RecognitionStatus string `json:"RecognitionStatus"` // RecognitionStatus field from API response.
	DisplayText       string `json:"DisplayText"`       // DisplayText field from API response.
	Offset            uint   `json:"Offset"`            // Offset field from API response.
	Duration          uint   `json:"Duration"`          // Duration field from API response.
}

/**
Function to convert speech to text using Azure speech API.

Input:
speech string : speech to be converted. Input is in JSON format and speech field is encoded in base64.

Returns:
[]byte : byte array of converted text or error message.
bool   : True if error.
int    : http sstatus code.

The converted speech is returned in JSON format with a field "text".
**/
func SpeechToText(speech string) ([]byte, bool, int) {
	var response Response        // Struct for Azure response.
	err := godotenv.Load(".env") // load environment variables from .env file.

	if err != nil {
		return []byte("Error loading required files."), true, http.StatusInternalServerError
	}

	body := parseSpeech(speech)

	if string(body) == "" {
		return []byte("Invalid JSON input."), true, http.StatusBadRequest
	}

	c := &http.Client{}
	req, err := http.NewRequest("POST", sttURI, bytes.NewReader(body))

	if err != nil {
		return []byte("Error preparing request."), true, http.StatusInternalServerError
	}

	req.Header.Set("Content-Type", "audio/wav; codecs=audio/pcm; samplerate=16000")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Except", "100-continue")
	req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("AZURE_KEY"))
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("User-Agent", "gocw/1.0")

	r, err := c.Do(req)
	if err != nil {
		return []byte("Error contacting external API."), true, http.StatusInternalServerError
	}

	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return []byte("Error parsing response from API."), true, http.StatusInternalServerError
		} else {
			err := json.Unmarshal(body, &response)
			if err != nil {
				return []byte("Could not parse API response as JSON."), true, http.StatusInternalServerError
			}
			return JSONify(response.DisplayText), false, http.StatusOK
		}
	} else {
		return []byte("Invalid external API response."), true, http.StatusInternalServerError
	}
}

func parseSpeech(s string) []byte {
	var p speechJSON
	err := json.Unmarshal([]byte(s), &p)

	if err != nil {
		p.Text = ""
		return []byte(p.Text)
	}

	o, _ := base64.StdEncoding.DecodeString(p.Text)

	return o
}
