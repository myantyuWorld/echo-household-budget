package domainmodel

type Information struct {
	ID          int
	Title       string
	Content     string
	Category    string
	IsPublished bool
}

func NewInformation(title string, content string, category string) *Information {
	return &Information{
		Title:       title,
		Content:     content,
		Category:    category,
		IsPublished: false,
	}
}
