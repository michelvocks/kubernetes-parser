package main

import (
	"fmt"
	"strconv"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

func (args scaleupRSArgs) execute() {
	// run command
	scaleUpReplicaSets(args.clientset, args.ns)
}

func scaleUpReplicaSets(clientset *kubernetes.Clientset, ns string) {
	// get repliction sets
	rs := clientset.ExtensionsV1beta1().ReplicaSets(ns)
	rsList, err := rs.List(v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// general debug output
	fmt.Printf("Replicas found: %d\n", len(rsList.Items))

	// iterate over all replication sets
	for _, rsObj := range rsList.Items {
		// only scale replication sets based on annotation setting
		if rsObj.Annotations == nil {
			continue
		}
		if _, ok := rsObj.Annotations[latestRS]; !ok {
			continue
		}

		// get last replication number from annotations
		lastRSNumber := rsObj.Annotations[latestRS]

		// general debug output
		fmt.Printf("Found latest replica number %s for %s\n", lastRSNumber, rsObj.Name)

		// scale replicas to given replica number from annotation
		lastRSNumberInt, err := strconv.Atoi(lastRSNumber)
		if err != nil {
			continue
		}
		lastRSNumberInt32 := int32(lastRSNumberInt)
		rsObj.Spec.Replicas = &lastRSNumberInt32

		// delete annotation
		delete(rsObj.Annotations, latestRS)

		// update replication controller
		rs.Update(&rsObj)
	}
}
