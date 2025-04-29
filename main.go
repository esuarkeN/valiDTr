package main

import (
	"fmt"
	"os"
	"valiDTr/cmd"
	"valiDTr/db"
)

func main() {
	db.InitDB()
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
