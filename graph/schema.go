package graph

import (
	"car-service/db"
	"car-service/models"
	"car-service/utils"

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

					// Validate
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
					id, _ := p.Args["id"].(int)
					// Simple implementation: Fetch, Update, Save (or direct update)
					// For simplicity in this demo, we'll try to update provided fields.
					// Building a dynamic query is safer.

					// 1. Check if exists
					var car models.Car
					err := db.DB.QueryRow("SELECT id, make, model, year, price, color, mileage FROM cars WHERE id=$1", id).
						Scan(&car.ID, &car.Make, &car.Model, &car.Year, &car.Price, &car.Color, &car.Mileage)
					if err != nil {
						return nil, err // Not found or DB error
					}

					// 2. Overwrite fields if provided
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

					// 3. Validate
					if err := utils.ValidateCar(car); err != nil {
						return nil, err
					}

					// 4. Update
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
