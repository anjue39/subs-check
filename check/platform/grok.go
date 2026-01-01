package platform

import (
	"io"
	"net/http"
	"strings"
)

// CheckGrok 检查 Grok 是否可用
// 当前逻辑：状态码 200 且页面包含 Grok 正常特征（如 "chat"、"Grok" 或其他核心元素），视为可用
// 如果有明确限制消息，视为不可用
func CheckGrok(httpClient *http.Client) (bool, error) {
	req, err := http.NewRequest("GET", "https://grok.x.ai/", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	// 可选：添加 Accept-Language: en-US 以模拟英文环境
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil // 非200，很可能被 block 或 redirect 到错误页
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	lowerBody := strings.ToLower(string(body))

	// 明确限制短语（当前最常见的）
	if strings.Contains(lowerBody, "not available in your region") ||
		strings.Contains(lowerBody, "this service is not available") ||
		strings.Contains(lowerBody, "regional restriction") ||
		strings.Contains(lowerBody, "coming soon") {
		return false, nil
	}

	// 正向判断：页面包含 Grok 正常元素（聊天界面特征，根据当前页面调整）
	if strings.Contains(lowerBody, "grok") &&
		strings.Contains(lowerBody, "chat") &&
		(strings.Contains(lowerBody, "xai") || strings.Contains(lowerBody, "ask anything")) {
		return true, nil
	}

	// 默认保守：有响应但无明确特征，视为不可用（避免误判）
	return false, nil
}
