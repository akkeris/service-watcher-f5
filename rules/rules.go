package rules

import (
	structs "github.com/akkeris/service-watcher-f5/structs"
	utils "github.com/akkeris/service-watcher-f5/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
)

func BuildRule(appname string, partition string, virtualhost string, nodeport string, unipool string) structs.Rulespec {
	var rule structs.Rulespec
	rule.Name = appname + "-rule"
	rule.Partition = partition
	rule.ApiAnonymous = "when HTTP_REQUEST { \n switch [string tolower [HTTP::host]] { \n\"" + virtualhost + "\" {\nset new_port \""+nodeport+"\"\npool /" + partition + "/" + unipool+"}}}"

	return rule
}

func AddRule(rule structs.Rulespec) {

	str, err := json.Marshal(rule)
	if err != nil {
		fmt.Println("Error preparing request")
	}
	jsonStr := []byte(string(str))
	urlStr := utils.F5url + "/mgmt/tm/ltm/rule"
	req, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonStr))
	req.Header.Add("X-F5-Auth-Token", utils.F5token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := utils.F5Client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Printf("%v : Add Rule %v\n", resp.Status, rule.Name)
	if resp.StatusCode > 299 {
		bodyj, _ := simplejson.NewFromReader(resp.Body)
		fmt.Println(bodyj)
	}
}
func DeleteRule(rule structs.Rulespec) {

	urlStr := utils.F5url + "/mgmt/tm/ltm/rule/~" + rule.Partition + "~" + rule.Name
	req, _ := http.NewRequest("DELETE", urlStr, nil)
	req.Header.Add("X-F5-Auth-Token", utils.F5token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := utils.F5Client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Printf("%v : Delete Rule %v\n", resp.Status, rule.Name)
	if resp.StatusCode > 299 {
		bodyj, _ := simplejson.NewFromReader(resp.Body)
		fmt.Println(bodyj)
	}

}
