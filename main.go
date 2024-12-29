package main

import (
	"fmt"
	"sftp_test/sftp_uploader"

	"golang.org/x/crypto/ssh"
)

func main() {
	client, err := sftp_uploader.Connect("tester", "password", "127.0.0.1", 22, ssh.Dial)
	if err != nil {
		fmt.Printf("Error connecting to SFTP server: %v\n", err)
		return
	}
	defer client.Close()

	// Example data to upload
	fileData := []byte("This is a test file content.")
	err = client.UploadFile(fileData, "test.txt")
	if err != nil {
		fmt.Printf("Error uploading file: %v\n", err)
		return
	}

	fmt.Println("File uploaded successfully")
}
