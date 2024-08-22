package services

import (
	"io/ioutil"
    "os"
    "path/filepath"
)

func (s *uploadService) UploadFileImage(ordercode,filename,fileType string , fileBytes []byte) error {
    filePath := filepath.Join(s.fileStoragePath, ordercode + "_"+ filename + "." + fileType)
    return ioutil.WriteFile(filePath, fileBytes, 0644)
}

func (s *uploadService) DownFileImage(ordercode,filename string) string {

    // List of supported file types
    supportedFileTypes := []string{"png", "jpg", "jpeg"}
    var filePath string

    // Find the file with the supported extension
    for _, fileType := range supportedFileTypes {
        potentialPath := filepath.Join(s.fileStoragePath, ordercode+"_"+ filename +"." + fileType)
        if _, err := os.Stat(potentialPath); err == nil {
            filePath = potentialPath
            break
        }
    }

    return filePath
}