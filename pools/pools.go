package pools

import (
	structs "github.com/akkeris/service-watcher-f5/structs"
	utils "github.com/akkeris/service-watcher-f5/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
)

func BuildPool(appname string, port string, partition string, monitor string, nodes []string) structs.Poolspec {

	var pool structs.Poolspec
	pool.Name = appname + "-pool"
	pool.Partition = partition
	pool.Monitor = monitor
	var members []structs.Memberspec

	for _, element := range nodes {
		var member structs.Memberspec
		member.Name = "/" + partition + "/" + element + ":" + port
		members = append(members, member)
	}
	pool.Members = members
	return pool

}

func AddPool(pool structs.Poolspec) {

	str, err := json.Marshal(pool)
	if err != nil {
		fmt.Println("Error preparing request")
	}
	jsonStr := []byte(string(str))
	urlStr := utils.F5url + "/mgmt/tm/ltm/pool"
	req, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonStr))
	req.Header.Add("X-F5-Auth-Token", utils.F5token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := utils.F5Client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Printf("%v : Add Pool %v\n", resp.Status, pool.Name)
	if resp.StatusCode > 299 {
		bodyj, _ := simplejson.NewFromReader(resp.Body)
		fmt.Println(bodyj)
	}
}
func DeletePool(pool structs.Poolspec) {

	urlStr := utils.F5url + "/mgmt/tm/ltm/pool/~" + pool.Partition + "~" + pool.Name
	req, _ := http.NewRequest("DELETE", urlStr, nil)
	req.Header.Add("X-F5-Auth-Token", utils.F5token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := utils.F5Client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Printf("%v : Delete Pool %v\n", resp.Status, pool.Name)
	if resp.StatusCode > 299 {
		bodyj, _ := simplejson.NewFromReader(resp.Body)
		fmt.Println(bodyj)
	}

}
