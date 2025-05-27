package repository

type FileStorageRepository interface {
	UploadFile(fileData []byte, fileName string) (string, error)
	GetFileURL(fileName string) (string, error)
	DeleteFile(fileName string) error
}
