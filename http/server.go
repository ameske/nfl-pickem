package http

import (
	"context"
	"encoding/gob"
	"errors"
	"log"
	"net/http"

	"github.com/ameske/nfl-pickem"
	"github.com/gorilla/securecookie"
)

// A Server exposes the NFL Pickem Service over HTTP
type Server struct {
	Address string
	router  *http.ServeMux
	sc      *securecookie.SecureCookie
	db      nflpickem.Service
}

func NewServer(address string, hashKey []byte, encryptKey []byte, db nflpickem.Service) (*Server, error) {
	sc := securecookie.New(hashKey, encryptKey)

	s := &Server{
		Address: address,
		router:  http.NewServeMux(),
		sc:      sc,
		db:      db,
	}

	gob.Register(nflpickem.User{})

	s.router.HandleFunc("/login", s.login)
	s.router.HandleFunc("/logout", s.logout)

	s.router.HandleFunc("/current", currentWeek(db))
	s.router.HandleFunc("/games", games(db))
	s.router.HandleFunc("/results", results(db))
	s.router.HandleFunc("/totals", weeklyTotals(db))

	s.router.HandleFunc("/picks", s.requireLogin(picks(db)))
	s.router.HandleFunc("/password", s.requireLogin(changePassword(db)))

	return s, nil
}

func (s *Server) Start() error {
	log.Printf("NFL Pick-Em Pool listening on %s", s.Address)
	return http.ListenAndServe(s.Address, s.router)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()
	if !ok {
		WriteJSONError(w, http.StatusBadRequest, "missing credentials")
		return
	}

	user, err := s.db.CheckCredentials(u, p)
	if err != nil {
		log.Println(err)
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	cookie, err := s.encodeCookie("nflpickem", user)
	if err != nil {
		log.Println(err)
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, cookie)

	WriteJSONSuccess(w, "successfully logged in")
}

func (s *Server) encodeCookie(name string, value interface{}) (*http.Cookie, error) {
	encoded, err := s.sc.Encode(name, value)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}, nil
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("nflpickem")
	if err != nil && err != http.ErrNoCookie {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	WriteJSONSuccess(w, "succesful logout")
}

func (s *Server) requireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := s.verifyLogin(w, r)
		if err != nil {
			// Regardless of the path here, let's just premptively clear this cookie out
			cookie := &http.Cookie{
				Name:   "nflpickem",
				MaxAge: -1,
			}
			http.SetCookie(w, cookie)
			WriteJSONError(w, http.StatusUnauthorized, "login required")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		next(w, r.WithContext(ctx))
	}
}

var errNoUser = errors.New("no user information stored in context")
var errNoLogin = errors.New("no login information found")

func retrieveUser(ctx context.Context) (nflpickem.User, error) {
	u, ok := ctx.Value("user").(nflpickem.User)
	if !ok {
		return nflpickem.User{}, errNoUser
	}

	return u, nil
}

func (s *Server) verifyLogin(w http.ResponseWriter, r *http.Request) (nflpickem.User, error) {
	cookie, err := r.Cookie("nflpickem")
	if err == nil {
		user := nflpickem.User{}
		if err := s.sc.Decode("nflpickem", cookie.Value, &user); err == nil {
			return user, nil
		}
	}

	u, p, ok := r.BasicAuth()
	if !ok {
		return nflpickem.User{}, errNoLogin
	}

	user, err := s.db.CheckCredentials(u, p)
	if err != nil {
		return nflpickem.User{}, err

	}

	cookie, err = s.encodeCookie("nflpickem", user)
	if err != nil {
		return nflpickem.User{}, err
	}

	http.SetCookie(w, cookie)

	return user, nil
}
