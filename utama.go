package main

import (
	"chatbot/chat"
	"chatbot/database"
	"chatbot/database/models"
	"chatbot/middleware"
	"chatbot/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

/*********             FUNCTIONS              *********/
func selectionSort(percakapan []*models.Percakapan, urutan string) {
	for i := 0; i < len(percakapan); i++ {
		maxIndex := i
		for j := i + 1; j < len(percakapan); j++ {
			kondisi := false
			if urutan == "asc" {
				kondisi = percakapan[j].UrgencyLevel < percakapan[maxIndex].UrgencyLevel
			} else if urutan == "desc" {
				kondisi = percakapan[j].UrgencyLevel > percakapan[maxIndex].UrgencyLevel
			}

			if kondisi {
				maxIndex = j
			}
		}

		temp := percakapan[i]
		percakapan[i] = percakapan[maxIndex]
		percakapan[maxIndex] = temp
	}
}

func insertionSort(percakapan []*models.Percakapan, urutan string) {
	for i := 1; i < len(percakapan); i++ {
		key := percakapan[i]
		j := i - 1

		kondisi := false
		if urutan == "asc" {
			kondisi = percakapan[j].CreatedAt.Before(key.CreatedAt)
		} else if urutan == "desc" {
			kondisi = percakapan[j].CreatedAt.After(key.CreatedAt)
		}

		for j >= 0 && kondisi {
			percakapan[j+1] = percakapan[j]
			j = j - 1
		}

		percakapan[j+1] = key
	}
}

func bubbleSort(percakapan []*models.Percakapan, urutan string) {
	for i := 0; i < len(percakapan)-1; i++ {
		for j := 0; j < len(percakapan)-i-1; j++ {
			urut := false

			if urutan == "asc" {
				urut = percakapan[j].Judul > percakapan[j+1].Judul
			} else if urutan == "desc" {
				urut = percakapan[j].Judul < percakapan[j+1].Judul
			}

			if urut {
				temp := percakapan[j+1]
				percakapan[j+1] = percakapan[j]
				percakapan[j] = temp
			}
		}
	}
}

func SequentialSearch(percakapan []*models.Percakapan, yangDiCari string) []*models.Percakapan {
	var FilterSemuaPercakapan []*models.Percakapan

	for i := 0; i < len(percakapan); i++ {
		percakapan := percakapan[i]
		if percakapan.Judul == yangDiCari || strings.Contains(percakapan.Judul, yangDiCari) {
			FilterSemuaPercakapan = append(FilterSemuaPercakapan, percakapan)
		}
	}

	return FilterSemuaPercakapan
}

func BinarySearch(percakapan []*models.Percakapan, yangDiCari string, urutan string) int {
	var kiri, kanan, tengah int
	ketemu := -1
	kiri = 0
	kanan = len(percakapan) - 1

	for kiri <= kanan && ketemu == -1 {
		tengah = (kanan + kiri) / 2

		if percakapan[tengah].Judul == yangDiCari {
			ketemu = tengah
		}

		if urutan == "a-z" {
			if percakapan[tengah].Judul < yangDiCari {
				kiri = tengah + 1
			} else {
				kanan = tengah - 1
			}
		} else if urutan == "z-a" {
			if percakapan[tengah].Judul > yangDiCari {
				kiri = tengah + 1
			} else {
				kanan = tengah - 1
			}
		}
	}

	return ketemu
}

/*********       PAGE HANDLER FUNCTIONS       *********/
func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Welcome to MindfulAI",
	})
}

func DashboardPage(c *gin.Context) {
	session := sessions.Default(c)

	akunID := session.Get("akun")
	percakapan_id := c.Query("chat")

	DataGin := utils.DefaultDataGin(c)
	DataGin["title"] = "Dashboard - MindfulAI"

	var percakapan models.Percakapan

	if percakapan_id != "" {
		database.Database.Preload("Omongan").Where("id = ? AND akun_id = ?", percakapan_id, akunID).First(&percakapan)

		n := len(percakapan.Omongan)
		for i := 0; i < n-1; i++ {
			for j := 0; j < n-i-1; j++ {
				if percakapan.Omongan[j].CreatedAt.After(percakapan.Omongan[j+1].CreatedAt) {
					temp := percakapan.Omongan[j]
					percakapan.Omongan[j] = percakapan.Omongan[j+1]
					percakapan.Omongan[j+1] = temp
				}
			}
		}

		percakapanJSON, err := json.Marshal(percakapan)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert conversation to JSON"})
			return
		}

		DataGin["dipilih"] = percakapanJSON
	}

	c.HTML(http.StatusOK, "dashboard.html", DataGin)
}

