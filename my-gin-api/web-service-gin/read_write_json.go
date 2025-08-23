package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// 1. Read the JSON file
	filePath := "data.json" // Replace with your JSON file path
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// 2. Unmarshal the JSON into a map[string]interface{}
	var data map[string]interface{}
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// 3. Access and process the data
	// You can access values by key and perform type assertions
	//if name, ok := data["name"].(string); ok {
	//	fmt.Printf("Name: %s\n", name)
	//}
	//
	//if age, ok := data["age"].(float64); ok { // JSON numbers are unmarshaled as float64
	//	fmt.Printf("Age: %.0f\n", age)
	//}
	//
	//if hobbies, ok := data["hobbies"].([]interface{}); ok {
	//	fmt.Println("Hobbies:")
	//	for _, hobby := range hobbies {
	//		if h, ok := hobby.(string); ok {
	//			fmt.Printf("- %s\n", h)
	//		}
	//	}
	//}

	// You can also iterate through the map
	fmt.Println("\nAll data %d:", len(data))

	for key, value := range data {
		fmt.Printf("%s: %v (Type: %T)\n", key, value, value)
	}

	// Marshal the map into JSON bytes
	jsonData, err := json.MarshalIndent(data, "", "  ") // Use MarshalIndent for pretty printing
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	fileList := data["files"].([]interface{})

	fmt.Println(fmt.Sprintf("\nAll data: %d", len(fileList)))

	for i, file := range fileList {
		if fileMap, ok := file.(map[string]interface{}); ok {
			// Extract some key fields from each file object
			songTitle := ""
			if title, exists := fileMap["song title"]; exists {
				if titleStr, ok := title.(string); ok {
					songTitle = titleStr
				}
			}
			
			originalFilename := ""
			if filename, exists := fileMap["original filename"]; exists {
				if filenameStr, ok := filename.(string); ok {
					originalFilename = filenameStr
				}
			}
			
			fmt.Printf("File %d: %s - %s (Type: %T)\n", i+1, songTitle, originalFilename, file)
		} else {
			fmt.Printf("File %d: Unexpected type %T\n", i+1, file)
		}
	}
	// Write the JSON bytes to a file
	err = ioutil.WriteFile("output.json", jsonData, 0644) // 0644 sets file permissions
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	fmt.Println("JSON data written to output.json")
}
