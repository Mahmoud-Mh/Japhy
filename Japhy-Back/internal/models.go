package internal

type Breed struct {
	ID                        int    `json:"id" db:"id"`
	Species                   string `json:"species" db:"species"`
	PetSize                   string `json:"pet_size" db:"pet_size"`
	Name                      string `json:"name" db:"name"`
	AverageMaleAdultWeight    int    `json:"average_male_adult_weight" db:"average_male_adult_weight"`
	AverageFemaleAdultWeight  int    `json:"average_female_adult_weight" db:"average_female_adult_weight"`
}

type BreedSearchParams struct {
	Species   string `json:"species"`
	PetSize   string `json:"pet_size"`
	MinWeight *int   `json:"min_weight"`
	MaxWeight *int   `json:"max_weight"`
} 