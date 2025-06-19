package internal

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type BreedsService struct {
	db *sql.DB
}

func NewBreedsService(db *sql.DB) *BreedsService {
	return &BreedsService{db: db}
}

// GetAll returns all breeds with optional search filters
func (s *BreedsService) GetAll(params BreedSearchParams) ([]Breed, error) {
	query := "SELECT id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight FROM breeds WHERE 1=1"
	args := []interface{}{}

	if params.Species != "" {
		query += " AND species = ?"
		args = append(args, params.Species)
	}

	if params.PetSize != "" {
		query += " AND pet_size = ?"
		args = append(args, params.PetSize)
	}

	if params.MinWeight != nil {
		query += " AND (average_male_adult_weight >= ? OR average_female_adult_weight >= ?)"
		args = append(args, *params.MinWeight, *params.MinWeight)
	}

	if params.MaxWeight != nil {
		query += " AND (average_male_adult_weight <= ? OR average_female_adult_weight <= ?)"
		args = append(args, *params.MaxWeight, *params.MaxWeight)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breeds []Breed
	for rows.Next() {
		var breed Breed
		err := rows.Scan(&breed.ID, &breed.Species, &breed.PetSize, &breed.Name, 
			&breed.AverageMaleAdultWeight, &breed.AverageFemaleAdultWeight)
		if err != nil {
			return nil, err
		}
		breeds = append(breeds, breed)
	}

	return breeds, nil
}

// GetByID returns a breed by ID
func (s *BreedsService) GetByID(id int) (*Breed, error) {
	var breed Breed
	query := "SELECT id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight FROM breeds WHERE id = ?"
	
	err := s.db.QueryRow(query, id).Scan(&breed.ID, &breed.Species, &breed.PetSize, &breed.Name,
		&breed.AverageMaleAdultWeight, &breed.AverageFemaleAdultWeight)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &breed, nil
}

// Create creates a new breed
func (s *BreedsService) Create(breed *Breed) error {
	query := `INSERT INTO breeds (id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight) 
			  VALUES (?, ?, ?, ?, ?, ?)`
	
	_, err := s.db.Exec(query, breed.ID, breed.Species, breed.PetSize, breed.Name,
		breed.AverageMaleAdultWeight, breed.AverageFemaleAdultWeight)
	
	return err
}

// Update updates an existing breed
func (s *BreedsService) Update(id int, breed *Breed) error {
	query := `UPDATE breeds SET species = ?, pet_size = ?, name = ?, 
			  average_male_adult_weight = ?, average_female_adult_weight = ? 
			  WHERE id = ?`
	
	result, err := s.db.Exec(query, breed.Species, breed.PetSize, breed.Name,
		breed.AverageMaleAdultWeight, breed.AverageFemaleAdultWeight, id)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("breed with id %d not found", id)
	}
	
	return nil
}

// Delete deletes a breed by ID
func (s *BreedsService) Delete(id int) error {
	query := "DELETE FROM breeds WHERE id = ?"
	
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("breed with id %d not found", id)
	}
	
	return nil
}

// ImportFromCSV imports breeds from the CSV file
func (s *BreedsService) ImportFromCSV(filename string) error {
	// Check if breeds already exist
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM breeds").Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return nil // Already imported
	}
	
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	
	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 6 {
			continue
		}
		
		// Remove quotes from fields
		for j := range record {
			record[j] = strings.Trim(record[j], "\"")
		}
		
		id, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}
		
		maleWeight, _ := strconv.Atoi(record[4])
		femaleWeight, _ := strconv.Atoi(record[5])
		
		breed := &Breed{
			ID:                       id,
			Species:                  record[1],
			PetSize:                  record[2],
			Name:                     record[3],
			AverageMaleAdultWeight:   maleWeight,
			AverageFemaleAdultWeight: femaleWeight,
		}
		
		err = s.Create(breed)
		if err != nil {
			return fmt.Errorf("failed to insert breed %s: %v", breed.Name, err)
		}
	}
	
	return nil
} 