package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// uploadFiles performs a multipart file upload to the given URL with the provided file paths.
func uploadFiles(url string, filePaths []string) error {
	// Create a buffer to write our multipart form data
	var requestBody bytes.Buffer
	// Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)
	// Loop through the file paths and add each file as a part
	for _, path := range filePaths {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		part, err := multipartWriter.CreateFormFile("my_files", filepath.Base(path))
		if err != nil {
			return err
		}
		if _, err = io.Copy(part, file); err != nil {
			return err
		}
	}
	// Close the multipart writer to set the terminating boundary
	if err := multipartWriter.Close(); err != nil {
		return err
	}
	// Create the HTTP request
	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return err
	}

	// Set up the username and password for basic auth.
	username := "user"
	password := "pass"
	request.SetBasicAuth(username, password)

	// Set the content type, this must be set to the value that includes the boundary.
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	// Perform the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check the response
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response status: %s", response.Status)
	}

	return nil
}

func main() {
	// Define the URL to upload the files
	uploadURL := "http://127.0.0.1:8080/upload"
	// Define the file paths to upload
	filePaths := []string{
		"path/to/your/second/file.png",
		// Add more files here
	}
	// Perform the file upload
	if err := uploadFiles(uploadURL, filePaths); err != nil {
		log.Fatalf("File upload failed: %v", err)
	}
	fmt.Println("Files uploaded successfully!")
}
