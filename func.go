package railwayNetwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

func SignUpAction(
	login string,
	passHash string,
	fName string,
	lName string,
	pNumber string) {

	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var exists bool
	err = db.QueryRow("select exists(select passport_number from client where passport_number = ?)", pNumber).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("error checking if row exists")
	}

	if exists {
		fmt.Println("This user already exists")
		return
	}

	fmt.Println("Continue")
	var key string
	fmt.Scan(&key)

	user := User{
		Login:        login,
		PasswordHash: passHash,
		FName:        fName,
		LName:        lName,
		PassportNum:  pNumber,
	}
	q := "INSERT INTO `client` (login, password_hash, first_name, last_name, passport_number) VALUES (?, ?, ?, ?, ?);"
	insert, err := db.Prepare(q)
	if err != nil {
		fmt.Println(err)
	}
	insert.Exec(user.Login, user.PasswordHash, user.FName, user.LName, user.PassportNum)

	var userID int
	db.QueryRow("select id from client where passport_number = ?", pNumber).Scan(&userID)
	db.QueryRow("insert into user_role (user_id, role) values (?, 'customer')", strconv.Itoa(userID))

	db.Close()
}

func CheckScheduleAction(stationName string) []Route {
	var routes []Route

	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("select train_id, arrival_time from station where station_name = ?", stationName)

	for rows.Next() {
		var route Route
		if err := rows.Scan(&route.RouteID, &route.arrivalTime); err != nil {
			db.Close()
			panic(err)
		}
		db.QueryRow("select route from train where id = ?", route.RouteID).Scan(&route.Stops)
		routes = append(routes, route)
	}

	//db.QueryRow("select route_name from train where id = ?", routes[0].RouteID).Scan(&(routes[0].RouteName))
	//fmt.Printf("%s %s\n", routes[0].RouteName, routes[0].arrivalTime)

	for i := 0; i < len(routes); i++ {
		db.QueryRow("select route_name from train where id = ?", routes[i].RouteID).Scan(&(routes[i].RouteName))
	}

	db.Close()
	return routes
}

func LoginAction(login string, pHash string) int {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var exists bool
	err = db.QueryRow("select exists(select id from client where login = ? and password_hash = ?)", login, pHash).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("error checking if row exists")
	}

	if !exists {
		fmt.Println("Incorrect login or password")
		return 0
	}

	var role string
	db.QueryRow("select role from user_role where user_id = (select id from client where login = ? and password_hash = ?)", login, pHash).Scan(&role)
	fmt.Printf("role: %s\n", role)
	switch role {
	case "customer":
		return 4
	}

	return 0

}

func Book(route string, departure string, arrival string, passNum string) {
	var routeID, departureID, arrivalID, userID int

	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	db.QueryRow("select train.id from train inner join station s on train.id = s.train_id where route_name = ?", route).Scan(&routeID)
	db.QueryRow("select station.id from station where station_name = ?", departure).Scan(&departureID)
	db.QueryRow("select station.id from station where station_name = ?", arrival).Scan(&arrivalID)
	db.QueryRow("select client.id from client where passport_number = ?", passNum).Scan(&userID)

	db.QueryRow("insert into ticket (user_id, train_id, departure_station_id, arrival_station_id)  values (?, ?, ?, ?)", userID, routeID, departureID, arrivalID)

	var buf1 []byte
	db.QueryRow("select places_available from train where route_name = ?", route).Scan(&buf1)
	var schedule map[string]interface{}
	var scheduleEdit map[string]string
	scheduleEdit = make(map[string]string)
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
			bufInt, _ := strconv.Atoi(fmt.Sprintf("%v", element))
			scheduleEdit[key] = strconv.Itoa(bufInt - 1)
			continue
		}
		bufInt, _ := strconv.Atoi(fmt.Sprintf("%v", element))
		scheduleEdit[key] = strconv.Itoa(bufInt)
	}
	fmt.Println(scheduleEdit["Kyiv"])
	fmt.Println("Press any key to proceed...")
	var key string
	fmt.Scan(&key)
	newJSON, _ := json.Marshal(scheduleEdit)

	db.QueryRow("update train set places_available = ? where id = ?", newJSON, routeID)
}
