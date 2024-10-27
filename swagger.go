package hive

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
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

type PathMethodConfig struct {
	Tags        []string                            `json:"tags"`
	Responses   map[string]PathMethodConfigResponse `json:"responses"`
	Security    []any                               `json:"security"`
	Parameters  []interface{}                       `json:"parameters"`
	RequestBody map[string]interface{}              `json:"requestBody"`
}

type PathMethodConfigResponse struct {
	Description string `json:"description"`
}

type SwaggerV2Config struct {
	Openapi    string                                 `json:"openapi"`
	Info       SwaggerV2ConfigInfo                    `json:"info"`
	Paths      map[string]map[string]PathMethodConfig `json:"paths"`
	Components SwaggerV2ConfigComponents              `json:"components"`
}

type SwaggerV2ConfigInfo struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type SwaggerV2ConfigComponents struct {
	SecuritySchemes SwaggerV2ConfigComponentsSecuritySchemes    `json:"securitySchemes"`
	Schemas         map[string]SwaggerV2ConfigComponentsSchemas `json:"schemas"`
}

type SwaggerV2ConfigComponentsSchemas struct {
	Type       string `json:"type"`
	Properties map[string]struct {
		Type  string `json:"type"`
		Items struct {
			Type string `json:"type"`
		} `json:"items"`
	} `json:"properties"`
}

type SwaggerV2ConfigComponentsSecuritySchemes struct {
	Bearer map[string]string `json:"bearer"`
}

