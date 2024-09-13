package utility

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
)

func UploadFile(file *multipart.FileHeader, destination string) (fileName string, err error) {
	src, err := file.Open()
	if err != nil {
		return
	}

	defer src.Close()

	err = os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		return
	}

	fileName = fmt.Sprintf("%s_%s", uuid.NewString(), file.Filename)
	filePath := fmt.Sprintf("%s/%s", destination, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return
	}

	return fileName, err
}

func DeleteFile(filePath string) (err error) {
	err = os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return
	}

	return
}
