#!/bin/bash

# Study Quest 系统启动脚本
# 用途：自动检查、停止旧服务、编译并启动新服务

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目目录
PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="${PROJECT_DIR}/backend"
WEB_DIR="${PROJECT_DIR}/web"
LOG_DIR="${PROJECT_DIR}/logs"

# 端口配置
BACKEND_PORT=8080

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   Study Quest 系统启动脚本${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 创建日志目录
mkdir -p "${LOG_DIR}"

# 1. 检查并停止旧服务
echo -e "${YELLOW}[1/4] 检查并停止旧服务...${NC}"

# 检查端口占用
OLD_PID=$(lsof -ti:${BACKEND_PORT} 2>/dev/null || true)

if [ -n "$OLD_PID" ]; then
    echo -e "${YELLOW}  发现端口 ${BACKEND_PORT} 被进程 ${OLD_PID} 占用${NC}"
    echo -e "${YELLOW}  正在停止旧服务...${NC}"
    
    # 尝试优雅停止
    kill $OLD_PID 2>/dev/null || true
    sleep 2
    
    # 检查是否还在运行
    if lsof -ti:${BACKEND_PORT} >/dev/null 2>&1; then
        echo -e "${RED}  优雅停止失败，强制终止...${NC}"
        kill -9 $OLD_PID 2>/dev/null || true
        sleep 1
    fi
    
    echo -e "${GREEN}  ✓ 旧服务已停止${NC}"
else
    echo -e "${GREEN}  ✓ 端口 ${BACKEND_PORT} 空闲${NC}"
fi

# 2. 检查并安装依赖
echo ""
echo -e "${YELLOW}[2/4] 检查依赖...${NC}"

# 检查 Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}  ✗ 未找到 Go，请先安装 Go 1.21+${NC}"
    exit 1
fi
echo -e "${GREEN}  ✓ Go 版本: $(go version | awk '{print $3}')${NC}"

# 检查后端依赖
cd "${BACKEND_DIR}"
if [ ! -f "go.mod" ]; then
    echo -e "${RED}  ✗ 未找到 go.mod 文件${NC}"
    exit 1
fi

echo -e "${YELLOW}  正在同步 Go 依赖...${NC}"
export GOPROXY=https://goproxy.cn,direct
go mod tidy
echo -e "${GREEN}  ✓ 依赖同步完成${NC}"

# 3. 编译后端
echo ""
echo -e "${YELLOW}[3/4] 编译后端代码...${NC}"

cd "${BACKEND_DIR}"
echo -e "${YELLOW}  正在编译 cmd/api/main.go...${NC}"

# 编译
if go build -o "${PROJECT_DIR}/study-quest-server" cmd/api/main.go; then
    echo -e "${GREEN}  ✓ 后端编译成功${NC}"
else
    echo -e "${RED}  ✗ 后端编译失败${NC}"
    exit 1
fi

# 4. 启动服务
echo ""
echo -e "${YELLOW}[4/4] 启动服务...${NC}"

cd "${PROJECT_DIR}"

# 后台启动服务并记录日志
nohup ./study-quest-server > "${LOG_DIR}/backend.log" 2>&1 &
BACKEND_PID=$!

# 等待服务启动
echo -e "${YELLOW}  等待服务启动...${NC}"
sleep 3

# 检查服务是否启动成功
if ps -p $BACKEND_PID > /dev/null; then
    echo -e "${GREEN}  ✓ 后端服务已启动 (PID: ${BACKEND_PID})${NC}"
    
    # 测试API
    sleep 1
    if curl -s http://localhost:${BACKEND_PORT}/api/v1/profile > /dev/null 2>&1; then
        echo -e "${GREEN}  ✓ API 响应正常${NC}"
    else
        echo -e "${YELLOW}  ⚠ API 暂时无响应，请稍后再试${NC}"
    fi
else
    echo -e "${RED}  ✗ 后端服务启动失败${NC}"
    echo -e "${RED}  请查看日志: ${LOG_DIR}/backend.log${NC}"
    exit 1
fi

# 输出访问信息
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}  ✓ 系统启动成功！${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${GREEN}访问地址:${NC}"
echo -e "  • Web Demo:  ${BLUE}http://localhost:${BACKEND_PORT}/web${NC}"
echo -e "  • API 文档:  ${BLUE}http://localhost:${BACKEND_PORT}/api/v1/profile${NC}"
echo ""
echo -e "${GREEN}日志文件:${NC}"
echo -e "  • 后端日志:  ${LOG_DIR}/backend.log"
echo ""
echo -e "${GREEN}停止服务:${NC}"
echo -e "  • 运行: ${YELLOW}./stop.sh${NC} 或 ${YELLOW}kill ${BACKEND_PID}${NC}"
echo ""
echo -e "${GREEN}查看日志:${NC}"
echo -e "  • 实时查看: ${YELLOW}tail -f ${LOG_DIR}/backend.log${NC}"
echo ""
echo -e "${BLUE}========================================${NC}"

# 保存 PID 到文件
echo $BACKEND_PID > "${PROJECT_DIR}/.backend.pid"

exit 0

