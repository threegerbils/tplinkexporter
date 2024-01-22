package clients

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type TPLINK1016DESwitch struct {
	host     string
	username string
	password string
}

func (client *TPLINK1016DESwitch) GetHost() string {
	return client.host
}

func (client *TPLINK1016DESwitch) GetPortStats() ([]portStats, error) {
	type allInfo struct {
		State      []int
		LinkStatus []int
		Pkts       []int
	}
	// "http://IP/logon.cgi"
	resp, err := http.PostForm(fmt.Sprintf("http://%s/logon.cgi", client.host), url.Values{"username": {client.username}, "password": {client.password}, "logon": {"Login"}})
	if err != nil {
		// handle error
		return nil, err
	}
	defer resp.Body.Close()
	// fmt.Println(resp, err)
	// "http://IP/PortStatisticsRpm.htm")
	resp2, err := http.Get(fmt.Sprintf("http://%s/PortStatisticsRpm.htm", client.host))
	if err != nil {
		// handle error
		return nil, err
	}
	defer resp2.Body.Close()
	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		// handle error
		return nil, err
	}
	// fmt.Println(string(body))
	// var jbody string = strings.ReplaceAll(
	// 	strings.ReplaceAll(
	// 		strings.ReplaceAll(
	// 			string(body), "link_status", `"linkStatus"`),
	// 		"state", `"State"`),
	// 	"pkts", `"Pkts"`)
	var jbody string = string(body)
	maxPort, err := strconv.Atoi(regexp.MustCompile(`max_port_num = ([^;]*);`).FindStringSubmatch(jbody)[1])
	fmt.Println(maxPort)
	res := regexp.MustCompile(`tmp_info = "([^"]*)";`).FindStringSubmatch(jbody)
	if res == nil {
		// fmt.Println(jbody)
		return nil, errors.New("unexpected response for port statistics http call: " + jbody)
	}
	info := strings.Split(strings.TrimRight(res[1], " "), " ")
	fmt.Println(info)
	columnRes := regexp.MustCompile(`<td class="TABLE_HEAD_BOTTOM" align=center width="78px">([^<]*)</td>`).FindAllStringSubmatch(jbody, -1)
	columns := make([]string, len(columnRes)-1)

	for i := 1; i < len(columnRes); i++ {
		columns[i-1] = columnRes[i][1]
	}
	fmt.Println(columns)
	var parsedInfo []map[string]int = make([]map[string]int, maxPort)
	for port := 0; port < maxPort; port++ {
		parsedInfo[port] = make(map[string]int)
		for col := 0; col < len(columns); col++ {
			parsedInfo[port][columns[col]], err = strconv.Atoi(info[port*len(columns)+col])
		}
	}
	//json.Unmarshal([]byte(res[1]), &jparsed)
	// fmt.Println(jparsed)
	var portsInfos []portStats
	portcount := maxPort
	for i := 0; i < portcount; i++ {
		var portInfo portStats
		portInfo.State = parsedInfo[i]["Status"]
		portInfo.LinkStatus = parsedInfo[i]["Link Status"]
		if portInfo.State == 1 {
			portInfo.PktCount = make(map[string]int)
			portInfo.PktCount["TxGoodPkt"] = parsedInfo[i]["TxGoodPkt"]
			portInfo.PktCount["TxBadPkt"] = parsedInfo[i]["TxBadPkt"]
			portInfo.PktCount["RxGoodPkt"] = parsedInfo[i]["RxGoodPkt"]
			portInfo.PktCount["RxBadPkt"] = parsedInfo[i]["RxBadPkt"]
		}
		portsInfos = append(portsInfos, portInfo)
	}
	fmt.Println(portsInfos)
	return portsInfos, nil
}

/*
sample output of PortStatisticsRpm.htm call:
<script>
var max_port_num = 8;
var port_middle_num  = 16;
var all_info = {
state:[1,1,1,1,1,1,1,1,0,0],
link_status:[6,6,0,6,0,0,0,5,0,0],
pkts:[1901830310,0,1338131260,33254,4291149014,0,2311488878,564,0,0,0,0,1814018004,0,33552310,0,0,0,0,0,0,0,0,0,0,0,0,0,1678459124,0,1866169392,0,0,0]
};
var tip = "";
</script>
*/

func NewTPLink1016DESwitch(host string, username string, password string) *TPLINK1016DESwitch {
	return &TPLINK1016DESwitch{
		host:     host,
		username: username,
		password: password,
	}
}
