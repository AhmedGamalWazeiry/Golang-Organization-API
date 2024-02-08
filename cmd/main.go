package main

import "org.com/org/pkg"
func main() {
	app := pkg.NewApp()
	app.Run(":8080")
}