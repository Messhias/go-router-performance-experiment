package main

import (
	"log"
	"messhias/router-expirement/internal/balancer"
	"messhias/router-expirement/internal/router"
	"os"
	"strings"
)

func main() {

	a := strings.TrimSpace(os.Getenv("UPSTREAM_A_URL"))
	b := strings.TrimSpace(os.Getenv("UPSTREAM_B_URL"))

	if a == "" || b == "" {
		log.Fatal("UPSTREAM_A_URL and UPSTREAM_B_URL must be set and non-empty")
	}

	bal, err := balancer.NewBalancer([]string{a, b})

	if err != nil {
		log.Fatal(err)
	}

	engine := router.NewEngine(bal, nil)

	serverAddress := os.Getenv("SERVER_ADDR")

	err = engine.Run(serverAddress)

	if err != nil {
		panic(err)
	}
}
