package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

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

func getMusicFileInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, musicFiles)
}

func main() {
	router := gin.Default()
	router.GET("/musicfileinfo", getMusicFileInfo)
	router.Run("localhost:8080")
}
