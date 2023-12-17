package main

import (
	"errors"
	"log"
	"os"

	"gofr.dev/pkg/gofr"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var db *sqlx.DB

type Car struct {
	ID                 int    `db:"id"`
	RegistrationNumber string `db:"RegistrationNumber"`
	Status             string `db:"Status"`
}

func init() {
	err := godotenv.Load("./configs/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
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

	// Create the "cars" table if it doesn't exist.
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS cars (
			id INT AUTO_INCREMENT PRIMARY KEY,
			RegistrationNumber VARCHAR(255) NOT NULL,
			Status VARCHAR(255) NOT NULL
		);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := gofr.New()

	// Get All cars in the Garage
	app.GET("/cars", func(ctx *gofr.Context) (interface{}, error) {
		ctx.Gofr.Logger.Info("Getting all cars in the garage")
		var cars []Car
		err := db.Select(&cars, "SELECT * FROM cars.cars")
		if err != nil {
			return nil, err
		}
		return cars, nil
	})

	// Add a car to the garage
	app.POST("/cars/add", func(ctx *gofr.Context) (interface{}, error) {
		ctx.Gofr.Logger.Info("Adding a new car to the garage")
		var registration_number = ctx.Param("registrationNumber")
		var status = ctx.Param("status")
		ctx.Gofr.Logger.Info("Registration Number: " + registration_number + " Status: " + status)

		// Check if the car with the same registration number already exists
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM cars WHERE RegistrationNumber = ?", registration_number)
		if err != nil {
			return nil, err
		}

		if count > 0 {
			return "Car already in garage", nil
		}

		// Check if the status is valid, allowing three values: entry, in_service, completed
		if status != "ENTRY" && status != "IN_SERVICE" && status != "COMPLETED" {
			return nil, errors.New("Invalid status")
		}

		// Insert the new car into the database
		result, err := db.Exec("INSERT INTO cars (RegistrationNumber, Status) VALUES (?, ?)", registration_number, status)
		if err != nil {
			return nil, err
		}

		// Get the ID of the newly inserted car
		newCarID, _ := result.LastInsertId()

		// Create a response object with the ID
		response := struct {
			ID                 int    `json:"id"`
			RegistrationNumber string `json:"registrationNumber"`
			Status             string `json:"status"`
		}{
			ID:                 int(newCarID),
			RegistrationNumber: registration_number,
			Status:             status,
		}

		return response, nil
	})

	// Update the status of a car
	app.PUT("/cars/update", func(ctx *gofr.Context) (interface{}, error) {
		ctx.Gofr.Logger.Info("Updating the status of a car")
		id := ctx.Param("id")
		status := ctx.Param("status")

		// Check if the car exists
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM cars WHERE id = ?", id)
		if err != nil {
			return nil, err
		}

		if count == 0 {
			return "Car not found", nil
		}

		// Check if the status is valid, allowing from entry,in_service to in_service or completed
		if status != "IN_SERVICE" && status != "COMPLETED" {
			return nil, errors.New("Invalid status")
		}

		// Update the status of the car
		_, err = db.Exec("UPDATE cars SET Status = ? WHERE id = ?", status, id)
		if err != nil {
			return nil, err
		}

		// Create a response object with the ID
		response := struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Status  string `json:"status"`
		}{
			ID:      id,
			Message: "Status updated successfully",
			Status:  status,
		}

		return response, nil
	})

	// Delete a car from the garage based on ID
	app.DELETE("/cars/delete/id", func(ctx *gofr.Context) (interface{}, error) {
		ctx.Gofr.Logger.Info("Deleting a car from the garage based on ID")
		id := ctx.Param("id")

		// Check if the car exists
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM cars WHERE id = ?", id)
		if err != nil {
			return nil, err
		}

		if count == 0 {
			return "Car not found", nil
		}

		// Delete the car from the database
		_, err = db.Exec("DELETE FROM cars WHERE id = ?", id)
		if err != nil {
			return nil, err
		}

		return map[string]string{"message": "Car deleted successfully", "id": id}, nil
	})

	// Delete a car from the garage based on Registration Number
	app.DELETE("/cars/delete/registration", func(ctx *gofr.Context) (interface{}, error) {
		ctx.Gofr.Logger.Info("Deleting a car from the garage based on Registration Number")
		registrationNumber := ctx.Param("registrationNumber")

		// Check if the car exists
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM cars WHERE RegistrationNumber = ?", registrationNumber)
		if err != nil {
			return nil, err
		}

		if count == 0 {
			return "Car not found", nil
		}

		// Delete the car from the database
		_, err = db.Exec("DELETE FROM cars WHERE RegistrationNumber = ?", registrationNumber)
		if err != nil {
			return nil, err
		}

		return map[string]string{"message": "Car deleted successfully", "registration": registrationNumber}, nil
	})

	// Root endpoint
	app.GET("/", func(ctx *gofr.Context) (interface{}, error) {
		return "Garage Management Application is Up!", nil
	})

	app.Start()
}
