package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const serverURL = "http://localhost:8080"

func main() {
	addCommand := flag.NewFlagSet("add", flag.ExitOnError)
	edgeCommand := flag.NewFlagSet("edge", flag.ExitOnError)
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)

	addR := addCommand.Bool("r", false, "Right-side key")
	addL := addCommand.Bool("l", false, "Left-side key")

	listR := listCommand.Bool("r", false, "List by right-side key")
	listL := listCommand.Bool("l", false, "List by left-side key")

	if len(os.Args) < 2 {
		fmt.Println("add, edge or list subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCommand.Parse(os.Args[2:])
		executeAdd(addCommand.Args(), *addR, *addL)
	case "edge":
		edgeCommand.Parse(os.Args[2:])
		executeEdge(edgeCommand.Args())
	case "list":
		listCommand.Parse(os.Args[2:])
		executeList(listCommand.Args(), *listR, *listL)
	default:
		fmt.Println("Only 'add', 'edge', or 'list' are valid commands")
		os.Exit(1)
	}
}

func executeAdd(args []string, rFlag, lFlag bool) {
	side := "r"
	if lFlag {
		side = "l"
	}

	var endpoint string
	switch side {
	case "r":
		endpoint = "/addr"
	case "l":
		endpoint = "/addl"
	default:
		endpoint = "/addmany"
	}

	var data []byte

	if endpoint == "/addmany" {
		payload := addValuelessPayload{
			Side: side,
			Keys: args[0:],
		}

		data, _ = json.Marshal(payload)
	} else {
		payload := addPayload{
			Key:   args[0],
			Value: args[1],
		}
		data, _ = json.Marshal(payload)
	}

	if sendRequest(endpoint, data)["status"] == "success" {
		fmt.Println("success")
	}
}

func executeEdge(args []string) {
	payload := edgePayload{
		RKey: args[0],
		LKey: args[1],
	}
	data, _ := json.Marshal(payload)

	if sendRequest("/edge", data)["status"] == "success" {
		fmt.Println("success")
	}
}

func executeList(args []string, rFlag, lFlag bool) {
	side := "r"
	if lFlag {
		side = "l"
	}

	payload := getPayload{
		Side: side,
		Key:  args[0],
	}

	data, _ := json.Marshal(payload)

	response := sendGetRequest("/list", data)

	for _, item := range response {
		fmt.Println(item)
	}
}

func sendRequest(endpoint string, data []byte) map[string]string {
	resp, err := http.Post(serverURL+endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Failed to send request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	json.Unmarshal(body, &result)

	return result
}

func sendGetRequest(endpoint string, data []byte) []string {
	resp, err := http.Post(serverURL+endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Failed to send request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result []string
	json.Unmarshal(body, &result)

	return result
}

type addPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type edgePayload struct {
	LKey string `json:"lKey"`
	RKey string `json:"rKey"`
}

type getPayload struct {
	Side string `json:"side"`
	Key  string `json:"key"`
}

type addValuelessPayload struct {
	Side string   `json:"side"`
	Keys []string `json:"keys"`
}
