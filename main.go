package main

import (
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	a := App{}
	a.Init()
	a.Start()

}
