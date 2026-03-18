# 代码审查报告

## 项目优化建议

---

## 一、严重问题 (Critical)

### 1. 数据库查询缺少错误处理

**问题位置**: `dao/*.go` 多个文件

**问题描述**: 所有 GORM 查询都没有检查错误，可能导致静默失败。

**示例代码**:
```go
// dao/storage.go
func GetStorage(source string) *model.Storage {
    storage := &model.Storage{}
    common.Db.Where("source = ? ", source).Find(storage)  // 没有检查 error
    if storage.ID == 0 {
        return nil
    }
    return storage
}
```

**建议修复**:
```go
func GetStorage(source string) *model.Storage {
    storage := &model.Storage{}
    result := common.Db.Where("source = ?", source).Find(storage)
    if result.Error != nil {
        logrus.Errorln("get storage error:", result.Error)
        return nil
    }
    if storage.ID == 0 {
        return nil
    }
    return storage
}
```

**影响范围**: 所有 DAO 层函数

---

### 2. 并发 goroutine 泄漏风险

**问题位置**: `agent/item.go:103`, `agent/user.go:96-104`

**问题描述**: `item.go` 中 worker 启动后立即发送 `done <- true`，可能导致 channel 阻塞或 goroutine 泄漏。

**问题代码**:
```go
// agent/item.go
for i := 0; i < concurrency; i++ {
    j := i + 1
    go func() {
        for {
            itemWorker(j, ch, done)
        }
    }()
    done <- true  // 这里会立即发送，但 worker 可能还没开始等待
}
```

**建议修复**: 移除初始的 `done <- true`，让 selector 和 worker 通过 channel 自然协调。

---

### 3. 图片下载 panic 导致服务崩溃

**问题位置**: `crawl/storage.go:56`

**问题描述**: 图片下载失败 5 次后直接 panic，会导致整个服务崩溃。

**问题代码**:
```go
if file == nil {
    panic("download file finally failed for : " + url)  // 直接 panic
}
```

**建议修复**: 返回错误而不是 panic，让调用层决定如何处理。

---

## 二、重要问题 (High)

### 4. SQL 注入风险

**问题位置**: `dao/*.go` 多个文件

**问题描述**: 虽然使用了参数化查询，但部分 WHERE 子句后有多余空格，且部分查询未验证输入。

**问题代码**:
```go
// dao/schedule.go
common.Db.Where("douban_id = ? AND type = ? ", doubanId, t)  // 末尾多余空格
```

**建议**: 统一 SQL 书写格式，移除多余空格。

---

### 5. 配置文件中硬编码敏感信息

**问题位置**: `application.yml.sample:32-38`

**问题描述**: S3 密钥、数据库密码等敏感信息出现在示例配置文件中。

**建议修复**:
```yaml
# application.yml.sample
s3:
  access_key: ${S3_ACCESS_KEY}  # 使用环境变量占位符
  secret_key: ${S3_SECRET_KEY}
```

---

### 6. Cookie 定期刷新可能失效

**问题位置**: `crawl/http.go:98-103`

**问题描述**: 每小时刷新一次 Cookie Jar，但如果 dbcl2 token 过期，整个客户端会失效。

**建议**: 添加 Cookie 过期检测机制，检测到 403 时尝试重新登录或切换账号。

---

## 三、一般问题 (Medium)

### 7. 日志输出过于冗长

**问题位置**: 全局

**问题描述**:
- 每次 HTTP 请求都打印日志
- storage 命中/未命中都打印
- 大量 Infoln 导致日志文件膨胀

**建议**:
- 区分日志级别 (Debug/Info/Warn/Error)
- 生产环境降低日志级别
- 关键路径使用结构化日志

```go
// 改进示例
logrus.WithFields(logrus.Fields{
    "url": url,
    "status": resp.StatusCode,
    "duration_ms": duration,
}).Debugln("HTTP request completed")
```

---

### 8. 数据库连接池未配置

**问题位置**: `common/database.go`

**问题描述**: 没有配置 GORM 底层的 SQL DB 连接池参数。

