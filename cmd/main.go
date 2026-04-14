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

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
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

	// ENDPOINT (/add-user)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var NewUser models.User

		err := json.NewDecoder(r.Body).Decode(&NewUser)
		if err != nil {
			http.Error(w, "Bad reqyest: invalid JSON", http.StatusBadRequest)
			return
		}
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
		var EditedUser models.User

		err := json.NewDecoder(r.Body).Decode(&EditedUser)
		if err != nil {
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		err = database.EditUser(conn, EditedUser)
		if err != nil {
			log.Printf("Error on edit user: %v", err)
			http.Error(w, "Can't edit user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "User edit successfully!")
	})

	// ENDPOINT (/delete-user)
	http.HandleFunc("/delete-user", func(w http.ResponseWriter, r *http.Request) {
		var User models.User

		err := json.NewDecoder(r.Body).Decode(&User)
		if err != nil {
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		err = database.DeleteUser(conn, User)
		if err != nil {
			log.Printf("Error on deleting: %v", err)
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
		var Item models.Purchase

		err := json.NewDecoder(r.Body).Decode(&Item)
		if err != nil {
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		err = database.DeletePurchase(conn, Item)
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
		var EditedPurchase models.Purchase

		err := json.NewDecoder(r.Body).Decode(&EditedPurchase)
		if err != nil {
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		err = database.EditPurchase(conn, EditedPurchase)
		if err != nil {
			log.Printf("Error on edit purchase: %v", err)
			http.Error(w, "Can't edit purchase", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
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
