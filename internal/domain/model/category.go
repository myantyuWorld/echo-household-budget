package domainmodel

// Category はカテゴリのドメインエンティティ
type Category struct {
	ID              uint
	HouseholdBookID uint
	Name            string
	Color           string
}

// CategoryRepository はカテゴリの永続化を担うリポジトリのインターフェース
type CategoryRepository interface {
	// Create は新しいカテゴリを作成します
	Create(category *Category) error

	// FindByID は指定されたIDのカテゴリを取得します
	FindByID(id uint) (*Category, error)

	// Update は既存のカテゴリを更新します
	Update(category *Category) error

	// Delete は指定されたIDのカテゴリを削除します
	Delete(id uint) error
}
