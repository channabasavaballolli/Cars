package graph

import (
	"car-service/db"
	"car-service/middleware"
	"car-service/models"
	"car-service/utils"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/graphql-go/graphql"
)

// CarType defines the GraphQL object for a Car
var CarType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Car",
		Fields: graphql.Fields{
			"id":      &graphql.Field{Type: graphql.Int},
			"make":    &graphql.Field{Type: graphql.String},
			"model":   &graphql.Field{Type: graphql.String},
			"year":    &graphql.Field{Type: graphql.Int},
			"price":   &graphql.Field{Type: graphql.Float},
			"color":   &graphql.Field{Type: graphql.String},
			"mileage": &graphql.Field{Type: graphql.Int},
		},
	},
)

// RootQuery defines the entry point for queries
var RootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"cars": &graphql.Field{
				Type: graphql.NewList(CarType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Fetch cars from DB (Duplicated logic from handlers for simplicity)
					rows, err := db.DB.Query("SELECT id, make, model, year, price, color, mileage FROM cars")
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var cars []models.Car
					for rows.Next() {
						var c models.Car
						if err := rows.Scan(&c.ID, &c.Make, &c.Model, &c.Year, &c.Price, &c.Color, &c.Mileage); err != nil {
							return nil, err
						}
						cars = append(cars, c)
					}
					return cars, nil
				},
			},
		},
	},
)

