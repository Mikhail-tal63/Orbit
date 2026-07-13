package main

import (
	"fmt"

	"github.com/Mikhail-Tal63/Orbit/configs"
	"github.com/Mikhail-Tal63/Orbit/internal/database"
)

func main() {

	cfg := configs.Load()

	db := database.Connect(*cfg)
	defer db.Close()

	fmt.Println("Orbit API starting...")
}
