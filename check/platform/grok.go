package platform

import (
	"net/http"
	"strings"
)

// CheckGrok 检测 Grok 解锁情况
func CheckGrok(client *http.Client) (bool, error) {
	// 1. 创建一个不跟随重定向的自定义 Client，用来捕捉 Location 头部
	testClient := &http.Client{
		Transport: client.Transport,
		Timeout:   client.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 强制停止重定向
		},
	}

	// 2. 访问 Grok 官网
	req, err := http.NewRequest("GET", "https://grok.com/", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	resp, err := testClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 3. 逻辑判定
	// 如果返回 200 OK，说明直接进入了首页，解锁成功
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	// 如果返回 301/302 重定向
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
		location := resp.Header.Get("Location")
		// 如果跳转地址包含 unavailable 或 restricted，说明该地区被封锁
		if strings.Contains(location, "unavailable") || strings.Contains(location, "restricted") {
			return false, nil
		}
		// 其他跳转（如跳转到 /login）视为解锁成功
		return true, nil
	}

	// 403 Forbidden 明确代表封锁
	return false, nil
}
