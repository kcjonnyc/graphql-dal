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
	Variant     *Variant   `json:"variant"`
	TopCategory int        `json:"topCategory"`
}

var data map[string]*Product

var variantType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "Variant",
		Fields: graphql.InputObjectConfigFieldMap{
			"status": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"upc": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"externalId": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"images": &graphql.InputObjectFieldConfig{
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
			"variant": &graphql.Field{
				Type: variantType,
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

/*var variantMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "variantMutation",
		Fields: graphql.Fields{
			"variant": &graphql.Field{
				Type: variantType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"status": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"upc": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"externalId": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"images": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if ok {
						product := data[id]
						status, ok := p.Args["status"].(int)
						if ok {
							product.Variant.Status = status
						}
						upc, ok := p.Args["upc"].(string)
						if ok {
							product.Variant.Upc = upc
						}
						externalId, ok := p.Args["externalId"].(string)
						if ok {
							product.Variant.ExternalId = externalId
						}
						images, ok := p.Args["images"].([]string)
						if ok {
							product.Variant.Images = images
						}
						return product.Variant, nil
					}
					return nil, nil
				},
			},
		},
	})*/

var productMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "productMutation",
		Fields: graphql.Fields{
			"product": &graphql.Field{
				Type: productType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"variant": &graphql.ArgumentConfig{
						Type: variantType,
					},
					"topCategory": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if ok {
						product := data[id]
						//variantMap := p.Args["variant"].(map[string]{})
						//status, ok := variantMap["stauts"].(int)
						log.Printf("%v", p.Args)
						/*if ok {
							log.Printf("test")
							product.Variant.Status = status
						}*/
						topCategory, ok := p.Args["topCategory"].(int)
						if ok {
							product.TopCategory = topCategory
						}
						return product, nil
					}
					return nil, nil
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
		Mutation: productMutation,
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
