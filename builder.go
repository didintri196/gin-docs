package ginswag

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-gonic/gin"
)

var doc *string

func init() {
	doc = flag.String("doc", "default", "value as swagger")
	flag.Parse()
	// fmt.Println(*gen)
}
func Use(router *gin.Engine) {
	// fmt.Println(*gen)
	dir := GetDir()
	router.Static("/readme", dir+"/readme/")
	url := ginSwagger.URL("/readme/readme.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	if *doc == "swagger" {
		allrouter := router.Routes()
		for _, route := range allrouter {
			generate(route.Path, route.Method, route.Handler)
		}
		Exec()
		os.Exit(3)
	}
}

//MENCARI
func getvalue(data, start, end string) (value string) {
	pecah1 := strings.Split(data, start)
	pecah2 := strings.Split(pecah1[1], end)
	return pecah2[0]
}

func GetDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func generate(url, method, handler string) {
	if strings.Contains(handler, "controllers") {
		findcontroll := getvalue(handler, "controllers.(*", ").")
		findfunc := getvalue(handler, ").", "-fm")
		fmt.Println("===================================")
		fmt.Println("URL : ", url)
		fmt.Println("METHOD : ", method)
		fmt.Println("HANDLER : ", handler)
		fmt.Println("CONTROLLER : ", findcontroll)
		fmt.Println("FUNCTION : ", findfunc)
		dir := GetDir()
		data, _ := ioutil.ReadFile(dir + "/controllers/" + findcontroll + ".go")
		// fmt.Println(err)
		parsingtahap1(findfunc, string(data), method, findcontroll)
	}
}

//LIST FUNGSI
func parsingtahap1(prefix, data, method, findcontroll string) {
	listfunc := strings.Split(string(data), "func")
	for i, func_row := range listfunc {
		if i > 0 {
			// fmt.Println(i, "->", func_row)
			status := parsingtahap2(prefix, func_row, method, findcontroll)
			if status == true {
				break
			}
		}
	}
}

//MENCARI FUNGSI
func parsingtahap2(prefix, data, method, findcontroll string) bool {
	if strings.Contains(data, prefix) {
		// fmt.Println(data) //FUNGSI DI TEMUKAN
		query := parsingquery(data)
		param := parsingparam(data)
		fmt.Println("QUERY :", query)
		fmt.Println("PARAM :", param)
		if method != "GET" {
			body := parsingbody(data, findcontroll)
			fmt.Println("BODY :", body)
		}
		parsingresponse(data, findcontroll)
		// fmt.Println("RESPONSE :", restponse)
		return true
	} else {
		return false
	}
}

//LIST QUERY
func parsingquery(data string) (query []string) {
	listquery := strings.Split(string(data), "Query(\"")
	for i, query_row := range listquery {
		if i > 0 {
			getvalue := strings.Split(query_row, "\")")
			// fmt.Println(i, "[QUERY]->", getvalue[0])
			query = append(query, getvalue[0])
		}
	}
	return
}

//LIST PARAM
func parsingparam(data string) (param []string) {
	listparam := strings.Split(string(data), "Param(\"")
	for i, param_row := range listparam {
		if i > 0 {
			getvalue := strings.Split(param_row, "\")")
			// fmt.Println(i, "[PARAM]->", getvalue[0])
			param = append(param, getvalue[0])
		}
	}
	return
}

//MENCARI BODY IF PUT,POST
func parsingbody(data, namecontrol string) (body string) {
	listbody := strings.Split(data, "var data request.")
	getvalue := strings.Split(listbody[1], "\n")
	namestructbody := getvalue[0]
	dir := GetDir()
	body = parsingstruct(dir+"/request/"+namecontrol+".go", namestructbody, namecontrol, "request")
	// fmt.Println(body)
	return
}

//MENCARI RESPONSE
func parsingresponse(data, namecontrol string) (resp string) {
	listbody := strings.Split(data, ".JSON(")
	for i, a := range listbody {
		if i > 0 {
			getvalue := strings.Split(a, "\n")
			oneline := getvalue[0]
			getvaluecode := strings.Split(oneline, ",")
			fmt.Println("CODE", getvaluecode[0])

			listvaluetext := strings.Split(oneline, "responses.")
			getvalueresponse := strings.Split(listvaluetext[1], "{")
			dir := GetDir()
			resp = parsingstruct(dir+"/responses/"+namecontrol+".go", getvalueresponse[0], namecontrol, "responses")
			fmt.Println("RESP", resp)
		}
	}

	return
}

