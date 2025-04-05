package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"john/models"
	"john/utils"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

var accountType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Account",
	Fields: graphql.Fields{
		"id":       &graphql.Field{Type: graphql.Int},
		"username": &graphql.Field{Type: graphql.String},
		"password": &graphql.Field{Type: graphql.String},
	},
})

var empType = graphql.NewObject(graphql.ObjectConfig{
	Name: "EmpDetails",
	Fields: graphql.Fields{
		"emp_id":    &graphql.Field{Type: graphql.String},
		"emp_title": &graphql.Field{Type: graphql.String},
		"address": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				token, ok := p.Context.Value("token").(string)
				if !ok || token == "" {
					return nil, errors.New("authorization required")
				}
				_, err := utils.ValidateToken(token)
				if err != nil {
					return nil, err
				}
				return p.Source.(models.EmpDetails).Address, nil
			},
		},
	},
})

func GraphQLHandler() http.Handler {
	// Create root mutation first
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"login": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
	
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					username := p.Args["username"].(string)
					password := p.Args["password"].(string)
				
					emp, err := models.VerifyEmployeeCredentials(username, password)
					if err != nil {
						log.Printf("Login failed for %s: %v", username, err)
						return nil, fmt.Errorf("login failed: %v", err) // More detailed error
					}
					
					token, err := utils.GenerateToken(*emp)
					if err != nil {
						log.Printf("Token generation failed: %v", err)
						return nil, errors.New("failed to generate token")
					}

					log.Printf("Login successful for user: %s", username)
					return token, nil
				},
			},
			"register": &graphql.Field{
				Type: accountType,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					username := p.Args["username"].(string)
					password := p.Args["password"].(string)
					return models.CreateEmployee(username, password)
				},
			},
		},
	})

	// Create root query
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"account": &graphql.Field{
				Type: accountType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					token, ok := p.Context.Value("token").(string)
					if !ok || token == "" {
						return nil, errors.New("authorization required")
					}
					emp, err := utils.ValidateToken(token)
					if err != nil {
						return nil, err
					}
					return models.GetEmployeeByUsername(emp.Username)
				},
			},
			"employeeDetails": &graphql.Field{
				Type: graphql.NewList(empType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					token, ok := p.Context.Value("token").(string)
					if !ok || token == "" {
						return nil, errors.New("authorization required")
					}
					_, err := utils.ValidateToken(token)
					if err != nil {
						return nil, err
					}
					return models.GetAllEmpDetails()
				},
			},
		},
	})

	// Create schema with both query and mutation
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// Debug: Print all available operations
	log.Println("Available Queries:")
	for name := range rootQuery.Fields() {
		log.Println("-", name)
	}
	log.Println("Available Mutations:")
	for name := range rootMutation.Fields() {
		log.Println("-", name)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}

		// Parse request body
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		log.Printf("Received GraphQL query: %s", reqBody.Query)

		// Execute GraphQL query


		
		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  reqBody.Query,
			VariableValues: reqBody.Variables,
			Context:        r.Context(),
		})

		if len(result.Errors) > 0 {
			log.Printf("GraphQL errors: %v", result.Errors)
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
}