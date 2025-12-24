package models

import "fmt"

func GetAllModels() []interface{} {
	var models []interface{}

	modelsMap := map[string]interface{}{
		"User": &User{},
	}

	for _, m := range modelsMap {
		models = append(models, m)
	}

	fmt.Printf("Found %d models for AutoMigrate\n", len(models))
	return models
}
