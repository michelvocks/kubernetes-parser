package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const createdField string = "created"
const expiresField string = "expires"
const createdLayout string = "20060102150405"

func main() {
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
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

	// get namespaces
	ns := clientset.CoreV1Client.Namespaces()
	nsList, err := ns.List(v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// iterate over all namespaces and get the annotatios
	for _, nsObj := range nsList.Items {
		// get annotations from namespace
		nsAnno := nsObj.Annotations

		// iterate annotations
		var createdTime time.Time
		for id, anno := range nsAnno {
			// did we found the created and expires tag?
			switch id {
			case createdField:
				t, err := time.ParseInLocation(createdLayout, anno, time.Local)
				if err != nil {
					fmt.Println(err)
				}
				createdTime = t
				fmt.Println("Got time:", createdTime.String())
				fmt.Println("Time now:", time.Now().Local().String())
			case expiresField:
				expiredTime := calculateExpireDate(createdTime, anno)
				if time.Now().Local().After(expiredTime) {
					fmt.Println("The following namespace has been expired:", nsObj.ObjectMeta.Name)
					fmt.Println("Expiry:", expiredTime.String())
					fmt.Println("Current:", time.Now().String())
				}
			}
		}
	}
}

func calculateExpireDate(t time.Time, addTime string) time.Time {
	// get the value of time (e.g. 12)
	i, err := strconv.Atoi(addTime[0 : len(addTime)-1])
	if err != nil {
		fmt.Println(err)
	}

	// get the time type (e.g. d for days)
	givenType := strings.ToLower(addTime[len(addTime)-1:])

	// find out the type
	var timeType time.Duration
	switch givenType {
	case "s":
		timeType = time.Second
	case "m":
		timeType = time.Minute
	case "h":
		timeType = time.Hour
	case "d":
		return t.AddDate(0, 0, i)
	case "w":
		return t.AddDate(0, 0, i*7)
	default:
		return t
	}

	// calculate new expired time
	t = t.Add(timeType * time.Duration(i))
	return t
}
