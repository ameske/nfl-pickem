package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ameske/nfl-pickem"

	"golang.org/x/crypto/bcrypt"
)

// ChangePassword processes the password change form, informing the user of any problems or success.
func changePassword(db nflpickem.PasswordUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := retrieveUser(r.Context())
		if err == errNoUser {
			WriteJSONError(w, http.StatusUnauthorized, "login required")
			return
		} else if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		p := r.FormValue("oldPassword")
		pN := r.FormValue("newPassword")

		bpass, err := bcrypt.GenerateFromPassword([]byte(pN), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			WriteJSONError(w, http.StatusInternalServerError, "contact admin")
			return
		}

		err = db.UpdatePassword(user.Email, p, string(bpass))
		if err != nil {
			log.Println(err)
			WriteJSONError(w, http.StatusInternalServerError, "contact admin")
			return
		}

		WriteJSONSuccess(w, fmt.Sprintf("Succesfully changed password for user %s", user.Email))
	}
}