// RootMutation defines the entry point for mutations
var RootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			// --- Auth Mutations ---
			"requestLogin": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					email, _ := p.Args["email"].(string)

					// 1. Ensure user exists (Upsert) - Default role is 'user' via DB default
					var userID int
					fmt.Printf("Attempting to login/register email: %s\n", email)
					err := db.DB.QueryRow("INSERT INTO users (email) VALUES ($1) ON CONFLICT (email) DO UPDATE SET email=EXCLUDED.email RETURNING id", email).Scan(&userID)
					if err != nil {
						fmt.Printf("Database error during user upsert: %v\n", err)
						return nil, fmt.Errorf("database error: %v", err)
					}
					fmt.Printf("User ID for %s is %d\n", email, userID)

					// 2. Generate 6-digit code
					rng := rand.New(rand.NewSource(time.Now().UnixNano()))
					code := fmt.Sprintf("%06d", rng.Intn(1000000))

					// 3. Save code to DB
					expiry := time.Now().Add(15 * time.Minute)
					_, err = db.DB.Exec("INSERT INTO verification_codes (user_id, code, expires_at) VALUES ($1, $2, $3)", userID, code, expiry)
					if err != nil {
						return nil, fmt.Errorf("failed to save code: %v", err)
					}

					// 4. Send Email
					err = utils.SendOTP(email, code)
					if err != nil {
						return nil, fmt.Errorf("failed to send email: %v", err)
					}

					return "Verification code sent to email", nil
				},
			},
			"verifyLogin": &graphql.Field{
				Type: graphql.String, // Returns JWT Token
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"code":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					email, _ := p.Args["email"].(string)
					code, _ := p.Args["code"].(string)

					// 1. Get User ID and Role
					var userID int
					var role string
					err := db.DB.QueryRow("SELECT id, role FROM users WHERE email=$1", email).Scan(&userID, &role)
					if err != nil {
						return nil, errors.New("user not found")
					}

					// 2. Verify Code
					var dbCode string
					var expiresAt time.Time
					err = db.DB.QueryRow("SELECT code, expires_at FROM verification_codes WHERE user_id=$1 AND code=$2 ORDER BY created_at DESC LIMIT 1", userID, code).Scan(&dbCode, &expiresAt)
					if err != nil {
						return nil, errors.New("invalid code")
					}

					if time.Now().After(expiresAt) {
						return nil, errors.New("code expired")
					}

					// 3. Generate JWT with Role
					token, err := utils.GenerateToken(userID, role)
					if err != nil {
						return nil, fmt.Errorf("failed to generate token: %v", err)
					}

					// 4. Clean up used codes (optional)
					_, _ = db.DB.Exec("DELETE FROM verification_codes WHERE user_id=$1", userID)

					return token, nil
				},
			},

			// --- Car Mutations (Protected) ---
			"createCar": &graphql.Field{
				Type: CarType,
				Args: graphql.FieldConfigArgument{
					"make":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"model":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"year":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"price":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
					"color":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"mileage": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Auth Check
					if p.Context.Value(middleware.UserIDKey) == nil {
						return nil, errors.New("unauthorized")
					}
					// RBAC Check
					if p.Context.Value(middleware.RoleKey) != "admin" {
						return nil, errors.New("forbidden: admins only")
					}

					make, _ := p.Args["make"].(string)
					model, _ := p.Args["model"].(string)
					year, _ := p.Args["year"].(int)
					price, _ := p.Args["price"].(float64)
					color, _ := p.Args["color"].(string)
					mileage, _ := p.Args["mileage"].(int)

					car := models.Car{
						Make:    make,
						Model:   model,
						Year:    year,
						Price:   price,
						Color:   color,
						Mileage: mileage,
					}

					if err := utils.ValidateCar(car); err != nil {
						return nil, err
					}

					err := db.DB.QueryRow(
						"INSERT INTO cars (make, model, year, price, color, mileage) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
						car.Make, car.Model, car.Year, car.Price, car.Color, car.Mileage).Scan(&car.ID)

					if err != nil {
						return nil, err
					}
					return car, nil
				},
			},
			"updateCar": &graphql.Field{
				Type: CarType,
				Args: graphql.FieldConfigArgument{
					"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"make":    &graphql.ArgumentConfig{Type: graphql.String},
					"model":   &graphql.ArgumentConfig{Type: graphql.String},
					"year":    &graphql.ArgumentConfig{Type: graphql.Int},
					"price":   &graphql.ArgumentConfig{Type: graphql.Float},
					"color":   &graphql.ArgumentConfig{Type: graphql.String},
					"mileage": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Auth Check
					if p.Context.Value(middleware.UserIDKey) == nil {
						return nil, errors.New("unauthorized")
					}
					// RBAC Check
					if p.Context.Value(middleware.RoleKey) != "admin" {
						return nil, errors.New("forbidden: admins only")
					}

					id, _ := p.Args["id"].(int)

					var car models.Car
					err := db.DB.QueryRow("SELECT id, make, model, year, price, color, mileage FROM cars WHERE id=$1", id).
						Scan(&car.ID, &car.Make, &car.Model, &car.Year, &car.Price, &car.Color, &car.Mileage)
					if err != nil {
						return nil, err
					}

					if val, ok := p.Args["make"].(string); ok {
						car.Make = val
					}
					if val, ok := p.Args["model"].(string); ok {
						car.Model = val
					}
					if val, ok := p.Args["year"].(int); ok {
						car.Year = val
					}
					if val, ok := p.Args["price"].(float64); ok {
						car.Price = val
					}
					if val, ok := p.Args["color"].(string); ok {
						car.Color = val
					}
					if val, ok := p.Args["mileage"].(int); ok {
						car.Mileage = val
					}

					if err := utils.ValidateCar(car); err != nil {
						return nil, err
					}

					_, err = db.DB.Exec("UPDATE cars SET make=$1, model=$2, year=$3, price=$4, color=$5, mileage=$6 WHERE id=$7",
						car.Make, car.Model, car.Year, car.Price, car.Color, car.Mileage, car.ID)
					if err != nil {
						return nil, err
					}

					return car, nil
				},
			},
			"deleteCar": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Auth Check
					if p.Context.Value(middleware.UserIDKey) == nil {
						return false, errors.New("unauthorized")
					}
					// RBAC Check
					if p.Context.Value(middleware.RoleKey) != "admin" {
						return false, errors.New("forbidden: admins only")
					}

					id, _ := p.Args["id"].(int)
					res, err := db.DB.Exec("DELETE FROM cars WHERE id=$1", id)
					if err != nil {
						return false, err
					}
					params, _ := res.RowsAffected()
					return params > 0, nil
				},
			},
		},
	},
)

// InitSchema creates and returns the GraphQL schema
func InitSchema() (graphql.Schema, error) {
	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    RootQuery,
			Mutation: RootMutation,
		},
	)
}
