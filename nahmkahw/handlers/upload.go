package handlers

import (
	"receipt-backend/nahmkahw/services"
    "receipt-backend/nahmkahw/errs"
	"net/http"
	"io/ioutil"
    "image/jpeg"
	"strings"
    "github.com/sirupsen/logrus"
    "runtime"
    "time"
    "fmt"
    "bytes"
    "image"

	"github.com/labstack/echo"
)

type (
	UploadtHandlers struct {
		uploadServices services.UploadServiceInterface
        logger *logrus.Logger
	}

    StudentForm struct {
		Std_code    string `json:"std_code" validate:"required"`
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

func (h *UploadtHandlers) GetPhoto(c echo.Context) error {
    
    student := new(StudentForm)

	if err := c.Bind(student); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, errs.NewBadRequestError())
	}

	if err := c.Validate(student); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, errs.NewBadRequestError())
	}

	 url := "http://10.2.1.155:9100/student/photo/" + student.Std_code

	client := &http.Client{
		Timeout: 60 * time.Second, // Set a higher timeout value
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	response, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	defer response.Body.Close()

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
        return c.JSON(http.StatusBadRequest, errs.NewMessageAndStatusCode(http.StatusBadRequest,"Unsupported image type."))
	}

	fmt.Println(contentType)

	// Decode the image
	var img image.Image
	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(response.Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
	case "image/png":
		img, err = jpeg.Decode(response.Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
	default:
        return c.JSON(http.StatusBadRequest, errs.NewMessageAndStatusCode(http.StatusBadRequest,"Unsupported image format."))
	}

	outputImg := bytes.NewBuffer(nil)

	if err := jpeg.Encode(outputImg, img, nil); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.Blob(http.StatusOK, contentType, outputImg.Bytes())
}

