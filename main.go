package main

import (
	"fmt"

	"github.com/wastewater-intelligence-network/win-api/core"
)

func main() {
	fmt.Print("main")
	app, err := core.NewWinApp()
	if err != nil {
		panic(err)
	}
	app.Run()
}
