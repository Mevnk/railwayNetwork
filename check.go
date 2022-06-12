package railwayNetwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func CheckDeparture(routeName string, station string) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var buf1 []byte
	db.QueryRow("select route from train where route_name = ?", routeName).Scan(&buf1)
	var route []string
	json.Unmarshal(buf1, &route)

	exists := false
	for i := 0; i < len(route); i++ {
		if route[i][:len(route[i])-6] == station {
			exists = true
			break
		}
	}

	return exists
}

func CheckArrival(routeName string, departure string, arrival string) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var buf1 []byte
	db.QueryRow("select route from train where route_name = ?", routeName).Scan(&buf1)
	var route []string
	json.Unmarshal(buf1, &route)

	exists := false
	flag := 0
	for i := 0; i < len(route); i++ {
		if route[i][:len(route[i])-6] == departure {
			flag = 1
			continue
		}
		if flag == 1 {
			if route[i][:len(route[i])-6] == arrival {
				exists = true
				break
			}
		}
	}

	return exists
}

func CheckUser(passNum string) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var exists bool
	db.QueryRow("select exists(select id from client where passport_number = ?)", passNum).Scan(&exists)

	return exists
}

func CheckPlaceAvailable(route string, departure string, arrival string) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var buf1 []byte
	db.QueryRow("select places_available from train where route_name = ?", route).Scan(&buf1)
	var schedule []string
	json.Unmarshal(buf1, &schedule)

	var flag int
	var key, value string
	for i := 0; i < len(schedule); i++ {
		key, value = ParseJSONBookedPlaces(schedule[i])
		fmt.Println("TEST4.11")
		if key == departure {
			flag = 1
		}
		if key == arrival {
			flag = 0
		}
		if flag == 1 {
			if value == "0" {
				return false
			}
		}
	}
	return true
}

func CheckFinalStation(route string, station string) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var buf1 []byte
	db.QueryRow("select route from train where route_name = ?", route).Scan(&buf1)
	var schedule []string
	json.Unmarshal(buf1, &schedule)

	var keys []string
	for i := 0; i < len(schedule); i++ {
		keys = append(keys, schedule[i][:len(schedule[i])-6])
	}

	if keys[len(keys)-1] == station {
		return true
	}
	return false
}

func CheckRoute(routeName string) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var exists bool
	db.QueryRow("select exists(select id from train where route_name = ?)", routeName).Scan(&exists)

	return exists
}

func CheckAdminPrivileges(id int) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var role string
	db.QueryRow("select role from client where id = ?", id).Scan(&role)

	if role == "admin" {
		return true
	}

	return false
}

func CheckStationPrivileges(id int) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var role string
	db.QueryRow("select role from client where id = ?", id).Scan(&role)

	if role == "station" {
		return true
	}

	return false
}
