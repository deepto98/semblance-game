package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
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

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

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

			//Call Cognitive Services API
			currentTags := TagRemoteImage(CreateComputerVisionClient(), url)

			var data ImageGame

			if statusOfAnswer == "correctResponse" || statusOfAnswer == "incorrectResponsesExhausted" {
				if statusOfAnswer == "incorrectResponsesExhausted" && lives == 0 {
					data =
						ImageGame{
							Score:         score,
							Lives:         lives,
							LivesAsHearts: strings.Repeat("❤️", lives),
							TotalGuesses:  totalGuesses,
							Accuracy:      (score / totalGuesses) * 100,
						}
					scoreTemplate := template.Must(template.ParseFiles("resources/Score.html"))
					scoreTemplate.Execute(w, data)
					return
				}

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

			} else if currentTags[strings.ToLower(guess)] {
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
			} else {
				var answer string

				if consecutiveWrongAnswers == 2 {
					lives--
					//Get answer
					//get the first key, since that has the highest confidence
					for key, _ := range currentTags {
						answer = key
						break
					}
				}
				data =
					ImageGame{
						CurrentImageSource:      url,
						Score:                   score,
						StatusOfAnswer:          "incorrect",
						Lives:                   lives,
						LivesAsHearts:           strings.Repeat("❤️", lives),
						ConsecutiveWrongAnswers: consecutiveWrongAnswers + 1,
						TotalGuesses:            totalGuesses + 1,
						Answer:                  answer,
					}
			}

			gameTemplate.Execute(w, data)
		}

	})
	http.ListenAndServe(":"+port, nil)
}
