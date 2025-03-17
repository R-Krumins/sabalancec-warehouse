package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	listenAddr string
	store      Storage
}

func NewServer(port string, store Storage) *Server {
	return &Server{
		listenAddr: ":" + port,
		store:      store,
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	mux.Handle("GET /api/product", makeHttpHandlerFunc(s.handleGetProduct))
	mux.Handle("GET /api/product/{id}", makeHttpHandlerFunc(s.handleGetProductById))
	mux.Handle("POST /api/product", makeHttpHandlerFunc(s.handleCreateProduct))

	mux.Handle("GET /api/allergen", makeHttpHandlerFunc(s.handleGetAllergen))
	mux.Handle("GET /api/allergen/{id}", makeHttpHandlerFunc(s.handleGetAllergenById))

	// Static file serving for images
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	}))

	fmt.Printf("Server listening on port %s...\n", s.listenAddr)
	http.ListenAndServe(s.listenAddr, mux)
}

func (s *Server) handleCreateProduct(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusForbidden, "Forbidden", nil)
	// product := new(ProductFull)

	// if err := json.NewDecoder(r.Body).Decode(product); err != nil {
	// 	return err
	// }

	// if err := s.store.CreateProduct(product); err != nil {
	// 	return err
	// }

	// return WriteJSON(w, http.StatusOK,
	// 	fmt.Sprintf("New product %s created", product.Name), product)
}

func (s *Server) handleGetProduct(w http.ResponseWriter, r *http.Request) error {
	products, err := s.store.GetProduct()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, "", products)
}

func (s *Server) handleGetAllergen(w http.ResponseWriter, r *http.Request) error {
	products, err := s.store.GetAllergen()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, "", products)
}

func (s *Server) handleGetProductById(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}

	product, err := s.store.GetProductById(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, "", product)
}

func (s *Server) handleGetAllergenById(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}

	product, err := s.store.GetAllergenByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, "", product)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHttpHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, msg string, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	resStatus := false
	if status >= 200 && status < 300 {
		resStatus = true
	}

	return json.NewEncoder(w).Encode(Res{
		Success: resStatus,
		Msg:     msg,
		Data:    data,
	})
}

func getId(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return 0, fmt.Errorf("invalid id: %s", r.PathValue("id"))
	}
	return id, nil
}
