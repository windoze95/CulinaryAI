package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// type User struct {
// 	gorm.Model
// 	Username         string `gorm:"unique"`
// 	Email            string `gorm:"unique"`
// 	HashedPassword   string
// 	Settings         UserSettings `gorm:"foreignKey:UserID"`
// 	CollectedRecipes []Recipe     `gorm:"many2many:user_collected_recipes;"`
// 	GuidingContent   GuidingContent
// }

type User struct {
	gorm.Model
	Username         string         `gorm:"unique;index"`
	Email            *string        `gorm:"unique;default:null"`
	FacebookID       *string        `gorm:"unique;default:null;index"`
	Settings         UserSettings   `gorm:"foreignKey:UserID"`
	CollectedRecipes []Recipe       `gorm:"many2many:user_collected_recipes;"`
	GuidingContent   GuidingContent `gorm:"foreignKey:UserID"`
	Auth             UserAuth       `gorm:"foreignKey:UserID"`
	Subscription     Subscription   `gorm:"foreignKey:UserID"`
}

type UserAuth struct {
	gorm.Model
	UserID         uint `gorm:"unique;index"`
	HashedPassword *string
	AuthType       string `gorm:"type:enum('standard','facebook');default:'standard'"`
}

type SubscriptionTier string

const (
	Free       SubscriptionTier = "Free"    // Free
	ThirtyUses SubscriptionTier = "30-Uses" // Basic
	NinetyUses SubscriptionTier = "90-Uses" // Premium
)

type Subscription struct {
	gorm.Model
	UserID           uint `gorm:"unique;index"`
	ExpiresAt        time.Time
	SubscriptionTier SubscriptionTier `gorm:"type:enum('Free','30-Uses','90-Uses');index"`
	RemainingUses    int
}

type UserSettings struct {
	gorm.Model
	UserID             uint `gorm:"unique;index"`
	EncryptedOpenAIKey string
}

type GuidingContent struct {
	gorm.Model
	UserID       uint `gorm:"unique;index"`
	UID          uuid.UUID
	UnitSystem   int    `gorm:"default:1"` // 1 = US Customary, 2 = Metric
	Requirements string // Additional instructions or guidelines
	// DietaryRestrictions string // Specific dietary restrictions
	// SupportingResearch string // Supporting research to help convey the user's expectations
}

func (gc *GuidingContent) GetUnitSystemName() string {
	switch gc.UnitSystem {
	case 1:
		return "US Customary"
	case 2:
		return "Metric"
	default:
		log.Println("Invalid Unit System used, defaulting to US Customary")
		return "US Customary"
	}
}

func (gc *GuidingContent) BeforeCreate(tx *gorm.DB) (err error) {
	if gc.UnitSystem != 1 && gc.UnitSystem != 2 {
		gc.UnitSystem = 1 // Default to 1 if not 1 or 2
	}
	return nil
}

func (gc *GuidingContent) BeforeUpdate(tx *gorm.DB) (err error) {
	if gc.UnitSystem != 1 && gc.UnitSystem != 2 {
		gc.UnitSystem = 1 // Default to 1 if not 1 or 2
	}
	return nil
}
