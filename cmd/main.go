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
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"

	_ "list-application-v2/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"reflect"

	"github.com/shopspring/decimal"
)

var validate *validator.Validate

// @title           List Application API
// @version         2.0
// @description     Сервер для управления списком покупок.
// @host            localhost:8081
// @BasePath        /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	/*
		┌-----------------------------------┐
		|          READ ENV FILE            |
		└-----------------------------------┘
	*/
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failer to load .env")
	}
	connStr := os.Getenv("DATABASE_URL")

	/*
		┌-----------------------------------┐
		|        CONNECT TO DATABASE        |
		└-----------------------------------┘
	*/
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Can't to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Succses to connect to PoostgreSQL!")

	/*
		┌-----------------------------------┐
		|            VALIDATOR              |
		└-----------------------------------┘
	*/
	validate = validator.New()

	validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if val, ok := field.Interface().(decimal.Decimal); ok {
			f, _ := val.Float64()
			return f
		}
		return nil
	}, decimal.Decimal{})

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
	http.HandleFunc("/ping", DBping(conn))

	/*
		┌-----------------------------------┐
		|     User managment endpoints      |
		└-----------------------------------┘
	*/

	// ENDPOINT (/login)
	http.HandleFunc("/login", LoginHandler(conn))

	// ENDPOINT (/register)
	http.HandleFunc("/register", RegisterHandler(conn))

	// ENDPOINT (/edit-user)
	http.HandleFunc("/edit-user", UserEditHandler(conn))

	// ENDPOINT (/delete-user)
	http.HandleFunc("/delete-user", UserDeleteHandler(conn))

	/*
		┌-----------------------------------┐
		|   Purchase managment endpoints    |
		└-----------------------------------┘
	*/

	// ENDPOINT (/delete-purchase)
	http.HandleFunc("/delete-purchase", PurchaseDeleteHandler(conn))

	// ENDPOINT (/edit-purchase)
	http.HandleFunc("/edit-purchase", PurchaseEditHandler(conn))

	// ENDPOINT (/add-purchase)
	http.HandleFunc("/add-purchase", PurchaseAddHandler(conn))

	/*
		┌-----------------------------------┐
		|     Get data lists endpoints      |
		└-----------------------------------┘
	*/

	// ENDPOINT (/users)
	http.HandleFunc("/users", UsersListHandler(conn))

	// ENDPOINT (/user/purchases)
	http.HandleFunc("/user/purchases", UserPurchasesListHandler(conn))

	/*
		┌-----------------------------------┐
		|            WEB SERVER             |
		└-----------------------------------┘
	*/

	// SWAGGER
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// START SERVER
	fmt.Println("Сервер запущен на http://localhost:8081")
	err = http.ListenAndServe(":8081", CORSMiddleware(http.DefaultServeMux))
	if err != nil {
		fmt.Printf("Ошибка сервера: %v\n", err)
	}
}

/*
┌-----------------------------------┐
|          JWT TOKEN CHECK          |
└-----------------------------------┘
*/
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int(claims["user_id"].(float64))

			ctx := context.WithValue(r.Context(), "userID", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
	}
}

/*
	┌-----------------------------------┐
	|      ENDPOINT FUNCTIONALITY       |
	└-----------------------------------┘
*/

// @Summary      Вход в систему
// @Description  Проверяет учетные данные и выдает JWT токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body models.User true "Данные для входа (username, password)"
// @Success      200 {object} map[string]string "token: string"
// @Router       /login [post]
func LoginHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Create a JWT (token)
		secret := []byte(os.Getenv("JWT_SECRET"))
		claims := jwt.MapClaims{
			"user_id": userFromDB.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен живет 24 часа
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(secret)
		if err != nil {
			http.Error(w, "Could not create token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	}
}

// @Summary      Регистрация
// @Description  Создает нового пользователя в системе
// @Tags         auth
// @Accept       json
// @Produce      plain
// @Param        input body models.User true "Данные пользователя"
// @Success      201 {string} string "User added successfully!"
// @Failure      400 {string} string "Validation failed / Invalid JSON"
// @Router       /register [post]
func RegisterHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var NewUser models.User

		err := json.NewDecoder(r.Body).Decode(&NewUser)
		if err != nil {
			http.Error(w, "Bad reqyest: invalid JSON", http.StatusBadRequest)
			return
		}
		if err := validate.Struct(NewUser); err != nil {
			http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
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
	}
}

// @Summary      Редактирование пользователя
// @Description  Обновляет данные профиля текущего пользователя (требует ID в Query)
// @Tags         auth
// @Accept       json
// @Produce      plain
// @Security     BearerAuth
// @Param        id query int true "User ID"
// @Param        input body models.User true "Новые данные"
// @Success      200 {string} string "User edit successfully!"
// @Router       /edit-user [post]
func UserEditHandler(conn *pgx.Conn) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)
		trustedID := r.Context().Value("userID").(int)
		if trustedID != id {
			http.Error(w, "Don't match the user ID", http.StatusBadRequest)
			return
		}

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
}

