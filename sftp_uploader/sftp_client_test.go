package sftp_uploader

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para a interface SFTPClientInterface
type MockSFTPClient struct {
	mock.Mock
}

func (m *MockSFTPClient) Create(path string) (io.Writer, io.Closer, error) {
	args := m.Called(path)
	return args.Get(0).(io.Writer), args.Get(1).(io.Closer), args.Error(2)
}

func (m *MockSFTPClient) Close() error {
	return m.Called().Error(0)
}

// Mock para io.Closer
type MockCloser struct {
	mock.Mock
}

func (m *MockCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func TestUploadFile_Success(t *testing.T) {
	// Criando o mock do cliente SFTP
	mockClient := new(MockSFTPClient)

	// Preparando o mock para o método Create
	mockWriter := &bytes.Buffer{} // Usamos um bytes.Buffer como mock de io.Writer
	mockCloser := new(MockCloser) // Criando um mock de io.Closer

	// Configurando as expectativas
	mockClient.On("Create", "remote/path").Return(mockWriter, mockCloser, nil)
	mockCloser.On("Close").Return(nil) // Expectativa de chamada do Close no MockCloser

	// Criando o SFTPClient com o mock
	sftpClient := SFTPClient{client: mockClient}

	// Dados para enviar
	data := []byte("hello world")

	// Chamando o método UploadFile
	err := sftpClient.UploadFile(data, "remote/path")
	assert.NoError(t, err)

	// Verificando se o método Create foi chamado corretamente
	mockClient.AssertExpectations(t)

	// Verificando se o método Close foi chamado corretamente
	mockCloser.AssertExpectations(t)

	// Verificando se os dados foram escritos no arquivo
	assert.Equal(t, "hello world", mockWriter.String())
}

func TestUploadFile_ErrorOnCreate(t *testing.T) {
	// Criando o mock do cliente SFTP
	mockClient := new(MockSFTPClient)

	// Criando mocks válidos para io.Writer e io.Closer
	mockWriter := new(MockWriter) // Um mock válido para io.Writer
	mockCloser := new(MockCloser) // Um mock válido para io.Closer

	// Preparando o mock para o método Create com erro
	mockClient.On("Create", "remote/path").Return(mockWriter, mockCloser, errors.New("create failed"))

	// Criando o SFTPClient com o mock
	sftpClient := SFTPClient{client: mockClient}

	// Dados para enviar
	data := []byte("hello world")

	// Chamando o método UploadFile e verificando o erro retornado
	err := sftpClient.UploadFile(data, "remote/path")
	assert.Error(t, err)
	assert.Equal(t, "create failed", err.Error()) // Verifica se o erro é o esperado

	// Verificando se o método Create foi chamado corretamente
	mockClient.AssertExpectations(t)
}

func TestUploadFile_ErrorOnWrite(t *testing.T) {
	// Criando o mock do cliente SFTP
	mockClient := new(MockSFTPClient)

	// Preparando o mock para o método Create
	mockWriter := new(MockWriter) // Usando um mock customizado para io.Writer
	mockCloser := new(MockCloser)
	mockClient.On("Create", "remote/path").Return(mockWriter, mockCloser, nil)

	// Simulando um erro ao escrever (mock do io.Writer)
	mockWriter.On("Write", []byte("hello world")).Return(0, errors.New("write failed"))

	// Configurando o comportamento esperado para o método Close
	mockCloser.On("Close").Return(nil).Once() // Garantindo que Close seja chamado uma vez

	// Criando o SFTPClient com o mock
	sftpClient := SFTPClient{client: mockClient}

	// Dados para enviar
	data := []byte("hello world")

	// Chamando o método UploadFile e verificando o erro retornado
	err := sftpClient.UploadFile(data, "remote/path")
	assert.Error(t, err)
	assert.Equal(t, "write failed", err.Error()) // Verifica se o erro é o esperado

	// Verificando se o método Create foi chamado corretamente
	mockClient.AssertExpectations(t)

	// Verificando se o método Close foi chamado corretamente
	mockCloser.AssertExpectations(t)
}

func TestUploadFile_ErrorOnClose(t *testing.T) {
	// Mock do cliente SFTP
	mockClient := new(MockSFTPClient)

	// Mock do Writer e Closer
	mockWriter := new(bytes.Buffer) // Usando bytes.Buffer como mock de Writer
	mockCloser := new(MockCloser)

	// Configurando os mocks
	mockClient.On("Create", "remote/path").Return(mockWriter, mockCloser, nil)
	mockCloser.On("Close").Return(errors.New("close failed")) // Simula falha ao fechar

	// Criando o cliente SFTP com o mock
	sftpClient := SFTPClient{client: mockClient}

	// Dados para upload
	data := []byte("hello world")

	// Chamando o método UploadFile
	err := sftpClient.UploadFile(data, "remote/path")

	// Verificando o erro
	assert.Error(t, err)
	assert.Equal(t, "close failed", err.Error()) // Esperado erro "close failed"

	// Verificando as chamadas nos mocks
	mockClient.AssertExpectations(t)
	mockCloser.AssertExpectations(t)
}

func TestUploadFile_EmptyData(t *testing.T) {
	// Criando o mock do cliente SFTP
	mockClient := new(MockSFTPClient)

	// Preparando o mock para o método Create
	mockWriter := &bytes.Buffer{}
	mockCloser := new(MockCloser)

	// Configurando as expectativas
	mockClient.On("Create", "remote/path").Return(mockWriter, mockCloser, nil)
	mockCloser.On("Close").Return(nil)

	// Criando o SFTPClient com o mock
	sftpClient := SFTPClient{client: mockClient}

	// Dados para enviar (vazio)
	data := []byte{}

	// Chamando o método UploadFile
	err := sftpClient.UploadFile(data, "remote/path")
	assert.NoError(t, err)

	// Verificando se o método Create foi chamado corretamente
	mockClient.AssertExpectations(t)

	// Verificando se o método Close foi chamado corretamente
	mockCloser.AssertExpectations(t)

	// Verificando se nada foi escrito no arquivo
	assert.Equal(t, "", mockWriter.String())
}
