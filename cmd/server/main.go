package main

import (
	"fmt"

	"github.com/evertras/sample-go-app/internal/server"
)

func main() {
	s := server.New("0.0.0.0:8080")

	fmt.Println("Starting")

	if err := s.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
