package main

import (
	"fmt"
	"time"
)

// StreamSession - сессия просмотра.
type StreamSession struct {
	UserName  string
	MovieName string
	StartedAt time.Time
	IsActive  bool
}

// NewStreamSession - конструктор StreamSession.
func NewStreamSession(user, movie string) *StreamSession {
	return &StreamSession{
		UserName:  user,
		MovieName: movie,
		StartedAt: time.Now(),
		IsActive:  true,
	}
}

func main() {
	session := NewStreamSession("Алексей Хувер", "Интерстеллар")
	fmt.Printf("%+v\n", *session)
}
