package main

import (
	"context"
	"encoding/json"
	"fmt"
	"list-application-v2/internal/database"
	"list-application-v2/internal/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Read secrets
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failer to load .env")
	}
	connStr := os.Getenv("DATABASE_URL")

	// Conection to database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Can't to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Succses to connect to PoostgreSQL!")

	/*
		┌-----------------------------------┐
		|        DATABASE MIGRATION         |
		└-----------------------------------┘
	*/
	m, err := migrate.New(
		"file://migrations",
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	/*
		┌-----------------------------------┐
		|      Status check endpoints       |
		└-----------------------------------┘
	*/

	// ENDPOINT (/ping)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Database connection is alive!")
	})

	/*
		┌-----------------------------------┐
		|     User managment endpoints      |
		└-----------------------------------┘
	*/

	// ENDPOINT (/login)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var input models.User
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		userFromDB, err := database.GetUserByUsername(conn, input.Username)
		if err != nil {
			http.Error(w, "User not found or Access Denied", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(input.Password))
		if err != nil {
			log.Printf("Wrong password for user %s", input.Username)
			http.Error(w, "Access Denied", http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, "Welcome, %s!", userFromDB.Username)
	})

	// ENDPOINT (/register)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var NewUser models.User

		err := json.NewDecoder(r.Body).Decode(&NewUser)
		if err != nil {
			http.Error(w, "Bad reqyest: invalid JSON", http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(NewUser.Password), 14)
		if err != nil {
			http.Error(w, "Ошибка безопасности", 500)
			return
		}
		NewUser.Password = string(hashedPassword)

		err = database.SaveNewUser(conn, NewUser)
		if err != nil {
			log.Printf("Error on saving user: %v", err)
			http.Error(w, "Can't save user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "User added successfully!")
	})

	// ENDPOINT (/edit-user)
	http.HandleFunc("/edit-user", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)

		var data models.User
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}
		data.ID = id

		err = database.EditUser(conn, data)
		if err != nil {
			http.Error(w, "Update failed", 500)
			return
		}

		fmt.Fprint(w, "User edit successfully!")
	})

	// ENDPOINT (/delete-user)
	http.HandleFunc("/delete-user", func(w http.ResponseWriter, r *http.Request) {

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Bad request: missing user ID", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id format", http.StatusBadRequest)
			return
		}

		err = database.DeleteUserByID(conn, id)
		if err != nil {
			log.Printf("Error deleting user: %v", err)
			http.Error(w, "Can't delete user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "User deleted successfully!")
	})

	/*
		┌-----------------------------------┐
		|   Purchase managment endpoints    |
		└-----------------------------------┘
	*/

	// ENDPOINT (/delete-purchase)
	http.HandleFunc("/delete-purchase", func(w http.ResponseWriter, r *http.Request) {

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Bad request: missing purchase ID", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id format", http.StatusBadRequest)
			return
		}

		err = database.DeletePurchaseByID(conn, id)
		if err != nil {
			log.Printf("Error on deleting: %v", err)
			http.Error(w, "Can't delete purchase", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Purchase deleted successfully!")
	})

	// ENDPOINT (/edit-purchase)
	http.HandleFunc("/edit-purchase", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)

		var EditedPurchase models.Purchase

		err := json.NewDecoder(r.Body).Decode(&EditedPurchase)
		if err != nil {
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}
		EditedPurchase.ID = id

		err = database.EditPurchase(conn, EditedPurchase)
		if err != nil {
			http.Error(w, "Update failed", 500)
			return
		}

		fmt.Fprint(w, "Purchase edit successfully!")
	})

	// ENDPOINT (/add-purchase)
	http.HandleFunc("/add-purchase", func(w http.ResponseWriter, r *http.Request) {
		var NewItem models.Purchase

		err := json.NewDecoder(r.Body).Decode(&NewItem)
		if err != nil {
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		err = database.SavePurchase(conn, NewItem)
		if err != nil {
			log.Printf("Error on saving: %v", err)
			http.Error(w, "Can't save item", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Purchase added successfully!")
	})

	/*
		┌-----------------------------------┐
		|     Get data lists endpoints      |
		└-----------------------------------┘
	*/

	// ENDPOINT (/purchases)
	http.HandleFunc("/purchases", func(w http.ResponseWriter, r *http.Request) {
		list, err := database.GetAllPurchases(conn)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	})

	// ENDPOINT (/users)
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		list, err := database.GetAllUsers(conn)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	})

	// START SERVER
	fmt.Println("Сервер запущен на http://localhost:8081")
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("Ошибка сервера: %v\n", err)
	}
}
