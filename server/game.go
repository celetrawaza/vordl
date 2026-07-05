package main

import (
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"slices"
	"strings"
	"time"
)

type GameParams struct {
	WordLength int       `json:"length"`
	MaxGuesses int       `json:"maxTries"`
	Letters    Alphabet  `json:"letters"`
	ResetTime  ResetTime `json:"resetTime"`
	WordNumber int       `json:"wordNumber"`
	word       string
	guesses    map[Player][]Guess
}

var Game GameParams

type ResetTime int64
type GuessStatus string

const (
	GuessCorrect GuessStatus = "correct"
	GuessPresent GuessStatus = "present"
	GuessWrong   GuessStatus = "wrong"
)

type LetterAnnotated struct {
	Letter  Letter      `json:"letter"`
	Correct GuessStatus `json:"correctness"`
}
type GuessAnnotated []LetterAnnotated // `json:"guess"`

type Player string

type Letter rune
type Guess []Letter

type Identity []Letter
type Alphabet []Identity

func makeParams() GameParams {
	now := time.Now().UTC()
	prevWordNumber := Game.WordNumber
	curWordNumber := prevWordNumber + 1
	if prevWordNumber == 0 {
		curWordNumber = Config.StartWordNumber
	}
	return GameParams{
		WordLength: 5,
		MaxGuesses: 7,
		guesses:    map[Player][]Guess{},
		ResetTime:  ResetTime(now.Add(Config.ResetInterval).Unix()),
		WordNumber: curWordNumber,
	}
}

func generatePlayer() (Player, error) {
	b := make([]byte, 16)
	var out Player
	for {
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		out = Player(hex.EncodeToString(b))
		if _, ok := Game.guesses[out]; !ok {
			Game.guesses[out] = []Guess{}
			break
		}
	}
	return out, nil
}

// func (id Identity) MarshalJSON() ([]byte, error) {
// 	// Convert the rune to its string representation and marshal that
// 	out := make([]string, len(id))
// 	for k, v := range id {
// 		out[k] = string(v)
// 	}
// 	return json.Marshal(out)
// }

func (l Letter) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(l))
}

// func (rt ResetTime) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(int64(rt))
// }

func normalizeWord(input string, allowed Alphabet) (ok bool, normalized string) {
	var sb strings.Builder
	letters := []Letter(input)
inputLoop:
	for _, r := range letters {
		// if strings.ContainsRune(allowed, r) {continue}
		for _, id := range allowed {
			for _, l := range id {
				if r != l {
					continue
				}
				sb.WriteRune(rune(id[0]))
				continue inputLoop
			}
		}
		return false, input
	}
	return true, sb.String()
}

func (g Guess) ToGuessAnnotated() (out GuessAnnotated) {
	out = make(GuessAnnotated, len(g))
	runeWord := []Letter(Game.word)
	for i, r := range g {
		out[i] = LetterAnnotated{
			Letter: r,
			Correct: func() (out GuessStatus) {
				switch {
				case r == runeWord[i]:
					out = GuessCorrect
				case slices.Contains(runeWord, r):
					out = GuessPresent
				default:
					out = GuessWrong
				}

				if out != GuessPresent {
					return
				}

				var letters, correct, incorrectIndex int
				for j, l := range runeWord {
					if l == r { // all occurences of the letter in runeWord
						letters++
					}
					if g[j] == l && g[j] == r { // all correct occurences of the letter in guess
						correct++
					}
					if j < i && g[j] == r && g[j] != l { // index of this misplaced occurence of the letter in guess
						incorrectIndex++
					}
				}
				if incorrectIndex >= letters-correct {
					return GuessWrong
				}
				return
			}(),
		}
	}
	return
}
