package railwayNetwork

import (
	"database/sql"
	"fmt"
	"github.com/manifoldco/promptui"
)

func (c Driver) ElevateUserWindow() int {
	var newRole, userPassport string

	fmt.Print("Enter user's passport: ")
	fmt.Scan(&userPassport)

	prompt := promptui.Select{
		Label: "Select user's new role: ",
		Items: []string{"Customer", "Station manager", "Admin"},
	}

	_, role, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}

	if role == "Station manager" {
		newRole = "station"
	} else {
		newRole = role
	}

	ElevateUser(userPassport, newRole)

	return 8
}

func ElevateUser(user string, newRole string) {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var exists bool
	db.QueryRow("select exists(select id from client where passport_number = ?)", user).Scan(&exists)
	if !exists {
		fmt.Println("No user with such passport")
		return
	}

	db.QueryRow("update client set role = ? where passport_number = ?", newRole, user)
}

func AssignToStation(userPassport string, station string) {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	if !CheckUser(userPassport) {
		fmt.Println("This user doesn't exist")
		return
	}

	id := GetIDFromPassport(userPassport)

	db.QueryRow("update station set manager_id = ? where station_name = ?", id, station)

}
