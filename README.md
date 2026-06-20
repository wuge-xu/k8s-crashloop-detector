# K8s CrashLoopBackOff Detector

一个用 Go + client-go 编写的命令行工具，自动扫描 Kubernetes 集群中处于 `CrashLoopBackOff` 状态的容器，并打印告警信息。

## 功能

- 遍历所有命名空间下的所有 Pod
- 检查每个 Pod 内每个容器的状态
- 发现 `CrashLoopBackOff` 状态时，打印告警（命名空间、Pod 名、容器名、重启次数）
- 统计本次检查共发现的异常数量

## 技术栈

- Go
- [client-go](https://github.com/kubernetes/client-go)：Kubernetes 官方 Go SDK

## 运行方式

确保本地已配置好可用的 kubeconfig（默认读取 `~/.kube/config`），然后：

\`\`\`bash
go mod tidy
go run main.go
\`\`\`

## 示例输出

\`\`\`
集群中共有 8 个 Pod，开始检查异常状态...

[WARN]
namespace: default
pod:       crash-test
container: crash-test
status:    CrashLoopBackOff
restarts:  5

检查完成，共发现 1 个 CrashLoopBackOff 异常。
\`\`\`

## 已知细节

`CrashLoopBackOff` 是容器在等待重试期间的瞬时状态（容器会在 Running → Terminated → Waiting(CrashLoopBackOff) → Running 之间循环），单次轮询可能会错过该状态窗口。本工具适合配合定时任务（如 cron 或 K8s CronJob）周期性运行，而不是只跑一次。该问题在开发过程中通过本地构造一个持续崩溃的测试 Pod（`busybox` + 非零退出码）进行过实际验证。

## 开发与测试环境

本项目在 WSL2 + K3s 单节点集群上开发和验证。
