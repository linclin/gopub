package jumpserver

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strings"
)

// type gid2hostinfo struct {
// 	gid ip2hostname
// }
type authinfo struct {
	Token string
	//	User interface{}
}

type asset struct {
	Ip       string `json:ip`
	Hostname string `json:hostname`
}

type node struct {
	Id    string `json:id`
	Value string `json:value`
}

func auth() (string, error) {
	auth_api_url := beego.AppConfig.String("jumpserver") + beego.AppConfig.String("jump_auth_api")
	param := "{\"username\": \"" + beego.AppConfig.String("jump_username") + "\", \"password\": \"" + beego.AppConfig.String("jump_password") + "\"}"
	resp, _ := http.Post(auth_api_url, "application/json", strings.NewReader(param))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	rs := authinfo{}
	err = json.Unmarshal([]byte(body), &rs)
	if err != nil {
		return "", nil
	}
	return string(rs.Token), nil
}
func GetGroups() (map[string]string, error) {
	token, err := auth()
	if err != nil {
		return nil, err
	}

	jumpserver_grouplist_api_url := beego.AppConfig.String("jumpserver") + beego.AppConfig.String("jump_grouplist_api")

	client := &http.Client{}
	request, err := http.NewRequest("GET", jumpserver_grouplist_api_url, strings.NewReader(""))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(request)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var rs []node
	err = json.Unmarshal(body, &rs)
	id2group := make(map[string]string)
	if len(rs) > 0 {
		for _, nodeinfo := range rs {
			if strings.HasPrefix(nodeinfo.Value, "up_") == true {
				id2group[nodeinfo.Id] = nodeinfo.Value
			}
		}
	}
	if err != nil {
		return nil, nil
	}
	return id2group, nil
}

func GetIpsByGroupid(group_id string) (map[string]string, error) {
	token, err := auth()
	if err != nil {
		return nil, err
	}

	jumpserver_grouplist_api_url := beego.AppConfig.String("jumpserver") + strings.Replace(beego.AppConfig.String("jump_groupid2ips_api"), "%id", string(group_id), -1)
	client := &http.Client{}
	request, err := http.NewRequest("GET", jumpserver_grouplist_api_url, strings.NewReader(""))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(request)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var rs []asset
	err = json.Unmarshal([]byte(body), &rs)
	ip2hostname := make(map[string]string)
	if len(rs) > 0 {
		for _, v := range rs {
			ip2hostname[v.Ip] = v.Hostname
		}
	}
	if err != nil {
		return nil, nil
	}
	return ip2hostname, nil
}
