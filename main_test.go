package main

import (
	"bytes"
	_ "fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_RejectGet(t *testing.T) {
	s := &server{}
	ts := httptest.NewServer(s)
	defer ts.Close()
	resp, err := http.Get(ts.URL)

	if err != nil {
		t.Errorf("Pinging root endpoint returned an error: %s", err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected MethodNotAllowed, got: %s", resp.Status)
	}
}

func Test_InsertData(t *testing.T) {
	s := &server{}
	ts := httptest.NewServer(s)
	defer ts.Close()

	b := bytes.NewReader([]byte(`{"oper": "INSERT", "k1": "TestSandwitch", "k2": "bread", "data": "sourdough"}`))
	resp, err := http.Post(ts.URL, "application/JSON", b)

	if err != nil {
		t.Errorf("Pinging root endpoint returned an error: %s", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected OK, got: %s", resp.Status)
	}
}

func Test_SelectAllData(t *testing.T) {
	s := &server{}
	ts := httptest.NewServer(s)
	defer ts.Close()

	b := bytes.NewReader([]byte(`{"oper": "INSERT", "k1": "TestSandwitch", "k2": "cheese", "data": "cheddar"}`))
	resp, err := http.Post(ts.URL, "application/JSON", b)

	if err != nil {
		t.Errorf("Pinging root endpoint returned an error: %s", err)
	}

	b = bytes.NewReader([]byte(`{"oper": "SELECT"}`))
	resp, err = http.Post(ts.URL, "application/JSON", b)

	if err != nil {
		t.Errorf("Pinging root endpoint returned an error: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected OK, got: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	expect := []byte(`{"TestSandwitch":{"bread":"sourdough","cheese":"cheddar"}}`)

	if !reflect.DeepEqual(body, expect) {
		t.Errorf("Bad output, got %s, expected %s", body, expect)
	}
}

func Test_SelectK1Data(t *testing.T) {
	s := &server{}
	ts := httptest.NewServer(s)
	defer ts.Close()

	b := bytes.NewReader([]byte(`{"oper": "INSERT", "k1": "TestSalad", "k2": "base", "data": "baby kale"}`))
	resp, err := http.Post(ts.URL, "application/JSON", b)

	if err != nil {
		t.Errorf("Pinging root endpoint returned an error: %s", err)
	}

	b = bytes.NewReader([]byte(`{"oper": "SELECT", "k1": "TestSalad"}`))
	resp, err = http.Post(ts.URL, "application/JSON", b)

	if err != nil {
		t.Errorf("Pinging root endpoint returned an error: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	expect := []byte(`{"TestSalad":{"base":"baby kale"}}`)

	if !reflect.DeepEqual(body, expect) {
		t.Errorf("Bad output, got %s, expected %s", body, expect)
	}
}

func Test_SelectK2Data(t *testing.T) {
	s := &server{}
	ts := httptest.NewServer(s)
	defer ts.Close()

	b := bytes.NewReader([]byte(`{"oper": "SELECT", "k2": "cheese"}`))
	resp, err := http.Post(ts.URL, "application/JSON", b)

	if err != nil {
		t.Errorf("Pinging root endpoint returned an error: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	expect := []byte(`{"TestSandwitch":{"cheese":"cheddar"}}`)

	if !reflect.DeepEqual(body, expect) {
		t.Errorf("Bad output, got %s, expected %s", body, expect)
	}
}

func Test_SelectSpecificData(t *testing.T) {
	s := &server{}
	ts := httptest.NewServer(s)
	defer ts.Close()

	b := bytes.NewReader([]byte(`{"oper": "SELECT", "k1": "TestSandwitch", "k2": "bread"}`))
	resp, err := http.Post(ts.URL, "application/JSON", b)

	if err != nil {
		t.Errorf("Pinging root endpoint returned an error: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	expect := []byte(`{"TestSandwitch":{"bread":"sourdough"}}`)

	if !reflect.DeepEqual(body, expect) {
		t.Errorf("Bad output, got %s, expected %s", body, expect)
	}
}

func Test_SelectNoData(t *testing.T) {
	s := &server{}
	ts := httptest.NewServer(s)
	defer ts.Close()

	b := bytes.NewReader([]byte(`{"oper": "SELECT", "k1": "TestSoup"}`))
	resp, err := http.Post(ts.URL, "application/JSON", b)

	if err != nil {
		t.Errorf("Pinging root endpoint returned an error: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	expect := []byte(`{"TestSoup":null}`)

	if !reflect.DeepEqual(body, expect) {
		t.Errorf("Bad output, got %s, expected %s", body, expect)
	}
}

func Test_DeleteData(t *testing.T) {
	s := &server{}
	ts := httptest.NewServer(s)
	defer ts.Close()

	b := bytes.NewReader([]byte(`{"oper": "DELETE", "k1": "TestSandwitch"}`))
	resp, err := http.Post(ts.URL, "application/JSON", b)

	if err != nil {
		t.Errorf("Pinging root endpoint returned an error: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected OK, got: %s", resp.Status)
	}

	b = bytes.NewReader([]byte(`{"oper": "SELECT", "k1": "TestSandwitch"}`))
	resp, err = http.Post(ts.URL, "application/JSON", b)

	if err != nil {
		t.Errorf("Pinging root endpoint returned an error: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	expect := []byte(`{"TestSandwitch":null}`)

	if !reflect.DeepEqual(body, expect) {
		t.Errorf("Bad output, got %s, expected %s", body, expect)
	}
}
