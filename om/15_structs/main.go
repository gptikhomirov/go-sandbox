package main

import "fmt"

type Coords struct {
	Lat  float64
	Long float64
}

func main() {
	s1 := Coords{Lat: 1, Long: 2}
	s2 := Coords{Lat: 1, Long: 2}
	fmt.Println(s1 == s2)

	s3 := Coords{Lat: 1, Long: 2}
	s4 := Coords{Long: 2, Lat: 1}
	fmt.Println(s3 == s4)

	copied := s1
	copied.Lat = 100
	fmt.Println(s1.Lat, copied.Lat)
}
