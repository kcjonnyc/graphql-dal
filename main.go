package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"io/ioutil"
	"log"
	"net/http"
)

type Variant struct {
	Status     int      `json:"status"`
	Upc        string   `json:"upc"`
	ExternalId string   `json:"externalId"`
	Images     []string `json:"images"`
}

type Product struct {
	Id          int        `json:"id"`
	Variants    []*Variant `json:"variants"`
	TopCategory int        `json:"topCategory"`
}

var data map[string]*Product

var variantType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Variant",
		Fields: graphql.Fields{
			"status": &graphql.Field{
				Type: graphql.Int,
			},
			"upc": &graphql.Field{
				Type: graphql.String,
			},
			"externalId": &graphql.Field{
				Type: graphql.String,
			},
			"images": &graphql.Field{
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
						Type: graphql.String,
					},
				},
				Resolve: func(product graphql.ResolveParams) (interface{}, error) {
					id, ok := product.Args["id"].(string)
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

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func openJsonFile(name string, p interface{}) (err error) {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, p)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := openJsonFile("example_product.json", &data)
	if err != nil {
		log.Fatalf("Error loading data")
	}
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	fmt.Println("Now server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
