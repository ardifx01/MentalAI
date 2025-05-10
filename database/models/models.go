package models

import "time"

type Akun struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	Nama      string    `json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PersonalSurvey struct {
	ID                          uint   `json:"id" gorm:"primaryKey"`
	PrimaryGoals                string `json:"primary_goals"`
	Experience                  string `json:"experience"`
	Frequency                   string `json:"frequency"`
	PrefferedCommunicationStyle string `json:"preffered_communication_style"`
	AkunID                      uint64 `gorm:"index"`
	Akun                        Akun   `gorm:"foreignKey:AkunID;constraint:OnDelete:CASCADE"`
}

type Percakapan struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	AkunID       uint64    `json:"akun_id"`
	Akun         Akun      `gorm:"foreignKey:AkunID;constraint:OnDelete:CASCADE"`
	Judul        string    `gorm:"type:varchar(100)" json:"judul"`
	UrgencyLevel int       `gorm:"type:int;check:urgency_level IN (1,2,3,4,5)" json:"urgency_level"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Omongan      []Omongan `json:"omongan" gorm:"foreignKey:PercakapanID"`
}

type Omongan struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Pesan        string     `json:"pesan"`
	Pengirim     string     `gorm:"type:varchar(10);check:pengirim IN ('user','bot')" json:"pengirim"`
	PercakapanID uint       `gorm:"index"`
	Percakapan   Percakapan `gorm:"foreignKey:PercakapanID;constraint:OnDelete:CASCADE"`
	CreatedAt    time.Time  `json:"created_at"`
}
