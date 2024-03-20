package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"whats-app-clone-service/models"
	"whats-app-clone-service/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	_ "golang.org/x/crypto/bcrypt"
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type loginResponse struct {
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db

func createConnection() *sql.DB {

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/wa_clone_db")

	if err != nil {
		panic(err.Error())
	}

	// defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connection successful!")
	return db
}

// CreateUser create a user in the postgres db
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty user of type models.User
	var user models.UserRegister

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert user function and pass the user
	insertID := register(user)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// GetAllUser will return all the users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the users in the db
	users, err := getJSON("SELECT id, username, phone, image from users")

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	// send all the users as response
	// json.NewEncoder(w).Encode(users)
	w.Header().Set("Content-Type", "application/json")
	w.Write(users)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty user of type models.User
	var user models.LoginPayload

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert user function and pass the user
	jwt := login(user)

	// format a response object
	res := loginResponse{
		Token:   jwt,
		Message: "User created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// ------------------------- handler functions ----------------

// register user
func register(user models.UserRegister) int64 {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	hassPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Fatalf("Error preparing statement", err)
		return 0
	}

	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO users (id, username, phone, image, password) VALUES (?, ?, ?, ?, ?) RETURNING id`
	stmt, err := db.Prepare(sqlStatement)

	if err != nil {
		log.Fatalf("Error preparing statement", err)
		return 0
	}

	res, err := stmt.Exec(uuid.New(), user.Username, user.Phone, user.Image, hassPassword)
	if err != nil {
		log.Fatalf("Error preparing statement %v", err)
		return 0
	}
	id, err := res.LastInsertId()

	if err != nil {
		log.Fatalf("error retrieving last inserted ID: %v", err)
	}
	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// register user
func login(user models.LoginPayload) string {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// declare query and get user from db
	query := "SELECT id, username, phone, image, password FROM users WHERE username = '" + user.Username + "'"
	exist, err := getJSON(query)

	//declare variable for store parses user data from db
	var data []models.User

	if err := json.Unmarshal([]byte(exist), &data); err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	isMatch := utils.CheckPasswordHash(user.Password, data[0].Password)
	if isMatch != true {
		log.Fatalf("Password is incorrect", err)
		return ""
	}

	jwt, err := utils.GenerateJWT(data[0])

	if err != nil {
		log.Fatalf("error generate jwt token", err)

	}

	return jwt
}

// getJSON select data from db then parse to json
func getJSON(sqlString string) ([]byte, error) {
	db := createConnection()
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return jsonData, err
	}

	return jsonData, nil
}
