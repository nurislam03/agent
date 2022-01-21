package main

import (
	"github.com/joho/godotenv"
	"github.com/nurislam03/agent/cmd"
)

func main() {
	_ = godotenv.Load(".env")
	cmd.RootCmd.Execute()
}
