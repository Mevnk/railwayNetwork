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
		fmt.Println("This user already exists")
		return
	}

	fmt.Println("Continue")
	var key string
	fmt.Scan(&key)

	user := User{
		Login:        login,
		PasswordHash: passHash,
		FName:        fName,
		LName:        lName,
		PassportNum:  pNumber,
	}
	q := "INSERT INTO `client` (login, password_hash, first_name, last_name, passport_number) VALUES (?, ?, ?, ?, ?);"
	insert, err := db.Prepare(q)
	if err != nil {
		fmt.Println(err)
	}
	insert.Exec(user.Login, user.PasswordHash, user.FName, user.LName, user.PassportNum)

	var userID int
	db.QueryRow("select id from client where passport_number = ?", pNumber).Scan(&userID)
	db.QueryRow("insert into user_role (user_id, role) values (?, 'customer')", strconv.Itoa(userID))

	db.Close()
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
