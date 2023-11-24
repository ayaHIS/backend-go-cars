package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Car struct {
	Model        string `json:"model"`
	Registration string `json:"registration"`
	Mileage      int    `json:"mileage"`
	IsRented     bool   `json:"isrented"`
}

var cars []Car

func listCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func addCar(w http.ResponseWriter, r *http.Request) {
	var newCar Car
	err := json.NewDecoder(r.Body).Decode(&newCar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the request
	if newCar.Model == "" || newCar.Registration == "" || newCar.Mileage < 0 {
		http.Error(w, "Invalid request. Make sure all fields are provided and mileage is non-negative.", http.StatusBadRequest)
		return
	}

	// Check if the car already exists by registration number
	for _, car := range cars {
		if car.Registration == newCar.Registration {
			http.Error(w, "Car with the same registration number already exists", http.StatusConflict)
			return
		}
	}

	// Save the new car to the MySQL database
	_, err = db.Exec("INSERT INTO cars (model, registration, mileage, is_rented) VALUES (?, ?, ?, ?)",
		newCar.Model, newCar.Registration, newCar.Mileage, newCar.IsRented)
	if err != nil {
		fmt.Printf("Error saving the new car to the database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the new car to your in-memory slice
	cars = append(cars, newCar)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCar)
}
func rentCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	registration := params["registration"]

	// Find the car with the given registration number
	for i, car := range cars {
		if car.Registration == registration {
			if car.IsRented {
				http.Error(w, "Car is already rented", http.StatusConflict)
				return
			}
			cars[i].IsRented = true
			json.NewEncoder(w).Encode(cars[i])
			return
		}
	}

	http.Error(w, "Car not found", http.StatusNotFound)
}

func returnCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	registration := params["registration"]

	// Find the car with the given registration number
	for i, car := range cars {
		if car.Registration == registration {
			if !car.IsRented {
				http.Error(w, "Car is not marked as rented", http.StatusConflict)
				return
			}

			// Simulate updating mileage and marking the car as available
			cars[i].Mileage += 100
			cars[i].IsRented = false

			json.NewEncoder(w).Encode(cars[i])
			return
		}
	}

	http.Error(w, "Car not found", http.StatusNotFound)
}

var db *sql.DB

func main() {

	var err error
	db, err = sql.Open("mysql", "pma:12341234@tcp(localhost:3306)/car_rental")
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}
	defer db.Close()
	r := mux.NewRouter()

	r.HandleFunc("/cars", listCars).Methods("GET")
	r.HandleFunc("/cars", addCar).Methods("POST")
	r.HandleFunc("/cars/{registration}/rentals", rentCar).Methods("POST")
	r.HandleFunc("/cars/{registration}/returns", returnCar).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
