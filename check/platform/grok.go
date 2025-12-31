package platform

import (
	"net/http"
)

func CheckGrok(client *http.Client) (bool, error) {
	// 访问 grok.com，不要拦截重定向，让 client 自动跟随
	req, _ := http.NewRequest("GET", "https://grok.com/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 只要状态码是 200，说明 IP 没被封锁，能够看到页面
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}
