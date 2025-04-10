package household

// UserAccount はユーザーアカウントを表すドメインモデルです
type UserAccount struct {
	ID     uint   `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// UserAccountRepository はユーザーアカウントのリポジトリインターフェースです
type UserAccountRepository interface {
	Create(userAccount *UserAccount) error
	Delete(id uint) error
}
