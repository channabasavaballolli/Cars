http://localhost:8000/graphql

1. Create (Mutation)
Creates a new car and returns the generated ID.

graphql
mutation {
  createCar(
    make: "BMW"
    model: "M3"
    year: 2024
    price: 85000.0
    color: "Green"
    mileage: 0
  ) {
    id
    make
    model
    id # You can ask for the ID to use in the update/delete steps
  }
}


mutation {
  createCar(
    make: "Honda"
    model: "Civic"
    year: 2026
    price: 25000.0
    color: "Blue"
    mileage: 12
  ) {
    id
    make
    model
    price
    id # You can ask for the ID to use in the update/delete steps
  }
}
2. Read (Query)
Lists all cars in the database.

graphql
query {
  cars {
    id
    make
    model
    year
    price
    color
    mileage
  }
}
3. Update (Mutation)
Updates specific fields for a car. Replace 1 with the actual ID you want to update.

graphql
mutation {
  updateCar(
    id: 1
    price: 82000.0
    mileage: 1500
  ) {
    id
    price
    mileage
  }
}
4. Delete (Mutation)
Deletes a car by ID. Returns true if satisfied.

graphql
mutation {
  deleteCar(id: 1)
}


POST http://localhost:8000/cars

{
    "make": "Tesla",
    "model": "Model 3",
    "year": 2024,
    "price": 45000.00,
    "color": "White",
    "mileage": 10
}
 This i s Json


GET http://localhost:8000/cars

