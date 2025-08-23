package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	c.IndentedJSON(http.StatusOK, musicFiles)
}

func read_write_json(inputFile, outputFile string) []FileInfo {
	// Read the input JSON file which has a top-level object with a "files" array
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("error reading %s: %v\n", inputFile, err)
		return nil
	}

	var fl fileList
	if err := json.Unmarshal(data, &fl); err != nil {
		fmt.Printf("error parsing %s: %v\n", inputFile, err)
		return nil
	}

	// Optionally write a pretty-printed flat array to outputFile for debugging/consumption
	out, err := json.MarshalIndent(fl.Files, "", "  ")
	if err != nil {
		fmt.Printf("error marshaling output: %v\n", err)
		return fl.Files
	}
	if err := os.WriteFile(outputFile, out, 0644); err != nil {
		fmt.Printf("error writing %s: %v\n", outputFile, err)
	}

	return fl.Files
}

func main() {
	musicFiles = read_write_json("data.json", "output.json")
	fmt.Printf("found %d\n", len(musicFiles))
	router := gin.Default()
	router.GET("/musicfileinfo", getMusicFileInfo)
	router.Run("localhost:8080")
}
