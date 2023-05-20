package engine

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/esgj/gochat/model"
)

var currentIntent model.IntentClass

func Run(intents []model.Intent, classes []model.IntentClass) {
	reader := bufio.NewReader(os.Stdin)
	class := classes[0]
	var previousInput string

	for {
		currentIntent := GetIntent(class, intents)

		if previousInput == "" {
			// Print out the first intent response
			fmt.Printf("%v %v   (Step: %d Topic: %v)\r\n", "🤖", currentIntent.Steps[0].Respones[0], class.CurrentStep, class.Intent)
			// Go to next step
			class.CurrentStep++
		}

		for {
			var data string
			var response string

			if previousInput == "" {
				if d, err := reader.ReadString('\n'); err != nil {
					log.Fatal(err)
				} else {
					data = d
				}
			} else {
				data = previousInput
				previousInput = ""
			}
			

			if len(currentIntent.Steps) == class.CurrentStep {
				rand.Seed(time.Now().Unix())
				class.CurrentStep = rand.Intn(len(currentIntent.Steps))
				if (class.CurrentStep == 0 && currentIntent.Name == intents[0].Name) {
					class.CurrentStep++
				}
			}

			for index, keyword := range currentIntent.Steps[class.CurrentStep].Match {
				if CompareTwoStrings(keyword, data) >= 0.5 {
					response = currentIntent.Steps[class.CurrentStep].Respones[index]
					break
				}
			}

			if response == "" {
				newIntent := MatchNewIntent(data, classes)
				if (newIntent.Intent != model.IntentClass{}.Intent) {
					class = newIntent
					previousInput = data
					break
				}
				rand.Seed(time.Now().Unix())
				randIndex := rand.Intn(len(currentIntent.Steps[class.CurrentStep].Fallback))
				response = currentIntent.Steps[class.CurrentStep].Fallback[randIndex]
			}

			fmt.Printf("%v %v   (Step: %d Topic: %v)\r\n", "🤖", response, class.CurrentStep, class.Intent)
		}
	}
}

func GetIntent(class model.IntentClass, intents []model.Intent) model.Intent {
	for _, intent := range intents {
		if intent.Name == class.Intent {
			return intent
		}
	}

	return intents[0]
}

func MatchNewIntent(word string, classes []model.IntentClass) model.IntentClass {
	for index, class := range classes {
		for _, classWord := range classes[index].Words {
			if CompareTwoStrings(classWord, word) > 0.5 {
				return class
			}
		}
	}

	return model.IntentClass{}
}