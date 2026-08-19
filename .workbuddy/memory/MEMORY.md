# 项目长期记忆

## 项目：blog（Go 后端，module 名 blog，原 tuyou）
- 技术栈：Gin + GORM(MySQL) + Redis + JWT + Zap + Viper。
- 前台原型：`prototype-final.html`；后台原型：`admin-prototype.html`（位于仓库根目录）。
- 实体由前台 UI 原型反推生成，放在 `internal/model/entity/`。

## 用户偏好 / 约定
- 删除某目录下的文件时，**只删文件、保留目录本身**（不要 `rm -rf` 整个目录）。
- 用户用中文沟通，回复用简体中文。

## 关键决策
- 删除实体时，基建层（app.go）的 AutoMigrate 改为集中注册当前所有实体，待业务模块成型后可拆为各模块自注册。

## 实体定稿（9 张表，2026-08-18）
- `users`（用户/博主，密码 json:"-" 不外露）、`articles`、`categories`（仅 name，按 id 排序，无 sort 字段）、`tags`、`article_tags`（多对多）、`comments`（parent_id 嵌套）、`comment_likes`、`article_likes`、`site_items`（关于页通用表，type 区分 技术栈分组/技能/联系方式/故事，parent_id 自引用）。
- `articles` 字段约定：`status` 1=已发布 0=下架(隐藏)；`hot` int 热门权重，可人工置顶或按规则计算；冗余计数 view_count/like_count/comment_count。
- 关于页 4 类内容（技术栈分组/技能/联系方式/故事）合并为单表 `site_items`（由最初 4 张表合并而来）。
- 点赞统一两张表：`article_likes`(user_id+article_id)、`comment_likes`(user_id+comment_id)，均复合唯一索引。
