package main

import (
	"fmt"
	"strconv"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

func (args scaledownRSArgs) execute() {
	// run command
	scaleDownReplicaSets(args.clientset, args.ns)
}

func scaleDownReplicaSets(clientset *kubernetes.Clientset, ns string) {
	// get repliction sets
	rs := clientset.ExtensionsV1beta1().ReplicaSets(ns)
	rsList, err := rs.List(v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// general debug output
	fmt.Printf("Replicas found: %d\n", len(rsList.Items))

	// iterate over all replication controllers
	for _, rsObj := range rsList.Items {
		// get actual replica number
		latestGivenRS := rsObj.Spec.Replicas

		// don't update replica sets which are already zero
		if *latestGivenRS == 0 {
			continue
		}

		// general debug output
		fmt.Printf("Found %d replicas for %s\n", *latestGivenRS, rsObj.Name)

		// create map if not exists
		if rsObj.Annotations == nil {
			rsObj.Annotations = make(map[string]string)
		}

		// set actual replica number as annotation
		rsObj.Annotations[latestRS] = strconv.Itoa(int(*latestGivenRS))

		// scale replicas to zero
		var zeroReplica int32
		rsObj.Spec.Replicas = &zeroReplica

		// update replication controller
		rs.Update(&rsObj)
	}
}
