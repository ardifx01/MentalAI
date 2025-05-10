package middleware

import (
	"chatbot/database"
	"chatbot/database/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func DapatinAkun(akun_id any) (models.Akun, bool) {
	akun := models.Akun{}
	if akun_id == nil {
		return akun, false
	}

	database.Database.Where("id = ?", akun_id).First(&akun)

	if akun.ID == 0 {
		return akun, false
	}

	return akun, true
}

func CekAutentikasiHandler(redirectToLogin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("akun") == nil && redirectToLogin {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		if session.Get("akun") != nil && !redirectToLogin {
			c.Redirect(http.StatusFound, "/dashboard")
			return
		}

		c.Next()
	}
}

func CekAkunValid() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		akunID := session.Get("akun")
		if akunID == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		akun, valid := DapatinAkun(akunID)
		if !valid {
			session.Delete("akun")
			session.Save()

			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.Set("akun", akun)
		c.Next()
	}
}
