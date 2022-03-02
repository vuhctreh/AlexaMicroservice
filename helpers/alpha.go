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
	URL        = "https://api.wolframalpha.com/v1/result?appid=2YKEVX-9RPE6GE88R&i"
	InputParam = "&i="
)

type TextJSON struct {
	Text string `json:"text"`
}

func ParseInput(i string) TextJSON {
	var p TextJSON
	err := json.Unmarshal([]byte(i), &p)

	if err != nil {
		log.Fatal(err)
	}

	return p
}

func JSONify(i string) string {
	o := TextJSON{i}

	byteArray, err := json.Marshal(o)

	if err != nil {
		log.Fatal(err)
	}

	return string(byteArray)
}

func GetAnswer(i string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	p := ParseInput(i).Text

	q := strings.Replace(p, " ", "+", -1)

	resp, err := http.Get(URL + os.Getenv("WOLFRAM_API") + InputParam + q)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return JSONify(string(body))
}
