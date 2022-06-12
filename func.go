package railwayNetwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func CheckTimeCorrect(timeTest string) bool {
	_, err := time.Parse("15:04", timeTest)
	if err != nil {
		return false
	}
	return true
}

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
	var schedule []string
	err = json.Unmarshal(buf1, &schedule)
	if err != nil {
		fmt.Println("Route book reset failed")
		return
	}
	var station string
	for i := 0; i < len(schedule); i++ {
		station, _ = ParseJSONBookedPlaces(schedule[i])
		schedule[i] = station + ":" + total
	}

	newJSON, _ := json.Marshal(schedule)
	db.QueryRow("update train set places_available = ? where route_name = ?", newJSON, route)

}

func GetIDFromPassport(pNumber string) int {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	if !CheckUser(pNumber) {
		return -1
	}

	var id int
	db.QueryRow("select id from client where passport_number = ?", pNumber).Scan(&id)

	return id
}

func GetTotalPlaces(route string) int {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var total int
	db.QueryRow("select total_places_available from train where route_name = ?", route).Scan(&total)

	return total
}

func GetPlaceAvailable(route string, departure string, arrival string) int {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var buf1 []byte
	db.QueryRow("select places_available from train where route_name = ?", route).Scan(&buf1)
	var schedule []string
	json.Unmarshal(buf1, &schedule)

	var flag int
	min := GetTotalPlaces(route)
	var elementInt int
	var station, places string
	for i := 0; i < len(schedule); i++ {
		station, places = ParseJSONBookedPlaces(schedule[i])
		if station == departure {
			flag = 1
		}
		if station == arrival {
			flag = 0
		}
		if flag == 1 {
			elementInt, _ = strconv.Atoi(fmt.Sprintf("%v", places))
			if min > elementInt {
				min = elementInt
			}
		}
	}
	return min
}

func ParseJSONBookedPlaces(unparsed string) (string, string) {
	var station, places string
	for i := 0; i < len(unparsed); i++ {
		if unparsed[i] == ':' {
			station = unparsed[:i]
			places = unparsed[i+1:]
			break
		}
	}
	return station, places
}
