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

func SpeechToText(speech string) []byte {
	var response Response
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	body := parseSpeech(speech)

	c := &http.Client{}
	req, err := http.NewRequest("POST", sttURI, bytes.NewReader([]byte(body)))

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

			return []byte(JSONify(response.DisplayText))
		}
	} else {
		log.Fatal(r.StatusCode)
	}
	return nil
}

func parseSpeech(s string) []byte {
	var p speechJSON
	err := json.Unmarshal([]byte(s), &p)

	if err != nil {
		log.Fatal(err)
	}

	o, _ := base64.StdEncoding.DecodeString(p.Text)

	return o
}
