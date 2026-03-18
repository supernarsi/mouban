# mouban

## 功能特性

- 用户数据抓取：自动获取豆瓣用户的书影音游标注信息
- 增量同步：基于 RSS 和评论时间的增量更新机制
- 条目发现：通过关联推荐自动发现新的条目和用户
- 定时任务：每日自动更新首页条目
- 数据持久化：MySQL 存储，支持 9 张数据表
- 监控指标：内置 Prometheus metrics 端点
- Docker 部署：支持容器化快速部署

## 架构设计

```
main.go (入口 + Gin 路由)
    │
    ├── common/        # 初始化：配置 (Viper)、数据库 (GORM)、日志 (logrus)
    ├── controller/    # HTTP 处理器：/admin/* 和 /guest/* 端点
    ├── agent/         # 后台 Worker：定时爬虫任务
    ├── crawl/         # 网页抓取：HTTP 客户端、HTML 解析、限流
    ├── dao/           # 数据访问层：GORM 操作
    ├── model/         # 数据模型：Book/Movie/Game/Song/User/Comment 等
    ├── consts/        # 常量定义：类型码、URL 模式
    └── util/          # 工具函数：JSON 解析、堆栈追踪
```

## 快速开始

### 构建

```bash
# 本地运行
go run main.go

# 生产构建（Linux amd64 二进制 + Docker）
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o main main.go
docker build -t mythsman/mouban -f Dockerfile --platform=linux/amd64 .
```

### 配置

配置文件 `application.yml`：

```yaml
server:
  cors: http://localhost
  port: 8080
  limit: 1h
agent:
  enable: true
  flow:
    discover: false
  user:
    concurrency: 3
  item:
    concurrency: 3
    max: 3000
datasource:
  driver: mysql
  host: localhost
  port: 3306
  database: mouban
  username: root
  password: your_password
http:
  timeout: 10000
  retry_max: 20
  auth: dbcl2_cookie_value,http://user:pass@proxy:port;
```

环境变量覆盖（`KEY__SUBKEY` 格式）：

```bash
export GIN_MODE=release
export server__cors=https://yourdomain.com
export datasource__host=localhost
export datasource__username=root
export datasource__password=secret
export agent__enable=true
```

### Docker Compose 部署

```yaml
services:
  mouban:
    image: mythsman/mouban
    container_name: mouban
    restart: always
    expose:
      - "8080"
    environment:
      - GIN_MODE=release
      - agent__enable=true
      - agent__flow__discover=false
      - agent__item__concurrency=5
      - agent__item__max=10000
      - http__timeout=30000
      - http__retry_max=20
      - http__interval__user=5000
      - http__interval__item=2000
      - server__cors=https://yourdomain.com
      - server__limit=30m
      - datasource__host=mysql-host
      - datasource__username=mysql-user
      - datasource__password=mysql-passwd
```

**重要**: `http__auth` 参数用于配置豆瓣登录态和 HTTP 代理，格式为：
```
<dbcl2_cookie>,http://<user>:<password>@<proxy_ip>:<proxy_port>;
```

`dbcl2` 需从豆瓣 Cookie 中获取。使用登录态可避免豆瓣对未登录用户的投毒策略。

## API 接口

### 访客接口

#### 用户录入/更新

```
GET /guest/check_user?id={douban_id}
```

响应示例：

```json
{
  "success": true,
  "result": {
    "id": 1000001,
    "domain": "ahbei",
    "name": "阿北",
    "thumbnail": "https://img1.doubanio.com/icon/u1000001-30.jpg",
    "book_wish": 81,
    "book_do": 61,
    "book_collect": 115,
    "movie_wish": 77,
    "movie_do": 17,
    "movie_collect": 218,
    "song_wish": 23,
    "song_do": 21,
    "song_collect": 24,
    "sync_at": 1667232000,
    "check_at": 1679646797,
    "publish_at": 1570409179
  }
}
```

时间戳说明：
- `publish_at`: 用户最近一次更新的时间
- `check_at`: 最近一次检测用户更新的时间
- `sync_at`: 最近一次同步用户信息的时间

#### 查询用户评论

| 类型 | 接口 |
|------|------|
| 读书 | `/guest/user_book?id={id}&action={wish\|do\|collect}` |
| 电影 | `/guest/user_movie?id={id}&action={wish\|do\|collect}` |
| 游戏 | `/guest/user_game?id={id}&action={wish\|do\|collect}` |
| 音乐 | `/guest/user_song?id={id}&action={wish\|do\|collect}` |

### 管理接口

#### 加载 Sitemap 数据

```
GET /admin/load_data?path={sitemap_file_path}
```

从豆瓣 [sitemap_index](https://www.douban.com/sitemap_index.xml) 离线加载存量数据。

#### 强制更新条目

```
GET /admin/refresh_item?type={1-book\|2-movie\|3-game\|4-song}&id={item_id}
```

#### 强制更新用户

```
GET /admin/refresh_user?id={douban_uid}
```

谨慎使用，会对系统造成较大压力。

### 监控端点

```
GET /metrics
```

Prometheus 格式的监控指标。

## 数据流

1. 用户调用 `/guest/check_user` 触发用户资料抓取
2. 抓取用户读书首页获取头像、域名等信息
3. 访问 RSS 页面获取最新更新时间用于去重
4. 滚动抓取用户书影音游评论页
5. 抓取条目详情页，发现关联的用户和条目
6. 新条目加入调度队列等待抓取
7. 每日定时任务更新首页条目

## 测试

```bash
# 运行所有测试
go test ./...

# 运行单个包测试
go test ./dao/
go test ./util/

# 运行特定测试
go test -run TestFunctionName ./path/to/package/

# 代码检查
golangci-lint run
```

## 数据库

自动创建 9 张表：

| 表名 | 说明 |
|------|------|
| users | 用户资料 |
| books | 图书条目 |
| movies | 电影条目 |
| games | 游戏条目 |
| songs | 音乐条目 |
| comments | 用户评论 |
| ratings | 聚合评分 |
| schedules | 爬虫任务队列 |
| access | 限流令牌 |

## 许可证

MIT License

## 相关项目

- [hexo-douban](https://github.com/mythsman/hexo-douban) - Hexo 豆瓣数据插件
