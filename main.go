package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "github.com/lib/pq" //import,but unuse
)

type Restaurant struct {
	Name     string
	distance string
}

func pg_conn(lat, lng float64) []Restaurant {

	connStr := "user=postgres dbname=*** password=*** sslmode=disable" //connection string
	db, err := sql.Open("postgres", connStr)                                //app---->pg
	if err != nil {
		panic(err)
	}
	//query to find top 10 restaurants within a 1km radius of a specified geographical point, ordered based on their distances.
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
	//executing the query with long and lat inputs
	rows, err := db.Query(query, lng, lat)
	if err != nil {
		panic(err)
	}
	defer rows.Next()
	//Append each SQL row to a collection of restaurants.
	var restaurants []Restaurant
	for rows.Next() {
		var rst Restaurant
		if err := rows.Scan(&rst.Name, &rst.distance); err != nil {
			panic(err)
		}
		restaurants = append(restaurants, rst)
	}
	//jsonData, err := json.Marshal(restaurants)
	if err != nil {
		panic(err)
	}
	return restaurants

}

func pong(c *gin.Context) {
	//c.JSON(http.StatusOK, gin.H{
	//	"ping": "pong",
	//})
	c.String(http.StatusOK, "pong")

}
func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/ping", pong)

	//Zafaraniyeh Office:   lat=35.803900&long=51.420431
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /http://localhost:8080/restaurants?lat=35.803900&long=51.420431

	r.GET("/restaurants", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)

		//converting string type to required float64 type
		lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
		long, err2 := strconv.ParseFloat(c.Query("long"), 64)

		if err1 == nil || err2 == nil {

			//get required data from postgres osm database
			restaurants := pg_conn(lat, long)

			//show the result with the help of restaurants.html template
			c.HTML(http.StatusOK, "restaurants.html", restaurants)

		} else {
			c.String(http.StatusOK, "Please enter latitude and longitude.")
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
