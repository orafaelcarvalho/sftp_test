package sftp_uploader

import (
	"fmt"
	"io"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTPClientInterface define a interface para o cliente SFTP
type SFTPClientInterface interface {
	Create(path string) (io.Writer, io.Closer, error) // Retorna io.Writer e io.Closer
	Close() error
}

// SFTPClient é a estrutura principal que usará o cliente SFTP.
type SFTPClient struct {
	client SFTPClientInterface
}

// SFTPClientWrapper é um wrapper para o *sftp.Client, permitindo que ele implemente a interface SFTPClientInterface
type SFTPClientWrapper struct {
	client *sftp.Client
}

// Create implementa a interface SFTPClientInterface para o wrapper de *sftp.Client
func (s *SFTPClientWrapper) Create(path string) (io.Writer, io.Closer, error) {
	// Retorna o arquivo SFTP e o seu closer
	file, err := s.client.Create(path)
	if err != nil {
		return nil, nil, err
	}
	return file, file, nil // file implementa io.Writer e io.Closer
}

// Close implementa a interface SFTPClientInterface para o wrapper de *sftp.Client
func (s *SFTPClientWrapper) Close() error {
	return s.client.Close()
}

// DialFunc é uma função que lida com a conexão SSH
type DialFunc func(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error)

// Connect estabelece uma conexão SFTP.
func Connect(user, password, host string, port int, dial DialFunc, newSFTPClient func(conn *ssh.Client) (*sftp.Client, error)) (*SFTPClient, error) {
	// Configuração do SSH
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Estabelece a conexão SSH
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %w", err)
	}

	// Cria o cliente SFTP
	sftpClient, err := newSFTPClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create SFTP client: %w", err)
	}

	// Retorna o SFTPClient usando o wrapper
	return &SFTPClient{client: &SFTPClientWrapper{client: sftpClient}}, nil
}

// UploadFile faz o upload de um arquivo para o SFTP.
func (s *SFTPClient) UploadFile(data []byte, remotePath string) error {
	writer, closer, err := s.client.Create(remotePath)
	if err != nil {
		return err
	}
	defer closer.Close() // Garantir que Close() seja chamado

	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// Close encerra a conexão SFTP.
func (s *SFTPClient) Close() error {
	return s.client.Close()
}
