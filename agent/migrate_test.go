package agent

import (
	"testing"
)

// TestDb 已废弃 - S3 迁移功能已移除
// 不再需要将图片从豆瓣迁移到 S3
func TestDb(t *testing.T) {
	t.Skip("S3 migration functionality has been removed")
}
