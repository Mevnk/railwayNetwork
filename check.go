package railwayNetwork

import (
	"database/sql"
	"encoding/json"
)

func CheckDeparture(routeName string, station string) bool {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var buf1 []byte
	db.QueryRow("select route from train where route_name = ?", routeName).Scan(&buf1)
	var route map[string]interface{}
	json.Unmarshal(buf1, &route)

	exists := false
	for key, _ := range route {
		if key == station {
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
	var route map[string]interface{}
	json.Unmarshal(buf1, &route)

	exists := false
	flag := 0
	for key, _ := range route {
		if key == departure {
			flag = 1
			continue
		}
		if flag == 1 {
			if key == arrival {
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
	var schedule map[string]interface{}
	json.Unmarshal(buf1, &schedule)

	var flag int
	for key, element := range schedule {
		if key == departure {
			flag = 1
		}
		if key == arrival {
			flag = 0
		}
		if flag == 1 {
			if element == "0" {
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
	var schedule map[string]interface{}
	json.Unmarshal(buf1, &schedule)

	keys := make([]string, len(schedule))
	i := 0
	for k := range schedule {
		keys[i] = k
		i++
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
