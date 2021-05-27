package handlers

import "net/http"

// CORS file for enabling https requests. This file should be exactly the same as our
// previous assignments.

type Cors struct {
	handler http.Handler
}

//Sets up the headers for our CORS enabled mux, with required API requests
func (cors *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
	w.Header().Set("Access-Control-Max-Age", "600")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	}
	cors.handler.ServeHTTP(w, r)
}

//returns CORS wrapped handler
func NewCors(handlerToWrap http.Handler) *Cors {
	return &Cors{handlerToWrap}
}
