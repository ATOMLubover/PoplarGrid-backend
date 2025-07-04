#!/bin/bash

# 获取当前脚本的目录（与执行路径无关）
ROOT_DIR="$(\
    (cd "$(dirname "${BASH_SOURCE[0]}")" \
        &> /dev/null) \
    && pwd)"
# 确保在脚本所在目录执行
cd "$ROOT_DIR" || exit 1

# 定义源代码路径和输出二进制文件路径
GO_SRC_DIR="$ROOT_DIR/cmd/update"
OUTPUT_DIR="$ROOT_DIR/bin"
GO_OUT="$OUTPUT_DIR/update" # 可执行文件的名称

# 确保输出目录存在
mkdir -p "$OUTPUT_DIR"

echo "[ 正在编译 $SOURCE_FILE... ]"

# 为了保证脚本的可移植性，进入源代码目录进行编译
cd "$GO_SRC_DIR" || {
    echo "[ 错误：无法切换到源代码目录 $GO_SRC_DIR ]"
    exit 1
}

if ! go build -o "$GO_OUT"; then
    echo "[ 编译失败 ]"
    exit 1
fi

# 检查编译是否成功
if [ $? -eq 0 ]; then
  echo "[ 编译成功，二进制文件输出路径：$GO_OUT ]"
  echo "[ 正在启动 $GO_OUT... ]"
  # 启动编译后的二进制文件
  if ! "$GO_OUT"; then
    echo "[ 启动失败 ]"
    exit 1
  fi
else
  echo "[ 编译失败 ]"
  exit 1
fi