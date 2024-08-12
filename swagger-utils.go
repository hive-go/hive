package hive

import (
	"fmt"
	"reflect"
	"strings"
)

func CreateMainOpenAiFile(
	configs ...string,
) string {

	configsMixed := ""

	if len(configs) > 0 {
		configsMixed += `,
`
	}

	for index, config := range configs {
		configsMixed += fmt.Sprintf(`%v`, config)
		if index != len(configs)-1 {
			configsMixed += `,
`
		}
	}

	text := fmt.Sprintf(`{
	"openapi": "3.0.0",
	"info": {
        "title": "Documentação API Template NestJS",
        "description": "",
        "version": "1.0.0",
        "contact": {}
  }%v
}`, configsMixed)

	return text
}

func CreateComponents(componentsProps ...string) string {

	fullStringWithSchemas := ""

	if len(componentsProps) > 0 {
		fullStringWithSchemas += `,
`
	}

	for _, component := range componentsProps {
		fullStringWithSchemas += fmt.Sprintf(`%v`, component)
	}

	var finalText = fmt.Sprintf(`	"components": {
		"securitySchemes": {
			"bearer": {
				"scheme": "Bearer",
				"bearerFormat": "Bearer",
				"description": "[just text field] Please enter token in following format: Bearer <JWT>",
				"name": "Authorization",
				"type": "http",
				"in": "Header"
			}
		}%v
	}`, fullStringWithSchemas)

	return finalText

}

func CreateSchemas(modules *[]Module) string {
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

	fullStringWithSchemas := ""

	for index, route := range filteredRoutesWithBody {
		fullStringWithSchemas += fmt.Sprintf(`%v`, GenerateBodyString(route.(map[string]interface{})["body"]))
		if index != len(filteredRoutesWithBody)-1 {
			fullStringWithSchemas += `,
`
		}
	}

	var finalText = fmt.Sprintf(`		"schemas": {
	 %v
		}`,
		fullStringWithSchemas)

	return finalText
}

func CreatePaths(modules *[]Module) string {

	groupMethodsByRoutePath := make(map[string][]interface{})

	for _, module := range *modules {
		for _, controller := range module.controllers {
			for _, route := range controller.routes {
				route.metadata["controller_tag"] = controller.metadata["tag"]

				if groupMethodsByRoutePath[route.metadata["path"].(string)] == nil {
					groupMethodsByRoutePath[route.metadata["path"].(string)] = []interface{}{route.metadata}
				} else {
					groupMethodsByRoutePath[route.metadata["path"].(string)] = append(groupMethodsByRoutePath[route.metadata["path"].(string)], route.metadata)
				}
			}
		}
	}

	fullStringWithRoutes := ""

	count := 0
	for path, routeMethods := range groupMethodsByRoutePath {

		var methodsString = ""

		for index, method := range routeMethods {

			method := method.(map[string]interface{})
			lowerCaseMethod := strings.ToLower(method["method"].(string))

			tagString := ""
			tag := method["controller_tag"]

			if tag != nil {
				tagString = tag.(string)
			}

			parameters := GetParametersString(method["parameters"].([]string))

			bodyString := ""
			body := method["body"]

			if body != nil {

				typeOfBody := reflect.TypeOf(body).Name()

				bodyString = fmt.Sprintf(`,
				"requestBody": {
					"required": true,
					"content": {
							"application/json": {
									"schema": {
											"$ref": "#/components/schemas/%v"
									}
							}
					}
        }`, typeOfBody)
			}

			bearerString := ""
			hasBearer := method["bearer"]

			if hasBearer != nil {
				bearerString = `{
					"bearer": []
				}`
			}

			methodsString += fmt.Sprintf(`			"%v": {
				"tags": [ "%v" ],
				"responses": {
					"201": {
						"description": ""
					}
				},
				"security": [
					%v
				],
				"parameters": [
					%v
				]%v
			}`, lowerCaseMethod, tagString, bearerString, parameters, bodyString)

			if index != len(routeMethods)-1 {
				methodsString += `,
`
			}

		}

		fullStringWithRoutes += fmt.Sprintf(`		"%v": {
%v
		}`, path, methodsString)

		if count != len(groupMethodsByRoutePath)-1 {
			fullStringWithRoutes += `,
`
		}

		count++
	}

	var finalText = fmt.Sprintf(`	"paths": {
%v
	}`, fullStringWithRoutes)

	return finalText
}

func GetParametersString(parameters []string) string {
	parametersString := ""
	for index, param := range parameters {
		parametersString += fmt.Sprintf(`{
						"in": "path",
						"name": "%v",
						"required": true,
						"schema": {
							"type": "string"
						}
					}`, param)

		if index != len(parameters)-1 {
			parametersString += `,`
		}
	}
	return parametersString
}

// func GenerateRouteString(routeMetadata map[string]interface{}) string {

// 	lowerCaseMethod := strings.ToLower(routeMetadata["method"].(string))

// 	text := `
// 				"` + lowerCaseMethod + `": {
// 			`

// 	if routeMetadata["tag"] != nil {
// 		text += AddTag(routeMetadata["tag"].(string))
// 	}

// 	text += `		` + GetResponse() + `
// 					` + GetSecurity() + `
// 					"parameters": [
// 	`
// 	if routeMetadata["parameters"] != nil {
// 		for index, param := range routeMetadata["parameters"].([]string) {
// 			text += GetPathParam(param)
// 			if index != len(routeMetadata["parameters"].([]string))-1 {
// 				text += ","
// 			}
// 		}
// 	}

// 	text += `
// 					]`

// 	if routeMetadata["body"] != nil {
// 		text += `
// 			"requestBody": {
// 				"required": true,
// 				"content": {
// 						"application/json": {
// 								"schema": {
// 										"$ref": "#/components/schemas/CreateGuideDto"
// 								}
// 						}
// 				}
// 		},
// 		`
// 	}

// 	text += `
// 				}`

// 	return text

// }

// routeString := fmt.Sprintf(`"%v": {
// 	%v
// }`, path)

// if indexPath != len(groupMethodsByRoutePath)-1 {
// 	routeString += `,`
// }
