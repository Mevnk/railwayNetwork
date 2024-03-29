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
	err = db.QueryRow("select exists(select role from client where id = ?)", userID).Scan(&exists)
	if err != nil {
		fmt.Println("SQL query failed")
		return false
	}

	return exists
}

func ReportDeparture(routeID string, actualDeparture string, manager int) {

	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var station string
	err = db.QueryRow("select station_name from station where manager_id = ?", manager).Scan(&station)
	if err != nil {
		fmt.Println("SQL query failed")
		return
	}
	if CheckFinalStation(routeID, station) {
		ClearBooked(routeID)
	}

	var buf1 []byte
	err = db.QueryRow("select route from train where route_name = ?", routeID).Scan(&buf1)
	if err != nil {
		fmt.Println("SQL query failed")
		return
	}
	var departureParse []string
	err = json.Unmarshal(buf1, &departureParse)
	if err != nil {
		fmt.Println("Unmarshaling failed")
		return
	}
	var timeSelect string
	for i := 0; i < len(departureParse); i++ {
		stationParse, _ := ParseJSONBookedPlaces(departureParse[i])
		if stationParse == station {
			timeSelect = departureParse[i][len(departureParse[i])-5:]
			break
		}
	}
	delay := TimeDiff(fmt.Sprintf("%v", timeSelect), actualDeparture)
	db.QueryRow("update train set delay = ? where route_name = ?", delay, routeID)
}
