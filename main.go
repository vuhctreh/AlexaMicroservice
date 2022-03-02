package main

import (
	"gocw/helpers"
	"io/ioutil"
)

func main() {
	//fmt.Println(helpers.GetAnswer(`{"text": "What is the melting point of ice"}`))
	//print(helpers.GetAnswer(string(helpers.SpeechToText(string(helpers.TextToSpeech("How tall is Mt Everest?"))))))
	output := helpers.AnswerQuestion(string(helpers.TextToSpeech("What is the capital of France")))
	ioutil.WriteFile("test.txt", []byte(output), 0644)
}
