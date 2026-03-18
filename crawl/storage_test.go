package crawl

import (
	"testing"
)

// TestStorage 已废弃 - S3 存储功能已移除
// 现在图片 URL 直接保存豆瓣原始地址，不再进行下载和上传
func TestStorage(t *testing.T) {
	t.Skip("S3 storage functionality has been removed")
}
