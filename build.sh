#!/usr/bin/env bash

# set -x

# 获取源码最近一次 git commit log，包含 commit sha 值，以及 commit message
GitCommitLog=$(git log --pretty=oneline -n 1)
# 将 log 原始字符串中的单引号替换成双引号
GitCommitLog=${GitCommitLog//\'/\"}
# 检查源码在git commit 基础上，是否有本地修改，且未提交的内容
GitStatus=$(git status -s)
# 查看当前的tag
GitTag=$(git describe --tag)
# 获取当前时间
BuildTime=$(date +'%Y.%m.%d.%H%M%S')
# 获取 Go 的版本
BuildGoVersion=$(go version)

# 将以上变量序列化至 LDFlags 变量中
# -s 去掉符号表,panic时候的stack trace就没有任何文件名/行号信息了，这个等价于普通C/C++程序被strip的效果
# -w 去掉DWARF调试信息，得到的程序就不能用gdb调试了。
LDFlags=" \
    -X 'github.com/q191201771/naza/pkg/bininfo.GitCommitLog=${GitCommitLog}' \
    -X 'github.com/q191201771/naza/pkg/bininfo.GitTag=${GitTag}' \
    -X 'github.com/q191201771/naza/pkg/bininfo.GitStatus=${GitStatus}' \
    -X 'github.com/q191201771/naza/pkg/bininfo.BuildTime=${BuildTime}' \
    -X 'github.com/q191201771/naza/pkg/bininfo.BuildGoVersion=${BuildGoVersion}' \
    -w -s \
"

ROOT_DIR=$(pwd)
BIN_DIR=${ROOT_DIR}/cmd/bin

# 如果可执行程序输出目录不存在，则创建
if [ ! -d ${BIN_DIR} ]; then
  cd ${ROOT_DIR} && mkdir -p cmd/bin
fi

# 编译多个可执行程序
# 根据cmd目录下的子程序列表来构建程序
# cd ${ROOT_DIR}/demo/add_blog_license && go build -ldflags "$LDFlags" -o ${ROOT_DIR}/bin/add_blog_license &&
#   cd ${ROOT_DIR}/demo/add_go_license && go build -ldflags "$LDFlags" -o ${ROOT_DIR}/bin/add_go_license &&
#   ls -lrt ${ROOT_DIR}/bin &&
#   cd ${ROOT_DIR} && ./bin/myapp -v &&
#   echo 'build done.'

# 获取构建程序列表
apps=$(cd ${ROOT_DIR}/cmd && find . -maxdepth 1 -type d ! -path ./bin -a ! -path . | cut -d"/" -f2)

for app in ${apps}; do
  echo "START BUILD APP: ${app}"
  cd ${ROOT_DIR}/cmd/${app}
  go build -ldflags "$LDFlags" -o ${BIN_DIR}/${app}
  echo "${app} BUILD DONE"
  echo "PRINT INFO  APP: ${app}"
  echo "-----------------INFO---------------------"
  cd ${BIN_DIR}
  # 压缩二进制文件
  upx ./${app}
  ./${app} -v
  echo "------------------------------------------"
  echo "******************************************"
done
echo -e "ALL APPS BUILD DONE\n"
