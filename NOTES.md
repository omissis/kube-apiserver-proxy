# Notes

## Commands to generate the graphql schema files starting from the kubernetes swagger.json

node_modules/.bin/openapi-to-graphql https://raw.githubusercontent.com/kubernetes/kubernetes/v1.27.1/api/openapi-spec/swagger.json --save internal/graph/swagger.graphqls
