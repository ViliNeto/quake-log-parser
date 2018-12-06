package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Test bundle set for RESTFul API's
func TestHTTPHandler(t *testing.T) {
	passGetTest(t)
	failGetTest(t)
	passPostTest(t)
	failPostTest(t)
}

//Testing simple GET HTTP request on /get route (PASSING TEST)
func passGetTest(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "http://localhost:9000/get", nil)
	if err != nil {
		fmt.Printf("FAILED GET: Could not do the request: %v", err)
	}
	recorder := httptest.NewRecorder()
	getHandler(recorder, request)
	if recorder.Code != http.StatusOK {
		fmt.Printf("FAILED GET: Expected status 200; got %d\n", recorder.Code)
	} else {
		fmt.Println("PASSED GET: Status Code 200")
	}
}

//Testing simple GET HTTP request on /get route (FAILING TEST)
func failGetTest(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "http://localhost:9000/get", nil)
	if err != nil {
		fmt.Printf("FAILED GET: Could not do the request: %v\n", err)
	}
	recorder := httptest.NewRecorder()
	getHandler(recorder, request)
	if recorder.Code != http.StatusOK {
		fmt.Printf("FAILED GET: Expected status 200; got %d\n", recorder.Code)
	} else {
		fmt.Println("PASSED GET: Status Code 200")
	}
}

//Testing simple POST HTTP request on /post route (PASSING TEST)
func passPostTest(t *testing.T) {
	jsonText := `[{"game_1":{"total_kills":0,"players":["Isgalamido"],"kills":null}},{"game_2":{"total_kills":2,"players":["Isgalamido","Dono da Bola","Mocinha","Zeh"],"kills":{"Isgalamido":0}}}]`
	err := json.Unmarshal([]byte(jsonText), &parsedJSON)
	payload := bytes.NewBufferString(`{"game": "game_1"}`)
	request, err := http.NewRequest("POST", "/post", payload)
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("FAILED POST: Could not do the request: %v\n", err)
	}
	recorder := httptest.NewRecorder()
	postHandler(recorder, request)
	if recorder.Code != http.StatusOK {
		fmt.Printf("FAILED POST: Expected status 200; got %d\n", recorder.Code)
	} else {
		fmt.Println("PASSED POST: Status Code 200")
	}
}

//Testing simple POST HTTP request on /post route (FAILING TEST)
func failPostTest(t *testing.T) {
	request, err := http.NewRequest("POST", "/post", nil)
	if err != nil {
		t.Fatalf("Could not do the request: %v\n", err)
	}
	recorder := httptest.NewRecorder()
	getHandler(recorder, request)
	if recorder.Code != http.StatusOK {
		fmt.Printf("Expected status 200; got %d\n", recorder.Code)
	} else {
		fmt.Println("PASSED: Status Code 200")
	}
}
