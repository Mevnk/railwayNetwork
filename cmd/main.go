package main

import (
	_ "github.com/go-sql-driver/mysql"
	railway "railwayNetwork"
)

func main() {
	menu := railway.Driver{}
	resp := 0

	menu.Actions = make(map[int]func() int)
	menu.Init()

	for resp != -1 {
		resp = menu.Actions[resp]()
	}
}
