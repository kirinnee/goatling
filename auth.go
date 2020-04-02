package goatling

import (
	"encoding/json"
	"errors"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

type (
	Jwks struct {
		Keys []JSONWebKeys `json:"keys"`
	}

	JSONWebKeys struct {
		Kty string   `json:"kty"`
		Kid string   `json:"kid"`
		Use string   `json:"use"`
		N   string   `json:"n"`
		E   string   `json:"e"`
		X5c []string `json:"x5c"`
	}
)

func getPemCert(token *jwt.Token, domain string) (string, error) {
	cert := ""
	resp, err := http.Get(domain + ".well-known/jwks.json")
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func jwtMiddleware(domain string) *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(
		jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

				// Verify 'iss' claim
				iss := domain
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss {
					log.Println("invalid issuer")
					return token, errors.New("invalid issuer")
				}

				cert, err := getPemCert(token, domain)
				if err != nil {
					panic(err.Error())
				}

				result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
				return result, nil
			},
			SigningMethod: jwt.SigningMethodRS256,
		})
}
