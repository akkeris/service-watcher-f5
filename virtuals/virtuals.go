package virtuals

import (
	structs "github.com/akkeris/service-watcher-f5/structs"
	utils "github.com/akkeris/service-watcher-f5/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
)

func GetVirtuals() {
	urlStr := utils.F5url + "/mgmt/tm/ltm/virtual"
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("X-F5-Auth-Token", utils.F5token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := utils.F5Client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println("Get Virtuals: " + resp.Status)
	if resp.StatusCode > 299 {
		bodyj, _ := simplejson.NewFromReader(resp.Body)
		fmt.Println(bodyj)
	}
}

func DetachRule(rule structs.Rulespec, virtual string) {
	rules := getRules(virtual)
	var newrules []string
	for _, element := range rules {
		if element != "/"+rule.Partition+"/"+rule.Name {
			newrules = append(newrules, element)
		}
	}
	var virtualo structs.Virtualspec
	virtualo.Rules = newrules

	str, err := json.Marshal(virtualo)
	if err != nil {
		fmt.Println("Error preparing request")
	}
	jsonStr := []byte(string(str))
	urlStr := utils.F5url + "/mgmt/tm/ltm/virtual/" + virtual
	req, _ := http.NewRequest("PATCH", urlStr, bytes.NewBuffer(jsonStr))
	req.Header.Add("X-F5-Auth-Token", utils.F5token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := utils.F5Client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Printf("%v : Detach Rule %v from %v\n", resp.Status, rule.Name, virtual)
	if resp.StatusCode > 299 {
		bodyj, _ := simplejson.NewFromReader(resp.Body)
		fmt.Println(bodyj)
	}

}

func AttachRule(rule structs.Rulespec, virtual string) {
	rules := getRules(virtual)
	rules = append(rules, "/"+rule.Partition+"/"+rule.Name)
	var virtualo structs.Virtualspec
	virtualo.Rules = rules

	str, err := json.Marshal(virtualo)
	if err != nil {
		fmt.Println("Error preparing request")
	}
	jsonStr := []byte(string(str))
	urlStr := utils.F5url + "/mgmt/tm/ltm/virtual/" + virtual
	req, _ := http.NewRequest("PATCH", urlStr, bytes.NewBuffer(jsonStr))
	req.Header.Add("X-F5-Auth-Token", utils.F5token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := utils.F5Client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Printf("%v : Attach Rule %v to %v\n", resp.Status, rule.Name, virtual)
	if resp.StatusCode > 299 {
		bodyj, _ := simplejson.NewFromReader(resp.Body)
		fmt.Println(bodyj)
	}

}

func getRules(virtual string) []string {

	urlStr := utils.F5url + "/mgmt/tm/ltm/virtual/" + virtual
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("X-F5-Auth-Token", utils.F5token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := utils.F5Client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Printf("%v : Get Rules from %v\n", resp.Status, virtual)
	bodyj, _ := simplejson.NewFromReader(resp.Body)
	if resp.StatusCode > 299 {
		bodyj, _ := simplejson.NewFromReader(resp.Body)
		fmt.Println(bodyj)
	}
	var rulesa []string
	rules, _ := bodyj.Get("rules").Array()
	for index, _ := range rules {
		value := rules[index]
		rulesa = append(rulesa, value.(string))
	}
	return rulesa

}
