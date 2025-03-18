package shared

import (
	"io"
	"net/http"
	"net/url"
)

// PostFormRequest は、URLエンコードされたフォームデータをPOSTする
func PostFormRequest(endpoint string, data url.Values) ([]byte, error) {
	resp, err := http.PostForm(endpoint, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// GetRequest は、Bearerトークン付きでGETリクエストを送る
func GetRequest(endpoint, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
