#!/bin/bash

# 获取当前脚本的目录
ROOT_DIR="$(dirname "${BASH_SOURCE[0]}")"

# 切换到脚本所在目录
cd "$ROOT_DIR" || exit 1

# 调用 serve.bash 脚本，并指定服务器名称为 "update"
# "$@" 会将所有传递给 serve_update.bash 的参数传递给 serve.bash
echo "[ 正在通过 serve_update.bash 调用 serve.bash，服务器名称: update ]"
source "./serve.bash" "update" "$@"