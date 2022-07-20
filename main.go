package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Form struct {
	Name  string
	email string
	Files []byte
}

func main() {
	e := echo.New()
	db := DbConnect()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "public")
	e.POST("/upload", db.upload)
	e.GET("/getformdata", db.getFormData)

	e.Logger.Fatal(e.Start(":1323"))
}

func (db *DBConnect) upload(c echo.Context) error {

	name := c.FormValue("name")
	email := c.FormValue("email")

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
		fileData, err := file.Open()
		if err != nil {
			return err
		}
		defer fileData.Close()

		file_data, err := ioutil.ReadAll(fileData)
		if err != nil {
			return err
		}

		sql := "INSERT INTO form_tab(Name,email,Files) VALUES(?,?,?)"
		data, err := db.conn.Prepare(sql)
		if err != nil {

			fmt.Print(err.Error())
		}
		defer data.Close()

		_, err2 := data.Exec(name, email, file_data)

		if err2 != nil {
			log.Println("this is an error", err2)
			panic(err2)

		}

	}

	return c.JSON(http.StatusOK, "success")
}

type DBConnect struct {
	conn *sql.DB
}

func DbConnect() *DBConnect {
	DB, err := sql.Open("mysql", "database:Aswad_database@123@tcp(127.0.0.1:3306)/form_db")
	if err != nil {
		panic(err.Error())
	}

	err = DB.Ping()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("connected")
	}

	return &DBConnect{
		conn: DB,
	}

}

func (db *DBConnect) getFormData(c echo.Context) error {
	ID := c.QueryParam("ID")
	fmt.Println(ID)
	rows, err := db.conn.Query("SELECT * FROM form_tab Where ID=? ", ID)
	if err != nil {

		panic(err)
	}

	var name string
	var emailid string
	var files []byte

	var data []Form
	for rows.Next() {
		err = rows.Scan(&name, &emailid, &files)
		if err != nil {
			panic(err)
		}
		data = append(data, Form{Name: name, email: emailid, Files: files})
	}
	fmt.Println(data)
	return c.JSON(http.StatusOK, "successs")
}
