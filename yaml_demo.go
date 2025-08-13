package main

import (
	"encoding/csv"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
func TestYAMLParsing() Config {
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
	return config
}

// Writing the config file after editing the config object
func writeConfig(config Config, filename string) Config {
	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("Error marshalling config: %v", err)
	}
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		log.Fatalf("Error writing config.yml: %v", err)
	}
	return config
}

func readCSV(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		log.Fatalf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		log.Fatalf("Error reading file: %v\n", err)
		return nil
	}

	fmt.Printf("Read %d records\n", len(records))
	for i, record := range records {
		fmt.Printf("Record %d: %v\n", i+1, record)
		// Also show as comma-separated values:
		fmt.Printf("  CSV format: %s\n", strings.Join(record, ","))

	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Printf("Error seeking file: %v\n", err)
		log.Fatalf("Error seeking file: %v\n", err)
		return nil
	}

	reader = csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			log.Fatalf("Error reading file: %v\n", err)
			return nil
		}
		fmt.Println(strings.Join(record, ","))
	}
	return records
}

func writeCSV(filename string, records [][]string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(records)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	}
	log.Printf("Wrote %d records to the file %s\n", len(records), filename)
}

// Main function to run the YAML test
func main() {
	fmt.Println("Testing YAML configuration parsing...")
	var conf Config
	conf = TestYAMLParsing()
	fmt.Println("\nYAML parsing test completed successfully!")
	fmt.Printf("Music Library Path: %s\n", conf.FilePaths.MusicLibraryPath)
	newMusicLibraryPath := "C:\\Users\\ggivl\\Documents\\PythonDevelopment\\FortyNinersDevelopment\\49ers-musiclibrary"
	conf.FilePaths.MusicLibraryPath = newMusicLibraryPath
	new_conf := writeConfig(conf, "new_config.yml")
	fmt.Printf("New Music Library Path %s\n", new_conf.FilePaths.MusicLibraryPath)
	records := readCSV("csv_output_full.csv")
	fmt.Printf("Found %d Records in CSV file %s\n", len(records), "csv_output_full.csv")
	writeCSV("new_csv_output_full.csv", records)
}
