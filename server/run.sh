set -e

go fmt
go build

./graphql_json_go -settings settings.development.json
