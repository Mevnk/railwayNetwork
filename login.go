package railwayNetwork

import (
	"database/sql"
	"fmt"
	"hash/fnv"
	"strconv"
)

func (c *Driver) LoginWindow() int {
	var login, password string

	fmt.Printf("\nEnter login: ")
	_, err := fmt.Scan(&login)
	if err != nil {
		fmt.Println("\nIncorrect input")
		return 0
	}
	fmt.Printf("\nEnter password: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		fmt.Println("\nIncorrect input")
		return 0
	}
	h := fnv.New32a()
	_, err = h.Write([]byte(password))
	if err != nil {
		fmt.Println("\nHashing error")
		return 0
	}
	passwordHash := strconv.Itoa(int(h.Sum32()))

	loginAttmp, userID := LoginAction(login, passwordHash)
	if userID != -1 {
		c.userID = userID
	}

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
	//db.QueryRow("select role from user_role where user_id = (select id from client where login = ? and password_hash = ?)", login, pHash).Scan(&role)
	err = db.QueryRow("select id, role from client where login = ? and password_hash = ?", login, pHash).Scan(&id, &role)
	if err != nil {
		return 0, 0
	}
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
