package handlers

import (
	"receipt-backend/nahmkahw/services"
    "receipt-backend/nahmkahw/errs"
	"net/http"
	"io/ioutil"
	"strings"

	"github.com/labstack/echo"
)

type (
	UploadtHandlers struct {
		uploadServices services.UploadServiceInterface
	}
)

func NewUploadtHandlers(uploadServices services.UploadServiceInterface) UploadtHandlers {
	return UploadtHandlers{uploadServices: uploadServices}
}

func (h *UploadtHandlers) UploadFileImage(c echo.Context) error {
	ordercode := c.FormValue("ordercode")
    filename := c.FormValue("filename")
    file, err := c.FormFile("file")
    if err != nil {
        return c.JSON(http.StatusBadRequest, errs.NewBadRequestError())
    }

    src, err := file.Open()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, errs.NewInternalServerError())
    }
    defer src.Close()

    fileBytes, err := ioutil.ReadAll(src)
    if err != nil {
       return c.JSON(http.StatusInternalServerError, errs.NewInternalServerError())
    }

	// Extract file type from Content-Type
    contentType := file.Header.Get("Content-Type")
    fileType := ""
    if strings.Contains(contentType, "/") {
        parts := strings.Split(contentType, "/")
        if len(parts) > 1 {
            fileType = parts[1]
        }
    }

    // Validate file type
    if fileType != "png" && fileType != "jpg" && fileType != "jpeg" {
        return c.JSON(http.StatusBadRequest, errs.NewMessageAndStatusCode(http.StatusBadRequest,"Unsupported file type. Only png, jpg, and jpeg are allowed."))
    }

    err = h.uploadServices.UploadFileImage(ordercode,filename,fileType,fileBytes)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, errs.NewInternalServerError())
    }

    return c.JSON(http.StatusOK , errs.NewMessageAndStatusCode(http.StatusOK,"File uploaded successfully."))

}

func (h *UploadtHandlers) GetFileImage(c echo.Context) error {
	ordercode := c.FormValue("ordercode")
    filename := c.FormValue("filename")
    if ordercode == "" && filename == "" {
        return c.JSON(http.StatusBadRequest, errs.NewBadRequestError())
    }

    filePath := h.uploadServices.DownFileImage(ordercode,filename)

    if filePath == "" {
        return c.JSON(http.StatusBadRequest, errs.NewMessageAndStatusCode(http.StatusBadRequest,"Image not found."))
    }

    return c.File(filePath)

}

