package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type server struct{}

type QUERY struct {
	OPER string `json:"oper"`
	K1   string `json:"k1"`
	K2   string `json:"k2"`
	DATA string `json:"data"`
}

type RESPONSE struct {
	DATA mytable `json:"data"`
}

type mytable map[string]map[string]string

var m = make(mytable)

func db_insert(K1 string, K2 string, D string) error {
	if K1 == "" || K2 == "" || D == "" {
		return errors.New("Missing Data")
	}

	mm, ok := m[K1]
    if !ok {
        mm = make(map[string]string)
        m[K1] = mm
    }
    mm[K2] = D

	return nil
}

func db_delete(K1 string, K2 string) error {
	if K1 == "" && K2 == "" {
		return errors.New("Missing Data")
	}

	if K1 != "" && K2 == "" {
		delete(m, K1)
	}

	if K1 == "" && K2 != "" {
		for k, _ := range m {
			_, ok := m[k][K2]
    		if ok {
    			delete(m[k], K2)
    		}
		}
	}

	if K1 != "" && K2 != "" {
		delete(m[K1], K2)
	}

	return nil
}

func db_select(K1 string, K2 string) ([]byte, error) {
	var temp mytable
	temp = make(mytable)

	if K1 == "" && K2 == "" {
		temp = m
	}

	if K1 != "" && K2 == "" {
		temp[K1] = m[K1]
	}

	if K1 == "" && K2 != "" {
		for k, _ := range m {
			i, ok := m[k][K2]
    		if ok {
    			temp[k] = make(map[string]string)
    			temp[k][K2] = i
    		}
		}
	}

	if K1 != "" && K2 != "" {
		temp[K1] = make(map[string]string)
		temp[K1][K2] = m[K1][K2]
	}

	b, err := json.Marshal(temp)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//handle JSON here
	if r.Method == http.MethodPost {
		var q QUERY

		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		err := json.NewDecoder(r.Body).Decode(&q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //bad request body
			return
		}

		switch q.OPER {
		case "INSERT":
			err := db_insert(q.K1, q.K2, q.DATA)
			if err != nil {
				reportError(w, err)
				return
			}
			w.WriteHeader(http.StatusCreated) //good request
		case "DELETE":
			err := db_delete(q.K1, q.K2)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		case "SELECT":
			data, err := db_select(q.K1, q.K2)
			if err != nil {
				reportError(w, err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		default:
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func reportError(w http.ResponseWriter, e error) {
	fmt.Println(e)
	http.Error(w, e.Error(), http.StatusInternalServerError)
	return
}

func main() {
	s := &server{}
	http.Handle("/", s)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
