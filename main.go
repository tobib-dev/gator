package main

import (
	"github.com/tobib-dev/gator/internal/config"
)

func main() {
	config.ReadJsonFile()
	config.SetUser("Tobi")
}