// @Summary      Удаление аккаунта
// @Description  Удаляет пользователя по ID
// @Tags         auth
// @Produce      plain
// @Security     BearerAuth
// @Param        id query int true "User ID"
// @Success      200 {string} string "User deleted successfully!"
// @Router       /delete-user [delete]
func UserDeleteHandler(conn *pgx.Conn) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		trustedID := r.Context().Value("userID").(int)

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
		if trustedID != id {
			http.Error(w, "Don't match the user ID", http.StatusBadRequest)
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
}

// @Summary      Проверка базы данных
// @Description  Проверяет, отвечает ли база данных
// @Tags         system
// @Produce      plain
// @Success      200 {string} string "Database connection is alive!"
// @Router       /ping [get]
func DBping(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Database connection is alive!")
	}
}

// @Summary      Добавить покупку
// @Description  Создает запись о покупке для текущего пользователя
// @Tags         purchases
// @Accept       json
// @Produce      plain
// @Security     BearerAuth
// @Param        input body models.Purchase true "Данные покупки"
// @Success      201 {string} string "Purchase added successfully!"
// @Router       /add-purchase [post]
func PurchaseAddHandler(conn *pgx.Conn) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var NewItem models.Purchase

		err := json.NewDecoder(r.Body).Decode(&NewItem)
		if err != nil {
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(NewItem); err != nil {
			http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		trustedUserID := r.Context().Value("userID").(int)

		NewItem.UserID = trustedUserID

		err = database.SavePurchase(conn, NewItem)
		if err != nil {
			log.Printf("Error on saving: %v", err)
			http.Error(w, "Can't save item", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Purchase added successfully!")
	})
}

// @Summary      Удалить покупку
// @Description  Удаляет покупку по её ID
// @Tags         purchases
// @Produce      plain
// @Security     BearerAuth
// @Param        id query int true "Purchase ID"
// @Success      200 {string} string "Purchase deleted successfully!"
// @Router       /delete-purchase [delete]
func PurchaseDeleteHandler(conn *pgx.Conn) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		trustedID := r.Context().Value("userID").(int)

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

		err = database.DeletePurchaseByID(conn, id, trustedID)
		if err != nil {
			log.Printf("Error on deleting: %v", err)
			http.Error(w, "Can't delete purchase", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Purchase deleted successfully!")
	})
}

// @Summary      Редактировать покупку
// @Description  Обновляет информацию о покупке
// @Tags         purchases
// @Accept       json
// @Produce      plain
// @Security     BearerAuth
// @Param        id query int true "Purchase ID"
// @Param        input body models.Purchase true "Новые данные покупки"
// @Success      200 {string} string "Purchase edit successfully!"
// @Router       /edit-purchase [post]
func PurchaseEditHandler(conn *pgx.Conn) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid purchase ID", http.StatusBadRequest)
			return
		}

		var EditedPurchase models.Purchase

		err = json.NewDecoder(r.Body).Decode(&EditedPurchase)
		if err != nil {
			http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
			return
		}
		EditedPurchase.ID = id

		trustedID := r.Context().Value("userID").(int)
		EditedPurchase.UserID = trustedID

		err = database.EditPurchase(conn, EditedPurchase)
		if err != nil {
			log.Printf("Edit error: %v", err)
			http.Error(w, "Update failed: item not found or access denied", http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Purchase edit successfully!")
	})
}

// @Summary      Список всех пользователей
// @Description  Возвращает список всех зарегистрированных юзеров (публично)
// @Tags         system
// @Produce      json
// @Success      200 {array} models.User
// @Router       /users [get]
func UsersListHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := database.GetAllUsers(conn)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	}
}

// @Summary      Список покупок пользователя
// @Description  Возвращает все покупки текущего авторизованного пользователя
// @Tags         purchases
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Purchase
// @Router       /user/purchases [get]
func UserPurchasesListHandler(conn *pgx.Conn) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		trustedUserID := r.Context().Value("userID").(int)
		list, err := database.GetPurchasesByID(conn, trustedUserID)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // В проде лучше указать конкретный домен
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
