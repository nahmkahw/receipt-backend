package services

type (
	 uploadService struct {
    	fileStoragePath string
	}

	ImageUpload struct {
		OrderID string
		File    []byte
	}

	UploadServiceInterface interface {
		UploadFileImage(ordercode,filename,fileType string , fileBytes []byte) error
		DownFileImage(ordercode,filename string) string
	}
)

func NewUploadServices(fileStoragePath string) UploadServiceInterface {
	return &uploadService{
		fileStoragePath: fileStoragePath,
	}
}