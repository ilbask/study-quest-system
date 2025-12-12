# MySQL 配置说明

## 快速开始

### 1. 安装 MySQL

#### macOS
```bash
brew install mysql
brew services start mysql
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install mysql-server
sudo systemctl start mysql
```

### 2. 配置 MySQL

#### 设置 root 密码
```bash
mysql -u root
```

在 MySQL 命令行中：
```sql
ALTER USER 'root'@'localhost' IDENTIFIED BY 'root';
FLUSH PRIVILEGES;
EXIT;
```

#### 创建数据库
```bash
mysql -u root -p < backend/init.sql
```

或手动创建：
```bash
mysql -u root -p
```

```sql
CREATE DATABASE study_quest CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

### 3. 配置连接信息

#### 方式 1: 环境变量（推荐）
```bash
export MYSQL_DSN="root:root@tcp(127.0.0.1:3306)/study_quest?charset=utf8mb4&parseTime=True&loc=Local"
```

#### 方式 2: 配置文件
创建 `config.yaml`:
```yaml
server:
  port: "8080"

database:
  dsn: "root:root@tcp(127.0.0.1:3306)/study_quest?charset=utf8mb4&parseTime=True&loc=Local"
```

### 4. 启动服务

```bash
./start.sh
```

## 数据库连接字符串格式

```
用户名:密码@tcp(主机:端口)/数据库名?charset=utf8mb4&parseTime=True&loc=Local
```

### 示例

#### 本地开发（默认）
```
root:root@tcp(127.0.0.1:3306)/study_quest?charset=utf8mb4&parseTime=True&loc=Local
```

#### 远程服务器
```
username:password@tcp(192.168.1.100:3306)/study_quest?charset=utf8mb4&parseTime=True&loc=Local
```

#### Docker MySQL
```
root:root@tcp(mysql:3306)/study_quest?charset=utf8mb4&parseTime=True&loc=Local
```

## 自动功能

### 自动迁移（Auto Migrate）
- 服务启动时自动创建/更新表结构
- 无需手动执行 SQL 脚本
- GORM 自动处理表结构变更

### 自动初始化数据（Seed Data）
- 首次启动时自动创建演示账号
- 包含：
  - 学生账号: `student1` / `123456`
  - 家长账号: `parent1` / `123456`
  - 2个演示任务

## 故障排查

### 问题 1: 无法连接 MySQL
**错误信息**: `Failed to connect to MySQL`

**解决方案**:
1. 检查 MySQL 是否运行:
   ```bash
   # macOS
   brew services list
   
   # Linux
   sudo systemctl status mysql
   ```

2. 检查连接信息是否正确
3. 检查防火墙设置

### 问题 2: 数据库不存在
**错误信息**: `Unknown database 'study_quest'`

**解决方案**:
```bash
mysql -u root -p -e "CREATE DATABASE study_quest CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### 问题 3: 权限不足
**错误信息**: `Access denied for user`

**解决方案**:
```sql
GRANT ALL PRIVILEGES ON study_quest.* TO 'root'@'localhost';
FLUSH PRIVILEGES;
```

## 数据备份

### 备份数据库
```bash
mysqldump -u root -p study_quest > backup_$(date +%Y%m%d).sql
```

### 恢复数据库
```bash
mysql -u root -p study_quest < backup_20251212.sql
```

## 生产环境建议

1. **修改默认密码**
   - 不要使用 `root` / `root`
   - 使用强密码

2. **创建专用用户**
   ```sql
   CREATE USER 'study_quest_user'@'localhost' IDENTIFIED BY '强密码';
   GRANT ALL PRIVILEGES ON study_quest.* TO 'study_quest_user'@'localhost';
   FLUSH PRIVILEGES;
   ```

3. **启用 SSL 连接**
   ```
   root:password@tcp(host:3306)/study_quest?charset=utf8mb4&parseTime=True&loc=Local&tls=true
   ```

4. **定期备份**
   - 设置自动备份脚本
   - 保留多个版本

5. **监控性能**
   - 启用慢查询日志
   - 监控连接数
   - 定期优化表

## Docker 部署（可选）

### docker-compose.yml
```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: study_quest
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  app:
    build: .
    environment:
      MYSQL_DSN: "root:root@tcp(mysql:3306)/study_quest?charset=utf8mb4&parseTime=True&loc=Local"
    ports:
      - "8080:8080"
    depends_on:
      - mysql

volumes:
  mysql_data:
```

### 启动
```bash
docker-compose up -d
```

## 回退到内存模式

如果 MySQL 连接失败，系统会自动回退到内存模式：
- 数据存储在内存中
- 重启后数据丢失
- 适合开发和演示

日志会显示：
```
Failed to connect to MySQL: ...
Falling back to In-Memory mode...
```

