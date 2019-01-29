package transport

import (
	"../books"
	"../security"
	"../utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Описание структуры ответа на list-запрос
type FetchResult struct {
	Name  string         `json:"name"`
	Total int            `json:"total"`
	Items []*books.Books `json:"items"`
}

// Описание структуры ответа при ошибке
type ErrorResponse struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

type Router interface {
	// Get
	Get(w http.ResponseWriter, req *http.Request)
	// Patch
	Patch(w http.ResponseWriter, req *http.Request)
	// Post
	Post(w http.ResponseWriter, req *http.Request)
	// Delete ...Handler
	Delete(w http.ResponseWriter, req *http.Request)
	// Put
	//Put(w http.ResponseWriter, req *http.Request)
	// Обработчик запросов
	Handle(http.ResponseWriter, *http.Request)
}

// Класс реализующий транспортный уровень
type HttpHandler struct {
	Sm    *security.SessionManager
	Books books.HandlerBooks
}

// Создание инстанса
func CreateInstance(sm *security.SessionManager, b books.HandlerBooks) *HttpHandler {
	return &HttpHandler{
		Sm:    sm,
		Books: b,
	}
}

// Get
func (s *HttpHandler) Get(w http.ResponseWriter, req *http.Request) {
	if id := req.URL.Query().Get("id"); id != "" {
		if err := s.fetchOneBook(id, w); err != nil {
			return
		}
	} else {
		if err := s.fetchListBooks(w); err != nil {
			return
		}
	}
}

func (s *HttpHandler) fetchListBooks(w http.ResponseWriter) error {
	books, count, err := s.Books.GetBooks()
	if err != nil {
		s.errorHandler(w, "Internal Server Error", http.StatusInternalServerError, err)
		return err
	}
	b, _ := json.Marshal(&FetchResult{
		Name:  "List books",
		Total: count,
		Items: books,
	})
	w.Write([]byte(b))
	w.WriteHeader(http.StatusOK)
	return nil
}

func (s *HttpHandler) fetchOneBook(id string, w http.ResponseWriter) error {
	book := &books.Books{}
	book, _, err := s.Books.GetBook(id)
	if err != nil {
		s.errorHandler(w, "Internal Server Error", http.StatusInternalServerError, err)
		return err
	}
	b, _ := json.Marshal(&book)
	w.Write([]byte(b))
	w.WriteHeader(http.StatusOK)
	return nil
}

// Post
func (s *HttpHandler) Post(w http.ResponseWriter, req *http.Request) {
	book := &books.Books{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.errorHandler(w, "Body read error", http.StatusInternalServerError, err)
		return
	}
	err = json.Unmarshal(body, book)
	if err != nil {
		s.errorHandler(w, "unmarshal error", http.StatusInternalServerError, err)
		return
	}

	if !book.IsValid() {
		s.errorHandler(w, "Not valid data", http.StatusBadRequest, nil)
		return
	}

	book.Id = utils.GenerateUUID()

	err = s.Books.AddBooks(book)

	if err != nil {
		s.errorHandler(w, "Internal Server Error", http.StatusInternalServerError, err)
		return
	}
	b, _ := json.Marshal(&book)
	w.Write([]byte(b))
	w.WriteHeader(http.StatusOK)
}

// Patch
func (s *HttpHandler) Patch(w http.ResponseWriter, req *http.Request) {
	book := &books.Books{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.errorHandler(w, "Body read error", http.StatusInternalServerError, err)
		return
	}
	err = json.Unmarshal(body, book)
	if err != nil {
		s.errorHandler(w, "unmarshal error", http.StatusInternalServerError, err)
		return
	}

	if !book.IsValid() {
		s.errorHandler(w, "Not valid data", http.StatusBadRequest, nil)
		return
	}

	bookFromDB, count, err := s.Books.GetBook(book.Id)

	if err != nil || count == 0 {
		s.errorHandler(w, "Record not found", http.StatusNotFound, err)
		return
	}
	if book.Name == "" {
		book.Name = bookFromDB.Name
	}
	if book.PageCount == 0 {
		book.PageCount = bookFromDB.PageCount
	}
	if book.ReleaseYear == 0 {
		book.ReleaseYear = bookFromDB.ReleaseYear
	}

	if err = s.Books.EditBooks(book); err != nil {
		s.errorHandler(w, "Internal Server Error", http.StatusInternalServerError, err)
		return
	}
	b, _ := json.Marshal(&book)
	w.Write([]byte(b))
	w.WriteHeader(http.StatusOK)
}

// Delete
func (s *HttpHandler) Delete(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	err := s.Books.DeleteBooks(id)
	if err != nil {
		s.errorHandler(w, "Internal Server Error", http.StatusInternalServerError, err)
		return
	}
	b, _ := json.Marshal(ErrorResponse{Message: id + " delete successful"})
	w.Write([]byte(b))
	w.WriteHeader(http.StatusOK)
}

func (s *HttpHandler) errorHandler(w http.ResponseWriter, message string, status int, err error) {
	b, _ := json.Marshal(&ErrorResponse{
		Message: message,
		Error:   err,
	})
	w.Write([]byte(b))
	w.WriteHeader(status)
}

// Http Handle
func (s *HttpHandler) Handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case http.MethodGet:
		s.Get(w, req)
		return
	case http.MethodPost:
		s.Post(w, req)
		return
	case http.MethodDelete:
		s.Delete(w, req)
	case http.MethodPatch:
		s.Patch(w, req)
	case http.MethodPut:
		// TODO:: create or load record
		//handler.Put(req)
	}
}
