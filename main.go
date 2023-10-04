package main

import (
	"fmt"
	"log"
	"os"
	bp "partdb/bipartite"
	"partdb/database"
)

const (
	protocol = "tcp"
	address  = ":1234"
)

func main() {
	os.Args = os.Args[1:]
	if len(os.Args) <= 0 {
		fmt.Println("Starting server... Ctrl+c to exit.")
		database.ServeRPC(protocol, address)
		return
	}

	if len(os.Args) <= 2 {
		log.Fatalln("Not enough arguments")
	}

	client, err := database.GetClient(protocol, address)
	if err != nil {
		log.Fatalln(err)
	}
	cmd := fmt.Sprintf("DB.%s", os.Args[0])

	switch os.Args[0] {
	case "Get":
		args := &database.TwoArg{
			First:  os.Args[1],
			Second: os.Args[2],
		}
		res := map[string]*bp.Part[string, string]{}
		client.Call(cmd, args, res)
		fmt.Println(res)
		for _, s := range res {
			fmt.Println(s)
		}
	case "Add":
		args := &database.AddArg{
			First:  os.Args[1],
			Second: os.Args[2],
			Third:  os.Args[3],
		}
		var err error
		_ = client.Call(cmd, args, err)
		if err != nil {
			log.Fatalln("something went wrong", err)
		}
		fmt.Println("Success!", args)
	default:
		args := &database.TwoArg{
			First:  os.Args[1],
			Second: os.Args[2],
		}
		var err error
		_ = client.Call(cmd, args, err)
		if err != nil {
			log.Fatalln("something went wrong", err)
		}
		fmt.Println("Success!", cmd, args)
	}
}
