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

	"github.com/jackc/pgx/v5/pgxpool"
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
	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	err = dbpool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Could not ping database: %v\n", err)
	}

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

	// ENDPOINT (/api/ping)
	http.HandleFunc("/api/ping", DBping(dbpool))

	/*
		┌-----------------------------------┐
		|     User managment endpoints      |
		└-----------------------------------┘
	*/

	// ENDPOINT (/api/login)
	http.HandleFunc("/api/login", LoginHandler(dbpool))

	// ENDPOINT (/api/register)
	http.HandleFunc("/api/register", RegisterHandler(dbpool))

	// ENDPOINT (/api/edit-user)
	http.HandleFunc("/api/edit-user", UserEditHandler(dbpool))

	// ENDPOINT (/api/delete-user)
	http.HandleFunc("/api/delete-user", UserDeleteHandler(dbpool))

	/*
		┌-----------------------------------┐
		|   Purchase managment endpoints    |
		└-----------------------------------┘
	*/

	// ENDPOINT (/api/delete-purchase)
	http.HandleFunc("/api/delete-purchase", PurchaseDeleteHandler(dbpool))

	// ENDPOINT (/api/edit-purchase)
	http.HandleFunc("/api/edit-purchase", PurchaseEditHandler(dbpool))

	// ENDPOINT (/api/add-purchase)
	http.HandleFunc("/api/add-purchase", PurchaseAddHandler(dbpool))

	/*
		┌-----------------------------------┐
		|     Get data lists endpoints      |
		└-----------------------------------┘
	*/

	// ENDPOINT (/api/users)
	http.HandleFunc("/api/users", UsersListHandler(dbpool))

	// ENDPOINT (/api/user/purchases)
	http.HandleFunc("/api/user/purchases", UserPurchasesListHandler(dbpool))

	// ENDPOINT (/api/me)
	http.HandleFunc("/api/me", GetCurrentUser(dbpool))

	/*
		┌-----------------------------------┐
		|            WEB SERVER             |
		└-----------------------------------┘
	*/

	// SWAGGER
	http.HandleFunc("/api/swagger/", httpSwagger.WrapHandler)

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
			SendError(w, http.StatusUnauthorized, "Missing token")
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			SendError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int(claims["user_id"].(float64))

			ctx := context.WithValue(r.Context(), "userID", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			SendError(w, http.StatusUnauthorized, "Invalid token claims")
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
// @Router       /api/login [post]
func LoginHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.User
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			SendError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		userFromDB, err := database.GetUserByUsername(r.Context(), db, input.Username)
		if err != nil {
			SendError(w, http.StatusUnauthorized, "User not found or Access Denied")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(input.Password))
		if err != nil {
			log.Printf("Wrong password for user %s", input.Username)
			SendError(w, http.StatusUnauthorized, "Access Denied")
			return
		}

		// Create a JWT (token)
		secret := []byte(os.Getenv("JWT_SECRET"))
		claims := jwt.MapClaims{
			"user_id": userFromDB.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(secret)
		if err != nil {
			SendError(w, http.StatusInternalServerError, "Could not create token")
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
func RegisterHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var NewUser models.User

		err := json.NewDecoder(r.Body).Decode(&NewUser)
		if err != nil {
			SendError(w, http.StatusBadRequest, "Bad request: invalid JSON")
			return
		}
		if err := validate.Struct(NewUser); err != nil {
			SendError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(NewUser.Password), 14)
		if err != nil {
			SendError(w, http.StatusInternalServerError, "Security error")
			return
		}
		NewUser.Password = string(hashedPassword)

		err = database.SaveNewUser(r.Context(), db, NewUser)
		if err != nil {
			log.Printf("Error on saving user: %v", err)
			SendError(w, http.StatusInternalServerError, "Can't save user")
			return
		}

		SendJSON(w, http.StatusOK, map[string]string{"message": "User added successfully!"})
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
func UserEditHandler(db *pgxpool.Pool) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			SendError(w, http.StatusBadRequest, "Invalid user ID format")
			return
		}
		trustedID, ok := r.Context().Value("userID").(int)
		if !ok {
			SendError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if trustedID != id {
			SendError(w, http.StatusBadRequest, "Don't match the user ID")
			return
		}

		var data models.User
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			SendError(w, http.StatusBadRequest, "Bad request: invalid JSON")
			return
		}
		data.ID = id

		if data.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
			if err != nil {
				SendError(w, http.StatusInternalServerError, "Security error")
				return
			}
			data.Password = string(hashedPassword)
		} else {
			currentUser, err := database.GetUserByID(r.Context(), db, id)
			if err != nil {
				SendError(w, http.StatusNotFound, "User not found")
				return
			}
			data.Password = currentUser.Password
		}

		err = database.EditUser(r.Context(), db, data)
		if err != nil {
			SendError(w, http.StatusInternalServerError, "Update failed")
			return
		}

		SendJSON(w, http.StatusOK, map[string]string{"message": "User edit successfully!"})
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
func UserDeleteHandler(db *pgxpool.Pool) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		trustedID := r.Context().Value("userID").(int)

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			SendError(w, http.StatusBadRequest, "Bad request: missing user ID")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			SendError(w, http.StatusBadRequest, "Invalid id format")
			return
		}
		if trustedID != id {
			SendError(w, http.StatusBadRequest, "Don't match the user ID")
			return
		}

		err = database.DeleteUserByID(r.Context(), db, id)
		if err != nil {
			log.Printf("Error deleting user: %v", err)
			SendError(w, http.StatusInternalServerError, "Can't delete user")
			return
		}

		SendJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully!"})
	})
}

