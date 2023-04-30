package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
)

type PreloadDepth int

const (
	PreloadAll = -1
)

func GetPreloads(ctx context.Context, maxDepth PreloadDepth) []string {
	return GetNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
		maxDepth,
		0,
	)
}

func GetNestedPreloads(
	ctx *graphql.OperationContext,
	fields []graphql.CollectedField,
	prefix string,
	maxDepth PreloadDepth,
	depth PreloadDepth,
) (preloads []string) {
	if maxDepth != PreloadAll && depth > maxDepth {
		return
	}

	for _, column := range fields {
		prefixColumn := GetPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(
			preloads,
			GetNestedPreloads(
				ctx,
				graphql.CollectFields(ctx, column.Selections, nil),
				prefixColumn,
				maxDepth,
				depth+1,
			)...,
		)
	}
	return
}

func GetPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}
