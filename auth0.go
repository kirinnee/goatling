package goatling

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"strings"
)

type Auth0Server struct {
	domain     string
	middleware *jwtmiddleware.JWTMiddleware
	*Server
}

func NewAuth(domain string, aud string) *Auth0Server {
	return &Auth0Server{
		domain:     domain,
		middleware: jwtMiddleware(domain, aud),
		Server:     New(),
	}
}

type AuthGoat interface {
	Claims(claims jwt.Claims) jwt.Claims
	Goat
}

type authGoat struct {
	domain string
	*goat
}

func (g *authGoat) Claims(claims jwt.Claims) jwt.Claims {
	authHeaderParts := strings.Split(g.req.Header.Get("Authorization"), " ")
	t := authHeaderParts[1]
	token, _ := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token, g.domain)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	return token.Claims
}

func (s *Auth0Server) ServePrivate(path string, handler func(AuthGoat) *ServerResponse) *mux.Route {
	return s.Handle(path, negroni.New(
		negroni.HandlerFunc(s.middleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if s.enableCors {
				w.Header().Set("Access-Control-Allow-Origin", s.cors)
				if r.Method == http.MethodOptions {
					return
				}
			}
			g := authGoat{
				domain: s.domain,
				goat: &goat{
					variables: mux.Vars(r),
					req:       r,
					resp:      w,
				},
			}
			response := handler(&g)
			if response == nil {
				return
			}
			w.WriteHeader(response.Status)
			if response.Content != nil {
				resp(w, response.Content)
			}
		}))))
}

func (s *Auth0Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
