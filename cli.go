package railwayNetwork

import (
	"database/sql"
	"fmt"
	"github.com/manifoldco/promptui"
	"hash/fnv"
	"strconv"
)

type Driver struct {
	LoggedIn bool
	userID   string
	role     string
	Actions  map[int]func() int
}

func (c Driver) Init() {
	c.LoggedIn = false
	c.role = ""

	c.Actions[0] = c.Index
	c.Actions[1] = c.SignUp
	c.Actions[2] = c.LoginWindow
	c.Actions[3] = c.StationSchedule
	c.Actions[4] = c.CustomerWindow
	c.Actions[5] = c.BookWindow
	c.Actions[6] = c.ViewTickets
}

//func (c Driver) Show() {
//	c.Actions[c.Current]()
//}

func (c Driver) Index() int {
	prompt := promptui.Select{
		Label: "Select Day",
		Items: []string{"1. Sign up", "2. Log in", "3. Check schedule"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}

	intVar, err := strconv.Atoi(result[0:1])
	return intVar
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

	fmt.Printf("Schedule role %s\n", c.role)
	fmt.Printf("Schedule login %t\n", c.LoggedIn)
	fmt.Println("Press any key to proceed...")
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

func (c Driver) SignUp() int {
	var login, password, fName, lName, pNumber string

	fmt.Printf("Enter login: ")
	fmt.Scan(&login)
	fmt.Printf("Enter password: ")
	fmt.Scan(&password)
	h := fnv.New32a()
	h.Write([]byte(password))
	passwordHash := strconv.Itoa(int(h.Sum32()))
	fmt.Printf("Enter your first name: ")
	fmt.Scan(&fName)
	fmt.Printf("Enter last name: ")
	fmt.Scan(&lName)
	fmt.Printf("Enter your passport number: ")
	fmt.Scan(&pNumber)

	if login != "" && password != "" {
		SignUpAction(login, passwordHash, fName, lName, pNumber)
	}
	return 0
}

func (c *Driver) LoginWindow() int {
	var login, password string

	fmt.Printf("Enter login: ")
	fmt.Scan(&login)
	fmt.Printf("Enter password: ")
	fmt.Scan(&password)
	h := fnv.New32a()
	h.Write([]byte(password))
	passwordHash := strconv.Itoa(int(h.Sum32()))

	loginAttmp, userID := LoginAction(login, passwordHash)
	c.userID = userID
	fmt.Println("TEST0 ", c.userID)

	fmt.Printf("Login1 role %s\n", c.role)
	fmt.Printf("Login1 login %t\n", c.LoggedIn)
	fmt.Println("Press any key to proceed...")
	var key string
	fmt.Scan(&key)

	if loginAttmp != 0 {
		c.LoggedIn = true

		switch loginAttmp {
		case 4:
			fmt.Printf("loginAttmp = %d", loginAttmp)
			c.role = "customer"
			break
		}
	}

	return loginAttmp
}

func (c Driver) CustomerWindow() int {
	prompt := promptui.Select{
		Label: "Select option",
		Items: []string{"Check schedule", "Book a ticket", "View your tickets"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}

	switch result {
	case "Check schedule":
		return 3
	case "Book a ticket":
		return 5
	case "View your tickets":
		return 6

	}
	return 0
}

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
