package railwayNetwork

import (
	"database/sql"
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
