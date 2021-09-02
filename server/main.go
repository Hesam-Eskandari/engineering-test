package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/flypay/engineering-test/pkg/api/alpha"
	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/api/beta"
	"github.com/flypay/engineering-test/pkg/api/orders"

	"github.com/gorilla/mux"
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
		RegisterMux(router, handler)
		RegisterHTTP(handler)
	}
	go launchHTTPServer()
	launchMuxServer(router)
}

func RegisterMux(router *mux.Router, handler apiHandler.Handler) *mux.Route {
	h := createHTTPHandler(handler)
	route := router.HandleFunc(handler.URL(), h.ServeHTTP)
	methods := handler.Methods()
	if len(methods) > 0 {
		route.Methods(methods...)
	}
	return route
}

func RegisterHTTP(handler apiHandler.Handler) {
	h := createHTTPHandler(handler)
	http.HandleFunc(handler.URL(), func(w http.ResponseWriter, r *http.Request) {
		for _, method := range handler.Methods() {
			if method == r.Method {
				h.ServeHTTP(w, r)
			}
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	})
}

func launchHTTPServer() {
	addrFlag := flag.String("addr", ":8086", "address to run unified server on")
	fmt.Println("http server ready")
	_ = http.ListenAndServe(*addrFlag, nil)
}

func launchMuxServer(router *mux.Router) {
	port := 8085
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: HTTPMiddleware(router),
	}
	conn, err := net.Listen("tcp", server.Addr)
	if err != nil {
		fmt.Printf("error listing to %v port, err: %s", port, err.Error())
	}
	fmt.Println("mux server ready")
	if err = server.Serve(conn); err != nil {
		fmt.Printf("server encountered err: %s", err)
	}
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
		return
	}
	return httpHandler{
		serveHTTP: hf,
	}
}

// httpHandler fulfills the http.Handler interface, allowing us to use logging http middleware
type httpHandler struct {
	serveHTTP http.HandlerFunc
}

// ServeHTTP calls through to the constructed function
func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveHTTP(w, r)
}

// HTTPMiddleware can provide logging/tracing for incoming http requests.
func HTTPMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		h.ServeHTTP(w, request)
	})
}
