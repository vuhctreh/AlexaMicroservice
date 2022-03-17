package helpers

import (
	"net/http"
)

func AnswerQuestion(speech string) ([]byte, bool, int) {
	q, b, e := SpeechToText(speech)

	if b != false {
		return q, b, e
	}

	a, b, _ := GetAnswer(string(q))

	o, b, _ := TextToSpeech(string(a))

	if b != false {
		return []byte("An error occurred whilst processing your request."), true, http.StatusInternalServerError
	}

	return o, false, http.StatusOK
}
