package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type PortStruct struct {
	id       int
	portType PortType
	value    int
	mu       sync.Mutex
}

type PortType int

const (
	IN PortType = iota
	OUT
)

var ports = make(map[int]*PortStruct)
var mu sync.Mutex

func writePortHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid port ID", http.StatusBadRequest)
		return
	}

	var data struct {
		Value int `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	mu.Lock()
	port, exists := ports[portID]
	mu.Unlock()

	if !exists || port.portType != OUT {
		http.Error(w, "Port not found or not an OUT port", http.StatusNotFound)
		return
	}

	port.mu.Lock()
	port.value = data.Value
	port.mu.Unlock()

	fmt.Printf("OUT Port %d: Write value %d\n", portID, port.value)
	w.WriteHeader(http.StatusOK)
}

func initPorts(numInPorts, numOutPorts int) {
	for i := 0; i < numInPorts; i++ {
		port := &PortStruct{id: i, portType: IN}
		ports[i] = port
		go port.operate()
	}
	for i := numInPorts; i < numInPorts+numOutPorts; i++ {
		port := &PortStruct{id: i, portType: OUT}
		ports[i] = port
		go port.operate()
	}
}

func (p *PortStruct) operate() {
	switch p.portType {
	case IN:
		for {
			p.mu.Lock()
			p.value = rand.Intn(2)
			fmt.Printf("IN Port %d: Read value %d\n", p.id, p.value)
			p.mu.Unlock()
			time.Sleep(time.Second)
		}
	case OUT:
		for {
			time.Sleep(time.Second)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	initPorts(3, 2)

	r := mux.NewRouter()
	r.HandleFunc("/port/{id}/write", writePortHandler).Methods("POST")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
