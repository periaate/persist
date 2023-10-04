package main

import (
	"fmt"
	"log"
	"os"
	"partdb/database"
	"partdb/structure"
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

	if os.Args[0] == "Edge" {
		args := [2]string{os.Args[1], os.Args[1]}
		var err error
		client.Call(os.Args[0], args, err)
		if err != nil {
			log.Fatalln("something went wrong", err)
		}
		fmt.Println("Success!", os.Args[0], os.Args[1], os.Args[2])
		return
	}

	side, err := Side(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	args := database.Args{
		Side: side,
		Key:  os.Args[2],
	}

	switch os.Args[0] {
	case "list":
		res := map[string]any{}
		client.Call("list", args, res)
		for _, s := range res {
			fmt.Println(s)
		}
	default:
		var err error
		client.Call(os.Args[0], args, err)
		if err != nil {
			log.Fatalln("something went wrong", err)
		}
		fmt.Println("Success!", os.Args[0], os.Args[1], os.Args[2])
	}
}

var r = map[string]any{
	"r":     "",
	"R":     "",
	"right": "",
	"Right": "",
	"RIGHT": "",
}

var l = map[string]any{
	"l":    "",
	"L":    "",
	"left": "",
	"Left": "",
	"LEFT": "",
}

func Side(s string) (structure.Side, error) {
	if _, ok := r[s]; ok {
		return structure.Right, nil
	}
	if _, ok := l[s]; ok {
		return structure.Left, nil
	}
	return structure.BadSide, fmt.Errorf("no such side")
}
