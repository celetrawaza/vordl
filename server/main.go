package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"time"
)

func getGuesses(w http.ResponseWriter, r *http.Request) {
	player := retrievePlayer(r)
	guesses := Game.guesses[player]
	out := make([]GuessAnnotated, len(guesses))
	for i, g := range guesses {
		out[i] = g.ToGuessAnnotated()
	}
	returnJSON(w, player, out)
}

func getParams(w http.ResponseWriter, r *http.Request) {
	player := retrievePlayer(r)
	returnJSON(w, player, Game)
}

func hasGuessed(player Player) bool {
	return slices.ContainsFunc(Game.guesses[player], func(g Guess) bool {
		return !slices.ContainsFunc(g.ToGuessAnnotated(), func(la LetterAnnotated) bool {
			return la.Correct != GuessCorrect
		})
	})
}

func makeGuess(w http.ResponseWriter, r *http.Request) {
	player := retrievePlayer(r)
	// if slices.ContainsFunc(Game.guesses[player], func(g Guess) bool {
	// 	return !slices.ContainsFunc(g.ToGuessAnnotated(), func(la LetterAnnotated) bool {
	// 		return la.Correct != GuessCorrect
	// 	})
	// })
	if hasGuessed(player) {
		// if already got correct guess
		http.Error(w, "Already guessed current word", http.StatusConflict)
		return
	}
	if len(Game.guesses[player]) >= Game.MaxGuesses {
		http.Error(w, "Too many guesses", http.StatusConflict) // todo fix statuses
		return
	}
	var guess Guess
	var guessWord string
	if err := json.NewDecoder(r.Body).Decode(&guessWord); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	guess = Guess(guessWord)
	if slices.ContainsFunc(Game.guesses[player], func(g Guess) bool {
		return slices.Equal(g, guess)
	}) {
		http.Error(w, "Already guessed this word", http.StatusConflict)
		return
	}
	// todo accept only real words
	if len(guess) != Game.WordLength {
		http.Error(w, "Wrong length", http.StatusBadRequest)
		return
	}
	ok, unnormGuess := normalizeWord(guessWord, Game.Letters)
	if !ok {
		http.Error(w, "Wrong letters", http.StatusBadRequest)
		return
	}
	normGuess := Guess(unnormGuess)
	Game.guesses[player] = append(Game.guesses[player], normGuess)
	getGuesses(w, r)
}

func getWord(w http.ResponseWriter, r *http.Request) {
	player := retrievePlayer(r)
	guessed := hasGuessed(player)
	failed := len(Game.guesses[player]) >= Game.MaxGuesses
	if !guessed && !failed {
		http.Error(w, "Not allowed to see the word", http.StatusForbidden)
		return
	}
	returnJSON(w, player, Game.word)
}

func resetGame() {
	Game = makeParams()
	// get letters
	Game.Letters = readLetters("letters.txt")
	// // parse words
	if _, err := os.Stat("data.txt"); err != nil {
		prepareFile("russian.txt", "data.txt", Game.Letters)
	}
	// pick word
	Game.word = pickRandomLine("data.txt")
	// Game.word = "емъёт"
	_, Game.word = normalizeWord(Game.word, Game.Letters)
	// todo timer
	// time.Timer
}

func main() {
	Config = LoadConfig()
	go func() {
		now := time.Now()
		targetTime := now.Truncate(Config.ResetInterval).Add(Config.ResetInterval)
		// init
		resetGame()
		Game.ResetTime = ResetTime(targetTime.Unix())
		// first reset
		delay := time.Until(targetTime)
		timer := time.NewTimer(delay)
		defer timer.Stop()
		<-timer.C
		resetGame()
		// next resets
		ticker := time.NewTicker(Config.ResetInterval)
		defer ticker.Stop()
		for range ticker.C {
			resetGame()
		}
	}()
	// serve and enjoy
	router := NewRouter()
	fmt.Printf("Serving on port %d\n", Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", Config.Port), router); err != nil {
		log.Fatal(err)
	}
}
