package railwayNetwork

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"hash/fnv"
	"strconv"
)

type Cli struct {
	Current int
	Actions map[int]func() int
}

func (c Cli) Init() {
	c.Actions[0] = c.Index
	c.Actions[1] = c.SignUp
	c.Actions[3] = StationSchedule
}

func (c Cli) Show() {
	c.Actions[c.Current]()
}

func (c Cli) Index() int {
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

func StationSchedule() int {
	prompt := promptui.Select{
		Label: "Select station",
		Items: []string{"Kyiv", "Zaporizhzhya", "Dnipro"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}

	schedule := CheckScheduleAction(result)
	for i, route := range schedule {
		fmt.Printf("%d: %s %s\n", i+1, route.RouteName, route.arrivalTime)
	}

	fmt.Println("Press any key to proceed...")
	var key string
	fmt.Scan(&key)

	return 0

}

func (c Cli) SignUp() int {
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
