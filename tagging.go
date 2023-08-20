package main

import (
	"strings"

	"github.com/jinzhu/gorm"
)

// func create() {
// 	// Define some tags
// 	parsedTagNames := []string{"Carnavore", "seedOil-free"}

// 	// Check if the tags exist, and create them if they don't
// 	var tags []Tag
// 	for _, tagName := range parsedTagNames {
// 		var tag Tag
// 		if err := db.Where("LOWER(name) = ?", strings.ToLower(tagName)).First(&tag).Error; err == gorm.ErrRecordNotFound {
// 			tag = Tag{Name: tagName}
// 			db.Create(&tag)
// 		}
// 		tags = append(tags, tag)
// 	}

// 	// Define a recipe with the tags
// 	recipe := Recipe{
// 		Title:   "Delicious Meat Salad without Lettuce",
// 		Content: "Recipe content here...",
// 		Tags:    tags,
// 	}

// 	// Save the recipe and its associated tags
// 	db.Create(&recipe)
// }

func appendTags(tagNames []string) error {
	for _, tagName := range tagNames {
		// Check if the tag exists, and create it if it does not
		canonicalName := strings.ToLower(tagName)
		var tag Tag
		if err := db.Where("LOWER(name) = ?", canonicalName).First(&tag).Error; err == gorm.ErrRecordNotFound {
			tag = Tag{Name: canonicalName}
			if err := db.Create(&tag).Error; err != nil {
				return err
			}
		}
		// Associate the tag with the recipe
		if err := db.Model(&Recipe{}).Association("Tags").Append(&tag).Error; err != nil {
			return err
		}
	}

	return nil
}
