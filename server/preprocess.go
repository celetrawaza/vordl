package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
)

func pickRandomLine(input string) string {
	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	rng := rand.New(rand.NewSource(int64(Game.WordNumber)))
	var result string
	scanner := bufio.NewScanner(file)

	for i := 1; scanner.Scan(); i++ {
		pick := rng.Intn(i)
		if pick == 0 {
			result = scanner.Text()
		}
	}

	return result
}

func prepareFile(input string, output string, allowed Alphabet) {
	inputFile, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	outputFile, err := os.OpenFile(output, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	for scanner.Scan() {
		line := scanner.Text()
		if len([]rune(line)) != Game.WordLength {
			continue
		}
		ok, norm := normalizeWord(line, allowed)
		if !ok {
			continue
		}
		fmt.Fprintln(outputFile, norm)
	}
}

func readLetters(input string) (out Alphabet) {
	out = make(Alphabet, 0)
	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		identity := make(Identity, 0)
		for _, r := range line {
			identity = append(identity, Letter(r))
		}
		out = append(out, identity)
	}
	return
}
