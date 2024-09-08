package utility

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
)

func UploadFile(file *multipart.FileHeader, destination string) (filePath string, err error) {
	src, err := file.Open()
	if err != nil {
		return
	}

	defer src.Close()

	err = os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		return
	}

	filePath = fmt.Sprintf("%s/%s_%s", destination, uuid.NewString(), file.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return
	}

	return filePath, err
}
