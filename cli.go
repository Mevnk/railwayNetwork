package railwayNetwork

import (
	"fmt"
	"github.com/manifoldco/promptui"
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
		Items: []string{"Sign up", "Log in", "Check schedule", "Exit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}

	if result == "Exit" {
		return -1
	}

	switch result {
	case "Sign up":
		return 1
	case "Log in":
		return 2
	case "Check schedule":
		return 3
	case "Exit":
		return -1
	}
	return 0
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
		_, err := fmt.Scan(&key)
		if err != nil {
			return 0
		}
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
		_, err := fmt.Scan(&trainID)
		if err != nil {
			fmt.Println("\nUnexpected input")
			return 7
		}
		fmt.Print("Enter actual departure time (format 00:00): ")
		_, err = fmt.Scan(&actualDeparture)
		if err != nil {
			fmt.Println("\nUnexpected input")
			return 7
		}
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
		_, err := fmt.Scan(&pNumber)
		if err != nil {
			fmt.Println("\nUnexpected input")
			return 8
		}

		AssignToStation(pNumber, selectStationResult)

		return 8
	case "Log out":
		return 0
	}

	return 0
}
