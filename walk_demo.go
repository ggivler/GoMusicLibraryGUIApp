package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	_ "github.com/marcboeker/go-duckdb"
)

// DuckDBDatabase represents a DuckDB database connection
type DuckDBDatabase struct {
	DbName     string
	Connection *sql.DB
}

// NewDuckDBDatabase creates a new DuckDB database instance
func NewDuckDBDatabase(dbName string) (*DuckDBDatabase, error) {
	if dbName == "" {
		dbName = "musiclibrary.duckdb"
	}

	db := &DuckDBDatabase{
		DbName: dbName,
	}

	err := db.connect()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DuckDBDatabase) logSQLCallback(statement string) {
	fmt.Printf("Executing SQL statement: %s\n", statement)
}

func (db *DuckDBDatabase) connect() error {
	var err error
	db.Connection, err = sql.Open("duckdb", db.DbName)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}
	return nil
}

func (db *DuckDBDatabase) Close() error {
	if db.Connection != nil {
		return db.Connection.Close()
	}
	return nil
}

func (db *DuckDBDatabase) CreateTable(tableName, columns string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, columns)
	_, err := db.Connection.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}
	return nil
}

func (db *DuckDBDatabase) InsertData(tableName string, data []interface{}) error {
	placeholders := strings.Repeat("?,", len(data))
	placeholders = strings.TrimSuffix(placeholders, ",")

	fmt.Printf("Data: %+v\n", data)

	query := fmt.Sprintf("INSERT INTO %s VALUES (%s)", tableName, placeholders)
	fmt.Printf("Query: %s\n", query)

	_, err := db.Connection.Exec(query, data...)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	return nil
}

func (db *DuckDBDatabase) FetchAll(tableName, condition string) ([][]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s %s", tableName, condition)
	rows, err := db.Connection.Query(query)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		result = append(result, values)
	}

	return result, nil
}

func (db *DuckDBDatabase) ExecuteQuery(query string, params ...interface{}) error {
	_, err := db.Connection.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}
	return nil
}

// FileMethods represents file processing methods
type FileMethods struct {
	BaseDir       string
	DbColumnNames string
	db            *DuckDBDatabase
}

// NewFileMethods creates a new FileMethods instance
func NewFileMethods(baseDir string) *FileMethods {
	return &FileMethods{
		BaseDir:       baseDir,
		DbColumnNames: "id INTEGER PRIMARY KEY, alphabetizing_letter TEXT, full_path_to_folder TEXT, original_filename TEXT, song_title TEXT, voicing TEXT, composer_or_arranger TEXT, file_type TEXT, file_create_date TEXT, library_type TEXT",
	}
}

func (fm *FileMethods) FindFilesRecursively(directoryPath string) ([]string, error) {
	var fileLst []string

	err := filepath.WalkDir(directoryPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileLst = append(fileLst, path)
		}
		return nil
	})

	return fileLst, err
}

func (fm *FileMethods) FindFilesInDirectory(directoryPath string) ([]string, error) {
	var filesList []string

	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("error: Directory not found at '%s': %v", directoryPath, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(directoryPath, entry.Name())
		fmt.Printf("Full path: %s\n", fullPath)
		filesList = append(filesList, fullPath)
	}

	return filesList, nil
}

func (fm *FileMethods) GetFullPathToFolder(directoryPath, fileName string) string {
	return filepath.Join(directoryPath, fileName)
}

func (fm *FileMethods) SplitFilename(filenameToSplit, separator string) []string {
	fmt.Printf("Filename to split: %s\n", filenameToSplit)

	// Use regex to split by space, underscore, or dot
	re := regexp.MustCompile(`[ _.]`)
	splitFileLst := re.Split(filenameToSplit, -1)

	fmt.Printf("Split file list: %+v\n", splitFileLst)
	fmt.Printf("Length of splitFileLst is: %d\n", len(splitFileLst))

	return splitFileLst
}

