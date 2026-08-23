package main

import (
	"log"
	"net/http"

	"github.com/Mikhail-Tal63/Orbit/configs"
	"github.com/Mikhail-Tal63/Orbit/internal/auth"
	"github.com/Mikhail-Tal63/Orbit/internal/database"
	db "github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/Mikhail-Tal63/Orbit/internal/driver"
	"github.com/Mikhail-Tal63/Orbit/internal/vehicle"
	"github.com/Mikhail-Tal63/Orbit/internal/websocket"
	"github.com/Mikhail-Tal63/Orbit/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// ── Config ──────────────────────────────────────────────
	cfg := configs.Load()

	// ── Database ────────────────────────────────────────────
	pool := database.Connect(*cfg)
	defer pool.Close()

	redisClient, err := database.NewReisClient(cfg.RedisAddrs)

	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}

	defer redisClient.Close()
	
	queries := db.New(pool)

	// ── Repositories ────────────────────────────────────────
	authRepo := auth.NewAuthRepository(queries)
	driverRepo := driver.NewDriverRepository(queries)
	vehicleRepo := vehicle.NewVechileRepository(queries)

	// ── Services ────────────────────────────────────────────
	authService := auth.NewAuthService(
		authRepo,
		[]byte(cfg.JWTSecret),
	)

	driverService := driver.NewDriverService(
		driverRepo,
		authRepo,
		vehicleRepo,
		pool,
	)

	// ── Handlers ────────────────────────────────────────────
	authHandler := auth.NewAuthHandler(authService)
	driverHandler := driver.NewDriverHandler(driverService)

	// ── WebSocket ───────────────────────────────────────────
	hub := websocket.NewHub()
	go hub.Run()

	// ── Router ──────────────────────────────────────────────
	router := mux.NewRouter()

	// Public routes
	public := router.PathPrefix("/api/v1").Subrouter()
	authHandler.AuthRouter(public)

	// Protected routes
	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	authHandler.ProtectedRouter(protected)
	driverHandler.DriverRouter(protected)

	// WebSocket
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServerWS(
			hub,
			[]byte(cfg.JWTSecret),
			w,
			r,
		)
	})

	// ── Start Server ────────────────────────────────────────
	addr := ":" + cfg.Port

	log.Printf("🚀 Orbit API listening on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