func SejarahPage(c *gin.Context) {
	session := sessions.Default(c)
	akunID := session.Get("akun")

	SearchParameter := c.Query("search")
	SortParameter := c.Query("sort")

	SemuaPercakapan := []*models.Percakapan{}
	if err := database.Database.Preload("Omongan").Where("akun_id = ?", akunID).Find(&SemuaPercakapan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
		return
	}

	if SortParameter == "urgency_lowest" {
		selectionSort(SemuaPercakapan, "asc")
	}

	if SortParameter == "urgency_highest" {
		selectionSort(SemuaPercakapan, "desc")
	}

	if SortParameter == "newest" {
		insertionSort(SemuaPercakapan, "asc")
	}

	if SortParameter == "oldest" {
		insertionSort(SemuaPercakapan, "desc")
	}

	if SortParameter == "a-z" {
		bubbleSort(SemuaPercakapan, "asc")
	}

	if SortParameter == "z-a" {
		bubbleSort(SemuaPercakapan, "desc")
	}

	if SearchParameter != "" {
		if SortParameter != "a-z" && SortParameter != "z-a" {
			SemuaPercakapan = SequentialSearch(SemuaPercakapan, SearchParameter)
		} else {
			ketemu := BinarySearch(SemuaPercakapan, SearchParameter, SortParameter)
			if ketemu == -1 {
				SemuaPercakapan = []*models.Percakapan{}
			} else {
				SemuaPercakapan = []*models.Percakapan{
					SemuaPercakapan[ketemu],
				}
			}
		}
	}

	var PercakapanDiPilih *models.Percakapan
	if len(SemuaPercakapan) > 0 {
		PercakapanDiPilih = SemuaPercakapan[0]
	}

	percakapanJSON, err := json.Marshal(PercakapanDiPilih)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert conversation to JSON"})
		return
	}

	c.HTML(http.StatusOK, "sejarah.html", gin.H{
		"title":      "History - MindfulAI",
		"akun":       akunID,
		"percakapan": SemuaPercakapan,
		"dipilih":    percakapanJSON,
	})
}

func LoginPage(c *gin.Context) {
	session := sessions.Default(c)

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login - MindfulAI",
		"error": session.Get("error"),
	})
}

func RegisterPage(ctx *gin.Context) {
	session := sessions.Default(ctx)

	ctx.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register - MindfulAI",
		"error": session.Get("error"),
	})

	session.Delete("error")
	session.Save()
}

/*********       REST API       *********/

func ChatHandler(c *gin.Context) {
	session := sessions.Default(c)

	akunID := session.Get("akun")
	akun, valid := middleware.DapatinAkun(akunID)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pesan := c.PostForm("pesan")
	fmt.Println(pesan)

	if pesan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message cannot be empty"})
		return
	}

	if len(pesan) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message is too long"})
		return
	}

	percakapanID := c.PostForm("percakapan_id")
	urgency_level := c.PostForm("urgency_level")
	mood := c.PostForm("mood")
	moodNotes := c.PostForm("mood_notes")

	var percakapan models.Percakapan
	if percakapanID != "" {
		if err := database.Database.Preload("Omongan", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).First(&percakapan, percakapanID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Conversation not found"})
			return
		}
	} else {
		if urgency_level == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Urgency level is required"})
			return
		}

		intUrgencyLevel, err := strconv.Atoi(urgency_level)
		if err != nil || intUrgencyLevel < 1 || intUrgencyLevel > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid urgency level"})
			return
		}

		percakapan = models.Percakapan{
			AkunID:       akun.ID,
			Judul:        chat.DapatinJudulPercakapan(pesan),
			UrgencyLevel: chat.DapatinUrgencyLevel(pesan),
		}

		if err := database.Database.Create(&percakapan).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
			return
		}
	}

	omongan := models.Omongan{
		Pesan:        pesan,
		Pengirim:     "user",
		PercakapanID: percakapan.ID,
	}

	if err := database.Database.Create(&omongan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	contentSejarah := []*genai.Content{}
	for _, v := range percakapan.Omongan {
		var role genai.Role = genai.RoleUser
		if v.Pengirim == "bot" {
			role = genai.RoleModel
		}

		fmt.Println("Text:", v.Pesan, "Role:", role, "- ID:", v.ID)
		contentSejarah = append(contentSejarah, genai.NewContentFromText(v.Pesan, role))
	}

	cs := chat.BuatChat(contentSejarah)

	// No Streaming, Only text | Yucky 🤢
	// hasil, err := chat.KirimPesan(cs, pesan+"\n\nUser's Mood: "+mood+"\nUser's Mood Notes: "+moodNotes)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
	// 	return
	// }

	// Streaming text, like ChatGPT | Sugoi 🤩
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	stream := chat.KirimPesanStream(cs, pesan+"\n\nUser's Mood: "+mood+"\nUser's Mood Notes: "+moodNotes)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "Streaming unsupported!")
		return
	}

	pesanFull := ""
	fmt.Println("SEBELUM STREAM")
	for chunk, err := range stream {
		if err != nil {
			fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			break
		}

		fmt.Fprintf(c.Writer, "data: %s\n\n", chunk.Candidates[0].Content.Parts[0].Text)
		flusher.Flush()
		pesanFull += chunk.Candidates[0].Content.Parts[0].Text
	}

	fmt.Println("SUDAH SELESAI", pesanFull)
	omongan = models.Omongan{
		Pesan:        pesanFull,
		Pengirim:     "bot",
		PercakapanID: percakapan.ID,
	}

	if err := database.Database.Create(&omongan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	final := fmt.Sprintf(`{"percakapan_id":"%d", "urgency_level":"%d", "judul":"%s"}`, percakapan.ID, percakapan.UrgencyLevel, percakapan.Judul)
	fmt.Fprintf(c.Writer, "event: done\ndata: %s\n\n", final)
	flusher.Flush()
}

