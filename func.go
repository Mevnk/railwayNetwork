package railwayNetwork

import (
	"database/sql"
	"fmt"
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
	} else {
		fmt.Printf("Success")
	}

	check := db.QueryRow("select exists(select passport_number from client where passport_number = ?)", pNumber)
	if check != nil {
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

	resp, err := insert.Exec(user.Login, user.PasswordHash, user.FName, user.LName, user.PassportNum)
	insert.Close()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
