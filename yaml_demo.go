package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// Config represents the structure of config.yml
type Config struct {
	SoftwareAuthor  string     `yaml:"SoftwareAuthor"`
	SoftwareCompany string     `yaml:"SoftwareCompany"`
	Phone           int        `yaml:"Phone"`
	Website         string     `yaml:"Website"`
	Database        Database   `yaml:"Database"`
	GoDatabase      GoDatabase `yaml:"GoDatabase"`
	CSV             CSV        `yaml:"CSV"`
	FilePaths       FilePaths  `yaml:"FilePaths"`
	Skills          []Skill    `yaml:"Skills"`
}

// Database represents the database configuration
type Database struct {
	SQLite           bool   `yaml:"sqlite"`
	DatabaseFilename string `yaml:"database_filename"`
}
type GoDatabase struct {
	DuckDB           bool   `yaml:"duckdb"`
	DatabaseFilename string `yaml:"go_database_filename"`
}

// CSV represents the CSV configuration
type CSV struct {
	CSVFilename string `yaml:"csv_filename"`
}

// FilePaths represents the file paths configuration
type FilePaths struct {
	MusicLibraryPath       string `yaml:"music_library_path"`
	PythonCodeRepoLocation string `yaml:"python_code_repo_location"`
	GoCodeRepoLocation     string `yaml:"go_code_repo_location"`
}

// Skill represents a skill entry
type Skill struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// TestYAMLParsing tests loading and parsing the config.yml file
func TestYAMLParsing() {
	// Read the YAML file
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Error reading config.yml: %v", err)
	}

	// Parse the YAML into our struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	// Print the parsed configuration
	fmt.Printf("Software Author: %s\n", config.SoftwareAuthor)
	fmt.Printf("Software Company: %s\n", config.SoftwareCompany)
	fmt.Printf("Phone: %d\n", config.Phone)
	fmt.Printf("Website: %s\n", config.Website)
	fmt.Printf("Database SQLite: %t\n", config.Database.SQLite)
	fmt.Printf("Database Filename: %s\n", config.Database.DatabaseFilename)
	fmt.Printf("GoDatabase DuckDB: %t\n", config.GoDatabase.DuckDB)
	fmt.Printf("GoDatabase Filename: %s\n", config.GoDatabase.DatabaseFilename)
	fmt.Printf("CSV Filename: %s\n", config.CSV.CSVFilename)
	fmt.Printf("Music Library Path: %s\n", config.FilePaths.MusicLibraryPath)
	fmt.Printf("Code Repo Location: %s\n", config.FilePaths.PythonCodeRepoLocation)
	fmt.Printf("Go Code Repo Location: %s\n", config.FilePaths.GoCodeRepoLocation)

	fmt.Println("\nSkills:")
	for i, skill := range config.Skills {
		fmt.Printf("  %d. %s: %s\n", i+1, skill.Name, skill.Description)
	}
}

// Main function to run the YAML test
func main() {
	fmt.Println("Testing YAML configuration parsing...")
	TestYAMLParsing()
	fmt.Println("\nYAML parsing test completed successfully!")
}
