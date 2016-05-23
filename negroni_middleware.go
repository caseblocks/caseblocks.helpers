package helpers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/negroni"
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
	if user, err := FindUserFromId(req, w, m.UserDB); err == nil {
		req.Header.Set(HDR_CASEBLOCKS_USER, strconv.Itoa(int(user.Id)))
		next(w, req)
	} else {
		req.Header.Del(HDR_CASEBLOCKS_USER)
	}
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
