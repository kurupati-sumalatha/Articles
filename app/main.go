package main

import (
	"art/article/delivery/http/middleware"
	repo "art/article/repository/mysql"
	"art/article/usecase"
	"art/author/repository/mysql"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	//"pkg/mod/github.com/spf13/viper@v1.8.0"
	"art/article/delivery/http"
	"time"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

var db *sql.DB

func main() {
	db, err := sql.Open("mysql", "root:Root@1234#@tcp(127.0.0.1:3306)/sakila")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	/*dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	*/
	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)
	authorRepo := mysql.NewMysqlAuthorRepository(db)
	ar := repo.NewMysqlArticleRepository(db)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := usecase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	http.NewArticleHandler(e, au)

	log.Fatal(e.Start(viper.GetString("server.address"))) //nolint
}
