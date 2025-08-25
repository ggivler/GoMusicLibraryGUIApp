package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"example/web-service-gin/logger"
	"github.com/gin-gonic/gin"
)

// musicFiles is a slice of FileInfo structs to hold data
var musicFiles []FileInfo

// FileInfo represents the structure for JSON output
type FileInfo struct {
	AlphabetizingLetter string `json:"alphabetizing letter"`
	FullPathToFolder    string `json:"full path to folder"`
	OriginalFilename    string `json:"original filename"`
	SongTitle           string `json:"song title"`
	Voicing             string `json:"voicing"`
	ComposerOrArranger  string `json:"composer or arranger"`
	FileType            string `json:"file type"`
	FileCreateDate      string `json:"file create date"`
	LibraryType         string `json:"library type"`
}

// wrapper for the top-level JSON object in data.json
type fileList struct {
	Files []FileInfo `json:"files"`
}

func getMusicFileInfo(c *gin.Context) {
	logger.Info("getMusicFileInfo endpoint called")
	c.IndentedJSON(http.StatusOK, musicFiles)
}

func read_write_json(inputFile, outputFile string) []FileInfo {
	// Read the input JSON file which has a top-level object with a "files" array
	data, err := os.ReadFile(inputFile)
	if err != nil {
		logger.Errorf("error reading %s: %v", inputFile, err)
		return nil
	}

	var fl fileList
	if err := json.Unmarshal(data, &fl); err != nil {
		logger.Errorf("error parsing %s: %v", inputFile, err)
		return nil
	}

	// Optionally write a pretty-printed flat array to outputFile for debugging/consumption
	out, err := json.MarshalIndent(fl.Files, "", "  ")
	if err != nil {
		logger.Errorf("error marshaling output: %v", err)
		return fl.Files
	}
	if err := os.WriteFile(outputFile, out, 0644); err != nil {
		logger.Errorf("error writing %s: %v", outputFile, err)
	}

	return fl.Files
}
func getWorkingDirectory() string {
	// Get the path of the current executable
	exePath, err := os.Executable()
	if err != nil {
		logger.Errorf("error getting executable path: %v", err)
		return ""
	}

	// Get the directory containing the executable
	dir := filepath.Dir(exePath)
	return dir
}

func main() {
	// Initialize the global logger
	err := logger.Init("music-api.log")
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close() // Ensure the log file is closed when the main function exits

	// Set Gin's output to use the same log file
	gin.DefaultWriter = logger.GetWriter()

	musicFiles = read_write_json("data.json", "output.json")
	logger.Infof("found %d music files", len(musicFiles))

	// Write log messages using the global logger
	logger.Info("Application started.")
	logger.Infof("Processing data for user: %s", "John Doe")
	workingdir := getWorkingDirectory()
	logger.Infof("executable directory: %s", workingdir)

	router := gin.Default()
	router.GET("/musicfileinfo", getMusicFileInfo)
	logger.Info("Starting server on localhost:8080")
	router.Run("localhost:8080")
}
