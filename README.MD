# Goatling, JSON HTTP server
A HTTP Server wrapper over gorilla/mux for better json support

# Getting Started
Import it as a library (go module)
```go
package main

import(
    goat "gitlab.com/kiringo/goatling"
)
```
or via go CLI
```bash
go get gitlab.com/kiringo/goatling
```

# Using Basic Server
```go
import (
	"flag"
	goat "gitlab.com/kiringo/goatling"
	"log"
	"net/http"
)

// Sample Response Model
type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

// Sample Endpoint
func Home(g goat.Goat) *goat.ServerResponse{
	user := User{
		Name: "kirin",
		Email: "kirinnee97@gmail.com",
	}
	return goat.OK(user)
}


func main() {

	var addr = flag.String("addr", ":9000", "http service address")
	flag.Parse()
	// Normal Server
	r := goat.New()

	// Now you can use r and a like gorilla mux server
	r.Serve("/", Home).Methods("GET")

	log.Println("Listening")
	log.Fatal(http.ListenAndServe(*addr, r))
}

```

## Accessing Body as JSON
```go
func Home(g goat.Goat) *goat.ServerResponse{
    var user User
    g.Body(&user) //Reads request into user struct
    return goat.OK(user)
}
```

## Accessing Headers
```go
func Home(g goat.Goat) *goat.ServerResponse{
    header := g.Header().Get("Header") //Retrieve request header
   	return goat.BadGateway(header)
}
```

## Accessing Path Fragments
```go
func Users(g goat.Goat) *goat.ServerResponse{
    pathFrag := g.Vars()["user"]
    return goat.OK(pathFrag)
}

func main() {
    //.. omitted for brevity
    r.Serve("/{user}", Users).Methods("GET")
    //.. omitted for brevity

}
```

## Returning Response Code
```go

import (
)

func Users(g goat.Goat) *goat.ServerResponse {
	typeToReturn := g.Vars()["user"]
	if typeToReturn == "200" {
		// Return with 200
		return goat.OK(typeToReturn)
	}
	if typeToReturn == "201" {
		// Return with 201
		return goat.Created(typeToReturn)
	}
	if typeToReturn == "202" {
		// Return with 202
		return goat.Accepted(typeToReturn)
	}

	if typeToReturn == "204" {
		// Return with 202
		return goat.NoContent(typeToReturn)
	}

	if typeToReturn == "400" {
		// Return with 400
		return goat.BadRequest(typeToReturn)
	}
	if typeToReturn == "401" {
		// Return with 401
		return goat.Unauthorized(typeToReturn)
	}
	if typeToReturn == "401" {
		// Return with 403
		return goat.Forbidden(typeToReturn)
	}
	if typeToReturn == "404" {
		// Return with 404
		return goat.NotFound(typeToReturn)
	}

	if typeToReturn == "409" {
		// Return with 409
		return goat.Conflict(typeToReturn)
	}
	if typeToReturn == "500" {
		// Return with 500
		return goat.InternalServerError(typeToReturn)
	}

	if typeToReturn == "502" {
		// Return with 502
		return goat.BadGateway(typeToReturn)
	}

	// Custom Code
	return &goat.ServerResponse{
		Status:  419,
		Content: typeToReturn,
	}

}
```

# Enabling CORS
```go
// Set CORS before serving any page 

func main() {

	var addr = flag.String("addr", ":9000", "http service address")
	flag.Parse()
	// Normal Server
	r := goat.New()

    r.SetCORS("*") 

	// Now you can use r and a like gorilla mux server
	r.Serve("/", Home).Methods("GET")

	log.Println("Listening")
	log.Fatal(http.ListenAndServe(*addr, r))
}


```

# Using Auth0 Integrated Server
```go
// Sample Response Model
type User struct {
    Name string `json:"name"`
    Email string `json:"email"`
}

// Sample Endpoint
func Home(g goat.AuthGoat) *goat.ServerResponse{  
       user := User{
            Name: "kirin",
            Email: "kirinnee97@gmail.com",
       }
    return goat.OK(user)
}


func main() {

    var addr = flag.String("addr", ":9000", "http service address")
    flag.Parse()
    // Auth Server 
    r := goat.NewAuth("https://domain.region.auth0.com/") //Auth0 issuing endpoint

    // Require Authenticated user
    r.ServePrivate("/", Home).Methods("GET")

    log.Println("Listening")
    log.Fatal(http.ListenAndServe(*addr, r))
}
```

## Obtain Claims
Create Custom Claim that extends jwt Claims 
```go
type CustomClaim struct {
    // Addition claims field
    Email string `json:"email"`
    // Extend basic claims
    jwt.StandardClaims
}
```

Cast Claims when using goat object 
```go
func Home(g goat.AuthGoat) *goat.ServerResponse{
	claim := g.Claims(&CustomClaim{}).(*CustomClaim)
    user := User{
        Name:  "kirin",
        Email: claim.Email,
    }
    return goat.OK(user)
}

```


## Contributing
Please read [CONTRIBUTING.md](CONTRIBUTING.MD) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning 
We use [SemVer](https://semver.org/) for versioning. For the versions available, see the tags on this repository.etting

## Authors
* [kirinnee](mailto:kirinnee97@gmail.com) 

## License
This project is licensed under MIT - see the [LICENSE.md](LICENSE.MD) file for details