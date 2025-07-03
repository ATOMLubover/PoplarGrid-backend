#!/bin/bash

# 获取当前脚本的目录
ROOT_DIR="$(dirname "${BASH_SOURCE[0]}")"

# 切换到脚本所在目录
cd "$ROOT_DIR" || exit 1

# 调用 serve.bash 脚本，并指定服务器名称为 "gateway"
# "$@" 会将所有传递给 serve_gateway.bash 的参数传递给 serve.bash
echo "[ 正在通过 serve_gateway.bash 调用 serve.bash，服务器名称: gateway ]"
source "./serve.bash" "gateway" "$@"


# #!/bin/bash

# # 获取当前脚本的目录（与执行路径无关）
# ROOT_DIR="$(\
#     (cd "$(dirname "${BASH_SOURCE[0]}")" \
#         &> /dev/null) \
#     && pwd)"
# # 确保在脚本所在目录执行
# cd "$ROOT_DIR" || exit 1

# # 变量宏定义
# BIN_DIR="$ROOT_DIR/bin" # 二进制文件输出目录
# GO_SRC_DIR="$ROOT_DIR/cmd/gateway" # Go 源码目录
# GO_OUT="$BIN_DIR/gateway" # 编译输出的二进制文件路径
# SWAG_OUTPUT_DIR="$ROOT_DIR/docs/gateway" # 定义 Swagger 文档的输出目录

# # 先确保输出目录存在
# mkdir -p "$BIN_DIR" || exit 1

# # 解析脚本参数以确定操作类型
# if (($# > 1)) \
#     || { (( $# == 1)) \
#         && [[ "$1" != "build" && "$1" != "run" ]]; }; then
#     echo "[ 参数错误，格式为: $0 [build|run] ]"
#     exit 1
# fi

# # 检查操作参数（默认为 ALL）
# ACTION="all"
# case "$1" in
#     build) ACTION="build" ;;
#     run)   ACTION="run" ;;
#     *)     echo "[ 无参数，已启用默认操作 ]" ;;
# esac

# # 检查 Go 环境
# if ! command -v go &> /dev/null; then
#     echo "[ 错误：无法使用 Go 命令，未检测到 Go 环境 ]"
#     exit 1
# fi

# # 检查 swag 命令
# if ! command -v swag &> /dev/null; then
#     echo "[ 错误：无法使用 swag 命令 ]"
#     exit 1
# fi

# # 定义 Swagger 初始化函数
# run_swag_init() {
#     # 确定 Swagger 的输出目录存在
#     mkdir -p "$SWAG_OUTPUT_DIR" || exit 1

#     # 确保在 Go 源码目录执行，保证 swag 命令成功
#     cd "$GO_SRC_DIR" || exit 1

#     echo "[ 开始生成 Swagger API 文档... ]"

#     if ! swag init \
#         -o "$SWAG_OUTPUT_DIR" \
#         -g "$GO_SRC/main.go"; then
#         echo "[ Swagger API 文档生成失败 ]"
#         exit 1
#     fi

#     echo "[ Swagger API 文档生成成功，输出路径：$SWAG_OUTPUT_DIR ]"

#     # 返回到脚本根目录
#     cd "$ROOT_DIR" || exit 1
# }

# # 定义编译函数
# compile_gateway() {
#     # 确保在 Go 源码目录执行编译
#     cd "$GO_SRC_DIR" || exit 1

#     echo "[ 开始编译 Gateway Server... ]"

#     if ! go build \
#         -o "$GO_OUT"; then
#         echo "[ 编译失败 ]"
#         exit 1
#     fi

#     echo "[ 编译成功，二进制文件输出路径：$GO_OUT ]"

#     # 返回到脚本根目录
#     cd "$ROOT_DIR" || exit 1
# }

# # 定义运行函数
# run_gateway() {
#     if [[ ! -f "$GO_OUT" ]]; then
#         echo "[ 未找到可执行文件，正在先行编译... ]"
#         compile_gateway || exit 1
#     fi

#     # 进入 GO_OUT 所在目录
#     cd "$(dirname "$GO_OUT")" || exit 1

#     echo "[ 开始启动 $GO_OUT ]"

#     if ! "$GO_OUT"; then
#         echo "[ 启动失败 ]"
#         exit 1
#     fi
# }

# # 最后根据操作类型执行对应命令
# case $ACTION in
#     build)
#         run_swag_init
#         compile_gateway
#         ;;

#     run)
#         run_gateway
#         ;;

#     all)
#         run_swag_init
#         compile_gateway
#         run_gateway
#         ;;
# esac