package gindocs

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var menumd = "## Daftar API\n"
var bodymd = ""
var countdata = 1

type BodyMd struct {
	Code     string
	Response string
}

var databodymd []BodyMd

func setmenumd(method, url string) {
	menumd += strconv.Itoa(countdata) + `.  [` + method + `] ` + url + "\n"
	countdata++
}
func seturl(method, url string) {
	body += "\n### [" + method + "] " + url + "\n"
}
func setrequestmd(jsondata string) {
	body += "*REQUEST JSON*\n"
	body += "\n```json\n"
	body += jsondata
	body += "\n```\n"
}

func setresponsemd(dataresp []BodyMd) {
	for _, a := range dataresp {
		body += "*RESPONSE CODE " + a.Code + "*\n"
		body += "\n```json\n"
		body += a.Response
		body += "\n```\n"
	}
}

func ExecMarkdown() {
	fmt.Println(menumd)
	dir := GetDir()
	os.Mkdir(dir+"/readme", 0755)
	err := ioutil.WriteFile("readme/readme.md", []byte(menumd+body), 0755)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}
}