//FIND STRUCT
func parsingstruct(path, namestructbody, namecontrol, pathserch string) (bodyjson string) {
	// fmt.Println("PATH CARI", pathserch)
	data, _ := ioutil.ReadFile(path)
	listvartext := strings.Split(string(data), "type "+namestructbody+" struct {")
	// fmt.Println("TOTAL CARI DATA STRUCT", len(listvartext))
	// fmt.Println("DATA STRUCT", listvartext)
	if len(listvartext) > 1 {
		getvalue := strings.Split(listvartext[1], "}")
		// fmt.Println(getvalue[0])
		listall := strings.Split(getvalue[0], "\n")
		totall := len(listall)
		// fmt.Println(totall)
		if (totall) > 2 {
			totall = totall - 2
		}
		i := 1
		for i_a, a := range listall {
			// fmt.Println("IA", i_a)
			space := regexp.MustCompile(`\s+`)
			a = space.ReplaceAllString(a, " ")
			listtext := strings.Split(a, "json:\"")
			if len(listtext) == 2 {
				getvaluejson := strings.Split(listtext[1], "\"")
				getvaluetype := strings.Split(listtext[0], " ")
				if i == 1 {
					i++
					if totall == 1 {
						bodyjson += "{\"" + getvaluejson[0] + "\":" + isstruct(getvaluetype[2], path, namestructbody, namecontrol) + "}"
					} else {
						bodyjson += "{\"" + getvaluejson[0] + "\":" + isstruct(getvaluetype[2], path, namestructbody, namecontrol)
					}
				} else if i_a == totall {
					if i_a == 1 {
						bodyjson += "{\"" + getvaluejson[0] + "\":" + isstruct(getvaluetype[2], path, namestructbody, namecontrol) + "}"
					} else {
						bodyjson += ", \"" + getvaluejson[0] + "\":" + isstruct(getvaluetype[2], path, namestructbody, namecontrol) + "}"
					}
				} else {
					bodyjson += ",\"" + getvaluejson[0] + "\":" + isstruct(getvaluetype[2], path, namestructbody, namecontrol)
				}
			}
		}
	} else {
		// path = strings.Replace(path, namecontrol, namestructbody, -1)
		// bodyjson = parsingstruct(path, namestructbody, namecontrol)
		files, err := ioutil.ReadDir("./" + pathserch)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			// fmt.Println(file.Name())
			dir := GetDir()
			path = dir + "/" + pathserch + "/" + file.Name()
			cek := cekstruct(path, namestructbody)
			if cek == true {
				fmt.Println("temu gess struct e")
				bodyjson = parsingstruct(path, namestructbody, namecontrol, pathserch)
				break
			}
		}
		// os.Exit(3)
	}
	return
}
func cekstruct(path, namestructbody string) (status bool) {
	data, _ := ioutil.ReadFile(path)
	if strings.Contains(string(data), "type "+namestructbody) {
		status = true
	} else {
		status = false
	}
	return
}
func isstruct(valuetipedata, path, namestructbody, namecontrol string) (value string) {
	// fmt.Println(valuetipedata)
	if strings.Contains(valuetipedata, "string") {
		if strings.Contains(valuetipedata, "[]") {
			value = "[\"\"]"
		} else {
			value = "\"\""
		}
	} else if strings.Contains(valuetipedata, "int") {
		if strings.Contains(valuetipedata, "[]") {
			value = "[0]"
		} else {
			value = "0"
		}
	} else {
		if strings.Contains(valuetipedata, "[]") {
			valuetipedata = strings.Replace(valuetipedata, "[]", "", -1)
			if strings.Contains(valuetipedata, ".") {
				getvaluetipepath := strings.Split(valuetipedata, ".")
				// fmt.Println(getvaluetipepath[1])
				dir := GetDir()
				path = dir + "/" + getvaluetipepath[0] + "/" + namecontrol + ".go"
				value = "[" + parsingstruct(path, getvaluetipepath[1], namecontrol, getvaluetipepath[0]) + "]"
			} else {
				pathsearch := ""
				dir := GetDir()
				// fmt.Println(dir)
				cekpathsearch := strings.Split(path+"/", dir)
				getvaluetipepathsearch := strings.Split(cekpathsearch[1], "/")
				if strings.Contains(path, getvaluetipepathsearch[0]) {
					// fmt.Println("KETEMU PATH NYA", getvaluetipepathsearch[1])
					pathsearch = getvaluetipepathsearch[1]
				}
				value = parsingstruct(path, valuetipedata, namecontrol, pathsearch)
			}
		} else {
			if strings.Contains(valuetipedata, ".") {
				getvaluetipepath := strings.Split(valuetipedata, ".")
				// fmt.Println(getvaluetipepath[1])
				dir := GetDir()
				path = dir + "/" + getvaluetipepath[0] + "/" + namecontrol + ".go"
				value = parsingstruct(path, getvaluetipepath[1], namecontrol, getvaluetipepath[0])
			} else {
				pathsearch := ""
				dir := GetDir()
				// fmt.Println(dir)
				cekpathsearch := strings.Split(path+"/", dir)
				getvaluetipepathsearch := strings.Split(cekpathsearch[1], "/")
				if strings.Contains(path, getvaluetipepathsearch[0]) {
					// fmt.Println("KETEMU PATH NYA", getvaluetipepathsearch[1])
					pathsearch = getvaluetipepathsearch[1]
				}
				value = parsingstruct(path, valuetipedata, namecontrol, pathsearch)
			}
		}
	}
	return
}
