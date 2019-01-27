package transport

import (
	"../security"
	"net/http"
)

type FetchResult struct {
	Header FetchResultHeader `json:"header"`
	Params FetchResultParams `json:"params"`
	Items  interface{}       `json:"items"`
}

type FetchResultHeader struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

type FetchResultParams struct {
	Filter string `json:"filter"`
	Sort   string `json:"sort"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type Router interface {
	// Get
	Get(w http.ResponseWriter, req *http.Request)
	// Patch
	//Patch(w http.ResponseWriter, req *http.Request)
	// Post
	//Post(w http.ResponseWriter, req *http.Request)
	// Put
	//Put(w http.ResponseWriter, req *http.Request)
	// Delete ...Handler
	//Delete(w http.ResponseWriter, req *http.Request)

	Handle(http.ResponseWriter, *http.Request)
}

type server struct {
	Sm     *security.SessionManager
	router *http.Request
	email  string
}

func CreateInstance(sm *security.SessionManager) Router {
	return &server{
		Sm: sm,
	}
}

func (s *server) Get(w http.ResponseWriter, req *http.Request) {

}

func (handler *server) Handle(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		handler.Get(w, req)
	case http.MethodPost:
		//handler.Post(req)
	case http.MethodDelete:
		//handler.Delete(req)
	case http.MethodPatch:
		//handler.Patch(req)
	case http.MethodPut:
		//handler.Put(req)
	}
}
