package helpers

import (
	"encoding/json"
	"log"
)

func AnswerQuestion(speech string) string {
	var Answer TextJSON

	q := SpeechToText(speech)
	a := GetAnswer(string(q))
	err := json.Unmarshal([]byte(a), &Answer)
	if err != nil {
		log.Fatal(err)
	}
	o := TextToSpeech(Answer.Text)

	return string(o)
}
