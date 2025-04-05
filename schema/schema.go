package schema

import (
	"encoding/json"
	"errors"
	"net/http"
	"john/models"
	"john/utils"

	"github.com/graphql-go/graphql"
)

var accountType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Account",
	Fields: graphql.Fields{
		"username": &graphql.Field{Type: graphql.String},
		"password": &graphql.Field{Type: graphql.String},
	},
})

var empType = graphql.NewObject(graphql.ObjectConfig{
	Name: "EmpDetails",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.String},
		"title": &graphql.Field{Type: graphql.String},
		"address": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				_, err := utils.ValidateToken(p.Context.Value("token").(string))
				if err != nil {
					return nil, err
				}
				return p.Source.(models.EmpDetails).Address, nil
			},
		},
	},
})

func GraphQLHandler() http.Handler {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"account": &graphql.Field{
				Type: accountType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					token := p.Context.Value("token").(string)
					if token == "" {
						return nil, errors.New("authorization required")
					}
					emp, err := utils.ValidateToken(token)
					if err != nil {
						return nil, err
					}
					return models.GetEmployeeByUsername(emp.Username)
				},
			},
			"EmployeeDetails": &graphql.Field{
				Type: graphql.NewList(empType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return models.GetAllEmpDetails()
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: rootQuery})
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "Missing query", http.StatusBadRequest)
			return
		}
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
			Context:       r.Context(),
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
}
