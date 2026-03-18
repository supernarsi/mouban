# 数据库表结构文档

豆瓣数据采集服务数据库设计文档

---

## 表概览

| 表名 | 说明 | 记录规模 |
|------|------|----------|
| `user` | 豆瓣用户信息 | 百万级 |
| `book` | 书籍条目信息 | ~447 万 |
| `movie` | 电影条目信息 | ~48 万 |
| `game` | 游戏条目信息 | ~3.8 万 |
| `song` | 音乐条目信息 | ~120 万 |
| `comment` | 用户评论/标注记录 | 千万级 |
| `rating` | 条目评分统计 | 千万级 |
| `schedule` | 爬虫任务调度队列 | 百万级 |
| `access` | 访问日志/限流记录 | 百万级 |
| `storage` | ~~图片/S3 存储映射~~ (已弃用) | 百万级 |

---

## 表详细结构

### 1. user - 用户信息表

存储豆瓣用户的基本信息和各类目数量统计。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `douban_uid` | BIGINT | NOT NULL, UNIQUE | 豆瓣用户 ID |
| `domain` | VARCHAR(64) | NOT NULL, INDEX | 豆瓣个人主页域名 (如 `ahbei`) |
| `name` | VARCHAR(512) | NOT NULL | 用户昵称 |
| `thumbnail` | VARCHAR(512) | | 用户头像 URL (已转存到 S3) |
| `book_wish` | INT UNSIGNED | DEFAULT 0 | 想读数量 |
| `book_do` | INT UNSIGNED | DEFAULT 0 | 在读数量 |
| `book_collect` | INT UNSIGNED | DEFAULT 0 | 读过数量 |
| `game_wish` | INT UNSIGNED | DEFAULT 0 | 想玩数量 |
| `game_do` | INT UNSIGNED | DEFAULT 0 | 在玩数量 |
| `game_collect` | INT UNSIGNED | DEFAULT 0 | 玩过数量 |
| `movie_wish` | INT UNSIGNED | DEFAULT 0 | 想看数量 |
| `movie_do` | INT UNSIGNED | DEFAULT 0 | 在看数量 |
| `movie_collect` | INT UNSIGNED | DEFAULT 0 | 看过数量 |
| `song_wish` | INT UNSIGNED | DEFAULT 0 | 想听数量 |
| `song_do` | INT UNSIGNED | DEFAULT 0 | 在听数量 |
| `song_collect` | INT UNSIGNED | DEFAULT 0 | 听过数量 |
| `sync_at` | DATETIME | | 最近同步时间 |
| `check_at` | DATETIME | | 最近检测时间 |
| `register_at` | DATETIME | | 注册时间 (首次抓取时间) |
| `publish_at` | DATETIME | | 用户最近发布时间 (用于判断是否变化) |
| `created_at` | DATETIME | | 记录创建时间 |
| `updated_at` | DATETIME | | 记录更新时间 |

**索引**:
- `PRIMARY KEY (id)`
- `UNIQUE KEY (douban_uid)`
- `KEY (domain)`

---

### 2. book - 书籍信息表

存储书籍条目的详细信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `douban_id` | BIGINT | NOT NULL, UNIQUE | 豆瓣书籍 ID |
| `title` | VARCHAR(1024) | NOT NULL | 书名 |
| `subtitle` | VARCHAR(1024) | | 副标题 |
| `orititle` | VARCHAR(1024) | | 原作名 (翻译书籍的原名) |
| `author` | VARCHAR(1024) | | 作者 (多人用分隔符) |
| `translator` | VARCHAR(512) | | 译者 |
| `press` | VARCHAR(512) | | 出版社 |
| `producer` | VARCHAR(512) | | 出品方 |
| `serial` | VARCHAR(512) | | 丛书名 |
| `publish_date` | VARCHAR(64) | | 出版年月 |
| `isbn` | VARCHAR(64) | | ISBN 号 |
| `framing` | VARCHAR(512) | | 装帧 (精装/平装等) |
| `page` | INT UNSIGNED | | 页数 |
| `price` | INT UNSIGNED | | 定价 (单位：分) |
| `book_intro` | MEDIUMTEXT | | 书籍简介 |
| `author_intro` | MEDIUMTEXT | | 作者简介 |
| `thumbnail` | VARCHAR(512) | | 封面图 URL (已转存到 S3) |
| `created_at` | DATETIME | | 记录创建时间 |
| `updated_at` | DATETIME | | 记录更新时间 |

