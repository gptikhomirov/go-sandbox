package main

import (
	"fmt"
	"strings"
)

type FavoriteList struct {
	Titles []string
}

func (l *FavoriteList) Add(title string) {
	l.Titles = append(l.Titles, title)
}

func (l *FavoriteList) PrintAll() {
	fmt.Println(strings.Join(l.Titles, " "))
}

func (l *FavoriteList) PrintUnique() {
	unique := make(map[string]struct{})
	for _, v := range l.Titles {
		unique[v] = struct{}{}
	}
	
	for title := range unique {
		fmt.Print(title)
	}
}

func main() {
	l := FavoriteList{}
	l.Add("item 1")
	l.Add("item 2")
	l.Add("item 3")
	l.Add("item 1")
	l.Add("item 3")
	l.Add("item 4")
	l.PrintAll()
	l.PrintUnique()
}
