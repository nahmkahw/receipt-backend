package services

import (
    "bytes"
    "image"
    "image/png"
    "os"
    "path/filepath"
    _ "image/gif"  // Support for GIF files
    _ "image/jpeg" // Support for JPEG files
    _ "image/png"  // Support for PNG files
)

func (s *uploadService) UploadFileImage(ordercode,filename,fileType string , fileBytes []byte) error {

        // Decode the image from fileBytes
    img, _, err := image.Decode(bytes.NewReader(fileBytes))
    if err != nil {
        return err
    }

    // Create the path for the PNG file
    filePath := filepath.Join(s.fileStoragePath, ordercode + "_"+ filename + ".png")

    // Save the image as PNG
    outFile, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer outFile.Close()

    err = png.Encode(outFile, img)
    if err != nil {
        return err
    }

    return nil

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