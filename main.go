package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"errors"
)

type server struct{}

type QUERY struct {
	OPER string `json:"oper"`
	K1 string `json:"k1"`
	K2 string `json:"k2"`
	DATA string `json:"data"`
}

type RESPONSE struct{
	DATA TABLE `json:"data"`
}

type ROW [3]string

type TABLE []ROW

var table TABLE

func match(A1 string, A2 string, C ROW) bool{
	if A1 != "" && A2 != "" && C[1] == A1 && C[2] == A2 {
		return true
	}

	if A1 != "" && A2 == "" && C[1] == A1 {
		return true
	}

	if A1 == "" && A2 == "" {
		return true
	}

	return false
}

func remove(s TABLE, i int) TABLE {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func db_insert(K1 string, K2 string, D string) error {
	if K1 == "" || K2 == "" || D == "" {
		return errors.New("Missing Data")
	}

	row := ROW{K1, K2, D}
	table = append(table, row)

	return nil
}

func db_delete(K1 string, K2 string) error {
	if K1 == "" && K2 == "" {
		return errors.New("Missing Data")
	}

	for i, s := range table {
		if match(K1, K2, s) {
			remove(table, i)
		}
	}

	return nil
}

func db_select(K1 string, K2 string) ([]byte, error) {
	var resp RESPONSE

	for _, s := range table {
		if match(K1, K2, s) {
			resp.DATA = append(resp.DATA, s)
		}
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request){
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
	}

    w.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func reportError(w http.ResponseWriter, e error){
	fmt.Println(e)
	http.Error(w, e.Error(), http.StatusInternalServerError)
	return
}

func main(){
	s := &server{}
	http.Handle("/", s)

	log.Fatal(http.ListenAndServe(":8080", nil))
}