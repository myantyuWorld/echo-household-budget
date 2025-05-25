package domainmodel

type ReceiptAnalyze struct {
	ID              int
	ImageURL        string
	AnalyzeStatus   string
	TotalPrice      int
	HouseholdBookID HouseHoldID
	Items           []ReceiptAnalyzeItem
}

type ReceiptAnalyzeReception struct {
	ImageURL        string
	ImageData       string
	HouseholdBookID HouseHoldID
}

type ReceiptAnalyzeItem struct {
	Name  string
	Price int
}

type ReceiptAnalyzeRepository interface {
	CreateReceiptAnalyzeReception(receiptAnalyze *ReceiptAnalyzeReception) error
	CreateReceiptAnalyzeResult(receiptAnalyze *ReceiptAnalyze) error
	FindByID(id HouseHoldID) (*ReceiptAnalyze, error)
}
