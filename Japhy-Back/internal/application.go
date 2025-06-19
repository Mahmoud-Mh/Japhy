package internal

import (
	"database/sql"

	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
)

type App struct {
	logger        *charmLog.Logger
	breedsHandler *BreedsHandler
}

func NewApp(logger *charmLog.Logger, db *sql.DB) *App {
	breedsService := NewBreedsService(db)
	breedsHandler := NewBreedsHandler(breedsService)

	// Import breeds from CSV on startup
	err := breedsService.ImportFromCSV("breeds.csv")
	if err != nil {
		logger.Error("Failed to import breeds from CSV", "error", err)
	} else {
		logger.Info("Breeds imported successfully from CSV")
	}

	return &App{
		logger:        logger,
		breedsHandler: breedsHandler,
	}
}

func (a *App) RegisterRoutes(r *mux.Router) {
	// Breeds CRUD endpoints
	r.HandleFunc("/breeds", a.breedsHandler.GetBreeds).Methods("GET")
	r.HandleFunc("/breeds/{id}", a.breedsHandler.GetBreed).Methods("GET")
	r.HandleFunc("/breeds", a.breedsHandler.CreateBreed).Methods("POST")
	r.HandleFunc("/breeds/{id}", a.breedsHandler.UpdateBreed).Methods("PUT")
	r.HandleFunc("/breeds/{id}", a.breedsHandler.DeleteBreed).Methods("DELETE")
}
