package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestWritePortHandler_Success(t *testing.T) {
	initPorts(3, 2)

	portID := 3
	value := 1

	data := map[string]int{"value": value}
	body, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", "/port/"+strconv.Itoa(portID)+"/write", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(writePortHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if ports[portID].value != value {
		t.Errorf("handler did not write correct value: got %v want %v", ports[portID].value, value)
	}
}

func TestWritePortHandler_PortNotFound(t *testing.T) {
	initPorts(3, 2)

	portID := 10
	value := 1

	data := map[string]int{"value": value}
	body, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", "/port/"+strconv.Itoa(portID)+"/write", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(writePortHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestWritePortHandler_WriteToInPort(t *testing.T) {
	initPorts(3, 2)
	portID := 1
	value := 1

	data := map[string]int{"value": value}
	body, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", "/port/"+strconv.Itoa(portID)+"/write", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(writePortHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
