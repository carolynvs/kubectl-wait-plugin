package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/pluginutils"

	_ "github.com/golang/glog"
)

func init() {
	// Initialize glog flags
	flag.CommandLine.Set("logtostderr", "true")
	flag.CommandLine.Set("v", os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_V"))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: kubectl plugin wait RESOURCE/NAME [--timeout 5m] [--interval 2s]")
		os.Exit(1)
	}

	resource := os.Args[1]
	resourceParts := strings.Split(resource, "/")
	if len(resourceParts) != 2 {
		fmt.Println("invalid resource specifier: expected RESOURCE/NAME, such as pod/mypod or deploy/mydeploy")
		os.Exit(1)
	}
	resourceType := resourceParts[0]
	resourceName := resourceParts[1]

	var timeout time.Duration
	timeoutFlag := os.Getenv("KUBECTL_PLUGINS_LOCAL_FLAG_TIMEOUT")
	if timeoutFlag == "-1" {
		notimeout := time.Duration(math.MaxInt64)
		timeout = notimeout
	} else {
		t, err := time.ParseDuration(timeoutFlag)
		if err != nil {
			fmt.Printf("invalid timeout: %s\n", err)
			os.Exit(1)
		}
		timeout = t
	}

	var interval time.Duration
	intervalFlag := os.Getenv("KUBECTL_PLUGINS_LOCAL_FLAG_INTERVAL")
	i, err := time.ParseDuration(intervalFlag)
	if err != nil {
		fmt.Printf("invalid timeout: %s\n", err)
		os.Exit(1)
	}
	interval = i

	waitOrDie(resourceType, resourceName, timeout, interval)
}

func waitOrDie(resourceType, resourceName string, timeout, interval time.Duration) {
	client, ns := loadConfig()

	switch resourceType {
	case "po", "pod":
		err := wait.PollImmediate(interval, timeout,
			func() (bool, error) {
				p, err := getPod(client, ns, resourceName)
				if err != nil {
					return false, err
				}
				return isPodReady(p), nil

			})
		if err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(-1)
		}
	case "deploy", "deployment":
		err := wait.PollImmediate(interval, timeout,
			func() (bool, error) {
				d, err := getDeployment(client, ns, resourceName)
				if err != nil {
					return false, err
				}
				return isDeploymentReady(d), nil

			})
		if err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(-1)
		}
	case "svc", "service":
		err := wait.PollImmediate(interval, timeout,
			func() (bool, error) {
				s, err := getService(client, ns, resourceName)
				if err != nil {
					return false, err
				}
				return isServiceReady(s), nil
			})
		if err != nil {
			fmt.Printf("error: %v", err)
			os.Exit(-1)
		}
	default:
		fmt.Printf("unsupported resource type: %s\n", resourceType)
		os.Exit(1)
	}
}

func loadConfig() (*kubernetes.Clientset, string) {
	restConfig, kubeConfig, err := pluginutils.InitClientAndConfig()
	if err != nil {
		panic(err)
	}
	c := kubernetes.NewForConfigOrDie(restConfig)
	ns, _, _ := kubeConfig.Namespace()
	return c, ns
}
