package helpers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type NegroniAuthenticatedSessionMiddleware struct {
	UserDB *sqlx.DB
}

// Middleware is a struct that has a ServeHTTP method
func NewAuthenticationSessionMiddleware(userdb *sqlx.DB) *NegroniAuthenticatedSessionMiddleware {
	return &NegroniAuthenticatedSessionMiddleware{UserDB: userdb}
}

const HDR_CASEBLOCKS_USER = "X-CASEBLOCKS-USER"

func (m *NegroniAuthenticatedSessionMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authorization, authOK := req.Header["Authorization"]
	jwtSessionKey := os.Getenv("JWT_SESSION_KEY")
	if authOK {
		tokenString := authorization[0][7:]
		if len(jwtSessionKey) == 0 {
			fmt.Println("WARNING: Unable to decode JWToken. Please specify JWT_SESSION_KEY envvar.")
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(jwtSessionKey), nil
		})
		if err == nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				var userID FKInt
				switch v := claims["id"].(type) {
				case int:
					userID = FKInt(v)
				case float64:
					userID = FKInt(v)
				}
				user, err := NewUserRepository(m.UserDB).FindById(userID)
				if err == nil {
					req.Header.Set(HDR_CASEBLOCKS_USER, strconv.Itoa(int(user.Id)))
					next(w, req)
				}
			}
		}
	} else if user, err := FindUserFromId(req, w, m.UserDB); err == nil {
		req.Header.Set(HDR_CASEBLOCKS_USER, strconv.Itoa(int(user.Id)))
		next(w, req)
	}
	req.Header.Del(HDR_CASEBLOCKS_USER)
}

func GetUserFromRequest(userRepository UserRepository, req *http.Request) (User, error) {
	id_str := req.Header.Get(HDR_CASEBLOCKS_USER)
	if id, err := strconv.Atoi(id_str); err != nil {
		return User{}, fmt.Errorf("Invalid user in header: %s", id_str)
	} else {
		return userRepository.FindById(FKInt(id))
	}

}

/* Generates a negroni handler for the route:
   pathType (base+string)
   which is handled by `f` and passed through
   middleware `mids` sequentially.

   ex:
   NegroniRoute(router, "/api/v1", "/users/update", "POST", UserUpdateHandler, LoggingMiddleware, AuthorizationMiddleware)

	 See https://gist.github.com/pagreczner/227913794e096953972c
*/
func NegroniRoute(m *mux.Router, base string, path string, pathType string, f http.Handler, mids ...func(http.ResponseWriter, *http.Request, http.HandlerFunc)) {
	_routes := mux.NewRouter()
	_routes.Handle(base+path, f).Methods(pathType)

	_n := negroni.New()
	for _, mid := range mids {
		_n.Use(negroni.HandlerFunc(mid))
	}
	_n.UseHandler(_routes)
	m.Handle(path, _n)
}

func GetPort(port int) string {
	if len(os.Getenv("PORT")) > 0 {
		return fmt.Sprintf(":%s", os.Getenv("PORT"))
	}
	return fmt.Sprintf(":%d", port)
}
