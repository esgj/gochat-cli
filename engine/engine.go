package engine

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/esgj/gochat/model"
)

var currentIntent model.IntentClass

func Run(intents []model.Intent, classes []model.IntentClass) {
	reader := bufio.NewReader(os.Stdin)
	class := classes[0]
	message := make(chan string)
	response := make(chan string)
	currentIntent := GetIntent(class, intents)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	var res string
	var prevMessage string
	var wordScore float32

	go func() {
		// Output the first intent default response
		response <- currentIntent.DefaultResponse

		for {
			wg.Wait()

			if currentIntent.Name != class.Intent {
				currentIntent = GetIntent(class, intents)
			}

			time.Sleep(time.Millisecond * 500)

			if res != "" {
				response <- res
				res = ""
				wg.Add(1)
				continue
			}

			if prevMessage != "" {
				message <- prevMessage
				prevMessage = ""
				wg.Add(1)
				continue
			}

			if data, err := reader.ReadString('\n'); err != nil {
				log.Fatal(err)
			} else {
				message <- data
				wg.Add(1)
			}
		}

	}()

	for {
		select {
		case data := <-message:
			for index, keyword := range currentIntent.Steps[class.CurrentStep].Match {
				if wordScore = CompareTwoStrings(keyword, data); wordScore >= 0.5 {
					res = currentIntent.Steps[class.CurrentStep].Respones[index]
					break
				}
			}

			if res == "" {
				newIntent := MatchNewIntent(data, classes)
				if (newIntent.Intent != model.IntentClass{}.Intent && currentIntent.Name != newIntent.Intent) {
					class = newIntent
					prevMessage = data
				} else {
					rand.Seed(time.Now().Unix())
					randIndex := rand.Intn(len(currentIntent.Steps[class.CurrentStep].Fallback))
					res = currentIntent.Steps[class.CurrentStep].Fallback[randIndex]
				}
			}
			wg.Done()
		case chatResponse := <-response:
			// fmt.Printf("%v %v   (Step: %d Topic: %v WordScore: %f)\r\n", "🤖", chatResponse, class.CurrentStep, class.Intent, wordScore)
			fmt.Printf("%v %v\r\n", "🤖", chatResponse)
			wg.Done()
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
