package main

import (
	"encoding/json"
	"log"
	"net/http"
	bp "partdb/bipartite"
	"time"
)

func logger(next http.Handler) http.Handler {
	// TODO: Add significantly better logging, including:
	// - Arguments
	// - Response status code and response values (or some simplification of them)
	// - other?
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.Error(w, "Internal Server Error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	db := bp.Make[string, string]()
	mux := http.NewServeMux()
	mux.HandleFunc("/addr", addRHandler(db))
	mux.HandleFunc("/addl", addLHandler(db))
	mux.HandleFunc("/addmany", addValuelessHandler(db))
	mux.HandleFunc("/edge", edgeHandler(db))
	mux.HandleFunc("/list", listHandler(db))

	// Apply middleware
	handler := logger(recoverPanic(mux))

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: handler,
	}

	log.Println("Server starting on localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
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

func addRHandler(b bp.Bipartite[string, string]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload addPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", 400)
			return
		}
		if err := b.AddR(payload.Key, payload.Value); err != nil {
			http.Error(w, "Failed to add", 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func addLHandler(b bp.Bipartite[string, string]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload addPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", 400)
			return
		}
		if err := b.AddL(payload.Key, payload.Value); err != nil {
			http.Error(w, "Failed to add", 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func addValuelessHandler(b bp.Bipartite[string, string]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload addValuelessPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", 400)
			return
		}
		if err := b.AddValueless(payload.Side, payload.Keys); err != nil {
			http.Error(w, "Failed to add", 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func edgeHandler(b bp.Bipartite[string, string]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload edgePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", 400)
			return
		}
		if err := b.Edge(payload.LKey, payload.RKey); err != nil {
			http.Error(w, "Failed to add", 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func listHandler(b bp.Bipartite[string, string]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload getPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", 400)
			return
		}
		list, err := b.List(payload.Side, payload.Key)
		if err != nil {
			http.Error(w, "Failed to add", 500)
			return
		}
		json.NewEncoder(w).Encode(list)
	}
}
