package main

import (
	"database/sql"
	"fmt"
	_ "github.com/marcboeker/go-duckdb"  // Import the duckdb driver
	"golang.org/x/text/encoding/charmap" // For common single-byte encodings like ISO-8859-1
	"golang.org/x/text/transform"        // For general encoding transformations
	"io"
	"os"
)

// ConvertCSVToUTF8 converts a CSV file from a specified source encoding to UTF-8.
func ConvertCSVToUTF8(inputFile, outputFile string, sourceEncoding string) error {
	// Open the input file
	in, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer in.Close()

	// Open the output file
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Determine the decoder based on the source encoding
	var decoder transform.Transformer
	switch sourceEncoding {
	case "iso-8859-1", "latin1":
		decoder = charmap.ISO8859_1.NewDecoder()
	// Add more cases for other encodings if needed (e.g., charmap.Windows1252.NewDecoder())
	default:
		return fmt.Errorf("unsupported source encoding: %s", sourceEncoding)
	}

	// Create a reader that decodes the input stream
	reader := transform.NewReader(in, decoder)

	// Copy the decoded (UTF-8) content to the output file
	_, err = io.Copy(out, reader)
	if err != nil {
		return fmt.Errorf("failed to copy and convert data: %w", err)
	}

	return nil
}

func main() {
	db, err := sql.Open("duckdb", "") // In-memory database
	if err != nil {
		panic(err)
	}
	defer db.Close()

	csvFilePath := "C:\\Users\\ggivl\\Documents\\GoDevelopment\\GoMusicLibraryGUIApp\\csv_output_full.csv"
	utf8CsvFilePath := "output_utf8.csv"
	tableName := "music_library_table"

	err = ConvertCSVToUTF8(csvFilePath, utf8CsvFilePath, "iso-8859-1")
	if err != nil {
		fmt.Printf("Error converting CSV: %v\n", err)
		return
	}
	fmt.Println("CSV file converted to UTF-8 successfully!")

	// Use the converted UTF-8 file and enable ignore_errors to skip problematic rows
	query := `CREATE TABLE ` + tableName + ` AS SELECT * FROM read_csv_auto('` + utf8CsvFilePath + `', ignore_errors=true);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	fmt.Println("Table created successfully from UTF-8 CSV!")

	// Verify the table was created and show some stats
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
	if err != nil {
		fmt.Printf("Error querying table: %v\n", err)
	} else {
		fmt.Printf("Successfully loaded %d rows into the table.\n", count)
	}

	// Show a few sample rows
	fmt.Println("\nFirst 3 rows:")
	rows, err := db.Query("SELECT * FROM " + tableName + " LIMIT 3")
	if err != nil {
		fmt.Printf("Error selecting sample rows: %v\n", err)
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		fmt.Printf("Error getting columns: %v\n", err)
		return
	}
	fmt.Printf("Columns: %v\n", columns)

	// Print sample rows
	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		fmt.Printf("Row: %v\n", values)
	}
}
