package errors

import "fmt"

// ErrorCode はエラーコードを表す型
type ErrorCode string

const (
	// エラーコードの定義
	ErrorCodeInvalidInput    ErrorCode = "INVALID_INPUT"
	ErrorCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrorCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrorCodeInternalError   ErrorCode = "INTERNAL_ERROR"
	ErrorCodeDatabaseError   ErrorCode = "DATABASE_ERROR"
	ErrorCodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"
)

// AppError はアプリケーション固有のエラーを表す構造体
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error はerrorインターフェースの実装
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError は新しいAppErrorを作成する
func NewAppError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// IsAppError はエラーがAppErrorかどうかを判定する
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError はエラーからAppErrorを取得する
func GetAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
