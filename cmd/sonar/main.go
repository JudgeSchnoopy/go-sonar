package main

import (
	"fmt"

	"github.com/JudgeSchnoopy/go-sonar/internal/server"
)

func main() {
	sonar, err := server.New()
	if err != nil {
		panic(err)
	}
	fmt.Printf("starting sonar")
	sonar.ListenAndServe()
	fmt.Println("ok i'm done")
}
