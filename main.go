package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/graphql-go/graphql"
	"encoding/json"
	"golang.org/x/net/context"
	"github.com/pressly/chi"
)


type Variant struct {
	Status int `json:"status"`
	Upc int `json:"upc"`
	ExternalId string `json:"externalId"`
	Images []string `json:"images"`
}

type Product struct {
	Id int `json:"id"`
	Variants []*Variant `json:"variants"`
	TopCategory int `json:"topCategory"`
}

var data map[int]*Product
// GraphQL

var variantType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Variant",
		Fields: graphql.Fields{
			"status": &graphql.Field{
				Type: graphql.Int,
			},
			"upc": &graphql.Field{
				Type: graphql.Int,
			},
			"externalId": &graphql.Field{
				Type: graphql.String,
			},
			"images:" &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
})

var productType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"variants": &graphql.Field{
				Type: graphql.NewList(variantType),
			},
			"topCategory": &graphql.Field{
				Type: graphql.Int,
			},
		},
})

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"product": &graphql.Field{
				Type: productType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(product graphql.ResolveParams) (interface{}, error) {
					id, ok := product.Args["id"].(int)
					if ok {
						return data[id], nil
					}
					return nil, nil
				},
		},
	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
		//Mutation: mutationType,
	},
)

func serveGraphQL (ctx context.Context, w http.ResponseWriter,r *http.Request) {
	// execute GraphQL query
	params := graphql.Params {
			Schema: schema,
			Context: ctx,
			RequestString: r.URL.Query().Get("query"),
	}
	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		fmt.Printf("errors : %v", result.Errors)
	}
	json.NewEncoder(w).Encode(result)
}


func main() {
	err := openJsonFile("example_product.json", &data)
	if err != nil {
		return
	}
	r := chi.NewRouter()
	r.Handle("/graphql", serveGraphQL)
	http.ListenAndServe(":3000", r)
}

// helper
func openJsonFile(name string, p interface{}) (err error) {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, p)
	if err != nil {
		return err
	}
	return err
}