func (fm *FileMethods) SplitSongTitle(songTitle string) string {
	fmt.Printf("Song title: %s\n", songTitle)

	// Add space before capital letters (except at the beginning)
	// Go doesn't support lookbehind, so we need to use a different approach
	if len(songTitle) <= 1 {
		return songTitle
	}

	// Start with the first character, then iterate through the rest
	result := string(songTitle[0])
	for _, r := range songTitle[1:] {
		if unicode.IsUpper(r) {
			// Add a space before uppercase letters
			result += " "
		}
		result += string(r)
	}

	fmt.Printf("Split song title: %s\n", result)
	return result
}

func (fm *FileMethods) GetAlphabetizerLetterFromFilename(filename string) string {
	if len(filename) == 0 {
		return ""
	}

	alphaLetter := string(filename[0])
	if unicode.IsDigit(rune(filename[0])) {
		for _, r := range filename {
			if unicode.IsLetter(r) {
				alphaLetter = string(r)
				break
			}
		}
	}

	fmt.Printf("Alpha letter: %s\n", alphaLetter)
	return strings.ToUpper(alphaLetter)
}

func (fm *FileMethods) GetExtensionFromFilename(composerArranger string) (string, string) {
	fmt.Printf("Composer/arranger: %s\n", composerArranger)

	ext := filepath.Ext(composerArranger)
	composer := strings.TrimSuffix(composerArranger, ext)

	fmt.Printf("Composer: %s\n", composer)
	fmt.Printf("Extension: %s\n", ext)

	return composer, ext
}

func (fm *FileMethods) GetFileCreationDateFromFilename(dtFilename string) string {
	formattedDt := ""

	info, err := os.Stat(dtFilename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Error: File not found at '%s'\n", dtFilename)
		return formattedDt
	}

	modTime := info.ModTime()
	formattedDt = modTime.Format("2006-01-02")

	fmt.Printf("Formatted date: %s\n", formattedDt)
	return formattedDt
}

func (fm *FileMethods) GetVoicingFromParsedFilename(listToParse interface{}) string {
	fmt.Printf("List to parse: %+v\n", listToParse)

	voicing := ""
	voicings := []string{
		"SATB", "SSATB", "SSAATTBB", "SAB", "SA", "SSA", "SAA",
		"SSAA", "TB", "TTB", "TBB", "TTBB",
	}

	switch v := listToParse.(type) {
	case []string:
		for _, entry := range voicings {
			for _, item := range v {
				if strings.Contains(strings.ToUpper(item), entry) {
					voicing = entry
					return voicing
				}
			}
		}
		voicing = "UNKNOWN"
	case string:
		upperStr := strings.ToUpper(v)
		for _, entry := range voicings {
			if strings.Contains(upperStr, entry) {
				voicing = entry
				return voicing
			}
		}
		voicing = fmt.Sprintf("Unable to find any of these voicings %+v in the following string '%s'", voicings, v)
	default:
		voicing = "Error: The variable listToParse must be a string or a slice"
	}

	return voicing
}

func (fm *FileMethods) GetFileTypeFromFilePath(filepath string) string {
	ext := strings.ToLower(filepath)

	switch {
	case strings.Contains(ext, "pdf"):
		return "PDF"
	case strings.Contains(ext, "mp3"):
		return "MP3"
	case strings.Contains(ext, "ogg"):
		return "OGG"
	case strings.Contains(ext, "wma"):
		return "WMA"
	case strings.Contains(ext, "mp4"):
		return "MP4"
	default:
		return "UNKNOWN"
	}
}

func (fm *FileMethods) GetLibraryTypeFromFilePath(filepath string) string {
	libraryType := strings.ToLower(filepath)

	switch {
	case strings.Contains(libraryType, "christmas"):
		return "Christmas"
	case strings.Contains(libraryType, "spring"):
		return "Spring"
	case strings.Contains(libraryType, "repertoire"):
		return "Repertoire"
	default:
		return "UNKNOWN"
	}
}

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

