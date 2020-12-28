package gindocs

import (
	"fmt"
	"io/ioutil"
	"os"
)

var json = ""
var body = ""
var definition = ""

func sethead() {
	json += `{
	"swagger": "2.0",
	"info": {
	  "version": "1.0.0",
	  "title": "Swagger Petstore",
	  "description": "A sample API that uses a petstore as an example to demonstrate features in the swagger-2.0 specification",
	  "termsOfService": "http://swagger.io/terms/",
	  "contact": {
		"name": "Swagger API Team"
	  },
	  "license": {
		"name": "MIT"
	  }
	},
	"host": "petstore.swagger.io",
	"basePath": "/api",
	"schemes": [
	  "http"
	],
	"consumes": [
	  "application/json"
	],
	"produces": [
	  "application/json"
	],`
}

func setbodystart() {
	body += `"paths": {`
}

type BodySwag struct {
	Code string
	Type string
	Ref  string
}

func setbody(url, method string, response []BodySwag) {
	body += `
		"` + url + `": {
		  "` + method + `": {
			"produces": [
			  "application/json"
			],
			"responses": {`
	for i, a := range response {
		if i > 0 {
			body += `
			,"` + a.Code + `": {
			  "schema": {
				"type": "` + a.Type + `",
				"items": {
				  "$ref": "` + a.Ref + `"
				}
			  }
			}`
		} else {
			body += `
			"` + a.Code + `": {
				"schema": {
				  "type": "` + a.Type + `",
				  "items": {
					"$ref": "` + a.Ref + `"
				  }
				}
			  }`
		}
	}

	body += `
			}
		  }
		}
	  `
}
func setbodyend() {
	body += `},`
}
func setdefinitionstart() {
	definition += `"definitions": {`
}

func setdefinition() {
	definition += `	
		"Pet": {
		  "type": "object",
		  "required": [
			"id",
			"name"
		  ],
		  "properties": {
			"id": {
			  "type": "integer",
			  "format": "int64"
			},
			"name": {
			  "type": "string"
			},
			"tag": {
			  "type": "string"
			}
		  }
		}
	  }`
}
func setdefinitionend() {
	definition += `}`
}

func Exec() {
	sethead()
	setbodystart()
	// setbody()
	// setbodyend()
	// setdefinitionstart()
	// setdefinition()
	// setdefinitionend()
	// fmt.Println(json)
	dir := GetDir()
	os.Mkdir(dir+"/readme", 0755)
	err := ioutil.WriteFile("readme/readme.json", []byte(json+body+definition), 0755)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}
}
