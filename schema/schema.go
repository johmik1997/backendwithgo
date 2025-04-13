package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"john/models"
	"john/utils"
	"log"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
)

var accountType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Account",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.Int},
			"username":  &graphql.Field{Type: graphql.String},
			"isAdmin":   &graphql.Field{Type: graphql.Boolean},  // Changed to match
			"createdAt": &graphql.Field{Type: graphql.DateTime},
		},
	})
	var employeeType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Employee",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.Int},
			"empId":       &graphql.Field{Type: graphql.Int},
			"empName":     &graphql.Field{Type: graphql.String},
			"department":  &graphql.Field{Type: graphql.String},
			"experience":  &graphql.Field{Type: graphql.Int},
			"address":     &graphql.Field{Type: graphql.String},
			"birthdate":   &graphql.Field{Type: graphql.String},
			"employePhoto": &graphql.Field{Type: graphql.String},
			"createdAt":   &graphql.Field{Type: graphql.String},
		},
	})
	

	var empType = graphql.NewObject(graphql.ObjectConfig{
		Name: "EmpDetails",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.Int},
			"empId":       &graphql.Field{Type: graphql.Int},
			"empName":     &graphql.Field{Type: graphql.String}, // Changed from emp_name to empName
			"department":  &graphql.Field{Type: graphql.String},
			"experience":  &graphql.Field{Type: graphql.Int},
			"address":     &graphql.Field{Type: graphql.String},
			"birthdate":   &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					emp := p.Source.(models.EmpDetails)
					return emp.Birthdate.Format("2006-01-02"), nil
				},
			},
			"employePhoto": &graphql.Field{Type: graphql.String},
		},
	})
	var participantType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Participant",
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.Int},
			"name":  &graphql.Field{Type: graphql.String},
			"class": &graphql.Field{Type: graphql.String},
		},
	})
	
var eventType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Event",
    Fields: graphql.Fields{
        "id":          &graphql.Field{Type: graphql.Int},
        "name":        &graphql.Field{Type: graphql.String},
        "date":        &graphql.Field{Type: graphql.DateTime},
        "location":    &graphql.Field{Type: graphql.String},
        "host":        &graphql.Field{Type: graphql.String},
        "description": &graphql.Field{Type: graphql.String},
        "participants": &graphql.Field{
            Type: graphql.NewList(participantType),
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                event := p.Source.(models.Event)
                return event.Participants, nil
            },
        },
    },
})
// schema.go
var employeType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Employees",
    Fields: graphql.Fields{
        "id": &graphql.Field{Type: graphql.Int},
        "username": &graphql.Field{Type: graphql.String},
        "isAdmin": &graphql.Field{Type: graphql.Boolean},
    },
})

var loginResponseType = graphql.NewObject(graphql.ObjectConfig{
    Name: "LoginResponse",
    Fields: graphql.Fields{
        "token": &graphql.Field{Type: graphql.String},
        "user": &graphql.Field{Type: employeType},
    },
})

// In your mutation fields


func GraphQLHandler() http.Handler {
	// Create root mutation first
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
		"login": &graphql.Field{
    Type: loginResponseType,
    Args: graphql.FieldConfigArgument{
        "username": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
        "password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
    },
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					username := p.Args["username"].(string)
					password := p.Args["password"].(string)

					user, err := models.VerifyEmployeeCredentials(username, password)
					if err != nil {
						return nil, errors.New("invalid credentials")
					}

					token, err := utils.GenerateToken(*user)
					if err != nil {
						return nil, errors.New("failed to generate token")
					}
					log.Printf("Login attempt for user: %s", username)

log.Printf("Login attempt for: %s", username)
log.Printf("Generated token: %s", token)
log.Printf("Returning: %+v", map[string]interface{}{
    "token": token,
    "user": map[string]interface{}{
        "id": user.ID,
        "username": user.Username,
        "isAdmin": user.IsAdmin,
    },
})
return map[string]interface{}{
	"token": token,
	"user": map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"isAdmin":  user.IsAdmin,
	},
}, nil
}},
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
			"addEmployee": &graphql.Field{
    Type: empType,
    Args: graphql.FieldConfigArgument{
        "empName":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
        "department":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
        "experience":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
        "address":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
        "birthdate":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
        "employePhoto": &graphql.ArgumentConfig{Type: graphql.String},
    },
    Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        // Get current user from context
        token, ok := p.Context.Value("token").(string)
        if !ok {
            return nil, errors.New("authentication required")
        }
        
        currentUser, err := utils.ValidateToken(token)
        if err != nil {
            return nil, fmt.Errorf("invalid token: %v", err)
        }

        // Parse input
        empName := p.Args["empName"].(string)
        department := p.Args["department"].(string)
        experience := p.Args["experience"].(int)
        address := p.Args["address"].(string)
        birthdateStr := p.Args["birthdate"].(string)
        
        birthdate, err := time.Parse("2006-01-02", birthdateStr)
        if err != nil {
            return nil, fmt.Errorf("invalid birthdate format: %v", err)
        }

        // Create employee details
        newEmp := models.EmpDetails{
            EmpID:        currentUser.ID, // Use logged-in user's ID
            EmpName:      empName,
            Department:   department,
            Experience:   experience,
            Address:      address,
            Birthdate:    birthdate,
            EmployePhoto: p.Args["employePhoto"].(string),
        }

        // Save to database
        if err := models.AddEmployee(&newEmp); err != nil {
            return nil, fmt.Errorf("failed to add employee: %v", err)
        }

        return newEmp, nil
    },
},

},
	})
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
            "account": &graphql.Field{
                Type: accountType,
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    claims, ok := p.Context.Value("user").(*utils.Claims)
                    if !ok {
                        return nil, errors.New("authentication required")
                    }
                    return models.GetEmployeeByUsername(claims.Username)
                },
            },
			"upcomingEvents": &graphql.Field{
    Type: graphql.NewList(eventType),
    Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        return models.GetUpcomingEvents()
    },
},
"teamMembers": &graphql.Field{
    Type: graphql.NewList(employeeType),
    Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        token, ok := p.Context.Value("token").(string)
        if !ok || token == "" {
            return nil, errors.New("authorization required")
        }
        
        if _, err := utils.ValidateToken(token); err != nil {
            return nil, errors.New("invalid token")
        }
        
        members, err := models.GetAllTeamMembers()
        if err != nil {
            log.Printf("Database error fetching team members: %v", err)
            return nil, errors.New("could not retrieve team members")
        }
        return members, nil
    },
},
			
	"employeeDetails": &graphql.Field{
	Type: graphql.NewList(empType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		// Get user from context
		user, ok := p.Context.Value("user").(*utils.Claims)
		if !ok {
			log.Println("No user found in context")
			return nil, errors.New("authentication required")
		}

		log.Printf("Fetching employee details for user %s", user.Username)
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
		// Reject non-POST requests
		if r.Method != "POST" {
			http.Error(w, "GraphQL only supports POST requests", http.StatusMethodNotAllowed)
			return
		}

		var reqBody struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}

		// Parse request with better error handling
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			log.Printf("Bad request: %v", err)
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return
		}

		log.Printf("GraphQL Query: %s", reqBody.Query)
		if reqBody.Variables != nil {
			log.Printf("Variables: %v", reqBody.Variables)
		}

		// Execute with timeout context
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  reqBody.Query,
			VariableValues: reqBody.Variables,
			Context:        ctx,
		})

		// Enhanced error handling
		if len(result.Errors) > 0 {
			log.Printf("GraphQL errors: %+v", result.Errors)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	})
}