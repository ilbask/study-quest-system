#!/bin/bash

# Study Quest 系统停止脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
PID_FILE="${PROJECT_DIR}/.backend.pid"
BACKEND_PORT=8080

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   Study Quest 系统停止脚本${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 方法1: 从 PID 文件读取
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p $PID > /dev/null 2>&1; then
        echo -e "${YELLOW}正在停止服务 (PID: ${PID})...${NC}"
        kill $PID 2>/dev/null || true
        sleep 2
        
        # 检查是否还在运行
        if ps -p $PID > /dev/null 2>&1; then
            echo -e "${RED}优雅停止失败，强制终止...${NC}"
            kill -9 $PID 2>/dev/null || true
        fi
        
        rm -f "$PID_FILE"
        echo -e "${GREEN}✓ 服务已停止${NC}"
    else
        echo -e "${YELLOW}⚠ PID 文件中的进程已不存在${NC}"
        rm -f "$PID_FILE"
    fi
fi

# 方法2: 通过端口查找
PORT_PID=$(lsof -ti:${BACKEND_PORT} 2>/dev/null || true)
if [ -n "$PORT_PID" ]; then
    echo -e "${YELLOW}发现端口 ${BACKEND_PORT} 被进程 ${PORT_PID} 占用${NC}"
    echo -e "${YELLOW}正在停止...${NC}"
    
    kill $PORT_PID 2>/dev/null || true
    sleep 2
    
    if lsof -ti:${BACKEND_PORT} >/dev/null 2>&1; then
        echo -e "${RED}优雅停止失败，强制终止...${NC}"
        kill -9 $PORT_PID 2>/dev/null || true
    fi
    
    echo -e "${GREEN}✓ 端口 ${BACKEND_PORT} 已释放${NC}"
else
    echo -e "${GREEN}✓ 端口 ${BACKEND_PORT} 空闲${NC}"
fi

# 清理编译文件（可选）
if [ -f "${PROJECT_DIR}/study-quest-server" ]; then
    echo -e "${YELLOW}清理编译文件...${NC}"
    rm -f "${PROJECT_DIR}/study-quest-server"
    echo -e "${GREEN}✓ 清理完成${NC}"
fi

echo ""
echo -e "${GREEN}所有服务已停止${NC}"
echo ""

