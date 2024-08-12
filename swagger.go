package hive

import (
	"io"
	"os"
	"reflect"
)

func generateSwagger(modules *[]Module) {

	stringToSave := generateString(modules)

	// for _, module := range *modules {
	// 	for _, controller := range module.controllers {
	// 		for _, route := range controller.routes {
	// 			if route.metadata["body"] != nil {
	// 				println("Body fields: ")
	// 				allFields := GetAllFieldsOfStruct(route.metadata["body"])
	// 				for _, field := range allFields {
	// 					println(field.name + " - " + field.type_field)
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	//save to file
	f, err := os.Create("swagger.json")
	if err != nil {
		panic(err)
	}

	data := []byte(stringToSave)

	_, err = io.WriteString(f, string(data))
	if err != nil {
		panic(err)
	}

	defer f.Close()
}

func GetAllFieldsOfStruct(structure interface{}) []struct {
	name       string
	type_field string
	json_name  string
} {

	var fields []struct {
		name       string
		type_field string
		json_name  string
	}

	t := reflect.TypeOf(structure)

	for i := 0; i < t.NumField(); i++ {
		name_of_field := t.Field(i).Name
		type_field := t.Field(i).Type

		fields = append(fields, struct {
			name       string
			type_field string
			json_name  string
		}{
			name:       name_of_field,
			type_field: type_field.String(),
			json_name:  t.Field(i).Tag.Get("json"),
		})
	}

	return fields
}

func generateString(modules *[]Module) string {

	// 	text := `{
	// 	"openapi": "3.0.0",
	// 		"paths": {
	// 	`

	// 	groupRoutesWithSamePath := make(map[string][]interface{})

	// 	for _, module := range *modules {
	// 		for _, controller := range module.controllers {
	// 			for _, route := range controller.routes {
	// 				if groupRoutesWithSamePath[route.metadata["path"].(string)] == nil {
	// 					groupRoutesWithSamePath[route.metadata["path"].(string)] = []interface{}{route.metadata}
	// 				} else {
	// 					groupRoutesWithSamePath[route.metadata["path"].(string)] = append(groupRoutesWithSamePath[route.metadata["path"].(string)], route.metadata)
	// 				}
	// 			}
	// 		}
	// 	}

	// 	// for _, module := range *modules {
	// 	// 	for _, controller := range module.controllers {
	// 	// 		for _, route := range controller.routes {
	// 	// 			text += GenerateRouteString(route.metadata)
	// 	// 		}
	// 	// 	}
	// 	// }

	// 	indexPath := 0

	// 	for path, routes := range groupRoutesWithSamePath {
	// 		text += `		"` + path + `": {`
	// 		for index, route := range routes {
	// 			text += GenerateRouteString(route.(map[string]interface{}))

	// 			if index != len(routes)-1 {
	// 				text += `,`
	// 			}
	// 		}

	// 		if indexPath != len(groupRoutesWithSamePath)-1 {
	// 			text += `
	// 		},
	// `
	// 		} else {
	// 			text += `
	// 		}
	// `
	// 		}

	// 		indexPath++

	// 	}

	// 	text += `	},`

	// 	text += `
	// 	"components": {
	// 		"schemas": {
	// 	`

	// 	filteredRoutesWithBody := []interface{}{}

	// 	for _, module := range *modules {
	// 		for _, controller := range module.controllers {
	// 			for _, route := range controller.routes {
	// 				if route.metadata["body"] != nil {
	// 					// text += GenerateBodyString(route.metadata["body"].(interface{}))
	// 					filteredRoutesWithBody = append(filteredRoutesWithBody, route.metadata)
	// 				}
	// 			}
	// 		}
	// 	}

	// 	for _, route := range filteredRoutesWithBody {
	// 		text += GenerateBodyString(route.(map[string]interface{})["body"].(interface{}))
	// 	}

	// 	text += `
	// 		}
	// 	}`

	// 	text += `
	// }`

	text := CreateMainOpenAiFile(
		CreatePaths(modules),
		CreateComponents(
			CreateSchemas(modules),
		),
	)

	return text
}

func GenerateBodyString(body interface{}) string {

	text := `		"` + reflect.TypeOf(body).Name() + `": {
				"type": "object",
				"properties": {
					`

	allFields := GetAllFieldsOfStruct(body)

	for index, field := range allFields {

		newType := ""

		switch field.type_field {
		case "string":
			newType = "string"
		case "int":
			newType = "integer"
		case "bool":
			newType = "boolean"
		default:
			newType = "object"
		}

		nameOfField := field.name

		if field.json_name != "" {
			nameOfField = field.json_name
		}

		text += `"` + nameOfField + `": {
						"type": "` + newType + `"
					}`

		if index != len(allFields)-1 {
			text += `,
					`
		}

	}

	text += `
				}
			}`

	return text
}
