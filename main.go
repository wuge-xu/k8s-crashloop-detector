package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Println("构建配置失败:", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("创建客户端失败:", err)
		os.Exit(1)
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("获取 Pod 列表失败:", err)
		os.Exit(1)
	}

	fmt.Printf("集群中共有 %d 个 Pod，开始检查异常状态...\n\n", len(pods.Items))

	problemCount := 0

	for _, pod := range pods.Items {
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.State.Waiting != nil && cs.State.Waiting.Reason == "CrashLoopBackOff" {
				problemCount++
				printCrashAlert(pod, cs)
			}
		}
	}

	fmt.Printf("\n检查完成，共发现 %d 个 CrashLoopBackOff 异常。\n", problemCount)
}

func printCrashAlert(pod corev1.Pod, cs corev1.ContainerStatus) {
	fmt.Println("[WARN]")
	fmt.Printf("namespace: %s\n", pod.Namespace)
	fmt.Printf("pod:       %s\n", pod.Name)
	fmt.Printf("container: %s\n", cs.Name)
	fmt.Printf("status:    CrashLoopBackOff\n")
	fmt.Printf("restarts:  %d\n", cs.RestartCount)
	fmt.Println()
}
