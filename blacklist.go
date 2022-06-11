package railwayNetwork

import (
	"database/sql"
	"fmt"
	"github.com/manifoldco/promptui"
)

func Blacklist() {
	prompt := promptui.Select{
		Label: "",
		Items: []string{"Add to blacklist", "Remove from blacklist"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {
	case "Add to blacklist":
		var pNumber string
		fmt.Printf("Enter user's passport number: ")
		fmt.Scan(&pNumber)
		userID := GetIDFromPassport(pNumber)
		BlacklistAdd(userID)
	case "Remove from blacklist":
		var pNumber string
		fmt.Printf("Enter user's passport number: ")
		fmt.Scan(&pNumber)
		userID := GetIDFromPassport(pNumber)
		BlacklistRemove(userID)
	}

}

func BlacklistAdd(userID int) {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	db.QueryRow("insert into blacklist (user_id) values (?)", userID)
}

func BlacklistRemove(userID int) {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var exists bool
	db.QueryRow("select exists(select user_id from blacklist where user_id = ?)", userID).Scan(&exists)
	if !exists {
		return
	}

	db.QueryRow("delete from blacklist where user_id = ?", userID)
}
