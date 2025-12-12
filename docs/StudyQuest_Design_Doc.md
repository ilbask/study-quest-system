# 小学生作业积分系统 (Study Quest) 详细技术设计文档 V1.0

## 1. 文档概述
### 1.1 项目背景
旨在通过游戏化（Gamification）机制解决家长辅导孩子作业的痛点。将枯燥的学习任务转化为积分挑战，建立“多劳多得、延迟满足”的正向反馈闭环，并通过积分兑换娱乐时间实现劳逸结合管理。

### 1.2 技术栈选型
*   **客户端 (App)**: React Native (TypeScript) —— 实现 iOS/Android 跨平台，支持热更新。
*   **服务端 (Backend)**: Golang (Gin Framework) —— 高性能，易于部署。
*   **数据库 (DB)**: MySQL 8.0 —— 存储核心业务数据（用户、任务、流水）。
*   **缓存 (Cache)**: Redis —— 存储 Session、实时排行榜、动态配置缓存。
*   **运维**: Docker + Nginx —— 容器化部署。

---

## 2. 系统架构设计

### 2.1 逻辑架构图
```
[iOS App] / [Android App]  (React Native)
       │
       ▼
[Nginx Load Balancer / HTTPS Gateway]
       │
       ▼
[Golang API Server] (Stateless) <---- [Config Management] (动态配置)
       │
       ├── [MySQL Master/Slave] (持久化存储)
       └── [Redis Cluster] (热数据/排行榜)
```

### 2.2 核心模块划分
1.  **用户中心 (User Center)**: 身份认证、角色管理（家长/学生/访客）、积分账户管理。
2.  **任务引擎 (Task Engine)**: 任务发布、任务流转（待办->审核->完成）、重复任务调度。
3.  **商城系统 (Marketplace)**: 商品管理（娱乐时间/实物）、库存控制、兑换交易。
4.  **配置中心 (Config Hub)**: **(关键)** 负责下发 UI 主题、功能开关（审核模式）、文案配置，实现 App Store 动态更新与过审控制。
5.  **数据分析 (Analytics)**: 学生行为报表、积分收支曲线。

---

## 3. 数据库详细设计 (MySQL)

### 3.1 用户表 `users`
| 字段名 | 类型 | 说明 | 索引 |
| :--- | :--- | :--- | :--- |
| `id` | BIGINT | 主键 | PK |
| `username` | VARCHAR(64) | 用户名/昵称 | |
| `role` | ENUM | 'parent', 'student', 'guest' | IDX |
| `points_balance` | INT | 当前可用积分余额 | |
| `total_points_earned` | INT | 历史累计获得积分（用于等级计算） | |
| `avatar_url` | VARCHAR(255) | 头像链接 | |
| `family_id` | BIGINT | 家庭组 ID (支持多子女扩展) | IDX |

### 3.2 任务表 `tasks`
| 字段名 | 类型 | 说明 | 索引 |
| :--- | :--- | :--- | :--- |
| `id` | BIGINT | 主键 | PK |
| `title` | VARCHAR(128) | 任务标题 | |
| `points` | INT | 奖励积分 | |
| `type` | TINYINT | 1:学习, 2:家务, 3:习惯 | IDX |
| `recurrence` | VARCHAR(32) | 重复规则 (cron表达式 或 'daily') | |
| `creator_id` | BIGINT | 创建者(家长) ID | |
| `is_active` | BOOLEAN | 是否启用 | |

### 3.3 任务流水表 `task_logs` (核心业务表)
| 字段名 | 类型 | 说明 | 索引 |
| :--- | :--- | :--- | :--- |
| `id` | BIGINT | 主键 | PK |
| `student_id` | BIGINT | 学生 ID | IDX |
| `task_id` | BIGINT | 关联任务 ID | |
| `status` | TINYINT | 0:进行中, 1:待审核, 2:已完成, 3:驳回 | IDX |
| `proof_img` | JSON | 凭证图片列表 | |
| `submitted_at` | DATETIME | 提交时间 | |
| `approved_at` | DATETIME | 审核时间 | |