func GenerateSwaggerV2(n *GoNest) {
	stringToSave := generateStringV2(n)
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

func generateStringV2(n *GoNest) string {

	config := SwaggerV2Config{
		Openapi: "3.0.0",
		Info: SwaggerV2ConfigInfo{
			Title:       "API title",
			Version:     "1.0.0",
			Description: "",
		},
		Paths: make(map[string]map[string]PathMethodConfig),
		Components: SwaggerV2ConfigComponents{
			SecuritySchemes: SwaggerV2ConfigComponentsSecuritySchemes{
				Bearer: map[string]string{
					"scheme":       "Bearer",
					"bearerFormat": "Bearer",
					"description":  `[just text field] Please enter token in following format: Bearer <JWT>`,
					"name":         "Authorization",
					"type":         "http",
					"in":           "Header",
				},
			},
		},
	}

	if n.config.SwaggerConfig.Title != "" {
		config.Info.Title = n.config.SwaggerConfig.Title
	}

	if n.config.SwaggerConfig.Version != "" {
		config.Info.Version = n.config.SwaggerConfig.Version
	}

	if n.config.SwaggerConfig.Description != "" {
		config.Info.Description = n.config.SwaggerConfig.Description
	}

	modules := &n.modules

	groupMethodsByRoutePath := []struct {
		path    string
		methods []interface{}
	}{}

	for _, module := range *modules {
		for _, controller := range module.controllers {

			for _, route := range controller.routes {
				route.metadata["controller_tag"] = controller.metadata["tag"]

				var groupAlreadyExists *struct {
					path    string
					methods []interface{}
				}

				for _, group := range groupMethodsByRoutePath {
					if group.path == route.metadata["path"].(string) {
						groupAlreadyExists = &group
					}
				}

				if groupAlreadyExists == nil {
					groupMethodsByRoutePath = append(groupMethodsByRoutePath, struct {
						path    string
						methods []interface{}
					}{
						path:    route.metadata["path"].(string),
						methods: []interface{}{route.metadata},
					})
				} else {
					for i, group := range groupMethodsByRoutePath {
						if group.path == route.metadata["path"].(string) {
							groupMethodsByRoutePath[i].methods = append(groupMethodsByRoutePath[i].methods, route.metadata)
						}
					}
				}
			}
		}
	}

	for _, pathConfig := range groupMethodsByRoutePath {

		pathMethods := make(map[string]PathMethodConfig)

		for _, method := range pathConfig.methods {

			method := method.(map[string]interface{})
			lowerCaseMethod := strings.ToLower(method["method"].(string))

			//if has method["controller_tag"].(string) then add if not add ""

			ControllerTag := ""

			if method["controller_tag"] != nil {
				ControllerTag = method["controller_tag"].(string)
			}

			Responses := make(map[string]PathMethodConfigResponse)

			Responses["200"] = PathMethodConfigResponse{
				Description: "OK",
			}

			Security := []any{}

			if method["bearer"] != nil {
				Responses["401"] = PathMethodConfigResponse{
					Description: "Unauthorized",
				}
				Security = append(Security, map[string][]any{
					"bearer": {},
				})
			}

			parametersBody := method["parameters"].([]string)

			Parameters := []interface{}{}

			for _, param := range parametersBody {
				Parameters = append(Parameters, map[string]interface{}{
					"in":       "path",
					"name":     param,
					"required": true,
					"schema": struct {
						Type string `json:"type"`
					}{
						Type: "string",
					},
				})
			}

			requestBody := map[string]interface{}{}

			if method["body"] != nil {
				requestBody = map[string]interface{}{
					"required": true,
					"content": map[string]struct {
						Schema struct {
							Ref string `json:"$ref"`
						} `json:"schema"`
					}{
						"application/json": {
							Schema: struct {
								Ref string `json:"$ref"`
							}{
								Ref: "#/components/schemas/" + reflect.TypeOf(method["body"]).Name(),
							},
						},
					},
				}
			}

			pathMethods[lowerCaseMethod] = PathMethodConfig{
				Tags: []string{
					ControllerTag,
				},
				Responses:   Responses,
				Security:    Security,
				Parameters:  Parameters,
				RequestBody: requestBody,
			}
		}

		//replace all cases of :param with {param} remember to match the case and also add the } at the end

		re := regexp.MustCompile(`:([a-zA-Z]+)`)
		path := re.ReplaceAllString(pathConfig.path, "{$1}")

		config.Paths[path] = pathMethods
	}

	filteredRoutesWithBody := []interface{}{}

	for _, module := range *modules {
		for _, controller := range module.controllers {
			for _, route := range controller.routes {
				if route.metadata["body"] != nil {
					filteredRoutesWithBody = append(filteredRoutesWithBody, route.metadata)
				}
			}
		}
	}

	schemas := make(map[string]SwaggerV2ConfigComponentsSchemas)

	for _, route := range filteredRoutesWithBody {
		route := route.(map[string]interface{})
		body := route["body"]

		typeOfBody := reflect.TypeOf(body).Name()

		schemas[typeOfBody] = SwaggerV2ConfigComponentsSchemas{
			Type: "object",
			Properties: make(map[string]struct {
				Type  string `json:"type"`
				Items struct {
					Type string `json:"type"`
				} `json:"items"`
			}),
		}

		allFields := GetAllFieldsOfStruct(body)

		for _, field := range allFields {

			newType := "object"

			if field.type_field == "string" {
				newType = "string"
			}

			if field.type_field == "int" {
				newType = "integer"
			}

			if field.type_field == "bool" {
				newType = "boolean"
			}

			if field.type_field == "float64" {
				newType = "number"
			}

			ItemsType := struct {
				Type string `json:"type"`
			}{}

			if strings.HasPrefix(field.type_field, "[]") {
				newType = "array"
				ItemsType = struct {
					Type string `json:"type"`
				}{
					Type: "string",
				}
			}

			nameOfField := field.name

			if field.json_name != "" {
				nameOfField = field.json_name

			}

			schemas[typeOfBody].Properties[nameOfField] = struct {
				Type  string `json:"type"`
				Items struct {
					Type string `json:"type"`
				} `json:"items"`
			}{
				Type:  newType,
				Items: ItemsType,
			}
		}
	}

	config.Components.Schemas = schemas

	jsonIndent, _ := json.MarshalIndent(config, "", "  ")

	jsonString := string(jsonIndent)

	return jsonString

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
