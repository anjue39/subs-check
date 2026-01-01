func CheckGrok(client *http.Client) (bool, error) {
	// 建议直接访问最新的独立域名 grok.com
	url := "https://grok.com/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	
	// 设置一个更现代的 UA
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	// 直接使用传入的 client 即可，不需要重新创建，保持全局超时一致
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 判定逻辑：
	// 1. 200 OK 代表可以直接访问（解锁）
	// 2. 301/302 重定向到其他页面通常也是解锁（因为未解锁 IP 通常直接给 403）
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
		// 排除掉明确的失败重定向
		location := resp.Header.Get("Location")
		if strings.Contains(location, "unavailable") || strings.Contains(location, "restricted") {
			return false, nil
		}
		return true, nil
	}

	// 403 Forbidden 明确代表地区封锁
	return false, nil 
}
