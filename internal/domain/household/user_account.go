//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package household

// UserAccount はユーザーアカウントを表すドメインモデルです
type UserAccount struct {
	ID     uint       `json:"id"`
	UserID LINEUserID `json:"user_id"`
	Name   string     `json:"name"`
}

type LINEUserInfo struct {
	UserID      LINEUserID `json:"userId"`
	DisplayName string     `json:"displayName"`
	PictureURL  string     `json:"pictureUrl"`
}

// UserAccountRepository はユーザーアカウントのリポジトリインターフェースです
type UserAccountRepository interface {
	Create(userAccount *UserAccount) error
	Delete(id uint) error
	FindByLINEUserID(userID LINEUserID) (*UserAccount, error)
}

type LINEUserID string

func NewUserAccount(lineUserID LINEUserID, name string) *UserAccount {
	return &UserAccount{
		UserID: lineUserID,
		Name:   name,
	}
}

func NewLINEUserInfo(lineUserID LINEUserID, displayName string, pictureURL string) *LINEUserInfo {
	return &LINEUserInfo{
		UserID:      lineUserID,
		DisplayName: displayName,
		PictureURL:  pictureURL,
	}
}
