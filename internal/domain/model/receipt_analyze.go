package domainmodel

type ReceiptAnalyze struct {
	ID              uint                 `json:"id"`
	TotalPrice      uint                 `json:"totalAmount"`
	CategoryID      CategoryID           `json:"categoryID"`
	S3FilePath      string               `json:"receiptImageURL"`
	HouseholdBookID HouseHoldID          `json:"householdID"`
	Items           []ReceiptAnalyzeItem `json:"items"`
}

type ReceiptAnalyzeReception struct {
	ImageURL        string
	ImageData       string
	HouseholdBookID HouseHoldID
	CategoryID      CategoryID
}

type ReceiptAnalyzeItem struct {
	Name  string `json:"name"`
	Price uint   `json:"amount"`
}

type ReceiptAnalyzeRepository interface {
	CreateReceiptAnalyzeReception(receiptAnalyze *ReceiptAnalyzeReception) error
	CreateReceiptAnalyzeResult(receiptAnalyze *ReceiptAnalyze) error
	FindReceiptAnalyzeByS3FilePath(s3FilePath string) (*ReceiptAnalyze, error)
	FindByID(id HouseHoldID) (*ReceiptAnalyze, error)
}
