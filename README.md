# PartDB

A key-value database which utilizes a bipartite graph as the main data structure, enabling powerful use cases.

## Overview

PartDB is a very light database, providing good usability for CLI, embedded, and API usage for systems requiring bipartite databases. Can also function as a typical KV database.

**What is a bipartite graph?**
A bipartite graph is a structure which has two distinct sets of nodes, where any node in either set may only ever have edges to nodes on the opposite set.

## Usage

To run the server: `go run server/server.go`

To run the client: `go run client/client.go {command} [flag] {arguments}`

Currently existing commands:
- `add` Add new keys and values to the database.
  - Flags:
    - `-r` Adds to the "right" partite of the graph.
	`partdb add -r rightKey rightValue`
    - `-l` Adds to the "left" partite of the graph.
	`partdb add -l leftKey leftValue`
    - If no flag is provided, the right partite is used, and the add can be given any number of arguments, adding each argument as its own key to right.
	`partdb add oneKey twoKey threeKey`
- `edge` Add an edge between a key on the left and on the right.
	`partdb edge leftKey rightKey`
- `list` List edge between a key on the left and on the right.
  - Flags:
    - `-r` Lists all of the edges `rightKey` has connected to the left.
	`partdb list -r rightKey`
    - `-l` Lists all of the edges `leftKey` has connected to the right.
	`partdb add -l leftKey`


