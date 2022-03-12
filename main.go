package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ImageGame struct {
	CurrentImageSource      string
	Score                   int
	Lives                   int
	LivesAsHearts           string
	ConsecutiveWrongAnswers int
	Answer                  string
	StatusOfAnswer          string
	TotalGuesses            int
	Accuracy                int
}

var myCache CacheItf

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	InitCache() // comment if want to use redis cache

	fs := http.FileServer(http.Dir("resources/Semblance Game_files"))
	http.Handle("/Semblance Game_files/", http.StripPrefix("/Semblance Game_files/", fs))
	gameTemplate := template.Must(template.ParseFiles("resources/Semblance Game.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			resp, _ := http.Get("https://random.imagecdn.app/1024/600")
			url := resp.Request.URL.String()

			data :=
				ImageGame{
					StatusOfAnswer:          "nonAnswer",
					CurrentImageSource:      url,
					Score:                   0,
					Lives:                   5,
					LivesAsHearts:           strings.Repeat("❤️", 5),
					ConsecutiveWrongAnswers: 0,
					TotalGuesses:            0,
				}

			gameTemplate.Execute(w, data)
		} else if r.Method == "POST" {
			r.ParseForm()

			var totalGuesses int
			if _, err := fmt.Sscanf(r.FormValue("totalGuesses"), "%d", &totalGuesses); err == nil {
				fmt.Println(totalGuesses)
			}
			var score int
			if _, err := fmt.Sscanf(r.FormValue("score"), "%d", &score); err == nil {
				fmt.Println(score)
			}
			var consecutiveWrongAnswers int
			if _, err := fmt.Sscanf(r.FormValue("consecutiveWrongAnswers"), "%d", &consecutiveWrongAnswers); err == nil {
				fmt.Println(consecutiveWrongAnswers)
			}

			var lives int
			if _, err := fmt.Sscanf(r.FormValue("lives"), "%d", &lives); err == nil {
				fmt.Println(lives)
			}
			url := r.FormValue("currentImageSource")
			statusOfAnswer := r.FormValue("statusOfAnswer")
			guess := r.FormValue("guess")

			var data ImageGame

			if statusOfAnswer == "correctResponse" || statusOfAnswer == "incorrectResponsesExhausted" {

				resp, _ := http.Get("https://random.imagecdn.app/1024/600")
				newurl := resp.Request.URL.String()
				data =
					ImageGame{
						CurrentImageSource:      newurl,
						Score:                   score,
						StatusOfAnswer:          "nonAnswer",
						Lives:                   lives,
						LivesAsHearts:           strings.Repeat("❤️", lives),
						ConsecutiveWrongAnswers: 0,
						TotalGuesses:            totalGuesses + 1,
					}
			} else {
				tagsFromCache, err := myCache.Get(url)
				var currentTags map[string]int
				json.Unmarshal(tagsFromCache, &currentTags)

				if err != nil {
					// error
					log.Fatal(err)
				}
				if currentTags == nil {
					currentTags = TagRemoteImage(CreateComputerVisionClient(), url)
					myCache.Set(url, currentTags, -1)
				}
				guess = strings.ToLower(guess)
				acceptableAnswers := getWordForms(guess)

				answerFound := false
				for _, acceptableAnswer := range acceptableAnswers {

					if currentTags[acceptableAnswer] > 0 {
						answerFound = true
						break
					}

				}
				if !answerFound {
					for currentTag := range currentTags {

						for _, acceptableAnswer := range acceptableAnswers {

							if checkEitherIsSubString(currentTag, acceptableAnswer) {
								answerFound = true
								break
							}

						}
					}
				}

				if answerFound {
					data =
						ImageGame{
							CurrentImageSource: url,
							Score:              score + 1,
							StatusOfAnswer:     "correct",
							Lives:              lives,
							LivesAsHearts:      strings.Repeat("❤️", lives),

							ConsecutiveWrongAnswers: 0,
							TotalGuesses:            totalGuesses + 1,
						}
				}
				if !answerFound {
					var answer string

					if consecutiveWrongAnswers == 2 {
						lives--
						//Get answer
						//get the first key, since that has the highest confidence
						for key, pos := range currentTags {
							if pos == 1 {
								answer = key
								break
							}

						}
					}
					if lives == 0 {
						data =
							ImageGame{
								Score:         score,
								Lives:         lives,
								LivesAsHearts: strings.Repeat("❤️", lives),
								TotalGuesses:  totalGuesses,
								Accuracy:      int(float64(score) / float64(totalGuesses) * 100),
							}
						scoreTemplate := template.Must(template.ParseFiles("resources/Score.html"))
						scoreTemplate.Execute(w, data)
						return
					}

					data =
						ImageGame{
							CurrentImageSource:      url,
							Score:                   score,
							StatusOfAnswer:          "incorrectResponse",
							Lives:                   lives,
							LivesAsHearts:           strings.Repeat("❤️", lives),
							ConsecutiveWrongAnswers: consecutiveWrongAnswers + 1,
							TotalGuesses:            totalGuesses + 1,
							Answer:                  answer,
						}
				}

			}

			gameTemplate.Execute(w, data)
		}

	})
	http.ListenAndServe(":"+port, nil)
}
