package utils

import (
	"io/ioutil"
	"mime/multipart"
)

// CreateFormFileFromContent creates file field named 'fieldName' with the value of 'content'
func CreateFormFileFromContent(writer *multipart.Writer, fieldName string, content []byte, fileName string) error {
	field, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return err
	}
	_, err = field.Write(content)
	return err
}

// CreateFormFileFromFileHeader creates file field named 'fieldName' with the contents of 'file'.
func CreateFormFileFromFileHeader(writer *multipart.Writer, fieldName string, file *multipart.FileHeader) error {
	contentFile, err := file.Open()
	if err != nil {
		return err
	}
	defer contentFile.Close()
	content, err := ioutil.ReadAll(contentFile)
	if err != nil {
		return err
	}
	return CreateFormFileFromContent(writer, fieldName, content, file.Filename)
}

// CreateFormFileFromFilePath creates file field named 'fieldName' with the content of a file on 'filePath'
func CreateFormFileFromFilePath(writer *multipart.Writer, fieldName, fileName, filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return CreateFormFileFromContent(writer, fieldName, content, fileName)
}

// CreateFormField creates a field named 'fieldName' with the content 'fieldContent'
func CreateFormField(writer *multipart.Writer, fieldName, fieldContent string) error {
	field, err := writer.CreateFormField(fieldName)
	if err != nil {
		return err
	}
	_, err = field.Write([]byte(fieldContent))
	if err != nil {
		return err
	}
	return nil
}
