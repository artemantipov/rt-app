package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	log "github.com/sirupsen/logrus"


	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

// User - main struct for Birthdays
type User struct {
	Name        string `json:"-"`
	DateOfBirth string `json:"dateOfBirth"`
}

// Birthdays - define table model for ORM
type Birthdays struct {
	gorm.Model
	Name        string `gorm:"column:name"`
	DateOfBirth string `gorm:"column:date_of_birth"`
}

// Change default log to proper JSON format for easy parsing by logstash/graylog/etc.
func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.PUT("/hello/:username", putData)
	e.GET("/hello/:username", getData)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler())) // expose default metrics for prometheus

	// Start server
	e.Logger.Fatal(e.Start(":8080"))

}

// Default message for R
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Test App for R")
}

// Put data into storage
func putData(c echo.Context) error {
	user := User{c.Param("username"), ""}

	if !onlyLetters(user.Name) {
		log.Printf("Non-letter(s) in username: %v", user.Name)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Non-letters in username: %v", user.Name))
	}

	err := c.Bind(&user)
	if err != nil {
		log.Printf("Failed to decode json: %v", err)
		return c.String(http.StatusInternalServerError, "")
	}

	valid, _, _ := dateParseCheck(user.DateOfBirth)
	if !valid {
		log.Printf("Date incorrect (more then today): %v", user.DateOfBirth)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Date incorrect: %v", user.DateOfBirth))
	}

	db, err := dbConnect()
	if err != nil {
		log.Printf("Failed connect to DB: %v", err)
		return c.String(http.StatusInternalServerError, "Failed connect to DB")
	}
	dbPut(user, db)

	return c.String(http.StatusNoContent, "")
}

// Get data from storage
func getData(c echo.Context) error {
	user := c.Param("username")

	if !onlyLetters(user) {
		log.Printf("Non-letter(s) in username: %v", user)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Non-letters in username: %v", user))
	}

	db, err := dbConnect()
	if err != nil {
		log.Printf("Failed connect to DB: %v", err)
		return c.String(http.StatusInternalServerError, "Failed connect to DB")
	}
	dateOfBirth, err := dbGet(user, db)
	if err != nil {
		log.Printf("Failed to get dateOfBirth: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to get dateOfBirth")
	}

	valid, today, diff := dateParseCheck(dateOfBirth)
	if !valid && !today {
		log.Printf("Failed to parse dateOfBirth: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to parse dateOfBirth")
	}

	switch diff {
	case 0:
		return c.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Hello %v! Happy birthday!", user)})
	default:
		return c.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Hello %v! Your birthday is in %v day(s)", user, diff)})
	}

}

// Check for only letter in username
func onlyLetters(u string) bool {
	for _, s := range u {
		if !unicode.IsLetter(s) {
			return false
		}
	}
	return true
}

// Date less then today and in expected format
func dateParseCheck(d string) (valid bool, today bool, diff int) {
	layout := "2006-01-02"
	today = false
	valid = true
	check, err := time.Parse(layout, d)
	if err != nil {
		log.Printf("Failed to parse date format: %v", err)
		diff = 0
		valid = false
		return
	}
	current := time.Now()
	diff = int(current.Sub(check).Hours() / 24)
	if diff < 1 && diff != 0 {
		log.Printf("Birthday date not correct, difference fromt today: %v", diff)
		valid = false
		return
	} else if diff == 0 {
		valid = false
		today = true
		return
	}
	return

}

func dbConnect() (db *gorm.DB, err error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbType := os.Getenv("DB_TYPE") // postgres|mysql|cloudsql

	connect := ""
	switch dbType {
	case "postgres":
		connect = fmt.Sprintf(`host=%v port=%v user=%v dbname=%v password=%v sslmode=disable`, dbHost, dbPort, dbUser, dbName, dbPass)
		log.Println("Database Type: postgres")
	case "mysql":
		connect = fmt.Sprintf(`%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local`, dbUser, dbPass, dbHost, dbPort, dbName)
		log.Println("Database Type: mysql")
	case "cloudsql":
		connect = fmt.Sprintf(`%v:%v@unix(/cloudsql/%v)/%v`, dbUser, dbPass, dbHost, dbName)
		dbType = "mysql"
		log.Println("Database Type: cloudsql")
		// default:
		// 	connect = "/tmp/gorm.db"
		// 	dbType = "sqlite3"
		// 	log.Println("Database Type: sqlite(default)")
	}

	db, err = gorm.Open(dbType, connect)
	if err != nil {
		log.Println("Connection Failed to Open")
	}
	log.Println("Connection Established")
	return
}

func dbPut(u User, db *gorm.DB) {
	dbType := os.Getenv("DB_TYPE")
	defer db.Close()
	db.AutoMigrate(&Birthdays{})
	user := &Birthdays{Name: strings.ToLower(u.Name), DateOfBirth: u.DateOfBirth}
	db.Where(Birthdays{Name: strings.ToLower(u.Name)}).Assign(Birthdays{Name: strings.ToLower(u.Name), DateOfBirth: u.DateOfBirth}).FirstOrCreate(&user)
	// Workaround for CloudSQL update record
	if dbType == "cloudsql" {
		db.Model(&user).Where(fmt.Sprintf("name = '%v'", strings.ToLower(u.Name))).Updates(Birthdays{Name: strings.ToLower(u.Name), DateOfBirth: u.DateOfBirth})
	}
}

func dbGet(username string, db *gorm.DB) (dateOfBirth string, err error) {
	defer db.Close()
	dateOfBirth = ""
	row := db.Table("birthdays").Where("name = ?", strings.ToLower(username)).Select("name, date_of_birth").Row()
	err = row.Scan(&username, &dateOfBirth)
	if err != nil {
		return dateOfBirth, err
	}
	return dateOfBirth, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
