package main

import (
	"log"

	"app/internal/config"
	"app/internal/database"
	"app/internal/server"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("gagal koneksi database: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("gagal migrasi database: %v", err)
	}
	if err := database.Seed(db, cfg); err != nil {
		log.Fatalf("gagal seeding database: %v", err)
	}
	r := server.NewRouter(cfg, db)
	log.Printf("Delcom Blog API berjalan di :%s (Swagger: /swagger/index.html)", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