**建议添加**:
```go
sqlDB, err := db.DB()
if err != nil {
    panic(err)
}
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

---

### 9. Prometheus 指标高基数风险

**问题位置**: `main.go:21-27`

**问题描述**: `serviceOpsHistogram` 使用 `ua` 和 `referer` 作为 label，可能导致高基数问题。

**问题代码**:
```go
serviceOpsHistogram = promauto.NewHistogramVec(..., []string{"method", "path", "ua", "referer"})
```

**建议**: 移除 `ua` 和 `referer`，或使用其哈希值/分类值。

---

### 10. 限流器全局单例

**问题位置**: `crawl/http.go:26-28`

**问题描述**: UserLimiter、ItemLimiter 是全局单例，无法根据不同类型的请求动态调整。

**建议**: 考虑使用令牌桶分层限流或基于用户的限流策略。

---

## 四、改进建议 (Low)

### 11. 代码复用问题

**问题位置**: `agent/processor.go` 180-420 行

**问题描述**: `syncCommentBook`、`syncCommentMovie`、`syncCommentGame`、`syncCommentSong` 函数逻辑高度相似，存在代码重复。

**建议**: 提取公共逻辑为通用函数，使用泛型或接口处理不同类型。

---

### 12. 缺少输入验证

**问题位置**: `controller/admin.go`

**问题描述**: 接口参数验证不完善。

**问题代码**:
```go
func RefreshUser(ctx *gin.Context) {
    idStr := ctx.Query("id")
    id := util.ParseNumber(idStr)
    if id == 0 {  // 0 也可能是合法 ID
        BizError(ctx, "参数错误")
        return
    }
```

**建议**: 使用更严格的验证逻辑，区分"解析失败"和"ID 为 0"。

---

### 13. 临时文件清理不彻底

**问题位置**: `crawl/storage.go:72`

**问题描述**: 临时文件在 panic 时可能不会被清理。

**建议**: 使用 defer 清理临时文件。

```go
defer func() {
    if file != nil {
        os.Remove(file.Name())
    }
}()
```

---

### 14. 缺少健康检查接口

**问题位置**: `main.go`

**问题描述**: 没有 `/health` 或 `/ready` 接口用于 Kubernetes 或负载均衡器。

**建议添加**:
```go
router.GET("/health", func(ctx *gin.Context) {
    // 检查数据库连接、S3 连接等
    ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
})
```

---

### 15. 依赖版本过时

**问题位置**: `go.mod`

**问题描述**:
- Go 版本 1.21 较旧
- gin v1.9.1 不是最新版
- 多个依赖可升级

**建议**: 定期执行 `go get -u ./...` 更新依赖。

---

## 五、架构建议

### 16. 爬虫调度策略优化

**当前问题**: 简单的 FIFO 队列，无法处理优先级。

**建议**:
- 添加优先级队列（新用户 > 老用户更新）
- 添加退避机制（失败的任务延迟更久再试）
- 添加任务去重逻辑

---

### 17. 缺少监控告警

**当前问题**: 只有 Prometheus 指标，没有告警配置。

**建议**:
- 配置 AlertManager 告警规则
- 监控爬虫失败率、队列积压量
- 监控 HTTP 403 比例（账号可能被封）

---

### 18. 图片存储策略单一

**当前问题**: 所有图片都上传 S3，增加成本。

**建议**:
- 添加 CDN 缓存层
- 对热门图片使用缓存
- 考虑使用豆瓣 CDN 直链（如果允许）

---

## 六、修复状态

### 已修复 (P0/P1)

| 优先级 | 问题编号 | 问题简述 | 状态 |
|--------|----------|----------|------|
| P0 | 1 | 数据库错误处理缺失 | ✅ 已修复 |
| P0 | 2 | Goroutine 泄漏风险 | ✅ 已修复 |
| P0 | 3 | Panic 导致服务崩溃 | ✅ 已修复 |
| P1 | 4 | SQL 格式问题 | ✅ 已修复 |
| P1 | 5 | 敏感信息配置 | ✅ 已修复 |
| P1 | 8 | 数据库连接池 | ✅ 已修复 |
| P1 | 13 | 临时文件清理 | ✅ 已修复 |
| P1 | 18 | 图片存储策略单一 | ✅ 已移除 S3 存储 |
| - | - | Go 版本升级 | ✅ 已升级至 1.25.0 |
| P2 | 7 | 日志级别优化 | ✅ 已修复 |
| P2 | 9 | Prometheus 指标优化 | ✅ 已修复 |
| P3 | 11 | 代码复用优化 | ✅ 已修复 |

---

## 修复摘要

### P0 问题修复

**1. 数据库错误处理** - 所有 DAO 层函数添加了错误检查和日志记录：
- `dao/storage.go`, `dao/schedule.go`, `dao/user.go`
- `dao/book.go`, `dao/movie.go`, `dao/game.go`, `dao/song.go`
- `dao/comment.go`, `dao/rating.go`, `dao/access.go`

**2. Goroutine 泄漏修复** - 重构了 `agent/item.go` 和 `agent/user.go`：
- 移除了 `done` channel 和不当的 `done <- true` 调用
- 使用 `for schedule := range ch` 替代单次接收
- selector 使用 `continue` 替代 `return`

**3. Panic 滥用修复** - 修改了 `crawl/storage.go`：
- 下载失败返回原 URL 而不是 panic
- 使用 `defer` 确保临时文件被清理
- `download` 函数返回 `nil` 而不是 panic

**4. S3 存储移除** - 完全移除了 S3 存储和图片下载逻辑：
- `crawl/storage.go` 简化为返回原始 URL 的透传函数
- 移除了 `application.yml.sample` 中的 S3 配置
- 移除了 `agent/processor.go` 中所有 `crawl.Storage()` 调用
- 简化了图片处理流程，直接保存豆瓣原始 URL

### P1 问题修复

**5. SQL 格式** - 移除了所有 WHERE 子句后的多余空格

**6. 敏感信息配置** - 更新了 `application.yml.sample`：
- 数据库密码使用 `${DB_PASSWORD:-changeme}` 占位符
- S3 密钥使用环境变量引用（已移除）
- 添加了配置说明注释

**7. 数据库连接池** - 在 `common/database.go` 中添加：
```go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

**8. 临时文件清理** - 使用 `defer` 确保文件被清理

### P2 问题修复

**9. 日志级别优化** - 实现了可配置的日志级别和结构化日志：
- 在 `common/log.go` 中添加日志级别配置支持
- 在 `application.yml.sample` 中添加 `log.level` 配置项
- 将 HTTP 请求日志从 `Info` 改为 `Debug` 级别 (`main.go`, `crawl/http.go`)
- 将爬虫发现日志从 `Info` 改为 `Debug` 级别 (`agent/item.go`, `agent/user.go`)
- 使用 `logrus.WithFields` 进行结构化日志输出

**10. Prometheus 指标优化** - 移除了高基数标签：
- 从 `serviceOpsHistogram` 中移除了 `ua` 和 `referer` 标签
- 仅保留 `method` 和 `path` 作为标签

### P3 问题修复

**11. 代码复用优化** - 重构了 `agent/processor.go` 中的重复代码：
- 创建了泛型函数 `syncComment[T]` 处理通用评论同步逻辑
- 将 `syncCommentBook/Movie/Game/Song` 四个函数简化为调用泛型函数
- 代码行数从 ~150 行减少到 ~50 行

### Go 版本升级

- Go 版本：1.21 → 1.25.0 (最新稳定版)
- 所有主要依赖已升级到最新版本

---

## 总结

该项目整体架构清晰，代码质量尚可。

**已解决的关键风险**:
1. ~~稳定性风险~~: 数据库查询无错误处理、panic 滥用 ✅
2. ~~并发风险~~: goroutine 泄漏、channel 使用不当 ✅
3. ~~安全风险~~: 敏感配置 ✅
4. ~~存储复杂度~~: S3 存储逻辑冗余、图片下载开销 ✅
5. ~~可维护性~~: 日志级别混乱、代码重复 ✅

**剩余建议**:
- P2/P3 级别问题可根据实际情况逐步优化
- 建议添加定期依赖更新机制
