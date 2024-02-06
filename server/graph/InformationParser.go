package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"graphql_json_go/graph/model"
	"reflect"
	"strconv"

	"github.com/rs/zerolog/log"
)

var models = []*model.Model{
	{
		Name: "Person",
		Fields: []*model.ModelField{
			{
				Name:       "name",
				Type:       model.ModelFieldEnumString,
				IsNullable: false,
			},
			{
				Name:       "age",
				Type:       model.ModelFieldEnumInt,
				IsNullable: false,
			},
			{
				Name:       "isSubscriber",
				Type:       model.ModelFieldEnumBoolean,
				IsNullable: false,
			},
			{
				Name:       "contacts",
				Type:       model.ModelFieldEnumArray,
				IsNullable: true,
				SubFields: []*model.ModelField{
					{
						Name:       "contact",
						Type:       model.ModelFieldEnumObject,
						IsNullable: false,
						SubFields: []*model.ModelField{
							{
								Name:       "email",
								Type:       model.ModelFieldEnumString,
								IsNullable: false,
							},
							{
								Name:       "phone",
								Type:       model.ModelFieldEnumString,
								IsNullable: true,
							},
						},
					},
				},
			},
		},
	},
	{
		Name: "Project",
		Fields: []*model.ModelField{
			{
				Name:       "title",
				Type:       model.ModelFieldEnumString,
				IsNullable: false,
			},
			{
				Name:       "description",
				Type:       model.ModelFieldEnumString,
				IsNullable: true,
			},
			{
				Name:       "members",
				Type:       model.ModelFieldEnumArray,
				IsNullable: false,
				SubFields: []*model.ModelField{
					{
						Name:       "member",
						Type:       model.ModelFieldEnumObject,
						IsNullable: false,
						SubFields: []*model.ModelField{
							{
								Name:       "name",
								Type:       model.ModelFieldEnumString,
								IsNullable: false,
							},
							{
								Name:       "role",
								Type:       model.ModelFieldEnumString,
								IsNullable: false,
							},
						},
					},
				},
			},
			{
				Name:       "tasks",
				Type:       model.ModelFieldEnumArray,
				IsNullable: true,
				SubFields: []*model.ModelField{
					{
						Name:       "task",
						Type:       model.ModelFieldEnumObject,
						IsNullable: false,
						SubFields: []*model.ModelField{
							{
								Name:       "description",
								Type:       model.ModelFieldEnumString,
								IsNullable: false,
							},
							{
								Name:       "status",
								Type:       model.ModelFieldEnumString,
								IsNullable: false,
							},
						},
					},
				},
			},
		},
	},
}

func CreateModel(ctx context.Context, name string, fields []*model.ModelFieldInput) (*model.Model, error) {
	exists := false
	for _, m := range models {
		if m.Name == name {
			exists = true
			break
		}
	}
	if exists {
		err := fmt.Errorf("Model with name %s already exists", name)
		return nil, err
	}
	newFields := []*model.ModelField{}
	for _, f := range fields {
		newFields = append(newFields, &model.ModelField{
			Type:       f.Type,
			IsNullable: f.IsNullable,
		})
	}
	newModel := model.Model{
		Name:   name,
		Fields: newFields,
	}
	models = append(models, &newModel)
	return &newModel, nil
}

func GetModels(ctx context.Context) ([]*model.Model, error) {
	return models, nil
}

func SendInformaton(ctx context.Context, info map[string]interface{}, modelName string) (bool, error) {
	// Get the model definition
	var modelDef *model.Model
	for _, m := range models {
		if m.Name == modelName {
			modelDef = m
			break
		}
	}
	if modelDef == nil {
		return false, fmt.Errorf("Model %s not found", modelName)
	}
	if !validateJSON(info, modelDef.Fields) {
		return false, fmt.Errorf("Invalid JSON for model %s", modelName)
	}
	return true, nil
}

func validateField(value interface{}, field *model.ModelField) bool {
	if value == nil {
		return field.IsNullable
	}

	valueType := reflect.TypeOf(value).Kind()

	switch field.Type {
	case model.ModelFieldEnumString:
		return valueType == reflect.String
	case model.ModelFieldEnumInt:
		switch v := value.(type) {
		case float64:
			return v == float64(int64(v)) // Check if the float is actually an int.
		case json.Number:
			// Try to parse the json.Number as an int64.
			if i, err := v.Int64(); err == nil {
				log.Info().Msgf("json.Number as Int64 value: %v", i)
				return true
			} else {
				log.Error().Msgf("Error converting json.Number to Int64: %v", err)
				return false
			}
		case string:
			// Try to parse the string as an int.
			_, err := strconv.Atoi(v)
			if err != nil {
				log.Error().Msgf("Error parsing string to int: %v", err)
				return false
			}
			return true
		default:
			// Handle other types as invalid.
			return false
		}
	case model.ModelFieldEnumFloat:
		return valueType == reflect.Float64
	case model.ModelFieldEnumBoolean:
		return valueType == reflect.Bool
	case model.ModelFieldEnumArray:
		if valueType != reflect.Slice {
			return false
		}
		elements, _ := value.([]interface{})
		if len(field.SubFields) == 0 {
			return true // If no subFields are defined, any array is valid.
		}
		for _, elem := range elements {
			// Validating each element of the array.
			if !validateField(elem, field.SubFields[0]) {
				return false
			}
		}
		return true
	case model.ModelFieldEnumObject:
		if valueType != reflect.Map {
			return false
		}
		obj, _ := value.(map[string]interface{})
		for _, subField := range field.SubFields {
			fieldValue, exists := obj[subField.Name]
			if !exists {
				if !subField.IsNullable {
					return false // Field is not nullable and is missing from the object.
				}
				continue // If the field is nullable and missing, it's valid.
			}
			if !validateField(fieldValue, subField) {
				return false
			}
		}
		// Verify that there are no extra fields in the object.
		if len(obj) > len(field.SubFields) {
			return false // If there are more fields in the object than in the model, it's invalid.
		}
		return true
	default:
		return false
	}
}

// validateJSON validates that the given JSON object is valid for the given model.
func validateJSON(info map[string]interface{}, fields []*model.ModelField) bool {
	log.Info().Msgf("Validating JSON: %v", info)
	// First, validate that all fields are present and valid.
	fieldNames := make(map[string]bool)
	for _, field := range fields {
		fieldNames[field.Name] = true
		value, exists := info[field.Name]
		if !exists {
			if !field.IsNullable {
				fmt.Printf("Missing non-nullable field: %s\n", field.Name)
				return false
			}
			continue // The field is nullable and missing, so it's valid.
		}

		if !validateField(value, field) {
			fmt.Printf("Validation failed for field: %s\n", field.Name)
			return false
		}
	}

	// Then, validate that there are no extra fields in the JSON.
	for key := range info {
		if !fieldNames[key] {
			fmt.Printf("Extra field in JSON not defined in the model: %s\n", key)
			return false
		}
	}

	return true
}
