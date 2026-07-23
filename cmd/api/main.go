package main

import (
	"fmt"

	"github.com/Mikhail-Tal63/Orbit/configs"
	"github.com/Mikhail-Tal63/Orbit/internal/auth"
	"github.com/Mikhail-Tal63/Orbit/internal/database"
	db "github.com/Mikhail-Tal63/Orbit/internal/db"
)

func main() {
	cfg := configs.Load()

	pool := database.Connect(*cfg)
	defer pool.Close()

	queries := db.New(pool)

	authRepo := auth.NewAuthRepository(queries)

authService := auth.NewAuthService(
    authRepo,
    []byte(cfg.JWTSecret),
)
print(authService)
	fmt.Println("Orbit API starting...")
}