**索引**:
- `PRIMARY KEY (id)`
- `UNIQUE KEY (douban_id)`

---

### 3. movie - 电影信息表

存储电影条目的详细信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `douban_id` | BIGINT | NOT NULL, UNIQUE | 豆瓣电影 ID |
| `title` | VARCHAR(512) | NOT NULL | 电影名称 |
| `director` | VARCHAR(512) | | 导演 |
| `writer` | VARCHAR(512) | | 编剧 |
| `actor` | VARCHAR(2048) | | 主演 (多人) |
| `style` | VARCHAR(512) | | 类型/风格 |
| `site` | VARCHAR(512) | | 官方网站 |
| `country` | VARCHAR(512) | | 制片国家/地区 |
| `language` | VARCHAR(512) | | 语言 |
| `publish_date` | VARCHAR(512) | | 上映日期 |
| `episode` | INT UNSIGNED | | 集数 (电视剧) |
| `duration` | INT UNSIGNED | | 片长 (分钟) |
| `alias` | VARCHAR(512) | | 又名 |
| `imdb` | VARCHAR(512) | | IMDb 链接 |
| `intro` | MEDIUMTEXT | | 简介 |
| `thumbnail` | VARCHAR(512) | | 海报 URL (已转存到 S3) |
| `created_at` | DATETIME | | 记录创建时间 |
| `updated_at` | DATETIME | | 记录更新时间 |

**索引**:
- `PRIMARY KEY (id)`
- `UNIQUE KEY (douban_id)`

---

### 4. game - 游戏信息表

存储游戏条目的详细信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `douban_id` | BIGINT | NOT NULL, UNIQUE | 豆瓣游戏 ID |
| `title` | VARCHAR(512) | NOT NULL | 游戏名称 |
| `platform` | VARCHAR(512) | | 平台 (PC/PS5/Switch 等) |
| `genre` | VARCHAR(512) | | 类型 |
| `alias` | VARCHAR(512) | | 又名 |
| `developer` | VARCHAR(512) | | 开发商 |
| `publisher` | VARCHAR(512) | | 发行商 |
| `publish_date` | VARCHAR(512) | | 发行日期 |
| `intro` | MEDIUMTEXT | | 游戏简介 |
| `thumbnail` | VARCHAR(512) | | 封面图 URL (已转存到 S3) |
| `created_at` | DATETIME | | 记录创建时间 |
| `updated_at` | DATETIME | | 记录更新时间 |

**索引**:
- `PRIMARY KEY (id)`
- `UNIQUE KEY (douban_id)`

---

### 5. song - 音乐信息表

存储音乐专辑的详细信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `douban_id` | BIGINT | NOT NULL, UNIQUE | 豆瓣音乐 ID |
| `title` | VARCHAR(512) | NOT NULL | 专辑名称 |
| `alias` | VARCHAR(512) | | 又名 |
| `musician` | VARCHAR(2048) | | 音乐人/乐队 |
| `album_type` | VARCHAR(512) | | 专辑类型 (录音室/现场/精选等) |
| `genre` | VARCHAR(512) | | 流派 |
| `media` | VARCHAR(512) | | 介质 (CD/黑胶/数字等) |
| `barcode` | VARCHAR(512) | | 条形码 |
| `publisher` | VARCHAR(512) | | 出版者 |
| `publish_date` | VARCHAR(512) | | 发行时间 |
| `isrc` | VARCHAR(512) | | ISRC 编码 |
| `album_count` | INT UNSIGNED | | 唱片数 |
| `intro` | MEDIUMTEXT | | 专辑简介 |
| `track_list` | MEDIUMTEXT | | 曲目列表 |
| `thumbnail` | VARCHAR(512) | | 封面图 URL (已转存到 S3) |
| `created_at` | DATETIME | | 记录创建时间 |
| `updated_at` | DATETIME | | 记录更新时间 |

