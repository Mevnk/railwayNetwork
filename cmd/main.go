package main

import (
	_ "github.com/go-sql-driver/mysql"
	"railwayNetwork/pkg"
)

func main() {
	menu := pkg.Cli{}
	resp := 0

	menu.Actions = make(map[int]func() int)
	menu.Init()

	for resp != -1 {
		resp = menu.Actions[resp]()
	}
}
