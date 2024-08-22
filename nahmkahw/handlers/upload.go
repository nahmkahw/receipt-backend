package handlers

import (
	"receipt-backend/nahmkahw/services"
    "receipt-backend/nahmkahw/errs"
	"net/http"
	"io/ioutil"
	"strings"
    "github.com/sirupsen/logrus"
    "runtime"
    "fmt"

	"github.com/labstack/echo"
)

type (
	UploadtHandlers struct {
		uploadServices services.UploadServiceInterface
        logger *logrus.Logger
	}
)

func NewUploadtHandlers(uploadServices services.UploadServiceInterface, logger *logrus.Logger) UploadtHandlers {
	return UploadtHandlers{uploadServices: uploadServices, logger : logger}
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
        param := fmt.Sprintf("%s,%s,%s",ordercode,filename,fileType)
		h.logError(errs.NewMessageAndStatusCode(http.StatusBadRequest,"Unsupported file type. Only png, jpg, and jpeg are allowed."),param)
        return c.JSON(http.StatusBadRequest, errs.NewMessageAndStatusCode(http.StatusBadRequest,"Unsupported file type. Only png, jpg, and jpeg are allowed."))
    }

    err = h.uploadServices.UploadFileImage(ordercode,filename,fileType,fileBytes)
    if err != nil {
        param := fmt.Sprintf("%s,%s,%s",ordercode,filename,fileType)
		h.logError(err,param)
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
        err := errs.NewMessageAndStatusCode(http.StatusBadRequest,"Image not found.")
        param := fmt.Sprintf("%s,%s",ordercode,filename)
		h.logError(err,param)
        return c.JSON(http.StatusBadRequest, err)
    }

    return c.File(filePath)

}

func (h *UploadtHandlers) logError(err error,param string) {
    pc, file, line, ok := runtime.Caller(1)
    if !ok {
        h.logger.Error("Failed to retrieve caller information")
    }
    funcName := runtime.FuncForPC(pc).Name()

    h.logger.WithFields(logrus.Fields{
        "func_name": funcName,
        "file":      file,
        "line":      line,
        "error":     err.Error(),
        "upload":  param,
    }).Error("Upload Error")
}

