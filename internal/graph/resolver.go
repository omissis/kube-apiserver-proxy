package graph

import "github.com/omissis/kube-apiserver-proxy/internal/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TodoStore []*model.Todo
}
