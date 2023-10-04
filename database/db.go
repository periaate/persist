package database

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	bp "partdb/bipartite"
)

type TwoArg struct {
	First  string
	Second string
}

type AddArg struct {
	First  string
	Second string
	Third  string
}

type DB struct {
	d bp.Bipartite[string, string]
}

func (d *DB) Get(args *TwoArg, ret *map[string]*bp.Part[string, string]) error {
	fmt.Println("GETTING:", args.Second, "; FROM SIDE:", args.First)
	res, err := d.d.Get(args.First, args.Second)
	*ret = res
	return err
}
func (d *DB) Add(args *AddArg, err *error) error {
	fmt.Println("ADDING TO:", args.First, "; KEY:", args.Second, "; VALUE:", args.Third)
	r := d.d.Add(args.First, args.Second, args.Third)
	fmt.Println("success")
	*err = r
	return r
}

func (d *DB) Edge(args *TwoArg, err *error) error {
	fmt.Println("ADDING EDGE BETWEEN RIGHT KEY:", args.First, " AND LEFT KEY:", args.Second)
	r := d.d.Edge(args.First, args.Second)
	fmt.Println("success")
	*err = r
	return r
}

func ServeRPC(prot string, addr string) {
	db := &DB{d: bp.Make[string, string]()}
	rpc.Register(db)
	rpc.HandleHTTP()
	l, err := net.Listen(prot, addr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	http.Serve(l, nil)
}

func GetClient(prot string, addr string) (*rpc.Client, error) { return rpc.DialHTTP(prot, addr) }
