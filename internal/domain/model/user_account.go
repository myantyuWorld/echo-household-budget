//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package domainmodel

// UserAccount はユーザーアカウントを表すドメインモデルです
type UserAccount struct {
	ID             UserID       `json:"id"`
	UserID         LINEUserID   `json:"user_id"`
	Name           string       `json:"name"`
	PictureURL     string       `json:"pictureURL"`
	HouseholdBooks []*HouseHold `json:"household_books"`
}

type LINEUserInfo struct {
	UserID      LINEUserID `json:"userId"`
	DisplayName string     `json:"displayName"`
	PictureURL  string     `json:"pictureUrl"`
}

// UserAccountRepository はユーザーアカウントのリポジトリインターフェースです
type UserAccountRepository interface {
	Create(userAccount *UserAccount) error
	Delete(id UserID) error
	FindByLINEUserID(userID LINEUserID) (*UserAccount, error)
	FetchMe(userID UserID) (*UserAccount, error)
}

type LINEUserID string
type UserID uint

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
