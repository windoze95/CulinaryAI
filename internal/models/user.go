package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// User is the model for a user.
type User struct {
	gorm.Model
	Username         string           `gorm:"unique;index"`
	FirstName        string           `gorm:"default:null"`
	Email            string           `gorm:"unique;default:null"`
	Auth             *UserAuth        `gorm:"foreignKey:UserID"`
	Subscription     *Subscription    `gorm:"foreignKey:UserID"`
	Settings         *UserSettings    `gorm:"foreignKey:UserID"`
	Personalization  *Personalization `gorm:"foreignKey:UserID"`
	CollectedRecipes []*Recipe        `gorm:"many2many:user_collected_recipes;"`
}

// UserAuth is the model for a user's authentication information.
type UserAuth struct {
	gorm.Model
	UserID         uint `gorm:"unique;index"`
	HashedPassword string
	AuthType       UserAuthType `gorm:"type:text"`
}

// UserAuthType is the type for the UserAuthType enum.
type UserAuthType string

// UserAuthType enum values.
const (
	Standard UserAuthType = "standard"
)

// IsValidAuthType checks if the AuthType is valid.
func (ua *UserAuth) IsValidAuthType() bool {
	switch ua.AuthType {
	case Standard:
		return true
	default:
		return false
	}
}

// BeforeCreate is a GORM hook that runs before creating a new UserAuth.
func (ua *UserAuth) BeforeCreate(tx *gorm.DB) (err error) {
	if !ua.IsValidAuthType() {
		// Cancel transaction
		return errors.New("invalid AuthType provided")
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a UserAuth.
func (ua *UserAuth) BeforeUpdate(tx *gorm.DB) (err error) {
	if !ua.IsValidAuthType() {
		// Cancel transaction
		return errors.New("invalid AuthType provided")
	}

	return nil
}

// SubscriptionTier is the type for the SubscriptionTier enum.
type SubscriptionTier string

// SubscriptionTier enum values.
const (
	Free    SubscriptionTier = "Free"    // Free
	Basic   SubscriptionTier = "Basic"   // Basic
	Premium SubscriptionTier = "Premium" // Premium
)

// Subscription is the model for a user's subscription.
type Subscription struct {
	gorm.Model
	UserID           uint             `gorm:"unique;index"`
	SubscriptionTier SubscriptionTier `gorm:"type:text;default:'Free'"`
	ExpiresAt        time.Time
	RemainingTokens  int `gorm:"default:50000"`
}

// IsValidSubscriptionTier checks if the SubscriptionTier is valid.
func (s *Subscription) IsValidSubscriptionTier() bool {
	switch s.SubscriptionTier {
	case Free, Basic, Premium:
		return true
	default:
		return false
	}
}

// BeforeCreate is a GORM hook that runs before creating a new user Subscription.
func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	if !s.IsValidSubscriptionTier() {
		// Set default
		s.SubscriptionTier = Free
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a user Subscription.
func (s *Subscription) BeforeUpdate(tx *gorm.DB) (err error) {
	if !s.IsValidSubscriptionTier() {
		// Cancel transaction
		return errors.New("invalid SubscriptionTier provided")
	}

	return nil
}

// UserSettings is the model for a user's settings.
type UserSettings struct {
	gorm.Model
	UserID          uint `gorm:"unique;index"`
	KeepScreenAwake bool `gorm:"default:true"`
}

// Personalization is the model for a user's personalization settings.
type Personalization struct {
	gorm.Model
	UserID       uint       `gorm:"unique;index"`
	UnitSystem   UnitSystem `gorm:"type:int"`
	Requirements string     // Additional instructions or guidelines
	UID          uuid.UUID
}

// UnitSystem is the type for the UnitSystem enum.
type UnitSystem int

// UnitSystem enum values.
const (
	USCustomary     UnitSystem       = iota // 0 - US Customary
	Metric                                  // 1 - Metric
	USCustomaryText = "US Customary"        // 0 - US Customary
	MetricText      = "Metric"              // 1 - Metric
)

// IsValidUnitSystem checks if the UnitSystem is valid.
func (p *Personalization) IsValidUnitSystem() bool {
	switch p.UnitSystem {
	case USCustomary, Metric:
		return true
	default:
		return false
	}
}

// GetUnitSystemText returns the text representation of the UnitSystem.
func (p *Personalization) GetUnitSystemText() string {
	switch p.UnitSystem {
	case USCustomary:
		return USCustomaryText
	case Metric:
		return MetricText
	default:
		return USCustomaryText
	}
}

// BeforeCreate is a GORM hook that runs before creating a new user Personalization.
func (p *Personalization) BeforeCreate(tx *gorm.DB) (err error) {
	if !p.IsValidUnitSystem() {
		// Set default
		p.UnitSystem = USCustomary
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a user Personalization.
func (p *Personalization) BeforeUpdate(tx *gorm.DB) (err error) {
	if !p.IsValidUnitSystem() {
		// Set default
		p.UnitSystem = USCustomary
	}

	return nil
}
