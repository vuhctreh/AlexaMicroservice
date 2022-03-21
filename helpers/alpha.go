package helpers

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	URL        = "https://api.wolframalpha.com/v1/result?appid="
	InputParam = "&i="
)

type TextJSON struct {
	Text string `json:"text"`
}

func ParseInput(i string) TextJSON {
	var p TextJSON
	err := json.Unmarshal([]byte(i), &p)

	if err != nil {
		p.Text = ""
	}

	return p
}

func JSONify(i string) []byte {
	o := TextJSON{i}

	ba, err := json.Marshal(o)

	if err != nil {
		return []byte("")
	}

	return ba
}

func GetAnswer(i string) ([]byte, bool, int) {
	err := godotenv.Load(".env")

	if err != nil {
		return []byte("Error loading required files."), true, http.StatusInternalServerError
	}

	p := ParseInput(i).Text

	if p == "" {
		return []byte("Invalid JSON input."), true, http.StatusBadRequest
	}

	q := strings.Replace(p, " ", "+", -1)

	resp, err := http.Get(URL + os.Getenv("WOLFRAM_API") + InputParam + q)
	if err != nil {
		log.Print(err)
		return []byte("Short Answers API could not be reached."), true, http.StatusBadGateway
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return []byte("Could not parse response from Wolfram Alpha."), true, http.StatusInternalServerError
	}

	o := JSONify(string(body))

	if string(o) == "" {
		return []byte("Error converting JSON response."), true, http.StatusInternalServerError
	}

	return JSONify(string(body)), false, http.StatusOK
}
