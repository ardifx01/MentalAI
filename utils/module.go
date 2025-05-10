package utils

import (
	"chatbot/database/models"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func DapatinLastValueOmongan(Percakapan models.Percakapan) string {
	Omongan := Percakapan.Omongan
	if len(Omongan) == 0 {
		return ""
	}

	return Omongan[len(Omongan)-1].Pesan
}

func ConvertWaktuKeString(waktu time.Time) string {
	// Convert Time to string with format "Today, 10:00 AM" or "Yesterday, 10:00 AM" or if it's older than 2 days, return the date in "May 26, 2005" format
	if waktu.After(time.Now().Add(-24 * time.Hour)) {
		return waktu.Format("Today, 15:04 PM")
	}
	// Check if the date is yesterday
	if waktu.After(time.Now().Add(-48 * time.Hour)) {
		return "Yesterday, " + waktu.Format("15:04 PM")
	}

	return waktu.Format("January 2, 2006")
}

func DefaultDataGin(c *gin.Context) gin.H {
	session := sessions.Default(c)

	err := session.Get("error")

	session.Delete("error")
	session.Save()

	return gin.H{
		"url":   strings.Split(c.Request.URL.Path, "/")[1],
		"akun":  session.Get("akun"),
		"error": err,
	}
}
