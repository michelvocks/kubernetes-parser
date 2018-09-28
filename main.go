package main

import (
	"flag"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type command interface {
	execute()
}

type expiredNSArgs struct {
	clientset   *kubernetes.Clientset
	expiredTime time.Duration
}

type scaledownRSArgs struct {
	ns        string
	clientset *kubernetes.Clientset
}

type scaleupRSArgs struct {
	ns        string
	clientset *kubernetes.Clientset
}

type noneNsArgs struct {
	clientset *kubernetes.Clientset
}

const latestRS string = "latestRS"
const expiresField string = "expires"
const userIdField string = "userId"

func main() {
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	command := flag.String("cmd", "", "command to execute: expiredns, scaledownrs, scaleuprs, nonens")
	namespace := flag.String("namespace", "", "namespace for replication controller scale down / up")
	expiretime := flag.Int("expiretime", 0, "minutes which will be added to the current time to calculate expiration")

	flag.Parse()
	if *kubeconfig == "" {
		panic("-kubeconfig not specified")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Setup command pattern
	expiredNS := expiredNSArgs{clientset: clientset, expiredTime: time.Duration(*expiretime)}
	scaledownRS := scaledownRSArgs{ns: *namespace, clientset: clientset}
	scaleupRS := scaleupRSArgs{ns: *namespace, clientset: clientset}
	noneNS := noneNsArgs{clientset: clientset}

	switch strings.ToLower(*command) {
	case "expiredns":
		expiredNS.execute()
	case "scaledownrs":
		scaledownRS.execute()
	case "scaleuprs":
		scaleupRS.execute()
	case "nonens":
		noneNS.execute()
	}
}
