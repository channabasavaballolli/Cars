package handlers

import (
	"car-app/db"
	"car-app/models"
	"car-app/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.DB.Query("SELECT id, make, model, year, price, color, mileage FROM cars")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var c models.Car

		if err := rows.Scan(&c.ID, &c.Make, &c.Model, &c.Year, &c.Price, &c.Color, &c.Mileage); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		cars = append(cars, c)
	}

	json.NewEncoder(w).Encode(cars)
}

func CreateCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c models.Car
	_ = json.NewDecoder(r.Body).Decode(&c)

	sqlStatement := `INSERT INTO cars (make, model, year, price, color, mileage) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	id := 0
	if err := utils.ValidateCar(c.Price, c.Year); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400 Bad Request
		return
	}
	err := db.DB.QueryRow(sqlStatement, c.Make, c.Model, c.Year, c.Price, c.Color, c.Mileage).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.ID = id
	json.NewEncoder(w).Encode(c)
}

func UpdateCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	var c models.Car
	_ = json.NewDecoder(r.Body).Decode(&c)

	sqlStatement := `UPDATE cars SET make=$1, model=$2, year=$3, price=$4, color=$5, mileage=$6 WHERE id=$7`

	res, err := db.DB.Exec(sqlStatement, c.Make, c.Model, c.Year, c.Price, c.Color, c.Mileage, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Car updated successfully"})
}

func DeleteCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	sqlStatement := `DELETE FROM cars WHERE id = $1`

	res, err := db.DB.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Car deleted successfully"})
}

func GetCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	var c models.Car
	sqlStatement := `SELECT id, make, model, year, price, color, mileage FROM cars WHERE id=$1`

	row := db.DB.QueryRow(sqlStatement, id)

	err := row.Scan(&c.ID, &c.Make, &c.Model, &c.Year, &c.Price, &c.Color, &c.Mileage)
	switch err {
	case sql.ErrNoRows:
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	case nil:
		json.NewEncoder(w).Encode(c)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
