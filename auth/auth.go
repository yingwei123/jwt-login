package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Authenticator struct {
	JWTKEY string
}
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateAuthenticator(JWTKEY string) Authenticator {
	return Authenticator{JWTKEY: JWTKEY}
}

func (a Authenticator) GenerateValidToken(userName string) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: userName,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(a.JWTKEY))
	if err != nil {
		return tokenString, expirationTime, err
	}

	println(tokenString)

	return tokenString, expirationTime, err
}

func (a Authenticator) ValidateRequest(r *http.Request) (string, error) {
	c, err := r.Cookie("access_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", errors.New("no access token")
			// If the cookie is not set, return an unauthorized status
		}
		// For any other type of error, return a bad request status

		return "", errors.New("could not get access token from cookie")
	}

	// Get the JWT string from the cookie

	content := strings.Split(c.Value, " ")

	tknStr := content[0]

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.JWTKEY), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New("unauthorized access")
		}

		return "", errors.New("internal server issue")
	}
	if !tkn.Valid {
		return "", errors.New("unauthorized access")
	}

	return content[1], nil
}
