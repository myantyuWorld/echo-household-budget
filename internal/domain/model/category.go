//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package domainmodel

// Category はカテゴリのドメインエンティティ
type Category struct {
	ID    CategoryID
	Name  string
	Color string
}

type HouseHoldCategory struct {
	ID              CategoryID
	Name            string
	Color           string
	HouseholdBookID HouseHoldID
}

type CategoryLimit struct {
	ID              CategoryLimitID `json:"categoryLimitID"`
	HouseholdBookID HouseHoldID     `json:"houseHoldID"`
	Category        Category        `json:"category"`
	LimitAmount     int             `json:"limitAmount"`
}

type CategoryID uint

const (
	CategoryIDFood CategoryID = iota + 1
	CategoryIDNecessary
)

type CategoryLimitID uint

// CategoryRepository はカテゴリの永続化を担うリポジトリのインターフェース
type CategoryRepository interface {
	// Create は新しいカテゴリを作成します
	CreateMasterCategory(category *Category) error
	CreateHouseHoldCategory(categoryLimit *CategoryLimit) error

	// FindByID は指定されたIDのカテゴリを取得します
	FindMasterCategoryByID(categoryID CategoryID) (*Category, error)
	FindHouseHoldCategoryByHouseHoldID(categoryLimitID CategoryLimitID) (*CategoryLimit, error)

	// Update は既存のカテゴリを更新します
	UpdateMasterCategory(category *Category) error
	UpdateHouseHoldCategory(categoryLimit *CategoryLimit) error

	// Delete は指定されたIDのカテゴリを削除します
	DeleteMasterCategory(id CategoryID) error
	DeleteHouseHoldCategory(categoryLimitID CategoryLimitID) error
}
