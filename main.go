package main

import (
	"fmt"
	"sftp_test/sftp_uploader"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	// Define a function to create a new SFTP client
	newSFTPClient := func(conn *ssh.Client) (*sftp.Client, error) {
		return sftp.NewClient(conn)
	}

	client, err := sftp_uploader.Connect("tester", "password", "127.0.0.1", 22, ssh.Dial, newSFTPClient)
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
