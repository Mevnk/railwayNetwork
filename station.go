package railwayNetwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func CheckStationAssignment(userID int) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var exists bool
	db.QueryRow("select exists(select role from user_role where user_id = ?)", userID).Scan(&exists)

	return exists
}

func ReportDeparture(routeID string, actualDeparture string, manager int) {

	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var station string
	db.QueryRow("select station_name from station where manager_id = ?", manager).Scan(&station)

	if CheckFinalStation(routeID, station) {
		ClearBooked(routeID)
		return
	}

	var buf1 []byte
	db.QueryRow("select route from train where route_name = ?", routeID).Scan(&buf1)
	var departureParse map[string]interface{}
	json.Unmarshal(buf1, &departureParse)
	fmt.Printf("Station: %v", departureParse)
	delay := TimeDiff(fmt.Sprintf("%v", departureParse[station]), actualDeparture)

	db.QueryRow("update train set delay = ? where route_name = ?", delay, routeID)
}
