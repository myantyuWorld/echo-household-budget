//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package domainmodel

// UserAccount はユーザーアカウントを表すドメインモデルです
type UserAccount struct {
	ID         uint       `json:"id"`
	UserID     LINEUserID `json:"user_id"`
	Name       string     `json:"name"`
	PictureURL string     `json:"pictureURL"`
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

func NewUserAccount(lineUserInfo *LINEUserInfo) *UserAccount {
	return &UserAccount{
		UserID:     lineUserInfo.UserID,
		Name:       lineUserInfo.DisplayName,
		PictureURL: lineUserInfo.PictureURL,
	}
}

func NewLINEUserInfo(lineUserID LINEUserID, displayName string, pictureURL string) *LINEUserInfo {
	return &LINEUserInfo{
		UserID:      lineUserID,
		DisplayName: displayName,
		PictureURL:  pictureURL,
	}
}
