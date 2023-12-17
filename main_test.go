package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func TestGetAllCars(t *testing.T) {
	// Set up a handler for your route
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, err := http.NewRequest("GET", "/cars", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func AddCar(car Car) error {
	// Add the car to the database
	_, err := db.Exec("INSERT INTO cars (RegistrationNumber, Status) VALUES (?, ?)", car.RegistrationNumber, car.Status)
	if err != nil {
		return err
	}
	return nil
}

func TestAddCars(t *testing.T) {
	// Set up a test database
	err := godotenv.Load("./configs/.env")
	if err != nil {
		t.Fatal("Failed to load .env file")
	}

	var (
		dbUser     = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbHost     = os.Getenv("DB_HOST")
		dbPort     = os.Getenv("DB_PORT")
		dbName     = os.Getenv("DB_NAME")
	)

	db, err = sqlx.Connect("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to database")
	}

	// Clean up the database after the test
	defer func() {
		db.Exec("DELETE FROM cars")
		db.Close()
	}()

	// Add a car
	car := Car{
		RegistrationNumber: "ABC123",
		Status:             "ENTRY",
	}

	err = AddCar(car)
	if err != nil {
		t.Fatalf("Failed to add car: %v", err)
	}

	// Check if the car was added successfully
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM cars where registrationNumber = ?", "ABC123")
	if err != nil {
		t.Fatalf("Failed to get car count: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected car count to be 1, got %d", count)
	}
}

func UpdateCar(car Car) error {
	// Add the car to the database
	_, err := db.Exec("UPDATE cars SET Status = ? WHERE RegistrationNumber = ?", car.Status, car.RegistrationNumber)
	if err != nil {
		return err
	}
	return nil
}

func TestUpdateCar(t *testing.T) {
	// Set up a test database
	err := godotenv.Load("./configs/.env")
	if err != nil {
		t.Fatal("Failed to load .env file")
	}

	var (
		dbUser     = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbHost     = os.Getenv("DB_HOST")
		dbPort     = os.Getenv("DB_PORT")
		dbName     = os.Getenv("DB_NAME")
	)

	db, err = sqlx.Connect("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	if err != nil {
		t.Fatal("Failed to connect to database")
	}

	// Create a car to update
	car := Car{
		RegistrationNumber: "MP08AB1234",
		Status:             "ENTRY",
	}

	// Clean up the database after the test
	defer func() {
		db.Exec("DELETE FROM cars where registrationNumber = ?", car.RegistrationNumber)
		db.Close()
	}()

	// Add the car to the database
	err = AddCar(car)
	if err != nil {
		t.Fatal("Failed to add car")
	}

	// Update the car
	car.Status = "IN_SERVICE"
	err = UpdateCar(car)
	if err != nil {
		t.Fatal("Failed to update car")
	}

	// Retrieve the car from the database
	var updatedCar Car
	err = db.Get(&updatedCar, "SELECT * FROM cars WHERE RegistrationNumber = ?", car.RegistrationNumber)
	if err != nil {
		t.Fatal("Failed to retrieve car")
	}

	// Check that the car was updated correctly
	if updatedCar.Status != "IN_SERVICE" {
		t.Errorf("Failed to update car: expected status to be 'IN_SERVICE', got '%s'", updatedCar.Status)
	}
}
