package railwayNetwork

import (
	"database/sql"
	"fmt"
	"hash/fnv"
	"strconv"
)

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
	if userID != -1 {
		c.userID = userID
	}
	fmt.Println("Press any key to proceed...")
	var key string
	fmt.Scan(&key)

	if loginAttmp != 0 {
		c.LoggedIn = true

		switch loginAttmp {
		case 4:
			c.role = "customer"
			break
		case 7:
			c.role = "station"
			break
		case 8:
			c.role = "admin"
			break
		}
	}

	return loginAttmp
}

func LoginAction(login string, pHash string) (int, int) {
	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var exists bool
	err = db.QueryRow("select exists(select id from client where login = ? and password_hash = ?)", login, pHash).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("error checking if row exists")
	}

	if !exists {
		fmt.Println("Incorrect login or password")
		return 0, -1
	}

	var role string
	var id int
	db.QueryRow("select role from user_role where user_id = (select id from client where login = ? and password_hash = ?)", login, pHash).Scan(&role)
	db.QueryRow("select id from client where login = ? and password_hash = ?", login, pHash).Scan(&id)
	switch role {
	case "customer":
		return 4, id
	case "station":
		return 7, id
	case "admin":
		return 8, id
	}

	return 0, -1

}