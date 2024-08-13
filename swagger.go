package hive

import (
	"io"
	"os"
	"reflect"
)

type SwaggerConfig struct {

	//default is false
	Enabled bool

	//default is "API title"
	Title string

	//default is "API description"
	Description string

	//default is 1.0.0
	Version string

	//default path is /api
	Path string
}

func generateSwagger(n *GoNest) {

	stringToSave := generateString(n)

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

func generateString(n *GoNest) string {

	modules := &n.modules

	text := CreateMainOpenAiFile(
		n.config.SwaggerConfig,
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
