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
		fmt.Printf("\nPrompt failed %v\n", err)
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
	fmt.Printf("\nEnter number of the route: ")
	_, err = fmt.Scan(&routeName)
	if err != nil {
		fmt.Println("Input failed")
		return
	}

	if CheckRoute(routeName) {
		fmt.Println("A route by this number already exists")
		return
	}

	var totalPlaces string
	fmt.Printf("\nEnter total number of available places: ")
	_, err = fmt.Scan(&totalPlaces)
	if err != nil {
		fmt.Println("Input failed")
		return
	}

	var time string
	var routeMap []string
	var routeBook []string
	prompt := promptui.Select{
		Label: "Select station",
		Items: []string{"Kyiv", "Zaporizhzhya", "Dnipro", "Donetsk", "Finish"},
	}
	i := 0
	for true {
		fmt.Printf("\nEnter station name")
		_, station, err := prompt.Run()
		if err != nil {
			fmt.Printf("\nPrompt failed %v\n", err)
			return
		}

		if station == "Finish" {
			break
		}

		fmt.Printf("\nEnter the time of departure from the station (format 00:00): ")
		_, err = fmt.Scan(&time)
		if err != nil {
			fmt.Println("Input failed")
			return
		}
		if !CheckTimeCorrect(time) {
			fmt.Println("Incorrect time format, try again...")
			continue
		}

		routeMap = append(routeMap, station+":"+time)
		routeBook = append(routeBook, station+":"+totalPlaces)
		i++
	}

	//for key, _ := range routeMap {
	//	fmt.Println(routeMap[key])
	//}

	routeMapJSON, _ := json.Marshal(routeMap)
	routeBookJSON, _ := json.Marshal(routeBook)
	db.QueryRow("insert into train (total_places_available, route_name, route, places_available, delay) values (?, ?, ?, ?, 0)", totalPlaces, routeName, routeMapJSON, routeBookJSON)

	var trainID int
	err = db.QueryRow("select id from train where route_name = ?", routeName).Scan(&trainID)
	if err != nil {
		fmt.Println("SQL query failed")
		return
	}
	for i := 0; i < len(routeMap); i++ {
		db.QueryRow("insert into station (station_name, train_id, arrival_time) values (?, ?, ?)", routeMap[i][:len(routeMap[i])-6], trainID, routeMap[i][len(routeMap[i])-5:])
	}

}

func RouteRemove() {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var routeName string
	fmt.Printf("\nEnter number of the route: ")
	_, err = fmt.Scan(&routeName)
	if err != nil {
		fmt.Println("Input failed")
		return
	}

	if !CheckRoute(routeName) {
		fmt.Println("A route by this number doesn't exist")
		return
	}

	var routeID int
	err = db.QueryRow("select id from train where route_name = ?", routeName).Scan(&routeID)
	if err != nil {
		fmt.Println("Query failed")
		return
	}
	var buf1 []byte
	err = db.QueryRow("select route from train where route_name = ?", routeName).Scan(&buf1)
	if err != nil {
		fmt.Println("Query failed")
		return
	}
	var schedule []string
	err = json.Unmarshal(buf1, &schedule)
	if err != nil {
		fmt.Println("Unmarshaling failed")
		return
	}

	for i := 0; i < len(schedule); i++ {
		db.QueryRow("delete from station where train_id = ? and station_name = ? and arrival_time = ?", routeID, schedule[i][:len(schedule[i])-6], schedule[i][len(schedule[i])-5:])
	}

	db.QueryRow("delete from train where id = ?", routeID)
}
