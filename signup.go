package railwayNetwork

import (
	"database/sql"
	"fmt"
	"hash/fnv"
	"strconv"
)

func SignUpAction(
	login string,
	passHash string,
	fName string,
	lName string,
	pNumber string) {

	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	}

	var exists bool
	err = db.QueryRow("select exists(select passport_number from client where passport_number = ?)", pNumber).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("error checking if row exists")
	}

	if exists {
		var checkRole string
		err := db.QueryRow("select role from client where passport_number = ?", pNumber).Scan(&checkRole)
		if err != nil {
			fmt.Println("\nSQL query failed")
			return
		}
		if checkRole == "customer" {
			fmt.Println("This user already exists")
			return
		}
	}

	fmt.Println("Press any key to continue")
	var key string
	_, err = fmt.Scan(&key)
	if err != nil {
		return
	}

	db.QueryRow("insert into client (login, password_hash, first_name, last_name, passport_number, role) values (?, ?, ?, ?, ?, 'customer')", login, passHash, fName, lName, pNumber)
}

func (c Driver) SignUp() int {
	var login, password, fName, lName, pNumber string

	fmt.Printf("Enter login: ")
	_, err := fmt.Scan(&login)
	if err != nil {
		fmt.Println("\nIncorrect input")
		return 0
	}
	fmt.Printf("Enter password: ")
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
	fmt.Printf("Enter your first name: ")
	_, err = fmt.Scan(&fName)
	if err != nil {
		fmt.Println("\nIncorrect input")
		return 0
	}
	fmt.Printf("Enter last name: ")
	_, err = fmt.Scan(&lName)
	if err != nil {
		fmt.Println("\nIncorrect input")
		return 0
	}
	fmt.Printf("Enter your passport number: ")
	_, err = fmt.Scan(&pNumber)
	if err != nil {
		fmt.Println("\nIncorrect input")
		return 0
	}

	if login != "" && password != "" {
		SignUpAction(login, passwordHash, fName, lName, pNumber)
	}
	return 0
}
