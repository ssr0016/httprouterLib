package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	db "github.com/ssr0016/librarySystem/model"
	"github.com/ssr0016/librarySystem/types"
)

type APIServer struct {
	listenAddr string
	store      db.Storage
	logger     *log.Logger
}

func NewAPIServer(listenAddr string, store db.Storage) *APIServer {
	logger := log.New(log.Writer(), "api server: ", log.LstdFlags|log.Lshortfile)
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
		logger:     logger,
	}
}

func (s *APIServer) Run() {
	router := httprouter.New()

	router.GET("/books", s.handlerGetBooks)
	router.GET("/books/:id", s.handlerGetBook)
	router.POST("/books", s.handleCreateBook)
	// router.PUT("/books/:id", updateBook)
	router.DELETE("/books/:id", s.handlerDeleteBook)

	log.Println("Server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleCreateBook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := new(types.CreateBookRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := types.NewBook(req.Title, req.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.store.CreateBook(r.Context(), *book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSONResponse(w, http.StatusCreated, book)
}

func (s *APIServer) handlerGetBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := s.store.GetBook(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSONResponse(w, http.StatusOK, book)
}

func (s *APIServer) handlerGetBooks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	books, err := s.store.GetBooks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSONResponse(w, http.StatusOK, books)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *APIServer) handlerDeleteBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.store.DeleteBook(r.Context(), int64(id)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "book deleted successfully"})
}