**索引**:
- `PRIMARY KEY (id)`
- `UNIQUE KEY (douban_id)`

---

### 6. comment - 用户评论表

存储用户对条目的评论/标注/评分记录。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `douban_uid` | BIGINT | NOT NULL, UNIQUE(uk_comment), INDEX | 豆瓣用户 ID |
| `douban_id` | BIGINT | NOT NULL, UNIQUE(uk_comment) | 条目 ID (书/影/音/游) |
| `type` | TINYINT | NOT NULL, UNIQUE(uk_comment), INDEX | 条目类型：1=书，2=电影，3=游戏，4=音乐 |
| `rate` | TINYINT | | 评分 (0-5 星，0 表示未评分) |
| `label` | VARCHAR(512) | | 标签/关键词 |
| `comment` | MEDIUMTEXT | | 评论内容 |
| `action` | TINYINT | NOT NULL, INDEX | 操作类型：0=do,1=wish,2=collect,3=hide |
| `mark_date` | DATETIME | NOT NULL, INDEX | 用户标记时间 (评论发布日期) |
| `created_at` | DATETIME | | 记录创建时间 |
| `updated_at` | DATETIME | | 记录更新时间 |

**索引**:
- `PRIMARY KEY (id)`
- `UNIQUE KEY uk_comment (douban_uid, douban_id, type)`
- `KEY idx_search (douban_uid, type, action, mark_date, updated_at)` - 用于用户评论列表查询

**action 字段枚举**:
| 值 | 说明 |
|----|------|
| 0 | do - 在读/在看/在玩/在听 |
| 1 | wish - 想读/想看/想玩/想听 |
| 2 | collect - 读过/看过/玩过/听过 |
| 3 | hide - 已隐藏 (软删除标记) |

---

### 7. rating - 评分统计表

存储条目的聚合评分数据。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `type` | TINYINT | NOT NULL, UNIQUE(uk_unique_id) | 条目类型：1=书，2=电影，3=游戏，4=音乐 |
| `douban_id` | BIGINT | NOT NULL, UNIQUE(uk_unique_id) | 条目 ID |
| `total` | INT UNSIGNED | | 评分总人数 |
| `rating` | FLOAT | | 平均分 (0-10 分) |
| `star5` | FLOAT | | 5 星占比 (百分比) |
| `star4` | FLOAT | | 4 星占比 (百分比) |
| `star3` | FLOAT | | 3 星占比 (百分比) |
| `star2` | FLOAT | | 2 星占比 (百分比) |
| `star1` | FLOAT | | 1 星占比 (百分比) |
| `status` | TINYINT | | 状态：0=normal,1=not enough,2=not allowed |
| `created_at` | DATETIME | | 记录创建时间 |
| `updated_at` | DATETIME | | 记录更新时间 |

**索引**:
- `PRIMARY KEY (id)`
- `UNIQUE KEY uk_unique_id (type, douban_id)`

**status 字段枚举**:
| 值 | 说明 |
|----|------|
| 0 | normal - 正常显示 |
| 1 | not enough - 评分人数不足 |
| 2 | not allowed - 不允许显示 (被豆瓣限制) |

---

### 8. schedule - 爬虫调度队列

存储待爬取任务队列，控制爬虫工作流程。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `douban_id` | BIGINT | NOT NULL, UNIQUE(uk_schedule) | 目标 ID (用户 ID 或条目 ID) |
| `type` | TINYINT | NOT NULL, UNIQUE(uk_schedule), INDEX | 类型：0=用户，1=书，2=电影，3=游戏，4=音乐 |
| `status` | TINYINT | NOT NULL, INDEX | 爬取状态 |
| `result` | TINYINT | NOT NULL, INDEX | 爬取结果 |
| `created_at` | DATETIME | | 任务创建时间 |
| `updated_at` | DATETIME | | 任务更新时间 (优先级排序字段) |

**索引**:
- `PRIMARY KEY (id)`
- `UNIQUE KEY uk_schedule (douban_id, type)`
- `KEY idx_status (type, status, updated_at)` - 按状态查询任务
- `KEY idx_result (type, result, updated_at)` - 按结果查询任务
- `KEY idx_search (status, result, updated_at)` - 综合查询

