package railwayNetwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
)

func RouteAdmin() {
	prompt := promptui.Select{
		Label: "",
		Items: []string{"Add route", "Remove route"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {
	case "Add route":
		RouteAdd()
	case "Remove route":
		RouteRemove()
	}

}

func RouteAdd() {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var routeName string
	fmt.Printf("Enter number of the route: ")
	fmt.Scan(&routeName)

	if CheckRoute(routeName) {
		fmt.Println("A route by this number already exists")
		return
	}

	var totalPlaces string
	fmt.Printf("Enter total number of available places: ")
	fmt.Scan(&totalPlaces)

	var time string
	var routeMap map[string]string
	var routeBook map[string]string
	routeMap = make(map[string]string)
	routeBook = make(map[string]string)
	prompt := promptui.Select{
		Label: "Select station",
		Items: []string{"Kyiv", "Zaporizhzhya", "Dnipro", "Donetsk", "Finish"},
	}
	for true {
		fmt.Printf("Enter station name")
		_, station, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if station == "Finish" {
			break
		}

		fmt.Printf("Enter the time of departure from the station (format 00:00): ")
		fmt.Scan(&time)
		if !CheckTimeCorrect(time) {
			fmt.Println("Incorrect time format, try again...")
			continue
		}

		routeMap[station] = time
		routeBook[station] = totalPlaces
	}

	routeMapJSON, _ := json.Marshal(routeMap)
	routeBookJSON, _ := json.Marshal(routeBook)
	db.QueryRow("insert into train (total_places_available, route_name, route, places_available, delay) values (?, ?, ?, ?, 0)", totalPlaces, routeName, routeMapJSON, routeBookJSON)

	var trainID int
	db.QueryRow("select id from train where route_name = ?", routeName).Scan(&trainID)
	for key, value := range routeMap {
		db.QueryRow("insert into station (station_name, train_id, arrival_time) values (?, ?, ?)", key, trainID, value)
	}

}

func RouteRemove() {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var routeName string
	fmt.Printf("Enter number of the route: ")
	fmt.Scan(&routeName)

	if !CheckRoute(routeName) {
		fmt.Println("A route by this number doesn't exist")
		return
	}

	var routeID int
	db.QueryRow("select id from train where route_name = ?", routeName).Scan(&routeID)
	var buf1 []byte
	db.QueryRow("select route from train where route_name = ?", routeName).Scan(&buf1)
	var schedule map[string]interface{}
	json.Unmarshal(buf1, &schedule)

	for key, value := range schedule {
		db.QueryRow("delete from station where train_id = ? and station_name = ? and arrival_time = ?", routeID, key, value)
	}

	db.QueryRow("delete from train where id = ?", routeID)
}
