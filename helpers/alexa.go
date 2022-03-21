package helpers

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func AnswerQuestion(s string) ([]byte, bool, int) {
	var body []byte

	c := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:3002/stt", bytes.NewReader([]byte(s)))

	a, b, d := checkReq(err)
	if b != false {
		return a, b, d
	}

	r, err := c.Do(req)

	x, y, z := checkResp(err)
	if y != false {
		return x, y, z
	}

	defer r.Body.Close()

	body, err = ioutil.ReadAll(r.Body)

	if err != nil {
		return []byte("Error parsing response."), true, http.StatusInternalServerError
	}

	req, err = http.NewRequest("POST", "http://127.0.0.1:3001/alpha", bytes.NewReader(body))

	a, b, d = checkReq(err)
	if b != false {
		return a, b, d
	}

	r, _ = c.Do(req)

	x, y, z = checkResp(err)
	if y != false {
		return x, y, z
	}

	defer r.Body.Close()

	body, err = ioutil.ReadAll(r.Body)

	if err != nil {
		return []byte("Error parsing response."), true, http.StatusInternalServerError
	}

	req, _ = http.NewRequest("POST", "http://127.0.0.1:3003/tts", bytes.NewReader(body))

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

	return body, false, http.StatusOK
}

func checkReq(e error) ([]byte, bool, int) {
	if e != nil {
		return []byte("Error preparing request."), true, http.StatusInternalServerError
	}

	return nil, false, http.StatusOK
}

func checkResp(e error) ([]byte, bool, int) {
	if e != nil {
		return []byte("Error contacting external API."), true, http.StatusInternalServerError
	}

	return nil, false, http.StatusOK
}
