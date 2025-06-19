CREATE TABLE IF NOT EXISTS breeds (
    id INT PRIMARY KEY,
    species VARCHAR(16) NOT NULL,
    pet_size VARCHAR(16) NOT NULL,
    name VARCHAR(64) NOT NULL,
    average_male_adult_weight INT,
    average_female_adult_weight INT
); 