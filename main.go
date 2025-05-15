package main

import (
	"fmt"

	"github.com/tobib-dev/gator/internal/config"
)

func main() {
	config.ReadJsonFile()
	config.SetUser("Tobi")
	file := config.ReadJsonFile()
	fmt.Println(file.DbUrl)
	fmt.Println(file.CurrentUserName)
}
