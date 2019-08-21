package utils

import (
	"fmt"
	"k8s.io/client-go/rest"
	"os"
	vault "github.com/akkeris/vault-client"
	"strings"
)

var Partition string
var InsidePartition string

var Monitor string
var InsideMonitor string

var Virtual string
var InsideVirtual string

var Defaultdomain string
var InsideDomain string

var Alamoapilocation string
var Alamoapiusername string
var Alamoapipassword string

var Client rest.Interface

var Blacklist map[string]bool

var Unipool string
func Variableinit() {

	fmt.Println("setting Partition to " + os.Getenv("F5_PARTITION"))
	Partition = os.Getenv("F5_PARTITION")

	fmt.Println("setting InsidePartition to " + os.Getenv("F5_INSIDE_PARTITION"))
	InsidePartition = os.Getenv("F5_INSIDE_PARTITION")

	fmt.Println("setting Monitor to " + os.Getenv("F5_MONITOR"))
	Monitor = os.Getenv("F5_MONITOR")

	fmt.Println("setting InsideMonitor to " + os.Getenv("F5_INSIDE_MONITOR"))
	InsideMonitor = os.Getenv("F5_INSIDE_MONITOR")

	fmt.Println("setting f5virtual to " + os.Getenv("F5_VIRTUAL"))
	f5virtual := os.Getenv("F5_VIRTUAL")

	fmt.Println("setting f5insidevirtual to " + os.Getenv("F5_INSIDE_VIRTUAL"))
	f5insidevirtual := os.Getenv("F5_INSIDE_VIRTUAL")

	fmt.Println("setting Virtual to " + "~" + Partition + "~" + f5virtual)
	Virtual = "~" + Partition + "~" + f5virtual

	fmt.Println("setting InsideVirtual to " + "~" + InsidePartition + "~" + f5insidevirtual)
	InsideVirtual = "~" + InsidePartition + "~" + f5insidevirtual

	fmt.Println("setting Defaultdomain to " + os.Getenv("DEFAULT_DOMAIN"))
	Defaultdomain = os.Getenv("DEFAULT_DOMAIN")

	fmt.Println("setting InsideDomain to " + os.Getenv("INSIDE_DOMAIN"))
	InsideDomain = os.Getenv("INSIDE_DOMAIN")

	alamoapisecret := os.Getenv("ALAMOAPI_SECRET")
	alamoapiusername := vault.GetField(alamoapisecret, "username")
	alamoapipassword := vault.GetField(alamoapisecret, "password")
	alamoapilocation := os.Getenv("REGIONAPI_LOCATION")
	Alamoapilocation = alamoapilocation
	Alamoapiusername = alamoapiusername
	Alamoapipassword = alamoapipassword
    
        Unipool = os.Getenv("UNIPOOL")

}

func InitBlacklist() {
	Blacklist = make(map[string]bool)
	blackliststring := os.Getenv("NAMESPACE_BLACKLIST")
	blacklistslice := strings.Split(blackliststring, ",")
	for _, element := range blacklistslice {
		Blacklist[element] = true
	}
	keys := make([]string, 0, len(Blacklist))
	for k := range Blacklist {
		keys = append(keys, k)
	}

	fmt.Printf("Setting blacklist to %v\n", strings.Join(keys, ","))

}
