package main

import "fmt"

func addFilm(catalog map[string]float64, title string, rating float64) {
	if _, ok := catalog[title]; ok {
		return
	}

	catalog[title] = rating
}

func topFilms(catalog map[string]float64, min float64) []string {
	s := make([]string, 0, len(catalog))
	for k, v := range catalog {
		if v >= min {
			s = append(s, k)
		}
	}
	return s
}

func avgRating(catalog map[string]float64) float64 {
	avg := 0.0
	for _, v := range catalog {
		avg += v
	}
	return avg / float64(len(catalog))
}

func main() {
	var catalog = map[string]float64{
		"Inception":                 8.2,
		"Dark Knight":               9.0,
		"Kolobok":                   3.2,
		"Spider man: Brand New Day": 7.0,
		"Top gun: Maverick":         8.5,
	}

	fmt.Println(catalog)
	addFilm(catalog, "Kolobok", 4.6)
	addFilm(catalog, "Se7en", 8.6)
	fmt.Println(catalog)

	fmt.Println(avgRating(catalog))

	fmt.Println(topFilms(catalog, 8))
	fmt.Println(topFilms(catalog, 9))
	fmt.Println(topFilms(catalog, 7))

}
