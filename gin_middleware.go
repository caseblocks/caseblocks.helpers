package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func DatabaseSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := NewConsoleLogger()
		db := NewSqlConnection(FindDBConnString(), logger)
		defer db.Close()
		c.Set("db", db)

		udb := NewSqlConnection(FindUserDBConnString(), logger)
		defer udb.Close()
		c.Set("userdb", udb)

		c.Next()
	}
}

func RailsSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("userdb").(*sqlx.DB)
		if user, err := FindUserFromId(c.Request, c.Writer, db); err == nil {
			c.Set("currentUser", user)
		}
		c.Next()
	}
}
