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
	err = db.Get(&count, "SELECT COUNT(*) FROM cars")
	if err != nil {
		t.Fatalf("Failed to get car count: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected car count to be 1, got %d", count)
	}
}
