package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql/api/_resolvers"
)

var handler http.Handler

func init() {
	schemaData, err := ioutil.ReadFile("schema.graphql")
	if err != nil {
		panic(err)
	}

	resolve := resolvers.New("https://www.speedrun.com/api/v1")

	schema, err := graphql.ParseSchema(string(schemaData), resolve,
		graphql.UseFieldResolvers())
	if err != nil {
		panic(err)
	}

	handler = &relay.Handler{Schema: schema}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r)
}
