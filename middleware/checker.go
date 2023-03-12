package middleware

import (
	"net/http"
	"strconv"
)

func RegisterChecker() Checker {
	var newCheck Checker
	newCheck.TotalAccess = 0
	newCheck.UserAccess = make(map[string]int)
	return newCheck
}

type Checker struct {
	TotalAccess int
	// It contains the how much times a user have access the route.
	// If a user is not yet accessed any of the route it will return 0.
	UserAccess map[string]int
}

func (check *Checker) Accessed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		check.TotalAccess++
		name := r.Header.Get("username")
		var count int = check.UserAccess[name]
		check.UserAccess[name] = count + 1
		next.ServeHTTP(w, r)
	})
}

func (check *Checker) GetAccessCount(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("username")
	w.Write([]byte("Route Access: " + strconv.Itoa(check.UserAccess[name])))
}
