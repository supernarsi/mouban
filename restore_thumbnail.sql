-- 恢复 thumbnail 字段为豆瓣原始 URL
-- 将 book、game、movie 表中的 S3 地址恢复为 storage 表中的 source 值

USE mouban;

-- ============================================
-- 1. 先查看需要恢复的数据量
-- ============================================

-- 查看 storage 表中有多少条映射记录
SELECT COUNT(*) AS storage_count FROM storage;

-- 查看 book 表中有多少条 thumbnail 是 S3 地址
SELECT COUNT(*) AS book_s3_count
FROM book b
JOIN storage s ON b.thumbnail = s.target;

-- 查看 movie 表中有多少条 thumbnail 是 S3 地址
SELECT COUNT(*) AS movie_s3_count
FROM movie m
JOIN storage s ON m.thumbnail = s.target;

-- 查看 game 表中有多少条 thumbnail 是 S3 地址
SELECT COUNT(*) AS game_s3_count
FROM game g
JOIN storage s ON g.thumbnail = s.target;

-- ============================================
-- 2. 备份当前数据（可选，建议执行）
-- ============================================

-- 备份 book 表的 thumbnail
CREATE TABLE IF NOT EXISTS book_thumbnail_backup AS
SELECT id, douban_id, thumbnail FROM book;

-- 备份 movie 表的 thumbnail
CREATE TABLE IF NOT EXISTS movie_thumbnail_backup AS
SELECT id, douban_id, thumbnail FROM movie;

-- 备份 game 表的 thumbnail
CREATE TABLE IF NOT EXISTS game_thumbnail_backup AS
SELECT id, douban_id, thumbnail FROM game;

-- ============================================
-- 3. 执行更新操作
-- ============================================

-- 更新 book 表的 thumbnail 为原始 URL
UPDATE book b
INNER JOIN storage s ON b.thumbnail = s.target
SET b.thumbnail = s.source;

-- 更新 movie 表的 thumbnail 为原始 URL
UPDATE movie m
INNER JOIN storage s ON m.thumbnail = s.target
SET m.thumbnail = s.source;

-- 更新 game 表的 thumbnail 为原始 URL
UPDATE game g
INNER JOIN storage s ON g.thumbnail = s.target
SET g.thumbnail = s.source;

-- ============================================
-- 4. 验证更新结果
-- ============================================

-- 验证 book 表更新结果
SELECT
    'book' AS table_name,
    COUNT(*) AS total_count,
    SUM(CASE WHEN thumbnail LIKE 'https://img%' THEN 1 ELSE 0 END) AS douban_url_count,
    SUM(CASE WHEN thumbnail LIKE '%s3%' OR thumbnail LIKE '%minio%' THEN 1 ELSE 0 END) AS s3_url_count
FROM book;

-- 验证 movie 表更新结果
SELECT
    'movie' AS table_name,
    COUNT(*) AS total_count,
    SUM(CASE WHEN thumbnail LIKE 'https://img%' THEN 1 ELSE 0 END) AS douban_url_count,
    SUM(CASE WHEN thumbnail LIKE '%s3%' OR thumbnail LIKE '%minio%' THEN 1 ELSE 0 END) AS s3_url_count
FROM movie;

-- 验证 game 表更新结果
SELECT
    'game' AS table_name,
    COUNT(*) AS total_count,
    SUM(CASE WHEN thumbnail LIKE 'https://img%' THEN 1 ELSE 0 END) AS douban_url_count,
    SUM(CASE WHEN thumbnail LIKE '%s3%' OR thumbnail LIKE '%minio%' THEN 1 ELSE 0 END) AS s3_url_count
FROM game;

-- 查看更新后的示例数据
SELECT 'book' AS table_name, douban_id, thumbnail
FROM book
WHERE thumbnail LIKE 'https://img%'
LIMIT 5;

SELECT 'movie' AS table_name, douban_id, thumbnail
FROM movie
WHERE thumbnail LIKE 'https://img%'
LIMIT 5;

SELECT 'game' AS table_name, douban_id, thumbnail
FROM game
WHERE thumbnail LIKE 'https://img%'
LIMIT 5;

-- ============================================
-- 5. 如果需要恢复备份数据（仅在需要时执行）
-- ============================================

-- 从备份恢复 book 表
-- UPDATE book b
-- INNER JOIN book_thumbnail_backup bk ON b.id = bk.id
-- SET b.thumbnail = bk.thumbnail;

-- 从备份恢复 movie 表
-- UPDATE movie m
-- INNER JOIN movie_thumbnail_backup mv ON m.id = mv.id
-- SET m.thumbnail = mv.thumbnail;

-- 从备份恢复 game 表
-- UPDATE game g
-- INNER JOIN game_thumbnail_backup gm ON g.id = gm.id
-- SET g.thumbnail = gm.thumbnail;
