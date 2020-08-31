package goatling

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
	cors       string
	enableCors bool
	*mux.Router
}

func New() *Server {
	s := &Server{enableCors: false, cors: "", Router: mux.NewRouter()}
	return s
}

func (s *Server) SetCORS(cors string) {
	if !s.enableCors {
		s.Use(mux.CORSMethodMiddleware(s.Router))
		s.enableCors = true
	}
	s.cors = cors
}

func resp(w http.ResponseWriter, any interface{}) {
	b, _ := json.Marshal(any)
	_, err := w.Write(b)
	if err != nil {
		log.Println(err)
	}
}

func respRaw(w http.ResponseWriter, any []byte) {
	_, err := w.Write(any)
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
	req       *http.Request
	resp      http.ResponseWriter
}

type Goat interface {
	Vars() map[string]string
	Body(any interface{}) *ServerResponse
	BodyString() string
	BodyBytes() []byte
	Header() *http.Header
	Request() *http.Request
	Response() http.ResponseWriter
}

func (g *goat) Request() *http.Request {
	return g.req
}

func (g *goat) Response() http.ResponseWriter {
	return g.resp
}

func (g *goat) Header() *http.Header {
	return &g.req.Header
}

func (g *goat) Vars() map[string]string {

	return g.variables
}

func (g *goat) Body(any interface{}) *ServerResponse {
	return getBody(g.req, any)
}

func (g *goat) BodyString() string {
	return string(g.BodyBytes())
}

func (g *goat) BodyBytes() []byte {
	b, _ := ioutil.ReadAll(g.req.Body)
	return b
}

func (s *Server) ServeString(path string, handler func(Goat) *ServerResponse) *mux.Route {
	return s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		g := goat{
			variables: mux.Vars(r),
			req:       r,
			resp:      w,
		}
		response := handler(&g)
		if response == nil {
			return
		}
		w.WriteHeader(response.Status)
		if response.Content != nil {
			resp(w, response.Content)
		}
	})
}

func (s *Server) ServeRaw(path string, handler func(Goat) *RawServerResponse) *mux.Route {
	return s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		g := goat{
			variables: mux.Vars(r),
			req:       r,
			resp:      w,
		}
		response := handler(&g)
		if response == nil {
			return
		}
		w.WriteHeader(response.Status)
		if response.Content != nil {
			respRaw(w, response.Content)
		}
	})
}

func (s *Server) Serve(path string, handler func(Goat) *ServerResponse) *mux.Route {
	return s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		g := goat{
			variables: mux.Vars(r),
			req:       r,
			resp:      w,
		}
		response := handler(&g)
		if response == nil {
			return
		}
		w.WriteHeader(response.Status)
		if response.Content != nil {
			resp(w, response.Content)
		}
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if s.enableCors {
		(w).Header().Set("Access-Control-Allow-Origin", s.cors)
		(w).Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			(w).Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	s.Router.ServeHTTP(w, r)
}
