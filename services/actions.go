package services

import (
	rules "github.com/akkeris/service-watcher-f5/rules"
	structs "github.com/akkeris/service-watcher-f5/structs"
	utils "github.com/akkeris/service-watcher-f5/utils"
	virtuals "github.com/akkeris/service-watcher-f5/virtuals"
	"encoding/json"
	"fmt"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"strconv"
	"strings"
)

func ProcessServiceDelete(obj interface{}) {
	servicenamepart := obj.(*corev1.Service).ObjectMeta.Name
	namespacepart := obj.(*corev1.Service).ObjectMeta.Namespace
	appname := servicenamepart + "-" + namespacepart
	if namespacepart == "default" {
		appname = servicenamepart
	}
	internal := IsInternal(namespacepart)
	if !internal {
		virtualhost := appname + "." + utils.Defaultdomain
		rule := rules.BuildRule(appname, utils.Partition, virtualhost, "0", utils.Unipool)
		remove(rule, utils.Virtual)
	}
        if internal {
                virtualhost := appname + "." + utils.InsideDomain
                rule := rules.BuildRule(appname, utils.InsidePartition, virtualhost, "0", utils.Unipool)
                remove(rule, utils.InsideVirtual)
        }

}

func ProcessServiceAdd(obj interface{}) {
	nodeporti := obj.(*corev1.Service).Spec.Ports[0].NodePort
        nodeport := strconv.Itoa(int(nodeporti))
	servicename := obj.(*corev1.Service).ObjectMeta.Name
	namespace := obj.(*corev1.Service).ObjectMeta.Namespace
	appname := servicename + "-" + namespace
	if namespace == "default" {
		appname = servicename
	}
	internal := IsInternal(namespace)
	if !internal {
		virtualhost := appname + "." + utils.Defaultdomain
		rule := rules.BuildRule(appname, utils.Partition, virtualhost, nodeport, utils.Unipool)
		add(rule, utils.Virtual)
	}
        if internal {
                virtualhost := appname + "." + utils.InsideDomain
                rule := rules.BuildRule(appname, utils.InsidePartition, virtualhost, nodeport, utils.Unipool)
                add(rule, utils.InsideVirtual)
        }


}

func IsInternal(space string) bool {

	client := &http.Client{}
	req, err := http.NewRequest("GET", utils.Alamoapilocation+"/v1/space/"+space, nil)
	req.SetBasicAuth(utils.Alamoapiusername, utils.Alamoapipassword)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	bb, err := ioutil.ReadAll(resp.Body)
	var spaceobject structs.Spacespec
	uerr := json.Unmarshal(bb, &spaceobject)
	if uerr != nil {
		fmt.Println(uerr)
	}
	return spaceobject.Internal
}

func getNodes() (n []string, e error) {
	var nodesids []string
	nodesresp := utils.Client.Get().Resource("nodes").Do()
	nodes, err := nodesresp.Raw()
	if err != nil {
		fmt.Println(err)
		return nodesids, err
	}
	var nodelist corev1.NodeList
	err = json.Unmarshal(nodes, &nodelist)
	if err != nil {
		fmt.Println(err)
		return nodesids, err
	}
	for _, element := range nodelist.Items {
		if !element.Spec.Unschedulable {
			uidparts := strings.Split(fmt.Sprintf("%v", element.ObjectMeta.UID), "-")
			nodeid := "uid" + uidparts[0]
			nodesids = append(nodesids, nodeid)
		}

	}
	return nodesids, nil
}
func add(rule structs.Rulespec, virtual string) {
	utils.NewToken()
	rules.AddRule(rule)
	virtuals.AttachRule(rule, virtual)

}

func remove(rule structs.Rulespec, virtual string) {
	utils.NewToken()
	virtuals.DetachRule(rule, virtual)
	rules.DeleteRule(rule)

}
