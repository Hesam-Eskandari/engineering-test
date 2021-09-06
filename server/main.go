package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/api/menus"
	"github.com/flypay/engineering-test/pkg/api/orders"
)

func main() {
	apiHandlers := make([]apiHandler.Handler, 0, 5)
	apiHandlers = append(apiHandlers, menus.NewGetAlphaMenu())
	apiHandlers = append(apiHandlers, menus.NewGetBetaMenu())
	apiHandlers = append(apiHandlers, orders.NewSetOrder())
	//router := mux.NewRouter().StrictSlash(false)  // uncomment to use mux server
	for _, handler := range apiHandlers {
		//RegisterMux(router, handler)
		RegisterHTTP(handler)
	}

	// go launchMuxServer(router) // uncomment to use mux server
	launchHTTPServer()

}

func RegisterHTTP(handler apiHandler.Handler) {
	h := createHTTPHandler(handler)
	http.HandleFunc(handler.URL(), func(w http.ResponseWriter, r *http.Request) {
		notAllowed := true
		for _, method := range handler.Methods() {
			if method == r.Method {
				h.ServeHTTP(w, r)
				notAllowed = false
			}
		}
		if notAllowed {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	})
}

func launchHTTPServer() {
	addrFlag := flag.String("addr", ":8086", "address to run unified server on")
	fmt.Println("http server ready")
	_ = http.ListenAndServe(*addrFlag, nil)
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
				log.Printf("createHTTPHandler: Error encoding response writer")
			}
			return
		}

		resp := handler.Process(request)
		w.WriteHeader(resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			if _, err = w.Write(body); err != nil {
				log.Fatalf("Error writing response body to writer. err: %s", err.Error())
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

//func launchMuxServer(router *mux.Router) {
//	port := 8085
//	server := &http.Server{
//		Addr:    fmt.Sprintf(":%v", port),
//		Handler: HTTPMiddleware(router),
//	}
//	conn, err := net.Listen("tcp", server.Addr)
//	if err != nil {
//		fmt.Printf("error listing to %v port, err: %s", port, err.Error())
//	}
//	fmt.Println("mux server ready")
//	if err = server.Serve(conn); err != nil {
//		fmt.Printf("server encountered err: %s", err)
//	}
//}

//func RegisterMux(router *mux.Router, handler apiHandler.Handler) *mux.Route {
//	h := createHTTPHandler(handler)
//	route := router.HandleFunc(handler.URL(), h.ServeHTTP)
//	methods := handler.Methods()
//	if len(methods) > 0 {
//		route.Methods(methods...)
//	}
//	return route
//}

// HTTPMiddleware can provide logging/tracing for incoming http requests.
//func HTTPMiddleware(h http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
//		h.ServeHTTP(w, request)
//	})
//}