func DapatinChatHandler(c *gin.Context) {
	session := sessions.Default(c)

	akunID := session.Get("akun")

	percakapanID := c.PostForm("percakapan_id")
	if percakapanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversation ID is required"})
		return
	}

	intPercakapanID, err := strconv.Atoi(percakapanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	intAkunID, ok := akunID.(uint64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var percakapan models.Percakapan
	if err := database.Database.Preload("Omongan").Where(&models.Percakapan{ID: uint(intPercakapanID), AkunID: intAkunID}).First(&percakapan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversation not found"})
		return
	}

	n := len(percakapan.Omongan)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if percakapan.Omongan[j].CreatedAt.After(percakapan.Omongan[j+1].CreatedAt) {
				temp := percakapan.Omongan[j]
				percakapan.Omongan[j] = percakapan.Omongan[j+1]
				percakapan.Omongan[j+1] = temp
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"response": percakapan,
	})
}

func UpdateChatHandler(c *gin.Context) {
	session := sessions.Default(c)

	akunID := session.Get("akun")

	UrgencyLevel := c.PostForm("urgency_level")
	percakapanID := c.PostForm("percakapan_id")

	if percakapanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversation ID is required"})
		return
	}

	var percakapan models.Percakapan
	if err := database.Database.Preload("Omongan").Where("id = ? AND akun_id = ?", percakapanID, akunID).First(&percakapan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversation not found"})
		return
	}

	if UrgencyLevel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Urgency level is required"})
		return
	}

	var UpdatePercakapan models.Percakapan
	intUrgencyLevel, err := strconv.Atoi(UrgencyLevel)
	if err != nil || intUrgencyLevel < 1 || intUrgencyLevel > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid urgency level"})
		return
	}
	UpdatePercakapan.UrgencyLevel = intUrgencyLevel

	if err := database.Database.Model(&percakapan).Updates(UpdatePercakapan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update conversation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sukses"})
}

