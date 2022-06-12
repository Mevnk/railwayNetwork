package railwayNetwork

import (
	"database/sql"
	"fmt"
	"github.com/manifoldco/promptui"
)

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

func (c *Driver) StationSchedule() int {
	prompt := promptui.Select{
		Label: "Select station",
		Items: []string{"Kyiv", "Zaporizhzhya", "Dnipro"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}

	schedule := CheckScheduleAction(result)
	for _, route := range schedule {
		fmt.Println("////////////////////////////")
		fmt.Println("Route number: ", route.RouteName)
		fmt.Println("Train arrival: ", route.arrivalTime)
		fmt.Println("Train route: ", route.Stops)
	}

	fmt.Println("Press any key to proceed...")
	var key string
	fmt.Scan(&key)
	if !c.LoggedIn {
		return 0
	} else {
		switch c.role {
		case "customer":
			return 4
		}
	}
	return 0
}
