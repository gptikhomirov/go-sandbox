package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Movie struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

func trackWorkTime(start time.Time) {
	end := time.Since(start).Round(time.Millisecond)
	fmt.Println("duration", end)
}

func loadFile(path string) (movies []Movie, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic at %s: %v", path, r)
		}
	}()

	file, errOpen := os.Open(path)
	if errOpen != nil {
		return nil, fmt.Errorf("error open file: %w", errOpen)
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			err = errors.Join(err, fmt.Errorf("error close file: %w", errClose))
		}
	}()

	if errDecode := json.NewDecoder(file).Decode(&movies); errDecode != nil {
		return nil, fmt.Errorf("error decode file: %w", errDecode)
	}

	return movies, nil
}

func loadCatalog(files []string) []Movie {
	defer trackWorkTime(time.Now())

	allMovies := make([]Movie, 0, len(files))
	for _, fileName := range files {
		movies, err := loadFile(fileName)

		if err != nil {
			fmt.Println(fmt.Errorf("err load file %s: %w", fileName, err))
		} else {
			allMovies = append(allMovies, movies...)
		}
	}

	return allMovies
}

func main() {
	movieNames := []string{
		"movies_empty.json",
		"movies_invalid.json",
		"movies_zero.json",
		"movies_valid.json",
	}

	movies := loadCatalog(movieNames)
	fmt.Println(movies)
}
