package railwayNetwork

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"strconv"
)

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
	c.Actions[7] = c.StationWindow
	c.Actions[8] = c.AdminWindow
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

func (c Driver) CustomerWindow() int {
	prompt := promptui.Select{
		Label: "Select option",
		Items: []string{"Check schedule", "Book a ticket", "View your tickets", "Log out"},
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
	case "Log out":
		return 0

	}
	return 0
}

func (c *Driver) StationWindow() int {
	if !CheckStationPrivileges(c.userID) {
		return 0
	}
	if !CheckStationAssignment(c.userID) {
		fmt.Println("You are not assigned to a station")
		fmt.Println("Press any key to proceed...")
		var key string
		fmt.Scan(&key)
		return 0
	}

	prompt := promptui.Select{
		Label: "",
		Items: []string{"Report departure", "Log out"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}

	switch result {
	case "Report departure":
		var trainID, actualDeparture string
		fmt.Print("Enter route number: ")
		fmt.Scan(&trainID)
		fmt.Print("Enter actual departure time (format 00:00): ")
		fmt.Scan(&actualDeparture)
		fmt.Println("TEST0")
		ReportDeparture(trainID, actualDeparture, c.userID)
		return 7
	case "Log out":
		return 0
	}

	return 0
}

func (c *Driver) AdminWindow() int {
	if !CheckAdminPrivileges(c.userID) {
		return 0
	}

	prompt := promptui.Select{
		Label: "",
		Items: []string{"Elevate user", "Blacklist", "Routes", "Assign to station", "Log out"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}

	switch result {
	case "Elevate user":
		resp := c.ElevateUserWindow()
		return resp
	case "Blacklist":
		Blacklist()
		return 8
	case "Routes":
		RouteAdmin()
		return 8
	case "Assign to station":
		selectStation := promptui.Select{
			Label: "",
			Items: []string{"Kyiv", "Zaporizhzhya", "Dnipro", "Donetsk", "Return"},
		}
		_, selectStationResult, _ := selectStation.Run()
		if selectStationResult == "Return" {
			return 8
		}
		var pNumber string
		fmt.Printf("Enter user's passport: ")
		fmt.Scan(&pNumber)

		AssignToStation(pNumber, selectStationResult)

		return 8
	case "Log out":
		return 0
	}

	return 0
}
