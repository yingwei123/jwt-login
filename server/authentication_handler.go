package server

import (
	"encoding/json"
	"net/http"

	"mongoTest.io/util"

	"golang.org/x/crypto/bcrypt"
)

type NewUser struct {
	Email     string `bson:"email" json:"email"`
	Password  string `bson:"password" json:"password"`
	TokenHash string `bson:"tokenHash" json:"tokenHash"`
	UserID    string `bson:"_id" json:"_id"`
}

func (rt Router) SignUpNewUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var newUser NewUser

		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		newUser.TokenHash = util.GenerateRandomString(15)
		newUser.Password = string(hashedPass)

		_, err = rt.MongoDBClient.CreateNewUser(newUser)

		if err != nil {
			response, _ := json.Marshal(Response{Status: http.StatusBadRequest, Message: "Failed to Create new User. Email is already in use"})
			w.Write(response)

			return
		}

		response, _ := json.Marshal(Response{Status: http.StatusOK, Message: "User Created"})
		w.Write(response)

		return
	}
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (rt Router) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var credentials Credentials

		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := rt.MongoDBClient.FindUserByEmail(credentials.Email)
		if err != nil {
			response, _ := json.Marshal(Response{Status: http.StatusBadRequest, Message: "Username does not exist. Sign Up!"})
			w.Write(response)

			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			response, _ := json.Marshal(Response{Status: http.StatusUnauthorized, Message: "Password is incorrect."})
			w.Write(response)

			return
		}

		token, expiration, err := rt.Authenticator.GenerateValidToken(user.Email)
		if err != nil {
			response, _ := json.Marshal(Response{Status: http.StatusUnauthorized, Message: "Could not create access token"})
			w.Write(response)

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "access_token",
			Value:   token + " " + user.UserID,
			Expires: expiration,
		})

		response, _ := json.Marshal(Response{Status: http.StatusOK, Message: ""})
		w.Write(response)

		return
	}
}
