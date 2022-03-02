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
	ttsURI = "https://uksouth.tts.speech.microsoft.com/cognitiveservices/v1"
)

type speechJSON struct {
	Text string `json:"speech"`
}

func TextToSpeech(text string) []byte {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	body := pb(text)

	c := &http.Client{}
	req, err := http.NewRequest("POST", ttsURI, bytes.NewReader([]byte(body)))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-Microsoft-OutputFormat", "riff-24khz-16bit-mono-pcm")
	req.Header.Set("Content-Type", "application/ssml+xml")
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
			return toJSON(body)
		}
	} else {
		log.Fatal(r.StatusCode)
	}
	return nil
}

func pb(t string) string {
	return "<speak version='1.0' xml:lang='en-US'><voice xml:lang='en-US' " +
		"xml:gender='Male' name='en-US-ChristopherNeural'>" + t + "</voice></speak>"
}

func toJSON(b []byte) []byte {
	s := speechJSON{base64.StdEncoding.EncodeToString(b)}
	o, _ := json.Marshal(s)
	return o
}
