package railwayNetwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"strconv"
)

func (c *Driver) BookWindow() int {
	if CheckBlacklist(c.userID) {
		fmt.Println("You are blacklisted")
		return 4
	}

	promptDeparture := promptui.Select{
		Label: "Select departure station",
		Items: []string{"Kyiv", "Zaporizhzhya", "Dnipro", "Donetsk", "Finish"},
	}

	_, DepStation, err := promptDeparture.Run()

	if err != nil {
		fmt.Printf("\nPrompt failed %v\n", err)
		return -1
	}

	promptArrival := promptui.Select{
		Label: "Select arrival station",
		Items: []string{"Kyiv", "Zaporizhzhya", "Dnipro", "Donetsk", "Finish"},
	}
	_, ArrStation, err := promptArrival.Run()
	if err != nil {
		fmt.Printf("\nPrompt failed %v\n", err)
		return -1
	}
	schedule := CheckScheduleAction(DepStation)
	if schedule == nil {
		return 4
	}
	for i, route := range schedule {
		totalAvailable := GetPlaceAvailable(route.RouteName, DepStation, ArrStation)
		if totalAvailable == -1 {
			continue
		}
		fmt.Println("//////////////////////")
		fmt.Println("#", i+1)
		fmt.Println("Route #", route.RouteName)
		fmt.Println("Arrival time: ", route.arrivalTime)
		fmt.Println("Places available: ", totalAvailable)
		fmt.Println("Stops: ", route.Stops)
	}
	promptTotal := promptui.Select{
		Label: "",
		Items: []string{"Continue", "Return"},
	}
	var route string
	fmt.Printf("\nEnter route number: ")
	_, err1 := fmt.Scan(&route)
	if err1 != nil {
		fmt.Println("Invalid route")
		return 4
	}
	if !CheckDeparture(route, DepStation) {
		fmt.Println("Given route doesn't arrive on this station")
		return 4
	}
	if !CheckArrival(route, DepStation, ArrStation) {
		fmt.Println("This route doesn't lead to this station")
		return 4
	}
	available := CheckPlaceAvailable(route, DepStation, ArrStation)
	if !available {
		fmt.Printf("\nAll is booked")
		return 4
	}
	_, selectTotal, _ := promptTotal.Run()
	if selectTotal == "Return" {
		return 4
	}

	var pNumber string
	fmt.Printf("\nEnter the passport number to book the ticket on: ")
	_, err2 := fmt.Scan(&pNumber)
	if err2 != nil {
		fmt.Println("Invalid passport")
		return 4
	}
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
			panic(err)
		}

		tickets = append(tickets, ticket)
	}

	for i := 0; i < len(tickets); i++ {
		err := db.QueryRow("select route_name from train where id = ?", tickets[i].train).Scan(&(tickets[i].train))
		if err != nil {
			fmt.Println("SQL query failed")
			return 4
		}
		err = db.QueryRow("select station_name from station where id = ?", tickets[i].departure).Scan(&(tickets[i].departure))
		if err != nil {
			fmt.Println("SQL query failed")
			return 4
		}
		err = db.QueryRow("select station_name from station where id = ?", tickets[i].arrival).Scan(&(tickets[i].arrival))
		if err != nil {
			fmt.Println("SQL query failed")
			return 4
		}
	}
	for i := 0; i < len(tickets); i++ {
		fmt.Println("Route number: ", tickets[i].train)
		fmt.Println("Departure station: ", tickets[i].departure)
		fmt.Println("Arrival station: ", tickets[i].arrival)
		fmt.Println("//////////////////////////////")
	}
	fmt.Println("Press any key to proceed...")
	var key string
	_, err = fmt.Scan(&key)
	if err != nil {
		return 4
	}

	return 4
}

func Book(route string, departure string, arrival string, passNum string) {
	var routeID, departureID, arrivalID, userID int

	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	err = db.QueryRow("select train.id from train inner join station s on train.id = s.train_id where route_name = ?", route).Scan(&routeID)
	if err != nil {
		fmt.Println("SQL query failed")
		return
	}
	err = db.QueryRow("select station.id from station where station_name = ? and train_id = ?", departure, routeID).Scan(&departureID)
	if err != nil {
		fmt.Println("SQL query failed")
		return
	}
	err = db.QueryRow("select station.id from station where station_name = ? and train_id = ?", arrival, routeID).Scan(&arrivalID)
	if err != nil {
		fmt.Println("SQL query failed")
		return
	}
	err = db.QueryRow("select client.id from client where passport_number = ?", passNum).Scan(&userID)
	if err != nil {
		fmt.Println("SQL query failed")
		return
	}

	db.QueryRow("insert into ticket (user_id, train_id, departure_station_id, arrival_station_id)  values (?, ?, ?, ?)", userID, routeID, departureID, arrivalID)

	var buf1 []byte
	err = db.QueryRow("select places_available from train where route_name = ?", route).Scan(&buf1)
	if err != nil {
		fmt.Println("SQL query failed")
		return
	}
	var schedule []string
	var scheduleEdit []string
	err = json.Unmarshal(buf1, &schedule)
	if err != nil {
		fmt.Println("Unmarshaling query failed")
		return
	}

	var flag int
	var station, placeN string
	for i := 0; i < len(schedule); i++ {
		station, placeN = ParseJSONBookedPlaces(schedule[i])
		if station == departure {
			flag = 1
		}
		if station == arrival {
			flag = 0
		}
		if flag == 1 {
			bufInt, _ := strconv.Atoi(fmt.Sprintf("%v", placeN))
			scheduleEdit = append(scheduleEdit, station+":"+strconv.Itoa(bufInt-1))
			continue
		}
		bufInt, _ := strconv.Atoi(fmt.Sprintf("%v", placeN))
		scheduleEdit = append(scheduleEdit, station+":"+strconv.Itoa(bufInt))
	}
	newJSON, _ := json.Marshal(scheduleEdit)

	db.QueryRow("update train set places_available = ? where id = ?", newJSON, routeID)
}
