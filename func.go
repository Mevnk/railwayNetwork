package railwayNetwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

func TimeDiff(departure string, actualDeparture string) float64 {
	departure1, _ := time.Parse("15:04", departure)
	actualDeparture1, _ := time.Parse("15:04", actualDeparture)

	difference := actualDeparture1.Sub(departure1)

	return difference.Minutes()
}

func ClearBooked(route string) {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var total string
	err = db.QueryRow("select total_places_available from train where route_name = ?", route).Scan(&total)
	if err != nil {
		fmt.Println("Getting total places failed")
		return
	}

	var buf1 []byte
	err = db.QueryRow("select places_available from train where route_name = ?", route).Scan(&buf1)
	if err != nil {
		fmt.Println("Getting booked places failed")
		return
	}
	var schedule map[string]interface{}
	var scheduleEdit map[string]string
	scheduleEdit = make(map[string]string)
	err = json.Unmarshal(buf1, &schedule)
	if err != nil {
		fmt.Println("Route book reset failed")
		return
	}

	for key, _ := range schedule {
		scheduleEdit[key] = total
	}

	newJSON, _ := json.Marshal(scheduleEdit)
	db.QueryRow("update train set places_available = ? where route_name = ?", newJSON, route)

}
