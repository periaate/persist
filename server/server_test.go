package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	bp "partdb/bipartite"
	"testing"
)

func TestAddRHandler(t *testing.T) {
	b := bp.Make[string, string]()
	payload := addPayload{"keyR", "valueR"}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/addr", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addRHandler(b))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response["status"] != "success" {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestAddLHandler(t *testing.T) {
	b := bp.Make[string, string]()
	payload := addPayload{"keyL", "valueL"}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/addl", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addLHandler(b))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response["status"] != "success" {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestAddValuelessHandler(t *testing.T) {
	b := bp.Make[string, string]()
	payload := addValuelessPayload{"R", []string{
		"One",
		"Two",
		"Three",
	}}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/addmany", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addValuelessHandler(b))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response["status"] != "success" {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestEdgeHandler(t *testing.T) {
	b := bp.Make[string, string]()
	b.AddR("Hello", ",")
	b.AddL("World", "!")
	payload := edgePayload{"World", "Hello"}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/edge", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(edgeHandler(b))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response["status"] != "success" {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestList(t *testing.T) {
	b := bp.Make[string, string]()
	b.AddR("Hello", "World")
	b.AddValueless("L", []string{"one", "two", "three"})
	b.Edge("one", "Hello")
	b.Edge("two", "Hello")
	b.Edge("three", "Hello")

	payload := getPayload{"R", "Hello"}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/list", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(listHandler(b))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []string
	json.Unmarshal(rr.Body.Bytes(), &response)
	if len(response) != 3 {
		t.Errorf("expected len 3")
	}
}
