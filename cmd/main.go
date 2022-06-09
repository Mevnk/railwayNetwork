package main

import (
	_ "github.com/go-sql-driver/mysql"
	railway "railwayNetwork"
)

func main() {
	menu := railway.Cli{}
	resp := 4

	menu.Actions = make(map[int]func() int)
	menu.Init()

	for resp != -1 {
		resp = menu.Actions[resp]()
	}
}
