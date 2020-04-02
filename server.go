package goatling

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
	*mux.Router
}

func New() *Server {
	return &Server{mux.NewRouter()}
}

func resp(w http.ResponseWriter, any interface{}) {
	b, _ := json.Marshal(any)
	_, err := w.Write(b)
	if err != nil {
		log.Println(err)
	}
}

func getBody(r *http.Request, any interface{}) *ServerResponse {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return BadRequest("Cannot read request")
	}
	log.Println(string(body))
	err = json.Unmarshal(body, &any)
	if err != nil {
		return BadRequest("Failed to unmarshal json")
	}
	return nil
}

type goat struct {
	variables map[string]string
	r         *http.Request
}

type Goat interface {
	Vars() map[string]string
	Body(any interface{}) *ServerResponse
	Header() *http.Header
	Request() *http.Request
}

func (g *goat) Request() *http.Request {
	return g.r
}

func (g *goat) Header() *http.Header {
	return &g.r.Header
}

func (g *goat) Vars() map[string]string {

	return g.variables
}

func (g *goat) Body(any interface{}) *ServerResponse {
	return getBody(g.r, any)
}

func (s *Server) Serve(path string, handler func(Goat) *ServerResponse) *mux.Route {
	return s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		g := goat{
			variables: mux.Vars(r),
			r:         r,
		}
		response := handler(&g)
		w.WriteHeader(response.Status)
		if response.Content != nil {
			resp(w, response.Content)
		}
	})
}
