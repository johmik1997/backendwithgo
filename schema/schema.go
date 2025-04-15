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
	"strings"
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
// 	var teamMemberType = graphql.NewObject(graphql.ObjectConfig{
//     Name: "TeamMember",
//     Fields: graphql.Fields{
//         "id":          &graphql.Field{Type: graphql.Int},
//         "empId":       &graphql.Field{Type: graphql.Int},
//         "empName":     &graphql.Field{Type: graphql.String},
//         "department":  &graphql.Field{Type: graphql.String},
//         "experience":  &graphql.Field{Type: graphql.Int},
//         "address":     &graphql.Field{Type: graphql.String},
//         "birthdate":   &graphql.Field{Type: graphql.String},
//         "employePhoto": &graphql.Field{Type: graphql.String},
//         "createdAt":   &graphql.Field{Type: graphql.String},
//     },
// })


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
        claims, ok := p.Context.Value("user").(*utils.Claims)
        if !ok {
            return nil, errors.New("authentication required")
        }

        // Parse input
        empName := p.Args["empName"].(string)
        department := p.Args["department"].(string)
        experience := p.Args["experience"].(int)
        address := p.Args["address"].(string)
        birthdateStr := p.Args["birthdate"].(string)
        
        birthdate, err := time.Parse("2006-01-02", birthdateStr)
        if err != nil {
            return nil, fmt.Errorf("invalid birthdate format (use YYYY-MM-DD): %v", err)
        }

        // Create employee details
        newEmp := models.EmpDetails{
            EmpID:       claims.ID, // Use logged-in user's ID
            EmpName:     empName,
            Department:  department,
            Experience:  experience,
            Address:     address,
            Birthdate:   birthdate,
            EmployePhoto: p.Args["employePhoto"].(string),
            CreatedAt:   time.Now(),
        }

        // Save to database
        if err := models.AddEmployee(&newEmp); err != nil {
            log.Printf("Failed to add employee: %v", err)
            return nil, fmt.Errorf("failed to add employee: %v", err)
        }

        log.Printf("Successfully added employee: %+v", newEmp)
        return newEmp, nil
    },
},
"updateEmployee": &graphql.Field{
				Type: empType,
				Args: graphql.FieldConfigArgument{
					"id":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"empName":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"department":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"experience":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"address":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"birthdate":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"employePhoto": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					birthdate, err := time.Parse("2006-01-02", p.Args["birthdate"].(string))
					if err != nil {
						return nil, fmt.Errorf("invalid birthdate format: %v", err)
					}
					

					emp := models.EmpUpdateDetails{
						ID:          p.Args["id"].(int),
						EmpName:     p.Args["empName"].(string),
						Department:  p.Args["department"].(string),
						Experience:  p.Args["experience"].(int),
						Address:     p.Args["address"].(string),
						Birthdate:   birthdate,
						EmployePhoto: p.Args["employePhoto"].(string),
					}

					if err := models.UpdateEmployee(&emp); err != nil {
						return nil, fmt.Errorf("failed to update employee: %v", err)
					}

					return emp, nil
				},
			},
		"deleteEmployee": &graphql.Field{
    Type: graphql.Boolean,
    Args: graphql.FieldConfigArgument{
        "id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
    },
    Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        // Enhanced error handling
        token, ok := p.Context.Value("token").(string)
        if !ok {
            return false, fmt.Errorf("missing authentication token")
        }

        claims, err := utils.ValidateToken(token)
        if err != nil {
            return false, fmt.Errorf("invalid token: %v", err)
        }

        log.Printf("User %d attempting to delete employee", claims.ID)

        id, ok := p.Args["id"].(int)
        if !ok {
            return false, fmt.Errorf("invalid employee ID")
        }

        if err := models.DeleteEmployee(id); err != nil {
            log.Printf("Delete failed: %v", err)
            return false, fmt.Errorf("failed to delete employee: %v", err)
        }

        return true, nil
    },
},

},
})
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			// In your rootQuery definition
"account": &graphql.Field{
    Type: accountType,
    Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        claims, ok := p.Context.Value("user").(*utils.Claims)
        if !ok {
            return nil, errors.New("authentication required")
        }
        user, err := models.GetEmployeeByUsername(claims.Username)
        if err != nil {
            return nil, fmt.Errorf("failed to fetch user: %v", err)
        }
        return user, nil
    },
},
			"employeeDetails": &graphql.Field{
				Type: graphql.NewList(empType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if _, ok := p.Context.Value("user").(*utils.Claims); !ok {
						return nil, errors.New("authentication required")
					}
					return models.GetAllEmpDetails() // Default pagination
				},
			},
			"health": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "OK", nil
				},
			},
			"upcomingEvents": &graphql.Field{
				Type: graphql.NewList(eventType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return models.GetUpcomingEvents()
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Handle OPTIONS for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only allow POST requests
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": []map[string]interface{}{
					{
						"message": "Only POST requests are supported",
					},
				},
			})
			return
		}

		// Parse request body
		var reqBody struct {
			Query         string                 `json:"query"`
			Variables     map[string]interface{} `json:"variables"`
			OperationName string                `json:"operationName"`
		}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": []map[string]interface{}{
					{
						"message": "Invalid request body: " + err.Error(),
					},
				},
			})
			return
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Add token to context if available
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenString)
			if err == nil {
				ctx = context.WithValue(ctx, "user", claims)
			}
			ctx = context.WithValue(ctx, "token", tokenString)
		}

		// Execute GraphQL operation
		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  reqBody.Query,
			VariableValues: reqBody.Variables,
			OperationName:  reqBody.OperationName,
			Context:        ctx,
		})

		// Handle errors
		if len(result.Errors) > 0 {
			log.Printf("GraphQL errors: %v", result.Errors)
		}

		// Return response
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("Failed to encode response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": []map[string]interface{}{
					{
						"message": "Failed to encode response",
					},
				},
			})
		}
	})
}