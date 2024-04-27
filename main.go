package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "github.com/lib/pq" //import,but unuse
)

//Db on hard---->slow
//important: conn db --->close
func pg_conn(id int64) string {
	//cache redis ---->ram
	//TODO
	//hard
	connStr := "user=postgres dbname=dB2 password=246 sslmode=disable" //connection string
	db, err := sql.Open("postgres", connStr)                           //app---->pg
	if err != nil {
		panic(err)
	}

	var name sql.NullString
	if err := db.QueryRow("SELECT name FROM student WHERE id = $1", id).Scan(&name); err != nil {
		panic(err)

	}
	return name.String

}

//65535 ports per ip - 1024 reserved
func pong(c *gin.Context) {
	// data databse

	time.Sleep(200 + time.Millisecond)
	c.JSON(http.StatusOK, gin.H{
		"id":   "1",
		"name": "maryh",
	})

}
func main() {
	r := gin.Default()
	r.GET("/ping", pong)
	r.GET("/user/:id", func(c *gin.Context) { //rest api
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			panic(err)
		}
		// query

		if id != 0 {
			c.JSON(http.StatusOK, gin.H{
				"id":   id,
				"name": pg_conn(id),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": "Not found",
			})
		}

	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
