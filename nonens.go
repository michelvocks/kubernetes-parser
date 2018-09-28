package main

import (
	"encoding/json"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
)

type nameSpaceNone struct {
	Name         string
	UserID       string
	CreationDate unversioned.Time
}

func (args noneNsArgs) execute() {
	// run command
	getNoneNS(args.clientset)
}

func getNoneNS(clientset *kubernetes.Clientset) {
	// get namespaces
	ns := clientset.CoreV1Client.Namespaces()
	nsList, err := ns.List(v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// iterate over all namespaces and get the annotatios
	noneNS := []*nameSpaceNone{}
	for _, nsObj := range nsList.Items {
		// get annotations from namespace
		nsAnno := nsObj.Annotations

		// iterate annotations
		var expiresFieldValue string
		var userID string
		var noneNamespace = nameSpaceNone{}
		for id, anno := range nsAnno {
			// did we found the expires tag?
			switch id {
			case expiresField:
				expiresFieldValue = anno
			case userIdField:
				userID = anno
			}
		}

		if expiresFieldValue == "none" && userID != "" {
			noneNamespace.Name = nsObj.ObjectMeta.Name
			noneNamespace.UserID = userID
			noneNamespace.CreationDate = nsObj.CreationTimestamp
			noneNS = append(noneNS, &noneNamespace)
		}
	}

	// print out all none namespaces as json
	b, err := json.Marshal(noneNS)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
