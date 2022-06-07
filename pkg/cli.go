package pkg

import (
	"database/sql"
	"fmt"
	"github.com/manifoldco/promptui"
	"hash/fnv"
	"railwayNetwork"
	"strconv"
)

type Cli struct {
	Current int
	Actions map[int]func() int
}

func (c Cli) Init() {
	c.Actions[0] = c.Index
	c.Actions[1] = c.SignUp
}

func (c Cli) Show() {
	c.Actions[c.Current]()
}

func (c Cli) Index() int {
	prompt := promptui.Select{
		Label: "Select Day",
		Items: []string{"1. Sign up", "2. Log in"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}

	intVar, err := strconv.Atoi(result[0:1])
	return intVar
}

func (c Cli) SignUp() int {
	var login, password, fName, lName, pNumber string

	db, err := sql.Open("mysql", "root:misha26105@tcp(127.0.0.1:3306)/railway")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Success")
	}

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

	user := railwayNetwork.User{
		Login:        login,
		PasswordHash: passwordHash,
		FName:        fName,
		LName:        lName,
		PassportNum:  pNumber,
	}

	check, err := db.Query("select * from train")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Success")
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
	defer check.Close()
	return 1
}
