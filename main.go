package main

import (
	"valiDTr/cmd"
	"valiDTr/db"
)

func main() {
	db.InitDB()
	cmd.Execute()
}
