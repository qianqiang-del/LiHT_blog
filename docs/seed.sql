-- ==========================================
-- LiHT 博客种子数据（完整版，覆盖全部 9 张表）
-- 导入前提：后端已启动过一次（AutoMigrate 建表完成）
-- 导入命令：mysql -u root -p liht < docs/seed.sql
-- 说明：所有用户密码均为 bcrypt("password")，即明文 password
-- ==========================================

SET NAMES utf8mb4;

USE `liht`;

-- -------------------------------------------
-- 清空旧数据（先删有外键/关联的表，再删主表）
-- -------------------------------------------
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE `comment_likes`;
TRUNCATE TABLE `article_likes`;
TRUNCATE TABLE `comments`;
TRUNCATE TABLE `article_tags`;
TRUNCATE TABLE `articles`;
TRUNCATE TABLE `tags`;
TRUNCATE TABLE `categories`;
TRUNCATE TABLE `site_items`;
TRUNCATE TABLE `users`;
SET FOREIGN_KEY_CHECKS = 1;

-- -------------------------------------------
-- 用户（博主 admin + 普通用户）
-- 注意：实体已含 github 列
-- -------------------------------------------
INSERT INTO `users` (`id`, `username`, `email`, `password`, `nickname`, `avatar`, `bio`, `github`, `role`, `status`, `created_at`, `updated_at`) VALUES
(1, 'admin',    'admin@liht.com',    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '牵强',   'https://api.dicebear.com/7.x/avataaars/svg?seed=admin',  '全栈开发者，记录技术与生活', 'https://github.com/example',   2, 1, '2026-01-01 10:00:00', '2026-08-01 10:00:00'),
(2, 'zhangsan', 'zhangsan@liht.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '张三',   'https://api.dicebear.com/7.x/avataaars/svg?seed=zhang',  '前端爱好者',                  'https://github.com/zhangsan', 1, 1, '2026-02-01 10:00:00', '2026-08-01 10:00:00'),
(3, 'lisi',     'lisi@liht.com',     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '李四',   'https://api.dicebear.com/7.x/avataaars/svg?seed=lisi',   '后端工程师',                  'https://github.com/lisi',     1, 1, '2026-03-01 10:00:00', '2026-08-01 10:00:00');

-- -------------------------------------------
-- 分类
-- -------------------------------------------
INSERT INTO `categories` (`id`, `name`, `created_at`, `updated_at`) VALUES
(1, '前端',   '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(2, '后端',   '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(3, '数据库', '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(4, 'DevOps', '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(5, '随笔',   '2026-01-01 10:00:00', '2026-01-01 10:00:00');

-- -------------------------------------------
-- 标签
-- -------------------------------------------
INSERT INTO `tags` (`id`, `name`, `created_at`, `updated_at`) VALUES
(1,  'Vue',         '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(2,  'React',       '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(3,  'Go',          '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(4,  'JavaScript',  '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(5,  'TypeScript',  '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(6,  'MySQL',       '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(7,  'Redis',       '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(8,  'Docker',      '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(9,  'Linux',       '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(10, 'CSS',         '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(11, 'Node.js',     '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(12, 'Git',         '2026-01-01 10:00:00', '2026-01-01 10:00:00');

-- -------------------------------------------
-- 文章（12 篇）
-- -------------------------------------------
INSERT INTO `articles` (`id`, `title`, `summary`, `content`, `cover`, `author_id`, `category_id`, `status`, `view_count`, `like_count`, `comment_count`, `hot`, `published_at`, `created_at`, `updated_at`) VALUES
(1,  'Vue 3 组合式 API 入门指南',
     '深入了解 Vue 3 Composition API 的核心概念与实战技巧。',
     '# Vue 3 组合式 API 入门\n\nVue 3 引入了 Composition API，让代码组织更加灵活...\n\n## 为什么选择组合式 API\n\n1. 更好的逻辑复用\n2. 更灵活的代码组织\n3. 更好的类型推导\n\n## 核心概念\n\n### ref 和 reactive\n\n~~~js\nimport { ref, reactive } from "vue"\n\nconst count = ref(0)\nconst state = reactive({ name: "Vue" })\n~~~\n\n### computed 和 watch\n\n~~~js\nimport { computed, watch } from "vue"\n\nconst double = computed(() => count.value * 2)\nwatch(count, (newVal) => console.log(newVal))\n~~~\n\n组合式 API 让我们可以按功能而非选项来组织代码，大型项目中尤为明显。',
     'https://picsum.photos/seed/vue3/800/400', 1, 1, 1, 1284, 36, 12, 90, '2026-08-01 09:30:00', '2026-08-01 09:30:00', '2026-08-01 09:30:00'),

(2,  'Go 语言并发编程实战',
     '通过实例学习 Goroutine 和 Channel 的使用技巧。',
     '# Go 并发编程\n\nGo 的并发模型基于 CSP（Communicating Sequential Processes）...\n\n## Goroutine\n\n~~~go\nfunc main() {\n    go func() {\n        fmt.Println("Hello from goroutine")\n    }()\n    time.Sleep(time.Second)\n}\n~~~\n\n## Channel\n\n~~~go\nch := make(chan int)\ngo func() {\n    ch <- 42\n}()\nfmt.Println(<-ch)\n~~~\n\n## select 多路复用\n\n~~~go\nselect {\ncase msg := <-ch1:\n    fmt.Println(msg)\ncase msg := <-ch2:\n    fmt.Println(msg)\ndefault:\n    fmt.Println("no message")\n}\n~~~\n\n记住：不要通过共享内存来通信，而要通过通信来共享内存。',
     'https://picsum.photos/seed/golang/800/400', 1, 2, 1, 2150, 58, 23, 120, '2026-07-28 14:00:00', '2026-07-28 14:00:00', '2026-07-28 14:00:00'),

(3,  'MySQL 索引优化实践',
     '深入理解 B+Tree 索引原理，掌握慢查询优化技巧。',
     '# MySQL 索引优化\n\n## B+Tree 索引结构\n\nInnoDB 使用 B+Tree 作为索引结构，聚簇索引存储完整行数据...\n\n## 最左前缀原则\n\n联合索引 (a, b, c) 可以匹配：\n- a\n- a, b\n- a, b, c\n\n不能匹配：\n- b\n- c\n- b, c\n\n## EXPLAIN 分析\n\n~~~sql\nEXPLAIN SELECT * FROM articles WHERE author_id = 1;\n~~~\n\n关注 type、key、rows、Extra 字段。\n\n## 常见优化\n\n1. 避免 SELECT *\n2. 避免在索引列上使用函数\n3. 小表驱动大表\n4. 覆盖索引减少回表',
     'https://picsum.photos/seed/mysql/800/400', 1, 3, 1, 1876, 42, 15, 100, '2026-07-25 10:00:00', '2026-07-25 10:00:00', '2026-07-25 10:00:00'),

(4,  'Docker 容器化部署 Node.js 应用',
     '从零开始学习 Dockerfile 编写与多阶段构建。',
     '# Docker 部署 Node.js\n\n## Dockerfile 编写\n\n~~~dockerfile\nFROM node:18-alpine AS builder\nWORKDIR /app\nCOPY package*.json ./\nRUN npm ci --only=production\nCOPY . .\nRUN npm run build\n\nFROM node:18-alpine\nWORKDIR /app\nCOPY --from=builder /app/dist ./dist\nCOPY --from=builder /app/node_modules ./node_modules\nEXPOSE 3000\nCMD ["node", "dist/index.js"]\n~~~\n\n## 最佳实践\n\n1. 使用 .dockerignore\n2. 多阶段构建减小镜像\n3. 合理利用缓存层\n4. 不要以 root 运行',
     'https://picsum.photos/seed/docker/800/400', 1, 4, 1, 965, 28, 8, 45, '2026-07-20 16:00:00', '2026-07-20 16:00:00', '2026-07-20 16:00:00'),

(5,  'TypeScript 高级类型体操',
     '掌握条件类型、映射类型、模板字面量类型等进阶用法。',
     '# TypeScript 高级类型\n\n## 条件类型\n\n~~~ts\ntype IsString<T> = T extends string ? true : false\ntype A = IsString<string> // true\ntype B = IsString<number> // false\n~~~\n\n## 映射类型\n\n~~~ts\ntype Readonly<T> = {\n  readonly [K in keyof T]: T[K]\n}\n~~~\n\n## 模板字面量类型\n\n~~~ts\ntype EventName = onCapitalize<string>\n~~~\n\n类型体操是 TypeScript 最强大的特性之一，掌握它可以让类型定义更加精确。',
     'https://picsum.photos/seed/typescript/800/400', 2, 1, 1, 1432, 35, 11, 78, '2026-07-18 08:30:00', '2026-07-18 08:30:00', '2026-07-18 08:30:00'),

(6,  'Redis 缓存策略与常见问题',
     '缓存穿透、缓存击穿、缓存雪崩的原理与解决方案。',
     '# Redis 缓存策略\n\n## 缓存穿透\n\n查询不存在的数据，每次都打到数据库。\n\n**解决方案：** 布隆过滤器 / 缓存空值\n\n## 缓存击穿\n\n热点 key 过期，大量请求打到数据库。\n\n**解决方案：** 互斥锁 / 永不过期\n\n## 缓存雪崩\n\n大量 key 同时过期。\n\n**解决方案：** 过期时间加随机值 / 多级缓存\n\n## 常用数据结构\n\n- String：计数器、分布式锁\n- Hash：对象存储\n- List：消息队列\n- Set：去重\n- ZSet：排行榜',
     'https://picsum.photos/seed/redis/800/400', 3, 2, 1, 1098, 31, 9, 55, '2026-07-15 11:00:00', '2026-07-15 11:00:00', '2026-07-15 11:00:00'),

(7,  'React 18 新特性解析',
     '深入理解 Suspense、useTransition、useDeferredValue 等新 API。',
     '# React 18 新特性\n\n## 自动批处理\n\nReact 18 默认对所有状态更新进行批处理，包括 setTimeout、Promise 中的更新。\n\n## useTransition\n\n~~~jsx\nconst [isPending, startTransition] = useTransition()\nstartTransition(() => {\n  setTab(\'comments\')\n})\n~~~\n\n## Suspense 改进\n\n支持服务端流式渲染，可以逐步发送 HTML。\n\n## useId\n\n生成唯一 ID，解决 SSR hydration 不匹配问题。',
     'https://picsum.photos/seed/react18/800/400', 2, 1, 1, 876, 22, 6, 38, '2026-07-12 15:30:00', '2026-07-12 15:30:00', '2026-07-12 15:30:00'),

(8,  'Linux 常用命令速查手册',
     '开发者必备的 Linux 命令整理，涵盖文件、进程、网络等。',
     '# Linux 命令速查\n\n## 文件操作\n\n~~~bash\nls -la          # 列出文件\ncp -r src dst   # 复制目录\nfind . -name "*.go" # 查找文件\n~~~\n\n## 进程管理\n\n~~~bash\nps aux          # 查看进程\ntop             # 实时监控\nkill -9 PID     # 强制终止\n~~~\n\n## 网络\n\n~~~bash\nnetstat -tlnp   # 查看端口\ncurl -I url     # 查看响应头\nss -tlnp        # 查看监听端口\n~~~\n\n## 文本处理\n\n~~~bash\ngrep -r "pattern" .  # 搜索文本\nawk \'{print $1}\'     # 提取列\nsed -i \'s/old/new/g\' # 替换文本\n~~~',
     'https://picsum.photos/seed/linux/800/400', 1, 4, 1, 2340, 65, 18, 130, '2026-07-10 09:00:00', '2026-07-10 09:00:00', '2026-07-10 09:00:00'),

(9,  'CSS Grid 布局完全指南',
     '从基础到实战，掌握 CSS Grid 二维布局系统。',
     '# CSS Grid 布局\n\n## 基本概念\n\n~~~css\n.container {\n  display: grid;\n  grid-template-columns: repeat(3, 1fr);\n  grid-template-rows: auto;\n  gap: 20px;\n}\n~~~\n\n## 常用属性\n\n- grid-template-columns / rows\n- grid-column / grid-row\n- grid-area\n- justify-items / align-items\n\n## 响应式布局\n\n~~~css\n.container {\n  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));\n}\n~~~\n\nGrid 是二维布局的最佳选择，Flexbox 适合一维布局。',
     'https://picsum.photos/seed/cssgrid/800/400', 2, 1, 1, 756, 19, 5, 30, '2026-07-08 13:00:00', '2026-07-08 13:00:00', '2026-07-08 13:00:00'),

(10, 'Go Gin 框架快速上手',
     '使用 Gin 构建 RESTful API 的完整教程。',
     '# Gin 框架入门\n\n## 安装\n\n~~~bash\ngo get github.com/gin-gonic/gin\n~~~\n\n## 基本用法\n\n~~~go\nfunc main() {\n    r := gin.Default()\n    r.GET("/ping", func(c *gin.Context) {\n        c.JSON(200, gin.H{"message": "pong"})\n    })\n    r.Run(":8080")\n}\n~~~\n\n## 路由分组\n\n~~~go\nv1 := r.Group("/api/v1")\n{\n    v1.GET("/articles", listArticles)\n    v1.POST("/articles", createArticle)\n}\n~~~\n\n## 中间件\n\n~~~go\nr.Use(authMiddleware())\n~~~\n\nGin 是 Go 生态中最流行的 Web 框架，性能优异且易于使用。',
     'https://picsum.photos/seed/gin/800/400', 1, 2, 1, 1120, 33, 10, 62, '2026-07-05 10:30:00', '2026-07-05 10:30:00', '2026-07-05 10:30:00'),

(11, 'Git 分支管理最佳实践',
     '团队协作中的 Git 工作流选择与分支策略。',
     '# Git 分支管理\n\n## 常见工作流\n\n### Git Flow\n\n- main：生产环境\n- develop：开发分支\n- feature/*：功能分支\n- release/*：发布分支\n- hotfix/*：热修复\n\n### Trunk Based\n\n所有人在 main 上开发，通过 feature flag 控制功能发布。\n\n## 提交规范\n\n~~~\nfeat: 新功能\nfix: 修复\ndocs: 文档\nstyle: 格式\nrefactor: 重构\ntest: 测试\nchore: 构建/工具\n~~~\n\n## 实用命令\n\n~~~bash\ngit rebase -i HEAD~3  # 交互式变基\ngit stash             # 暂存修改\ngit cherry-pick abc   # 摘取提交\n~~~',
     'https://picsum.photos/seed/git/800/400', 3, 4, 1, 654, 15, 4, 22, '2026-07-02 17:00:00', '2026-07-02 17:00:00', '2026-07-02 17:00:00'),

(12, '写代码之外的思考',
     '关于技术成长、工作与生活平衡的一些随想。',
     '# 写代码之外的思考\n\n## 技术深度 vs 广度\n\n初学者常犯的错误是什么都学，什么都不精。\n\n建议：\n1. 先精通一门语言\n2. 理解计算机基础\n3. 再横向扩展\n\n## 保持学习\n\n- 每周读一篇技术文章\n- 每月完成一个小项目\n- 每季度学一项新技术\n\n## 工作与生活\n\n写代码是工作，不是生活的全部。\n\n适当休息、运动、社交，反而能提高编程效率。\n\n> "最好的调试工具是充足的睡眠。"',
     'https://picsum.photos/seed/think/800/400', 1, 5, 1, 432, 12, 3, 15, '2026-06-28 20:00:00', '2026-06-28 20:00:00', '2026-06-28 20:00:00');

-- -------------------------------------------
-- 文章-标签关联（每篇 0~3 个标签）
-- -------------------------------------------
INSERT INTO `article_tags` (`article_id`, `tag_id`) VALUES
(1, 1), (1, 4), (1, 5),
(2, 3),
(3, 6),
(4, 8), (4, 9), (4, 11),
(5, 5), (5, 4),
(6, 7), (6, 6),
(7, 2), (7, 4), (7, 5),
(8, 9),
(9, 10),
(10, 3),
(11, 12);

-- -------------------------------------------
-- 评论（含嵌套回复，parent_id 指向一级评论）
-- -------------------------------------------
INSERT INTO `comments` (`id`, `article_id`, `parent_id`, `user_id`, `nickname`, `content`, `like_count`, `status`, `created_at`, `updated_at`) VALUES
-- 文章1：Vue 3
(1, 1, NULL, 2, '张三', '写得很清楚，组合式 API 终于看懂了，ref 那段讲得最好。', 3, 1, '2026-08-02 10:00:00', '2026-08-02 10:00:00'),
(2, 1, 1,    3, '李四', '回复张三：同感，希望多来点实战例子。', 1, 1, '2026-08-02 11:00:00', '2026-08-02 11:00:00'),
(3, 1, 1,    2, '张三', '补充：watch 默认是浅层的，深层次要加 deep 选项。', 0, 1, '2026-08-02 12:00:00', '2026-08-02 12:00:00'),
(4, 1, NULL, 3, '李四', '收藏了，期待下一篇关于 setup 语法的文章。', 2, 1, '2026-08-03 09:00:00', '2026-08-03 09:00:00'),
(5, 1, NULL, 1, '牵强', '感谢支持，下一篇会写状态管理。', 1, 1, '2026-08-03 10:00:00', '2026-08-03 10:00:00'),
-- 文章2：Go 并发
(6, 2, NULL, 2, '张三', 'channel 的示例很经典，CSP 那节解释得很清楚。', 4, 1, '2026-07-29 10:00:00', '2026-07-29 10:00:00'),
(7, 2, 6,    1, '牵强', '感谢支持，后续会更新 select 的实战场景。', 2, 1, '2026-07-29 11:00:00', '2026-07-29 11:00:00'),
(8, 2, NULL, 3, '李四', 'goroutine 泄漏的问题能再讲讲吗？', 1, 1, '2026-07-30 09:00:00', '2026-07-30 09:00:00'),
-- 文章3：MySQL 索引
(9, 3, NULL, 3, '李四', 'B+Tree 那节可以再配一张图就更直观了。', 1, 1, '2026-07-26 10:00:00', '2026-07-26 10:00:00'),
-- 文章5：TypeScript
(10, 5, NULL, 1, '牵强', '条件类型的例子很简洁，适合入门。', 0, 1, '2026-07-19 09:00:00', '2026-07-19 09:00:00'),
-- 文章8：Linux
(11, 8, NULL, 2, '张三', 'mark 一下，常用命令整理得很全。', 1, 1, '2026-07-11 10:00:00', '2026-07-11 10:00:00');

-- -------------------------------------------
-- 文章点赞（user_id + article_id 唯一）
-- -------------------------------------------
INSERT INTO `article_likes` (`id`, `article_id`, `user_id`, `created_at`, `updated_at`) VALUES
(1, 1, 1, '2026-08-02 10:00:00', '2026-08-02 10:00:00'),
(2, 1, 2, '2026-08-02 11:00:00', '2026-08-02 11:00:00'),
(3, 1, 3, '2026-08-02 12:00:00', '2026-08-02 12:00:00'),
(4, 2, 1, '2026-07-29 10:00:00', '2026-07-29 10:00:00'),
(5, 2, 2, '2026-07-29 11:00:00', '2026-07-29 11:00:00'),
(6, 3, 1, '2026-07-26 10:00:00', '2026-07-26 10:00:00'),
(7, 3, 3, '2026-07-26 11:00:00', '2026-07-26 11:00:00'),
(8, 5, 2, '2026-07-19 09:00:00', '2026-07-19 09:00:00'),
(9, 8, 3, '2026-07-11 10:00:00', '2026-07-11 10:00:00');

-- -------------------------------------------
-- 评论点赞（user_id + comment_id 唯一）
-- -------------------------------------------
INSERT INTO `comment_likes` (`id`, `comment_id`, `user_id`, `created_at`, `updated_at`) VALUES
(1, 1, 1, '2026-08-02 10:00:00', '2026-08-02 10:00:00'),
(2, 1, 2, '2026-08-02 10:05:00', '2026-08-02 10:05:00'),
(3, 1, 3, '2026-08-02 10:10:00', '2026-08-02 10:10:00'),
(4, 2, 3, '2026-08-02 11:00:00', '2026-08-02 11:00:00'),
(5, 4, 1, '2026-08-03 09:00:00', '2026-08-03 09:00:00'),
(6, 4, 2, '2026-08-03 09:05:00', '2026-08-03 09:05:00'),
(7, 5, 2, '2026-08-03 10:00:00', '2026-08-03 10:00:00'),
(8, 6, 1, '2026-07-29 10:00:00', '2026-07-29 10:00:00'),
(9, 6, 2, '2026-07-29 10:05:00', '2026-07-29 10:05:00'),
(10, 6, 3, '2026-07-29 10:10:00', '2026-07-29 10:10:00'),
(11, 7, 3, '2026-07-29 11:00:00', '2026-07-29 11:00:00'),
(12, 8, 2, '2026-07-30 09:00:00', '2026-07-30 09:00:00');

-- -------------------------------------------
-- 关于页内容（site_items）
-- type: skill_group | skill | contact | story
-- -------------------------------------------
INSERT INTO `site_items` (`id`, `type`, `title`, `content`, `parent_id`, `code`, `icon`, `sort`, `created_at`, `updated_at`) VALUES
-- 技术栈分组
(1,  'skill_group', '后端',   '', NULL, '', '', 1, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(2,  'skill_group', '前端',   '', NULL, '', '', 2, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(3,  'skill_group', 'DevOps', '', NULL, '', '', 3, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
-- 技术栈项（parent_id 指向分组）
(4,  'skill', 'Go',          '', 1, '', '', 1, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(5,  'skill', 'Gin',         '', 1, '', '', 2, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(6,  'skill', 'MySQL',       '', 1, '', '', 3, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(7,  'skill', 'Redis',       '', 1, '', '', 4, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(8,  'skill', 'Vue',         '', 2, '', '', 1, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(9,  'skill', 'React',       '', 2, '', '', 2, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(10, 'skill', 'TypeScript',  '', 2, '', '', 3, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(11, 'skill', 'Docker',      '', 3, '', '', 1, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(12, 'skill', 'Linux',       '', 3, '', '', 2, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
-- 联系方式
(13, 'contact', 'GitHub', 'https://github.com/example',   NULL, 'github', 'github', 1, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(14, 'contact', '邮箱',   'admin@liht.com',               NULL, 'email',  'email',  2, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(15, 'contact', '微信',   'qianqiang-dev',                NULL, 'wechat', 'wechat', 3, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
-- 博客故事时间线（title=时间文本，content=故事）
(16, 'story', '2020-03', '第一次搭建个人博客，写下第一篇文章。', NULL, '', '', 1, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(17, 'story', '2021-06', '系统学习 Go 语言，开始记录后端学习笔记。', NULL, '', '', 2, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(18, 'story', '2022-09', '转型全栈，前端与后端一起深耕。', NULL, '', '', 3, '2026-01-01 10:00:00', '2026-01-01 10:00:00'),
(19, 'story', '2024-01', '重写博客，采用前后端分离架构。', NULL, '', '', 4, '2026-01-01 10:00:00', '2026-01-01 10:00:00');