func RegisterHandler(c *gin.Context) {
	session := sessions.Default(c)

	nama := c.PostForm("nama")
	email := c.PostForm("email")
	password := c.PostForm("password")
	konfirmasipassword := c.PostForm("konfirmasipassword")
	goalstress := c.PostForm("goal-stress")
	goalproductivity := c.PostForm("goal-productivity")
	goalreflection := c.PostForm("goal-reflection")
	goalhabits := c.PostForm("goal-habits")
	experience := c.PostForm("experience")
	frequency := c.PostForm("frequency")
	communicationstyle := c.PostForm("communication-style")
	terms := c.PostForm("terms")

	if terms != "on" {
		c.Redirect(http.StatusFound, "/register")
		return
	}

	if nama == "" || email == "" || password == "" || konfirmasipassword == "" || experience == "" || frequency == "" || communicationstyle == "" || terms == "" {
		session.Set("error", "All fields are required")
		session.Save()

		c.Redirect(http.StatusFound, "/register")
		return
	}

	if len(password) < 6 {
		session.Set("error", "Password must be at least 6 characters")
		session.Save()

		c.Redirect(http.StatusFound, "/register")
		return
	}

	if password != konfirmasipassword {
		session.Set("error", "Passwords do not match")
		session.Save()

		c.Redirect(http.StatusFound, "/register")
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		session.Set("error", "Failed to hash password")
		session.Save()

		c.Redirect(http.StatusFound, "/register")
		return
	}

	akun := models.Akun{
		Nama:     nama,
		Email:    email,
		Password: hashedPassword,
	}

	if err := database.Database.Create(&akun).Error; err != nil {
		session.Set("error", "Failed to create account: "+err.Error())
		session.Save()

		c.Redirect(http.StatusFound, "/register")
		return
	}

	primaryGoal := ""
	if goalstress == "on" {
		primaryGoal += "Stress,"
	}
	if goalproductivity == "on" {
		primaryGoal += "Productivity,"
	}
	if goalreflection == "on" {
		primaryGoal += "Reflection,"
	}
	if goalhabits == "on" {
		primaryGoal += "Habits,"
	}

	if err := database.Database.Create(&models.PersonalSurvey{
		PrimaryGoals:                primaryGoal,
		Frequency:                   frequency,
		PrefferedCommunicationStyle: communicationstyle,
		Experience:                  experience,
		AkunID:                      akun.ID,
	}).Error; err != nil {
		session.Set("error", "Failed to create personal survey. please try again")
		session.Save()

		c.Redirect(http.StatusFound, "/register")
		return
	}

	session.Set("akun", akun.ID)
	session.Save()

	c.Redirect(http.StatusFound, "/dashboard")
}

func HapusChatHandler(c *gin.Context) {
	session := sessions.Default(c)

	akunID := session.Get("akun")
	percakapanID := c.PostForm("percakapan_id")

	if percakapanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversation ID is required"})
		return
	}

	intPercakapanID, err := strconv.Atoi(percakapanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	intAkunID, ok := akunID.(uint64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	if err := database.Database.Where("percakapan_id = ?", intPercakapanID).Delete(&models.Omongan{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete messages"})
		return
	}

	if err := database.Database.Where("id = ? AND akun_id = ?", intPercakapanID, intAkunID).Delete(&models.Percakapan{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete conversation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted successfully"})
}

func LoginHandler(c *gin.Context) {
	session := sessions.Default(c)

	email := c.PostForm("email")
	password := c.PostForm("password")

	var akun models.Akun
	database.Database.Where("email = ?", email).First(&akun)

	if err := database.Database.Where("email = ?", email).First(&akun).Error; err != nil {

		session.Set("error", "Invalid email or password")
		session.Save()

		c.Redirect(http.StatusFound, "/login")
		return
	}

	if !utils.VerifyPassword(password, akun.Password) {
		session.Set("error", "Invalid email or password")
		session.Save()

		c.Redirect(http.StatusFound, "/login")
		return
	}

	session.Set("akun", akun.ID)
	session.Save()

	c.Redirect(http.StatusFound, "/dashboard")
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("akun")
	session.Save()

	c.Redirect(http.StatusFound, "/login")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()
	chat.ClientGenAI()

	router := gin.Default()

	store := cookie.NewStore([]byte(os.Getenv("KUE")))
	router.Use(sessions.Sessions("session", store))
	router.SetFuncMap(template.FuncMap{
		"DapatinLastValueOmongan": utils.DapatinLastValueOmongan,
		"ConvertWaktuKeString":    utils.ConvertWaktuKeString,
	})

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", IndexPage)
	router.GET("/logout", LogoutHandler)

	AuthRedirectLogin := router.Group("/")
	AuthRedirectLogin.Use(middleware.CekAutentikasiHandler(true))
	AuthRedirectLogin.Use(middleware.CekAkunValid())
	{
		AuthRedirectLogin.GET("/dashboard", DashboardPage)
		AuthRedirectLogin.GET("/history", SejarahPage)

		AuthAPI := AuthRedirectLogin.Group("/api")
		AuthAPI.POST("/chat", ChatHandler)
		AuthAPI.POST("/dapatinChat", DapatinChatHandler)
		AuthAPI.POST("/chat/hapus", HapusChatHandler)
		AuthAPI.POST("/chat/update", UpdateChatHandler)
	}

	AuthRedirectDashboard := router.Group("/")
	AuthRedirectDashboard.Use(middleware.CekAutentikasiHandler(false))
	{
		AuthRedirectDashboard.GET("/login", LoginPage)
		AuthRedirectDashboard.GET("/register", RegisterPage)

		AuthRedirectDashboard.POST("/api/register", RegisterHandler)
		AuthRedirectDashboard.POST("/api/login", LoginHandler)
	}

	router.Run(":8080")
}
