#!/bin/bash

# 获取当前脚本的目录（与执行路径无关）
ROOT_DIR="$(\
    (cd "$(dirname "${BASH_SOURCE[0]}")" \
    &> /dev/null) \
&& pwd)"
# 确保在脚本所在目录执行
cd "$ROOT_DIR" || exit 1

# 解析脚本参数以确定操作类型和服务器名称
# 调用形式为 source serve.bash "server_name" [run|build]
if (($# < 1)) || (($# > 2)); then
    echo "[ 参数错误，格式为: source $0 \"server_name\" [build|run] ]"
    echo "[ 例如: source $0 \"gateway\" build ]"
    exit 1
fi

SERVER_NAME="$1" # 服务器的文件夹名和最终可执行文件的名字

# 变量宏定义
BIN_DIR="$ROOT_DIR/bin"                             # 二进制文件输出目录
GO_MAIN_DIR="$ROOT_DIR/cmd/${SERVER_NAME}server"    # main.go 所在目录
GO_OUT="$BIN_DIR/$SERVER_NAME"                      # 编译输出的二进制文件路径
SWAG_OUTPUT_DIR="$ROOT_DIR/docs/$SERVER_NAME"       # 定义 Swagger 文档的输出目录
CONFIG_FILE_NAME="${SERVER_NAME}_config.yaml"       # 配置文件名
CONFIG_SOURCE_PATH="$GO_MAIN_DIR/$CONFIG_FILE_NAME" # 配置文件源路径
CONFIG_DEST_PATH="$BIN_DIR/$CONFIG_FILE_NAME"       # 配置文件目标路径 (软链接)

# 先确保输出目录存在
mkdir -p "$BIN_DIR" || exit 1

# 检查操作参数（默认为 ALL）
ACTION="all"
if (($# == 2)); then
    case "$2" in
        build) ACTION="build" ;;
        run)   ACTION="run" ;;
        *)     echo "[ 参数错误，第二个参数必须是 'build' 或 'run' ]"
            exit 1
        ;;
    esac
else
    echo "[ 未指定操作参数，已启用默认操作: all (生成 Swagger, 编译, 运行) ]"
fi


# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "[ 错误：无法使用 Go 命令，未检测到 Go 环境 ]"
    exit 1
fi

# 检查 swag 命令
if ! command -v swag &> /dev/null; then
    echo "[ 错误：无法使用 swag 命令，请确保已安装并配置 PATH ]"
    exit 1
fi

# 定义 Swagger 初始化函数
run_swag_init() {
    # 确定 Swagger 的输出目录存在
    mkdir -p "$SWAG_OUTPUT_DIR" || exit 1
    
    if [[ ! -d "$GO_MAIN_DIR" ]]; then
        echo "[ 错误：Go 源码目录 '$GO_MAIN_DIR' 不存在。请检查 server_name 是否正确。 ]"
        exit 1
    fi
    
    # 确保返回到脚本根目录，以扫描到所有包
    cd "$ROOT_DIR" || exit 1
    
    echo "[ 开始生成 Swagger API 文档... ]"
    
    if ! swag init \
    -o "$SWAG_OUTPUT_DIR" \
    --dir "$GO_MAIN_DIR","$ROOT_DIR/internal/${SERVER_NAME}server/handlers","$ROOT_DIR/internal/${SERVER_NAME}server/dtos"; then
        echo "[ Swagger API 文档生成失败 ]"
        exit 1
    fi
    
    echo "[ Swagger API 文档生成成功，输出路径：$SWAG_OUTPUT_DIR ]"
}

# 定义编译函数
compile_server() {
    # 确保在 Go 源码目录执行编译
    if [[ ! -d "$GO_MAIN_DIR" ]]; then
        echo "[ 错误：Go 源码目录 '$GO_MAIN_DIR' 不存在。请检查 server_name 是否正确。 ]"
        exit 1
    fi
    cd "$GO_MAIN_DIR" || exit 1
    
    echo "[ 开始编译 $SERVER_NAME Server... ]"
    
    if ! go build \
    -o "$GO_OUT"; then
        echo "[ 编译失败 ]"
        exit 1
    fi
    
    echo "[ 编译成功，二进制文件输出路径：$GO_OUT ]"

    if [[ ! -f "$CONFIG_SOURCE_PATH" ]]; then
        echo "[ 警告：未找到配置文件 '$CONFIG_SOURCE_PATH'，跳过创建软链接。 ]"
    else
        echo "[ 准备将配置文件 '$CONFIG_FILE_NAME' 以软链接形式复制到 '$BIN_DIR' ]"
        # 移除旧的软链接或文件（如果存在）
        if [[ -L "$CONFIG_DEST_PATH" ]]; then
            rm "$CONFIG_DEST_PATH"
        elif [[ -f "$CONFIG_DEST_PATH" ]]; then
            echo "[ 警告：目标目录 '$BIN_DIR' 中存在同名文件 '$CONFIG_FILE_NAME'，将被覆盖。 ]"
            rm "$CONFIG_DEST_PATH"
        fi

        # 创建新的软链接
        if ! ln -s "$CONFIG_SOURCE_PATH" "$CONFIG_DEST_PATH"; then
            echo "[ 错误：创建配置文件软链接失败。 ]"
            exit 1
        fi
        echo "[ 配置文件软链接创建成功：$CONFIG_DEST_PATH -> $CONFIG_SOURCE_PATH ]"
    fi
    
    # 返回到脚本根目录
    cd "$ROOT_DIR" || exit 1
}

# 定义运行函数
run_server() {
    if [[ ! -f "$GO_OUT" ]]; then
        echo "[ 未找到可执行文件 '$GO_OUT'，正在先行编译... ]"
        compile_server || exit 1
    fi
    
    # 进入 GO_OUT 所在目录
    cd "$(dirname "$GO_OUT")" || exit 1
    
    echo "[ 开始启动 $GO_OUT ]"
    
    # 运行二进制文件
    if ! "$GO_OUT"; then
        echo "[ 启动失败 ]"
        exit 1
    fi
}

# 最后根据操作类型执行对应命令
case $ACTION in
    build)
        run_swag_init
        compile_server
    ;;
    
    run)
        run_server
    ;;
    
    all)
        run_swag_init
        compile_server
        run_server
    ;;
esac
