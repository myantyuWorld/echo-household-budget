package domainmodel

type ReceiptAnalyze struct {
	ID         uint
	TotalPrice uint
	CategoryID CategoryID
	S3FilePath string
	Items      []ReceiptAnalyzeItem
}

type ReceiptAnalyzeReception struct {
	ImageURL        string
	ImageData       string
	HouseholdBookID HouseHoldID
	CategoryID      CategoryID
}

type ReceiptAnalyzeItem struct {
	Name  string
	Price uint
}

type ReceiptAnalyzeRepository interface {
	CreateReceiptAnalyzeReception(receiptAnalyze *ReceiptAnalyzeReception) error
	CreateReceiptAnalyzeResult(receiptAnalyze *ReceiptAnalyze) error
	FindReceiptAnalyzeByS3FilePath(s3FilePath string) (*ReceiptAnalyze, error)
	FindByID(id HouseHoldID) (*ReceiptAnalyze, error)
}
