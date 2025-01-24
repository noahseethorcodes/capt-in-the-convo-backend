package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Define the API endpoint
	apiURL := "http://localhost:8080/admin/tags" // Replace with your actual endpoint

	// Open the preset_tags.txt file
	file, err := os.Open("preset_tags.txt")
	if err != nil {
		// Fallback to the secondary path
		file, err = os.Open("utils/preset_tags.txt")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Get the tag name from the line
		tagName := scanner.Text()

		// Skip empty lines
		if tagName == "" {
			continue
		}

		// Create the request body
		requestBody, err := json.Marshal(map[string]string{
			"name": tagName,
		})
		if err != nil {
			fmt.Println("Error creating request body:", err)
			continue
		}

		// Send the POST request
		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Printf("Error sending request for tag '%s': %v\n", tagName, err)
			continue
		}
		defer resp.Body.Close()

		// Check the response status
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("Tag '%s' created successfully.\n", tagName)
		} else {
			fmt.Printf("Failed to create tag '%s'. Status: %d\n", tagName, resp.StatusCode)
		}
	}

	// Check for errors during file scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
