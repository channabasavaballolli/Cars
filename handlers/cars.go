package handlers

import (
	"car-service/db"
	"car-service/models"
	"car-service/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cars = append(cars, c)
	}
	json.NewEncoder(w).Encode(cars)
}

func CreateCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var c models.Car
	_ = json.NewDecoder(r.Body).Decode(&c)

	if err := utils.ValidateCar(c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.DB.QueryRow("INSERT INTO cars (make, model, year, price, color, mileage) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		c.Make, c.Model, c.Year, c.Price, c.Color, c.Mileage).Scan(&c.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(c)
}

func GetCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var c models.Car
	err := db.DB.QueryRow("SELECT id, make, model, year, price, color, mileage FROM cars WHERE id=$1", id).
		Scan(&c.ID, &c.Make, &c.Model, &c.Year, &c.Price, &c.Color, &c.Mileage)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Car not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(c)
}

func UpdateCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var c models.Car
	_ = json.NewDecoder(r.Body).Decode(&c)

	_, err := db.DB.Exec("UPDATE cars SET make=$1, model=$2, year=$3, price=$4, color=$5, mileage=$6 WHERE id=$7",
		c.Make, c.Model, c.Year, c.Price, c.Color, c.Mileage, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.ID = id
	json.NewEncoder(w).Encode(c)
}

func DeleteCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := db.DB.Exec("DELETE FROM cars WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}
