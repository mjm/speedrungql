package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql/api/_resolvers"
)

func main() {
	schemaData, err := ioutil.ReadFile("api/schema.graphql")
	if err != nil {
		panic(err)
	}

	resolve := resolvers.New("https://www.speedrun.com/api/v1")

	schema, err := graphql.ParseSchema(string(schemaData), resolve,
		graphql.UseFieldResolvers())
	if err != nil {
		panic(err)
	}

	handler := &relay.Handler{Schema: schema}
	http.Handle("/graphql", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
