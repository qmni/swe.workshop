package model

import "time"

type PlayerClass string

const (
	PlayerClassWarrior PlayerClass = "WARRIOR"
	PlayerClassMage    PlayerClass = "MAGE"
	PlayerClassRogue   PlayerClass = "ROGUE"
	PlayerClassPriest  PlayerClass = "PRIEST"
	PlayerClassHunter  PlayerClass = "HUNTER"
)

type PlayerStatus string

const (
	PlayerStatusActive  PlayerStatus = "ACTIVE"
	PlayerStatusBanned  PlayerStatus = "BANNED"
	PlayerStatusDeleted PlayerStatus = "DELETED"
)

type Player struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Username    string       `json:"username" gorm:"not null;unique"`
	Email       string       `json:"email" gorm:"not null;unique"`
	Level       int          `json:"level" gorm:"not null;default:1"`
	Experience  int          `json:"experience" gorm:"not null;default:0"`
	PlayerClass PlayerClass  `json:"playerClass" gorm:"column:playerClass;not null"`
	Status      PlayerStatus `json:"status" gorm:"not null;default:ACTIVE"`
	GuildID     *uint        `json:"guildId,omitempty" gorm:"column:guildId"`
	CreatedAt   time.Time    `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt" gorm:"column:updatedAt"`
	Version     int          `json:"version" gorm:"not null;default:0"`
}

func (Player) TableName() string {
	return "player"
}

type Guild struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique"`
	Description string    `json:"description"`
	FoundedAt   time.Time `json:"foundedAt" gorm:"column:foundedAt"`
	Version     int       `json:"version" gorm:"not null;default:0"`
}

func (Guild) TableName() string {
	return "guild"
}