### 3.4 商品/奖励表 `rewards`
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `id` | BIGINT | 主键 |
| `title` | VARCHAR(128) | 商品名称 (如:看电视30分钟) |
| `cost` | INT | 消耗积分 |
| `category` | TINYINT | 1:虚拟权益(时间), 2:实物奖励 |
| `stock` | INT | 库存 (-1 为无限) |

### 3.5 应用配置表 `app_configs` (动态更新基石)
| 字段名 | 类型 | 说明 | 索引 |
| :--- | :--- | :--- | :--- |
| `key` | VARCHAR(64) | 配置键 (如 'ios_review_mode') | UNIQUE |
| `value` | TEXT | 配置值 (JSON 格式) | |
| `platform` | VARCHAR(10) | 'ios', 'android', 'all' | |
| `min_version` | VARCHAR(10) | 最低生效版本号 | |

---

## 4. API 接口设计 (RESTful v1)

所有接口需带 `Authorization: Bearer <token>` 头。

### 4.1 通用/配置
*   `GET /api/v1/config/init`
    *   **功能**: App 启动时调用，获取全局配置。
    *   **返回**: 
        ```json
        {
          "review_mode": false, // iOS审核开关
          "theme": "default",
          "features": ["game_center", "mall"], // 启用的功能模块
          "daily_quote": "今日名言..."
        }
        ```

### 4.2 任务模块
*   `GET /api/v1/tasks/today`
    *   **功能**: 获取今日任务列表（包含已完成和未完成的状态）。
*   `POST /api/v1/tasks/submit`
    *   **参数**: `{ "task_id": 101, "proof_img": "url..." }`
    *   **功能**: 学生提交任务，状态变更为“待审核”。
*   `POST /api/v1/tasks/approve` (家长权限)
    *   **参数**: `{ "log_id": 505, "action": "approve" }`
    *   **功能**: 审核通过，**触发数据库事务**：更新 log 状态 -> 增加用户积分 -> 写入积分流水。

### 4.3 积分与商城
*   `GET /api/v1/mall/items`
    *   **功能**: 获取可兑换商品列表。
*   `POST /api/v1/mall/redeem`
    *   **参数**: `{ "item_id": 202 }`
    *   **功能**: 兑换商品，扣除积分。若为“娱乐时间”，服务端记录开始时间戳。

---

## 5. App Store 上架合规与动态更新策略

### 5.1 审核模式 (Review Mode)
这是确保包含“游戏化”、“娱乐控制”功能 App 能顺利过审的关键设计。

1.  **服务端控制**: 在 `app_configs` 表中设置 `ios_review_mode = true`。
2.  **API 响应**: 当 App 启动调用 `/config/init` 时，返回 `review_mode: true`。
3.  **前端表现**:
    *   隐藏“娱乐游戏”入口。
    *   隐藏“锁屏控制”等敏感描述。
    *   界面展示为纯粹的“任务清单管理工具”或“习惯养成工具”。
4.  **过审后**: 服务端修改配置为 `false`，用户重启 App 即可看到完整功能。

### 5.2 隐私合规
*   **儿童隐私保护 (COPPA)**:
    *   注册时需由家长操作。
    *   不收集精确地理位置。
    *   头像上传仅支持系统预设或本地处理，避免上传真人照片（或需严格审核）。
*   **权限申请**:
    *   **相机**: "用于拍摄作业完成情况以便家长审核"。
    *   **通知**: "用于接收任务提醒和审核结果"。

### 5.3 动态更新 (CodePush)
*   集成 `react-native-code-push`。
*   对于 UI 调整、Bug 修复、非原生代码的功能迭代，直接通过 CodePush 下发 JS Bundle，无需重新走 App Store 审核流程。

---

## 6. 未来扩展设计 (Phase 2)

1.  **AI 作业批改**: 引入 OCR SDK，后端对接 GPT-4o 或垂直教育模型，实现数学题自动判卷。
2.  **IoT 协议**: 定义 MQTT Topic，例如 `device/tv/control`，配合家庭网关实现积分耗尽自动关闭电视。
3.  **SaaS 化**: 支持多家庭租户，为学校或培训机构提供“班级作业管理端”，老师布置作业，家长监督任务。

