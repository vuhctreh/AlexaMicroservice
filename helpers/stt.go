package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	sttURI = "https://uksouth.stt.speech.microsoft.com/speech/recognition/conversation/cognitiveservices/v1?language=en-GB"
)

type Response struct {
	RecognitionStatus string `json:"RecognitionStatus"`
	DisplayText       string `json:"DisplayText"`
	Offset            uint   `json:"Offset"`
	Duration          uint   `json:"Duration"`
}

func SpeechToText(speech string) ([]byte, bool, int) {
	var response Response
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	body := parseSpeech(speech)

	if string(body) == "" {
		return []byte("Invalid JSON input."), true, http.StatusBadRequest
	}

	c := &http.Client{}
	req, err := http.NewRequest("POST", sttURI, bytes.NewReader(body))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "audio/wav; codecs=audio/pcm; samplerate=16000")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Except", "100-continue")
	req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("AZURE_KEY"))
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("User-Agent", "gocw/1.0")

	r, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			err := json.Unmarshal(body, &response)
			if err != nil {
				log.Fatal(err)
			}
			return JSONify(response.DisplayText), false, http.StatusOK
		}
	} else {
		log.Fatal(r.StatusCode)
	}
	return []byte("An unknown error occurred."), true, http.StatusTeapot
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
