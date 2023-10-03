package database

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"partdb/structure"
)

type DB struct{ d structure.Bipartite[string] }

type Args struct {
	Side structure.Side
	Key  string
}

func (d *DB) List(args *Args, r *map[string]any) error {
	res, err := d.d.Get(args.Side, args.Key)
	if err != nil {
		*r = nil
		return err
	}
	*r = res
	return err
}

func (d *DB) Find(args *Args, err *error) error      { *err = d.d.Find(args.Side, args.Key); return nil }
func (d *DB) Add(args *Args, err *error) error       { *err = d.d.Add(args.Side, args.Key); return nil }
func (d *DB) Edge(args *[2]string, err *error) error { *err = d.d.Edge(args[0], args[1]); return nil }

func ServeRPC(prot string, addr string) {
	db := new(DB)
	rpc.Register(db)
	rpc.HandleHTTP()
	l, err := net.Listen(prot, addr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	http.Serve(l, nil)
}

func GetClient(prot string, addr string) (*rpc.Client, error) { return rpc.DialHTTP(prot, addr) }
