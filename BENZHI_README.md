# BENZHI_README

## 项目说明

- 项目：zhanglei10281852-gif/gogo-58
- 项目用途：CaveLoop is an offline command line tool for speleological survey data reduction and traverse network adjustment. It takes raw cave survey field notes (tape distance, compass azimuth, clinometer inclination, optional backsights), reduces them into three dimensional station coordinates, analyses the passage network, finds the independent loop basis, measures how well each loop closes, distributes the closure error and flags the gross errors that typically spoil a survey.
- Go 工具链：`golang:1.22`
- 前端工具链：无

## 标准构建、运行和测试命令

进入容器后执行：

```bash
# 编译
cd '/app' && GOTOOLCHAIN=local go build ./...

# 启动
# 该项目为 Go 库或未提供可识别的可执行入口，无独立启动命令

# 测试
cd '/app' && GOTOOLCHAIN=local go test ./...
```

## Docker 构建和进入容器

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh benzhi-task-58-amd64 linux/amd64
./build_benzhi_docker.sh benzhi-task-58-arm64 linux/arm64
docker run -it benzhi-task-58-amd64:latest
docker run -it --platform linux/arm64 benzhi-task-58-arm64:latest
```

## 题目验证命令

1. 预期退出码 0：`go test ./internal/store -run "^TestVerifyAuditRejectsSplicedChain$" -count=1 -v`
2. 预期退出码 0：`go test -buildvcs=false -count=1 ./...`

## Bug 复现

Bug 现象、触发步骤和完整错误信息见 `BUG_REPRO.md`。
