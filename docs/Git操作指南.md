# Git 操作快速参考

## 当前仓库信息
- **本地路径**: `/Users/it/debao.huang/private_project/study-quest-system`
- **远程仓库**: https://github.com/ilbask/study-quest-system.git
- **当前分支**: main
- **最新提交**: 4073cad - first version

---

## 常用 Git 操作

### 1. 查看当前状态
```bash
cd /Users/it/debao.huang/private_project/study-quest-system
git status
```

### 2. 添加新文件或修改
```bash
# 添加所有变更
git add -A

# 或添加指定文件
git add <文件路径>

# 或添加指定目录
git add backend/
```

### 3. 提交代码
```bash
# 提交暂存区的所有更改
git commit -m "提交说明"

# 示例
git commit -m "feat: 添加任务审核功能"
git commit -m "fix: 修复积分计算错误"
git commit -m "docs: 更新 README"
```

### 4. 推送到远程
```bash
# 推送到 main 分支
git push origin main

# 如果已经设置了上游分支（upstream），可以简化为
git push
```

### 5. 拉取远程最新代码
```bash
# 拉取并合并
git pull origin main

# 或简化为
git pull
```

### 6. 查看提交历史
```bash
# 查看简洁历史
git log --oneline -10

# 查看详细历史
git log -5

# 查看图形化历史
git log --graph --oneline --all
```

### 7. 创建和切换分支
```bash
# 创建新分支
git branch feature/new-feature

# 切换到新分支
git checkout feature/new-feature

# 或一步完成
git checkout -b feature/new-feature

# 推送新分支到远程
git push -u origin feature/new-feature
```

### 8. 查看远程仓库信息
```bash
git remote -v
git remote show origin
```

---

## 推荐的开发流程

### 功能开发流程
```bash
# 1. 确保在最新的 main 分支
git checkout main
git pull origin main

# 2. 创建功能分支
git checkout -b feature/积分兑换功能

# 3. 开发并提交
git add backend/internal/service/reward.go
git commit -m "feat: 实现积分兑换服务"

# 4. 推送到远程
git push -u origin feature/积分兑换功能

# 5. 在 GitHub 上创建 Pull Request

# 6. 合并后切回 main 并删除本地分支
git checkout main
git pull origin main
git branch -d feature/积分兑换功能
```

### Bug 修复流程
```bash
# 1. 创建修复分支
git checkout -b fix/积分计算错误

# 2. 修复并提交
git add backend/internal/service/task.go
git commit -m "fix: 修复任务审核时积分未正确增加的问题"

# 3. 推送并创建 PR
git push -u origin fix/积分计算错误
```

---

## Commit Message 规范

推荐使用 Conventional Commits 规范：

### 格式
```
<类型>(<范围>): <简短描述>

<详细描述>（可选）

<Footer>（可选）
```

### 常用类型
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式调整（不影响功能）
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建/工具链相关

### 示例
```bash
feat(backend): 添加用户积分查询接口
fix(frontend): 修复任务列表刷新问题
docs: 更新 App Store 上架指南
refactor(service): 重构任务审核逻辑
```

---

## .gitignore 建议

确保项目根目录有 `.gitignore` 文件，避免提交不必要的文件：

```gitignore
# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
.go/
go.work

# Node / React Native
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*
.expo/
.expo-shared/

# iOS
ios/Pods/
ios/build/
*.pbxuser
*.mode1v3
*.mode2v3
*.perspectivev3
*.xcuserstate
project.xcworkspace/
xcuserdata/

# Android
android/build/
android/.gradle/
*.apk
*.ap_
*.aab

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# macOS
.DS_Store

# Environment
.env
.env.local

# Database
*.db
*.sqlite
</ignore>
```

---

## 常见问题

### Q1: 如何撤销未提交的修改？
```bash
# 撤销指定文件
git checkout -- <文件名>

# 撤销所有修改
git checkout -- .
```

### Q2: 如何撤销已暂存的文件？
```bash
# 取消暂存指定文件
git reset HEAD <文件名>

# 取消所有暂存
git reset HEAD
```

### Q3: 如何修改最后一次提交？
```bash
# 修改提交信息
git commit --amend -m "新的提交信息"

# 添加遗漏的文件到上次提交
git add 遗漏的文件
git commit --amend --no-edit
```

### Q4: 如何查看某个文件的修改历史？
```bash
git log -p <文件名>
```

### Q5: 推送时提示 "rejected"？
```bash
# 先拉取远程代码
git pull origin main

# 解决冲突后再推送
git push origin main
```

---

## 团队协作建议

1. **保持 main 分支稳定**：不要直接在 main 上开发
2. **功能分支开发**：每个功能使用独立分支
3. **频繁提交**：小步快跑，每天至少提交一次
4. **Pull Request 审查**：大功能需要代码审查
5. **定期同步**：每天开始工作前先 `git pull`

---

**祝开发顺利！🚀**

