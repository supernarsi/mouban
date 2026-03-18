package crawl

// Storage 直接返回原始 URL，不再进行 S3 存储
// 该函数保留是为了兼容现有代码调用
func Storage(url string) string {
	return url
}
