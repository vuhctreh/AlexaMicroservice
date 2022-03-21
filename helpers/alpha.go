/**
The helpers package contains alexa.go, alpha.go, stt.go and tts.go.

alpha.go uses the wolframalpha short answers API to generate an answer
to a given query.
**/
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
	URL        = "https://api.wolframalpha.com/v1/result?appid=" // the base URL for wolframalpha API
	InputParam = "&i="                                           // input parameter to appand ot URL
)

type TextJSON struct {
	Text string `json:"text"` // text key of a JSON
}

/**
Function to parse JSON.

Input:
i string : the JSON input.

Returns:
TextJSON : "text" field of the JSON mapped to struct.

If the JSON is invalid, aka does not containt a "text" field,
maps struct to an empty string instead.
**/
func ParseInput(i string) TextJSON {
	var p TextJSON // struct for storing parsed query
	err := json.Unmarshal([]byte(i), &p)

	if err != nil {
		p.Text = "" // sets query to empty string if JSON is invalid.
	}

	return p
}

/**
Function to convert string to JSON.

Input:
i string : string to be converted.

Returns:
[]byte : byte array of the converted string.

If an error occurs trying to Marshal the struct, returns
byte array of an empty string.
**/
func JSONify(i string) []byte {
	o := TextJSON{i} // The struct used to build the JSON

	ba, err := json.Marshal(o) // converts struct to JSON format.

	if err != nil {
		return []byte("") // Returns byte array of an empty string if error.
	}

	return ba
}

/**
Function to get answer for query "i" from wolframalpha short answers API.

Input:
i string : query string.

Returns:
[]byte : byte array of the converted string.
bool   : True if there was an error.
int    : http status code.

If I is not in JSON  format, a 400 error is returned. Many other errors are returned
and marked as internal server errors.
**/
func GetAnswer(i string) ([]byte, bool, int) {
	err := godotenv.Load(".env") // load environment variables from .env file.

	// returns internal server if .env cannot be opened.
	if err != nil {
		return []byte("Error loading required files."), true, http.StatusInternalServerError
	}

	p := ParseInput(i).Text // call ParseInput on JSON

	if p == "" {
		return []byte("Invalid JSON input."), true, http.StatusBadRequest // return 400 error if JSON is invalid.
	}

	q := strings.Replace(p, " ", "+", -1) // replace all spaces in query with a "+"

	resp, err := http.Get(URL + os.Getenv("WOLFRAM_API") + InputParam + q) // Send request to API with query
	if err != nil {
		log.Print(err)
		return []byte("Short Answers API could not be reached."), true, http.StatusBadGateway // Return 502 if API unreachable.
	}

	body, err := ioutil.ReadAll(resp.Body) // Read response
	if err != nil {
		log.Print(err)
		return []byte("Could not parse response from Wolfram Alpha."), true, http.StatusInternalServerError // Return 500 if error parsing response.
	}

	o := JSONify(string(body)) // Convert response to JSON.

	if string(o) == "" {
		return []byte("Error converting JSON response."), true, http.StatusBadGateway // Return 502 if converted response is invalid.
	}

	return JSONify(string(body)), false, http.StatusOK // Return 200 and answer JSON.
}
