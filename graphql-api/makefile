dev-server:
	reflex -r '\.go$$' -s -- sh -c 'DEBUG=true LOCAL=true go run main.go'

gen:
	go generate ./...
doc:
	@-cat `find . -name '*.graphql'` > temp.gql
	@-graphql-markdown temp.gql > ../../docs/gql_schema.md
	@-rm temp.gql