package main

import (
	"encoding/json"
	"net/http"
)

func schedulizeFunc(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})
		GameScheduler <- func() {
			f(w, r)
			close(done)
		}
		<-done
	}
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello!"))
	})
	mux.HandleFunc("GET /api/tries", schedulizeFunc(getGuesses))
	mux.HandleFunc("GET /api/params", schedulizeFunc(getParams))
	mux.HandleFunc("POST /api/guess", schedulizeFunc(makeGuess))
	mux.HandleFunc("GET /api/word", schedulizeFunc(getWord))

	return mux
}

func retrievePlayer(r *http.Request) Player {
	cookie, err := r.Cookie("player_id")
	if err != nil {
		player, _ := generatePlayer()
		return player
	}
	out := Player(cookie.Value)
	if _, ok := Game.guesses[out]; !ok {
		Game.guesses[out] = []Guess{}
	}
	return out
}

func updateCookie(w http.ResponseWriter, player Player) {
	cookie := &http.Cookie{
		Name:   "player_id",
		Value:  string(player),
		Path:   "/",
		MaxAge: int(Config.ResetInterval.Seconds()) * 2,
	}
	http.SetCookie(w, cookie)
}

func returnJSON(w http.ResponseWriter, player Player, data any) {
	w.Header().Set("Content-Type", "application/json")
	updateCookie(w, player)
	json.NewEncoder(w).Encode(data)
}