// MasterJSONFile represents the complete JSON structure
type MasterJSONFile struct {
	Files []FileInfo `json:"files"`
}

func (fm *FileMethods) WriteCSVOutputFile(inputJSON MasterJSONFile, outputPath, outputFilename string, fieldnames []string) error {
	csvFilename := filepath.Join(outputPath, outputFilename)
	fmt.Printf("CSV filename: %s\n", csvFilename)

	file, err := os.Create(csvFilename)
	if err != nil {
		fmt.Printf("I/O error: %v\n", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write(fieldnames); err != nil {
		return err
	}

	// Write data rows
	fmt.Printf("Input JSON files: %+v\n", inputJSON.Files)
	for _, fileInfo := range inputJSON.Files {
		row := []string{
			fileInfo.AlphabetizingLetter,
			fileInfo.FullPathToFolder,
			fileInfo.OriginalFilename,
			fileInfo.SongTitle,
			fileInfo.Voicing,
			fileInfo.ComposerOrArranger,
			fileInfo.FileType,
			fileInfo.FileCreateDate,
			fileInfo.LibraryType,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	fmt.Printf("Wrote CSV output to %s\n", csvFilename)
	return nil
}

func (fm *FileMethods) ImportCSVFileIntoDB(csvFilename, databaseFilename, tableName string) error {
	var err error
	fm.db, err = NewDuckDBDatabase(databaseFilename)
	if err != nil {
		return err
	}
	defer fm.db.Close()

	err = fm.db.CreateTable("music_library", fm.DbColumnNames)
	if err != nil {
		return err
	}

	file, err := os.Open(csvFilename)
	if err != nil {
		fmt.Printf("Error: The file '%s' was not found: %v\n", csvFilename, err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) > 0 {
		header := records[0]
		fmt.Printf("Header: %+v\n", header)

		for i, row := range records[1:] {
			fmt.Printf("Row: %+v\n", row)

			// Insert ID at the beginning
			data := make([]interface{}, len(row)+1)
			data[0] = i + 1
			for j, v := range row {
				data[j+1] = v
			}
			fmt.Printf("Data with ID: %+v\n", data)

			err = fm.db.InsertData(tableName, data)
			if err != nil {
				fmt.Printf("Database Error has occurred: %v\n", err)
				continue
			}
		}

		allRows, err := fm.db.FetchAll("music_library", "")
		if err != nil {
			return err
		}
		fmt.Printf("All rows: %+v\n", allRows)
		fmt.Printf("Number of rows: %d\n", len(allRows))
	}

	return nil
}

func main() {
	// Command line arguments
	dirpath := flag.String("d", `C:\Users\ggivl\Documents\PythonDevelopment\FortyNinersDevelopment\49ersMusicLibrary`, "Path to the directory of the files to parsed")
	extension := flag.String("e", ".pdf", "Extension of the files to parse")
	outputCSV := flag.String("o", "csv_output_full.csv", "CSV output file")
	dbname := flag.String("b", "musiclibrary.duckdb", "Database filename")
	flag.Parse()

	fileExt := *extension

	keywords := []string{
		"alphabetizing letter",
		"full path to folder",
		"original filename",
		"song title",
		"voicing",
		"composer or arranger",
		"file type",
		"file create date",
		"library type",
	}

	fmt.Printf("Directory path: %s\n", *dirpath)

	fileMethods := NewFileMethods(*dirpath)

	// Find files recursively
	fileLst, err := fileMethods.FindFilesRecursively(*dirpath)
	if err != nil {
		log.Fatalf("Error finding files: %v", err)
	}

	fmt.Printf("File list: %+v\n", fileLst)

	var pdfFileLst []string
	for _, filename := range fileLst {
		if strings.HasSuffix(filename, fileExt) {
			fmt.Printf("PDF filename: %s\n", filename)
			pdfFileLst = append(pdfFileLst, filename)
		}
	}

	var jsonFileLst []FileInfo
	for _, pdfFilepath := range pdfFileLst {
		var jsonFileInfo FileInfo
		pdfFilename := filepath.Base(pdfFilepath)
		splitFilename := fileMethods.SplitFilename(pdfFilename, "_")

		fmt.Printf("Split filename: %+v\n", splitFilename)

		// Initialize default values
		dirPath := "UNKNOWN"
		rawTitle := "UNKNOWN"
		rawVoicing := "UNKNOWN"
		rawComposerArranger := "UNKNOWN"

		k := len(splitFilename)
		if k > 0 {
			dirPath = splitFilename[0]
		}
		if k > 1 {
			rawTitle = splitFilename[1]
		}
		if k > 2 {
			rawVoicing = splitFilename[2]
		}
		if k > 3 {
			// Extract composer/arranger, but remove file extension if it's the last element
			composerWithExt := splitFilename[3]
			if k > 4 {
				// If there are more elements, this is not the last one, so it shouldn't have extension
				rawComposerArranger = composerWithExt
			} else {
				// This might be the last element, so remove extension
				composer, _ := fileMethods.GetExtensionFromFilename(composerWithExt)
				if composer != "" {
					rawComposerArranger = composer
				} else {
					rawComposerArranger = composerWithExt
				}
			}
		}

		_ = dirPath
		_ = rawVoicing

		jsonFileInfo.AlphabetizingLetter = fileMethods.GetAlphabetizerLetterFromFilename(filepath.Base(pdfFilepath))
		jsonFileInfo.FullPathToFolder = filepath.Dir(pdfFilepath)
		jsonFileInfo.OriginalFilename = filepath.Base(pdfFilepath)
		jsonFileInfo.FileCreateDate = fileMethods.GetFileCreationDateFromFilename(pdfFilepath)
		jsonFileInfo.SongTitle = fileMethods.SplitSongTitle(rawTitle)
		jsonFileInfo.Voicing = fileMethods.GetVoicingFromParsedFilename(splitFilename)
		jsonFileInfo.ComposerOrArranger = rawComposerArranger
		jsonFileInfo.FileType = fileMethods.GetFileTypeFromFilePath(pdfFilepath)
		// jsonFileInfo.LibraryType = fileMethods.GetLibraryTypeFromFilePath(pdfFilepath)
		jsonFileInfo.LibraryType = fmt.Sprintf("%s,", fileMethods.GetLibraryTypeFromFilePath(pdfFilepath))

		fmt.Printf("PDF filepath: %s\n", pdfFilepath)
		fmt.Printf("JSON file info: %+v\n", jsonFileInfo)

		jsonFileLst = append(jsonFileLst, jsonFileInfo)
	}

	masterJSONFile := MasterJSONFile{
		Files: jsonFileLst,
	}

	fmt.Printf("Number of PDF files: %d\n", len(pdfFileLst))
	fmt.Printf("Master JSON file: %+v\n", masterJSONFile)

	// Write CSV output
	err = fileMethods.WriteCSVOutputFile(masterJSONFile, ".", *outputCSV, keywords)
	if err != nil {
		log.Printf("Error writing CSV: %v", err)
	}

	// Import CSV to database
	err = fileMethods.ImportCSVFileIntoDB(*outputCSV, *dbname, "music_library")
	if err != nil {
		log.Printf("Error importing to database: %v", err)
	}

	// Write JSON output
	jsonOutPath := "output_file_full.json"
	jsonData, err := json.MarshalIndent(masterJSONFile, "", "    ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}

	err = os.WriteFile(jsonOutPath, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing JSON file: %v", err)
		return
	}

	fmt.Printf("JSON output written to: %s\n", jsonOutPath)
}
