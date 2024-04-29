package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "github.com/lib/pq" //import,but unuse
)

type Restaurant struct {
	Name     string `json:"name"`
	distance string `json:"distance"`
	// Add more fields here as needed
}

//Db on hard---->slow
//important: conn db --->close
func pg_conn(lat, lng float64) string {
	//cache redis ---->ram
	//TODO
	//hard
	connStr := "user=postgres dbname=postgres password=246 sslmode=disable" //connection string
	db, err := sql.Open("postgres", connStr)                                //app---->pg
	if err != nil {
		panic(err)
	}

	query := `
	SELECT
	name,
	ST_Distance(
        ST_Transform(way, 4326),
        ST_SetSRID(ST_MakePoint($1,$2), 4326)
    ) AS distance
    --,osm_id, amenity,    
	FROM planet_osm_point
	where amenity = 'restaurant' AND 
	ST_DWithin(
		ST_Transform(way, 4326),
		ST_SetSRID(ST_MakePoint($1,$2), 4326),
		1000
	)
	ORDER BY distance
	limit 10;

	`
	rows, err := db.Query(query, lng, lat)
	if err != nil {
		panic(err)
	}
	defer rows.Next()
	//var name string
	//var distance string
	var restaurants []Restaurant
	for rows.Next() {
		var rst Restaurant
		if err := rows.Scan(&rst.Name, &rst.distance); err != nil {
			panic(err)
		}
		restaurants = append(restaurants, rst)
	}
	jsonData, err := json.Marshal(restaurants)
	if err != nil {
		panic(err)
	}
	return string(jsonData)

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
				"id": id,
				//"name": pg_conn(id),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": "Not found",
			})
		}

	})

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	r.GET("/point", func(c *gin.Context) {
		lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
		long, err2 := strconv.ParseFloat(c.Query("long"), 64) // shortcut for c.Request.URL.Query().Get("lastname")
		fmt.Print(lat, long)
		if err1 == nil || err2 == nil {
			c.JSON(http.StatusOK, gin.H{
				"res": pg_conn(lat, long),
			})
			//restaurants := pg_conn(lat, long)
			//c.JSON(http.StatusOK, restaurants)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": "lat or Long Not found",
			})
		}

		//c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
