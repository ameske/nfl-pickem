package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ameske/nfl-pickem/jsonhttp"

	"golang.org/x/crypto/bcrypt"
)

type PasswordUpdater interface {
	UpdatePassword(username string, oldPassword string, newPassword string) error
}

// ChangePassword processes the password change form, informing the user of any problems or success.
func ChangePassword(db PasswordUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _ := currentUser(r)
		if u == "" {
			jsonhttp.WriteError(w, http.StatusUnauthorized, "login required")
			return
		}

		r.ParseForm()

		p := r.FormValue("oldPassword")
		pN := r.FormValue("newPassword")

		bpass, err := bcrypt.GenerateFromPassword([]byte(pN), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			jsonhttp.WriteError(w, http.StatusInternalServerError, "contact admin")
			return
		}

		err = db.UpdatePassword(u, p, string(bpass))
		if err != nil {
			log.Println(err)
			jsonhttp.WriteError(w, http.StatusInternalServerError, "contact admin")
			return
		}

		jsonhttp.WriteSuccess(w, fmt.Sprintf("Succesfully changed password for user %s", u))
	}
}
