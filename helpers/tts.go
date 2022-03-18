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
	ttsURI = "https://uksouth.tts.speech.microsoft.com/cognitiveservices/v1"
)

type speechJSON struct {
	Text string `json:"speech"`
}

func TextToSpeech(text string) ([]byte, bool, int) {
	err := godotenv.Load(".env")

	if err != nil {
		return []byte("Error loading required files."), true, http.StatusInternalServerError
	}

	body := pb(text)

	if body == "" {
		return []byte("Invalid JSON input."), true, http.StatusBadRequest
	}

	c := &http.Client{}
	req, err := http.NewRequest("POST", ttsURI, bytes.NewReader([]byte(body)))

	if err != nil {
		return []byte("Error preparing request."), true, http.StatusInternalServerError
	}

	req.Header.Set("X-Microsoft-OutputFormat", "riff-24khz-16bit-mono-pcm")
	req.Header.Set("Content-Type", "application/ssml+xml")
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
			return []byte("Error parsing API response."), true, http.StatusInternalServerError
		} else {
			return toJSON(body), false, http.StatusOK
		}
	} else {
		return []byte("Error retrieving converted text."), true, http.StatusInternalServerError
	}
}

func pb(t string) string {
	var q TextJSON
	err := json.Unmarshal([]byte(t), &q)

	if err != nil {
		return q.Text
	}

	return "<speak version='1.0' xml:lang='en-US'><voice xml:lang='en-US' " +
		"xml:gender='Male' name='en-US-ChristopherNeural'>" + q.Text + "</voice></speak>"
}

func toJSON(b []byte) []byte {
	s := speechJSON{base64.StdEncoding.EncodeToString(b)}
	o, _ := json.Marshal(s)
	return o
}