**type 字段枚举**:
| 值 | 说明 |
|----|------|
| 0 | user - 用户 |
| 1 | book - 书籍 |
| 2 | movie - 电影 |
| 3 | game - 游戏 |
| 4 | song - 音乐 |

**status 字段枚举**:
| 值 | 说明 |
|----|------|
| 0 | to crawl - 待爬取 |
| 1 | crawling - 爬取中 |
| 2 | crawled - 已爬取 |
| 3 | can crawl - 可爬取 (新发现待处理) |

**result 字段枚举**:
| 值 | 说明 |
|----|------|
| 0 | unready - 未就绪 (等待中) |
| 1 | ready - 就绪 (爬取成功) |
| 2 | invalid - 无效 (条目不存在或被封禁) |

---

### 9. access - 访问日志表

记录 API 访问日志，用于限流和审计。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `douban_uid` | BIGINT | NOT NULL, INDEX | 豆瓣用户 ID (请求参数) |
| `path` | VARCHAR(64) | NOT NULL | 请求路径 |
| `ip` | VARCHAR(64) | NOT NULL, INDEX | 请求 IP |
| `user_agent` | VARCHAR(512) | NOT NULL | User-Agent |
| `referer` | VARCHAR(512) | NOT NULL | Referer 来源 |
| `created_at` | DATETIME | | 访问时间 |
| `updated_at` | DATETIME | | 记录更新时间 |

**索引**:
- `PRIMARY KEY (id)`
- `KEY (douban_uid)`
- `KEY (ip)`

---

### 10. storage - 存储映射表（已弃用）

**注意**: 该表已弃用，不再使用。图片 URL 现在直接保存豆瓣原始地址。

~~记录原始图片 URL 到 S3 存储的映射关系，用于图片转存和去重。~~

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `source` | VARCHAR(256) | NOT NULL, UNIQUE | 原始图片 URL |
| `target` | VARCHAR(256) | NOT NULL, INDEX | S3 存储后的 URL |
| `md5` | VARCHAR(64) | NOT NULL, INDEX | 图片内容 MD5 值 |
| `created_at` | DATETIME | | 记录创建时间 |
| `updated_at` | DATETIME | | 记录更新时间 |

**索引**:
- `PRIMARY KEY (id)`
- `UNIQUE KEY (source)`
- `KEY (target)`
- `KEY (md5)` - 用于图片去重

---

## 核心业务流程

### 用户数据抓取流程

```
1. 接收用户 ID → schedule 表插入任务 (type=0, status=3)
   ↓
2. 爬取用户主页 → user 表 (基本信息 + 各类目数量)
   ↓
3. 爬取 RSS 页面 → 获取 publish_at 用于增量判断
   ↓
4. 爬取评论列表 → comment 表 (用户标注记录)
   ↓
5. 发现新条目 → schedule 表插入任务 (type=1/2/3/4)
   ↓
6. 爬取条目详情 → book/movie/game/song 表 + rating 表
```

### 调度队列处理逻辑

```sql
-- 查询待处理任务 (按 updated_at 排序，优先处理旧的)
SELECT * FROM schedule
WHERE type = ? AND status = 0 AND result = 0
ORDER BY updated_at ASC
LIMIT ?;

-- 任务完成后更新状态
UPDATE schedule SET status = 2, result = 1 WHERE id = ?;

-- 发现新条目前先尝试插入，避免重复
INSERT IGNORE INTO schedule (douban_id, type, status, result)
VALUES (?, ?, 3, 0);
```

### 图片处理逻辑（已简化）

**注意**: S3 存储功能已移除，图片 URL 现在直接保存豆瓣原始地址。

`crawl.Storage()` 函数现在直接返回原始 URL，不再进行下载和上传操作。

---

## 数据量统计

根据 README.md 所述 (截至 2025 年 06 月):

| 数据类型 | 有效条目数 |
|----------|------------|
| 书籍 | ~447 万 |
| 电影 | ~48 万 |
| 音乐 | ~120 万 |
| 游戏 | ~3.8 万 |
