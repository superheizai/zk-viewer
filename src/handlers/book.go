package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "models"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Handler for the books Create action
// POST /books

func BookCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	book := &Book{}
	if err := populateModelFromHandler(w, r, params, book); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}
	Bookstore[book.ISDN] = book
	writeOKResponse(w, book)
}

// Handler for the books index action
// GET /books
func BookIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	books := []*Book{}
	for _, book := range Bookstore {
		books = append(books, book)
	}
	writeOKResponse(w, books)
}

// Handler for the books Show action
// GET /books/:isdn
func BookShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isdn := params.ByName("isdn")
	book, ok := Bookstore[isdn]
	if !ok {
		// No book with the isdn in the url has been found
		writeErrorResponse(w, http.StatusNotFound, "Record Not Found")
		return
	}
	writeOKResponse(w, book)
}

// Writes the response as a standard JSON response with StatusOK
func writeOKResponse(w http.ResponseWriter, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length, Authorization, Accept,X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods","PUT,POST,GET,DELETE,OPTIONS")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&JsonResponse{Data: m}); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// Writes the error response as a Standard API JSON response with a response code
func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length, Authorization, Accept,X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods","PUT,POST,GET,DELETE,OPTIONS")
	w.WriteHeader(errorCode)
	json.
		NewEncoder(w).
		Encode(&JsonErrorResponse{Error: &ApiError{Status: errorCode, Title: errorMsg}})
}

//Populates a model from the params in the Handler
func populateModelFromHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params, model interface{}) error {
	//body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	decoder:=json.NewDecoder(r.Body)
	err:=decoder.Decode(model)
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	//if err := json.Unmarshal(body, model); err != nil {
	//	return err
	//}
	return nil
}
