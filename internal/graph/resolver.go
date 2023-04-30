package graph

import "k8s.io/client-go/kubernetes"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Clientset *kubernetes.Clientset
}
