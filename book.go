package railwayNetwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"strconv"
)

func (c Driver) BookWindow() int {
	prompt := promptui.Select{
		Label: "Select departure station",
		Items: []string{"Kyiv", "Zaporizhzhya", "Dnipro"},
	}

	_, DepStation, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}

	schedule := CheckScheduleAction(DepStation)
	for i, route := range schedule {
		fmt.Printf("%d: %s %s\n", i+1, route.RouteName, route.arrivalTime)
	}

	var route string
	fmt.Printf("Enter route number: ")
	fmt.Scan(&route)
	if !CheckDeparture(route, DepStation) {
		fmt.Println("Given route doesn't arrive on this station")
		return 4
	}

	_, ArrStation, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}
	if !CheckArrival(route, DepStation, ArrStation) {
		fmt.Println("This route doesn't lead to this station")
		return 4
	}

	available := CheckPlaceAvailable(route, DepStation, ArrStation)
	if !available {
		fmt.Printf("All is booked")
		return 4
	}

	var pNumber string
	fmt.Printf("Enter the passport number to book the ticket on: ")
	fmt.Scan(&pNumber)
	if !CheckUser(pNumber) {
		fmt.Println("There is no user with this passport in the database")
		return 4
	}

	Book(route, DepStation, ArrStation, pNumber)

	return 4
}

func (c *Driver) ViewTickets() int {
	var tickets []Ticket

	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("select train_id, departure_station_id, arrival_station_id from ticket where user_id = ?", c.userID)

	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(&ticket.train, &ticket.departure, &ticket.arrival); err != nil {
			db.Close()
			panic(err)
		}

		tickets = append(tickets, ticket)
	}

	for i := 0; i < len(tickets); i++ {
		db.QueryRow("select route_name from train where id = ?", tickets[i].train).Scan(&(tickets[i].train))
		db.QueryRow("select station_name from station where id = ?", tickets[i].departure).Scan(&(tickets[i].departure))
		db.QueryRow("select station_name from station where id = ?", tickets[i].arrival).Scan(&(tickets[i].arrival))
	}
	fmt.Println("TEST1 ", c.userID)
	for i := 0; i < len(tickets); i++ {
		fmt.Println("Route number: ", tickets[i].train)
		fmt.Println("Departure station: ", tickets[i].departure)
		fmt.Println("Arrival station: ", tickets[i].arrival)
	}
	fmt.Println("TEST2")
	fmt.Println("Press any key to proceed...")
	var key string
	fmt.Scan(&key)

	return 4
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
