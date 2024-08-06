package fileutil

import (
	"log"
	"mime/multipart"
	"net/http"
)

var ImageContentTypes = []string{
	"image/apng",
	"image/avif",
	"image/bmp",
	"image/gif",
	"image/vnd.microsoft.icon",
	"image/jpeg",
	"image/png",
	"image/svg+xml",
	"image/tiff",
	"image/webp",
}

func DetectMultipartFileType(file *multipart.FileHeader) (string, error) {
	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(fileContent multipart.File) {
		err := fileContent.Close()
		if err != nil {
			log.Println("Error closing file: ", err)
		}
	}(fileContent)

	buffer := make([]byte, 512)
	_, err = fileContent.Read(buffer)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}
