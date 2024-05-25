cd rpc-service  # 进入当前目录下的 rpc_service 文件夹

for file in *.go; do
    # 使用 go run 命令运行每个文件
    go run "$file"
done
