package railwayNetwork

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"hash/fnv"
	"strconv"
)

type Cli struct {
	LoggedIn bool
	role     string
	Actions  map[int]func() int
}

func (c Cli) Init() {
	c.LoggedIn = false
	c.role = ""

	c.Actions[0] = c.Index
	c.Actions[1] = c.SignUp
	c.Actions[2] = c.LoginWindow
	c.Actions[3] = c.StationSchedule
	c.Actions[4] = c.CustomerWindow
	c.Actions[5] = c.BookWindow
}

//func (c Cli) Show() {
//	c.Actions[c.Current]()
//}

func (c Cli) Index() int {
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

func (c *Cli) StationSchedule() int {
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
	for i, route := range schedule {
		fmt.Printf("%d: %s %s\n", i+1, route.RouteName, route.arrivalTime)
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

func (c Cli) SignUp() int {
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

func (c *Cli) LoginWindow() int {
	var login, password string

	fmt.Printf("Enter login: ")
	fmt.Scan(&login)
	fmt.Printf("Enter password: ")
	fmt.Scan(&password)
	h := fnv.New32a()
	h.Write([]byte(password))
	passwordHash := strconv.Itoa(int(h.Sum32()))

	loginAttmp := LoginAction(login, passwordHash)

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

	fmt.Printf("Login2 role %s\n", c.role)
	fmt.Printf("Login2 login %t\n", c.LoggedIn)
	fmt.Println("Press any key to proceed...")
	fmt.Scan(&key)

	return loginAttmp
}

func (c Cli) CustomerWindow() int {
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

	}
	return 0
}

func (c Cli) BookWindow() int {
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
	fmt.Println("BEFORECHECK")
	if !CheckArrival(route, DepStation, ArrStation) {
		fmt.Println("This route doesn't lead to this station")
		return 4
	}

	fmt.Println("ALLCHECK")

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
