package main

import "log"
import "golang.org/x/tools/go/loader"

//import "golang.org/x/tools/go/types"

func main() {
	srccode := `
	package main

	func main(){
	  log.Println("Hello")
	}
	`

	//fset := token.NewFileSet() // positions are relative to fset

	var conf loader.Config
	f, err := conf.ParseFile("thing.go", srccode)

	log.Println(f)
	log.Println(err)
}
