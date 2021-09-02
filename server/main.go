package main

import (
	"encoding/json"
	"fmt"
	"github.com/flypay/engineering-test/pkg/api/alpha"
	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/api/beta"
	"github.com/flypay/engineering-test/pkg/api/orders"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func main() {
	apiHandlers := make([]apiHandler.Handler, 0, 5)
	apiHandlers = append(apiHandlers, alpha.NewGetMenuCategories())
	apiHandlers = append(apiHandlers, alpha.NewGetMenuIngredients())
	apiHandlers = append(apiHandlers, alpha.NewGetMenuProductsHandler())
	apiHandlers = append(apiHandlers, beta.NewGetMenu())
	apiHandlers = append(apiHandlers, orders.NewSetOrder())
	router := mux.NewRouter().StrictSlash(true)
	for _, handler := range apiHandlers {
		Register(router, handler)
	}
	fmt.Println("Reached 2")
	launchServer(router)
}

func Register(router *mux.Router, handler apiHandler.Handler) *mux.Route {
	h := createHTTPHandler(handler)
	route := router.HandleFunc(handler.URL(), h.ServeHTTP)
	methods := handler.Methods()
	if len(methods) > 0 {
		route.Methods(methods...)
	}
	return route
}

func createHTTPHandler(handler apiHandler.Handler) http.Handler {

	hf := func(w http.ResponseWriter, r *http.Request) {

		request, err := handler.ParseArgs(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err2 := json.NewEncoder(w).Encode(http.Response{
				Status:     err.Error(),
				StatusCode: http.StatusBadRequest,
			})
			if err2 != nil {
				// this should be a log instead of print in real application
				fmt.Printf("createHTTPHandler: Error encoding response writer")
			}
			return
		}

		resp := handler.Process(request)
		fmt.Println("rawResponse:", resp)
		if resp.StatusCode >= http.StatusBadRequest {
			// Todo
		}
		w.WriteHeader(resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			// fmt.Println("body", body)
			if _, err := w.Write(body); err != nil {
				log.Fatalf("Error writing response body to writer")
			}
		}
		fmt.Println("last response")
		return
	}
	return httpHandler{
		serveHTTP: hf,
	}
}

// httpHandler fulfills the http.Handler interface, allowing us to use logging http middleware
//
type httpHandler struct {
	serveHTTP http.HandlerFunc
}

// ServeHTTP calls through to the constructed function
//
func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveHTTP(w, r)
}

func launchServer(router *mux.Router) {
	port := 8085
	server := &http.Server{
		Addr:    fmt.Sprintf("http://localhost:%v", port),
		Handler: HTTPMiddleware(router),
	}
	fmt.Println("listening")
	conn, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		fmt.Printf("error listing to %v port, err: %s", port, err.Error())
	}
	if err = server.Serve(conn); err != nil {
		fmt.Printf("server encountered err: %s", err)
	}
}

// HTTPMiddleware provides logging/tracing for incoming http requests.
func HTTPMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		h.ServeHTTP(w, request)
	})
}
