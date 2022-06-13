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
	err = db.QueryRow("select route from train where route_name = ?", routeName).Scan(&buf1)
	if err != nil {
		fmt.Println("SQL query failed")
		return false
	}
	var route []string
	err = json.Unmarshal(buf1, &route)
	if err != nil {
		fmt.Println("Unmarshaling failed")
		return false
	}

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
	err = db.QueryRow("select route from train where route_name = ?", routeName).Scan(&buf1)
	if err != nil {
		fmt.Println("SQL query failed")
		return false
	}
	var route []string
	err = json.Unmarshal(buf1, &route)
	if err != nil {
		fmt.Println("Unmarshaling failed")
		return false
	}

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
	err = db.QueryRow("select exists(select id from client where passport_number = ?)", passNum).Scan(&exists)
	if err != nil {
		fmt.Println("SQL query failed")
		return false
	}

	return exists
}

func CheckPlaceAvailable(route string, departure string, arrival string) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var buf1 []byte
	err = db.QueryRow("select places_available from train where route_name = ?", route).Scan(&buf1)
	if err != nil {
		fmt.Println("SQL query failed")
		return false
	}
	var schedule []string
	err = json.Unmarshal(buf1, &schedule)
	if err != nil {
		fmt.Println("Unmarshaling failed")
		return false
	}

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
	err = db.QueryRow("select route from train where route_name = ?", route).Scan(&buf1)
	if err != nil {
		fmt.Println("SQL query failed")
		return false
	}
	var schedule []string
	err = json.Unmarshal(buf1, &schedule)
	if err != nil {
		fmt.Println("Unmastshal failed")
		return false
	}

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
	err = db.QueryRow("select exists(select id from train where route_name = ?)", routeName).Scan(&exists)
	if err != nil {
		fmt.Println("SQL query failed")
		return false
	}

	return exists
}

func CheckAdminPrivileges(id int) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var role string
	err = db.QueryRow("select role from client where id = ?", id).Scan(&role)
	if err != nil {
		fmt.Println("SQL query failed")
		return false
	}

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
	err = db.QueryRow("select role from client where id = ?", id).Scan(&role)
	if err != nil {
		fmt.Println("SQL query failed")
		return false
	}

	if role == "station" {
		return true
	}

	return false
}