// @Summary      Проверка базы данных
// @Description  Проверяет, отвечает ли база данных
// @Tags         system
// @Produce      plain
// @Success      200 {string} string "Database connection is alive!"
// @Router       /ping [get]
func DBping(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SendJSON(w, http.StatusOK, map[string]string{"message": "Database connection is alive!"})
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
func PurchaseAddHandler(db *pgxpool.Pool) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var NewItem models.Purchase

		err := json.NewDecoder(r.Body).Decode(&NewItem)
		if err != nil {
			SendError(w, http.StatusBadRequest, "Bad request: invalid JSON")
			return
		}

		if err := validate.Struct(NewItem); err != nil {
			SendError(w, http.StatusBadRequest, "Validation failed", err.Error())
			return
		}

		trustedUserID, ok := r.Context().Value("userID").(int)
		if !ok {
			SendError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		NewItem.UserID = trustedUserID

		err = database.SavePurchase(r.Context(), db, NewItem)
		if err != nil {
			log.Printf("Error on saving: %v", err)
			SendError(w, http.StatusInternalServerError, "Can't save item")
			return
		}

		SendJSON(w, http.StatusCreated, map[string]string{"message": "Purchase added successfully!"})
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
func PurchaseDeleteHandler(db *pgxpool.Pool) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		trustedID, ok := r.Context().Value("userID").(int)
		if !ok {
			SendError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			SendError(w, http.StatusBadRequest, "Bad request: missing purchase ID")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			SendError(w, http.StatusBadRequest, "Invalid id format")
			return
		}

		err = database.DeletePurchaseByID(r.Context(), db, id, trustedID)
		if err != nil {
			log.Printf("Error on deleting: %v", err)
			SendError(w, http.StatusInternalServerError, "Can't delete purchase")
			return
		}

		SendJSON(w, http.StatusOK, map[string]string{"message": "Purchase deleted successfully!"})
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
func PurchaseEditHandler(db *pgxpool.Pool) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			SendError(w, http.StatusBadRequest, "Invalid purchase ID")
			return
		}

		var EditedPurchase models.Purchase

		err = json.NewDecoder(r.Body).Decode(&EditedPurchase)
		if err != nil {
			SendError(w, http.StatusBadRequest, "Bad request: invalid JSON")
			return
		}
		EditedPurchase.ID = id

		trustedID, ok := r.Context().Value("userID").(int)
		if !ok {
			SendError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		EditedPurchase.UserID = trustedID

		err = database.EditPurchase(r.Context(), db, EditedPurchase)
		if err != nil {
			log.Printf("Edit error: %v", err)
			SendError(w, http.StatusForbidden, "Update failed: item not found or access denied")
			return
		}

		SendJSON(w, http.StatusOK, map[string]string{"message": "Purchase edit successfully!"})
	})
}

// @Summary      Список всех пользователей
// @Description  Возвращает список всех зарегистрированных юзеров (публично)
// @Tags         system
// @Produce      json
// @Success      200 {array} models.User
// @Router       /users [get]
func UsersListHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := database.GetAllUsers(r.Context(), db)
		if err != nil {
			SendError(w, http.StatusInternalServerError, "DB error")
			return
		}

		SendJSON(w, http.StatusOK, list)
	}
}

// @Summary      Список покупок пользователя
// @Description  Возвращает все покупки текущего авторизованного пользователя
// @Tags         purchases
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Purchase
// @Router       /user/purchases [get]
func UserPurchasesListHandler(db *pgxpool.Pool) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		trustedUserID, ok := r.Context().Value("userID").(int)
		if !ok {
			SendError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		list, err := database.GetPurchasesByID(r.Context(), db, trustedUserID)
		if err != nil {
			SendError(w, http.StatusInternalServerError, "DB error")
			return
		}

		SendJSON(w, http.StatusOK, list)
	})
}

func GetCurrentUser(db *pgxpool.Pool) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			SendError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		ctx := r.Context()

		user, err := database.GetUserByID(ctx, db, userID)
		if err != nil {
			log.Printf("Error: %v", err)
			SendError(w, http.StatusNotFound, "User not found")
			return
		}

		user.Password = ""

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			SendError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

/*
┌-----------------------------------┐
|          Error Handling           |
└-----------------------------------┘
*/
type APIError struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func SendError(w http.ResponseWriter, statusCode int, message string, details ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := APIError{Error: message}
	if len(details) > 0 {
		resp.Details = details[0]
	}

	json.NewEncoder(w).Encode(resp)
}

func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
