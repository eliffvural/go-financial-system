package api

import (
	"fmt"
	"log"
	"net/http"
)

// HandlerFunc, custom router için handler fonksiyon tipidir
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Route, bir endpoint ve ona karşılık gelen handler'ı tutar
type Route struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

// Router, custom router yapısıdır
// Route'ları ve middleware'leri tutar
type Router struct {
	routes     []Route
	middleware []func(HandlerFunc) HandlerFunc
}

// Yeni bir Router oluşturur
func NewRouter() *Router {
	return &Router{}
}

// Route ekler
func (r *Router) Handle(method, path string, handler HandlerFunc) {
	r.routes = append(r.routes, Route{Method: method, Path: path, Handler: handler})
}

// Middleware ekler
func (r *Router) Use(mw func(HandlerFunc) HandlerFunc) {
	r.middleware = append(r.middleware, mw)
}

// HTTP isteklerini karşılar
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if req.Method == route.Method && req.URL.Path == route.Path {
			h := route.Handler
			// Middleware zincirini uygula
			for i := len(r.middleware) - 1; i >= 0; i-- {
				h = r.middleware[i](h)
			}
			h(w, req)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 Not Found")
}

// Server, HTTP sunucusunu başlatır
func StartServer(addr string, router *Router) {
	fmt.Printf("Sunucu %s adresinde başlatılıyor...\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}
