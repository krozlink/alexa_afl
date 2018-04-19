package main

import (
	"fmt"
	"github.com/Nyarum/betting"
)

func main() {
	bet := betting.NewBet("test")
	err := bet.GetSession("", "", "", "")

	if err != nil {
		fmt.Printf(err.Error())
	}
}
