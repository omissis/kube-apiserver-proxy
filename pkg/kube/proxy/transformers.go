package proxy

import (
	"encoding/json"

	"github.com/itchyny/gojq"
)

type ResponseBodyTransformer interface {
	Name() string
	Run([]byte, map[string]any) ([]byte, error)
}

func NewJqResponseBodyTransformer() *JqResponseBodyTransformer {
	return &JqResponseBodyTransformer{}
}

type JqResponseBodyTransformer struct{}

func (jq *JqResponseBodyTransformer) Name() string {
	return "jq"
}

func (jq *JqResponseBodyTransformer) Run(body []byte, opts map[string]any) ([]byte, error) {
	query, err := gojq.Parse(opts["src"].(string))
	if err != nil {
		return nil, err
	}

	data := map[string]any{}
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	vvv := make([]byte, 0)

	iter := query.Run(data)

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			return nil, err
		}

		vv, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		vvv = append(vvv, vv...)
	}

	return vvv, nil
}
