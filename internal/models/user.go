package models

import (
	"errors"
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
	FirstName        *string        `gorm:"default:null"`
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
	AuthType       string `gorm:"type:text"`
}

type UserAuthType string

const (
	Standard UserAuthType = "standard"
	Facebook UserAuthType = "facebook"
)

func (ua *UserAuth) IsValidAuthType() bool {
	switch ua.AuthType {
	case "standard", "facebook":
		return true
	default:
		return false
	}
}

func (ua *UserAuth) BeforeCreate(tx *gorm.DB) (err error) {
	if !ua.IsValidAuthType() {
		// Cancel transaction
		return errors.New("invalid AuthType provided")
	}
	return nil
}

func (ua *UserAuth) BeforeUpdate(tx *gorm.DB) (err error) {
	if !ua.IsValidAuthType() {
		// Cancel transaction
		return errors.New("invalid AuthType provided")
	}
	return nil
}

type SubscriptionTier string

const (
	Free       SubscriptionTier = "Free"    // Free
	ThirtyUses SubscriptionTier = "30-Uses" // Basic
	NinetyUses SubscriptionTier = "90-Uses" // Premium
)

type Subscription struct {
	gorm.Model
	UserID           uint             `gorm:"unique;index"`
	SubscriptionTier SubscriptionTier `gorm:"type:text;default:'Free'"`
	ExpiresAt        time.Time
	RemainingUses    int `gorm:"default:5"`
}

func (s *Subscription) IsValidSubscriptionTier() bool {
	switch s.SubscriptionTier {
	case Free, ThirtyUses, NinetyUses:
		return true
	default:
		return false
	}
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	if !s.IsValidSubscriptionTier() {
		// Set default
		s.SubscriptionTier = Free
	}
	return nil
}

func (s *Subscription) BeforeUpdate(tx *gorm.DB) (err error) {
	if !s.IsValidSubscriptionTier() {
		// Cancel transaction
		return errors.New("invalid SubscriptionTier provided")
	}
	return nil
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
	UnitSystem   GuidingContentUnitSystem `gorm:"type:text"`
	Requirements string                   // Additional instructions or guidelines
	// DietaryRestrictions string // Specific dietary restrictions
	// SupportingResearch string // Supporting research to help convey the user's expectations
}

type GuidingContentUnitSystem string

const (
	USCustomary GuidingContentUnitSystem = "US Customary"
	Metric      GuidingContentUnitSystem = "Metric"
)

func (gc *GuidingContent) IsValidUnitSystem() bool {
	switch gc.UnitSystem {
	case USCustomary, Metric:
		return true
	default:
		return false
	}
}

func (gc *GuidingContent) BeforeCreate(tx *gorm.DB) (err error) {
	if !gc.IsValidUnitSystem() {
		// Set default
		gc.UnitSystem = USCustomary
	}
	return nil
}

func (gc *GuidingContent) BeforeUpdate(tx *gorm.DB) (err error) {
	if !gc.IsValidUnitSystem() {
		// Set default
		gc.UnitSystem = USCustomary
	}
	return nil
}
