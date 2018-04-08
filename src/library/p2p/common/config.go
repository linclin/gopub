package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	nettool "github.com/toolkits/net"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// 定义配置映射的结构体
type Config struct {
	Server bool //是否为服务端

	Name string

	DownDir string //只有客户端才配置

	Log string

	Net struct {
		IP       string
		MgntPort int
		DataPort int

		AgentMgntPort int
		AgentDataPort int
	}

	Auth struct {
		Username string
		Password string
	}

	Control *Control
}

type Control struct {
	Speed     int
	MaxActive int
	CacheSize int
}

func normalFile(dir string) string {
	if !filepath.IsAbs(dir) {
		pwd, _ := os.Getwd()
		dir = filepath.Join(pwd, dir)
		dir, _ = filepath.Abs(dir)
		dir = filepath.Clean(dir)
		return dir
	}
	return dir
}

func (c *Config) defaultValue() {
	c.DownDir = normalFile(c.DownDir)
	f, err := os.Stat(c.DownDir)
	if err == nil || !os.IsExist(err) {
		os.MkdirAll(c.DownDir, os.ModePerm)
	} else {
		if !f.IsDir() {
			fmt.Printf("DownDir is not a directory")
			os.Exit(6)
		}
	}

	if c.Log != "" {
		c.Log = normalFile(c.Log)
	}

	if c.Control == nil {
		c.Control = &Control{Speed: 20, MaxActive: 10, CacheSize: 50}
	}
	if c.Net.IP == "" {
		//获取本机ip
		LocalIps, err := nettool.IntranetIP()
		if err != nil {
			fmt.Printf("get LocalIp error")
			os.Exit(6)
		}
		daemon := LocalIps[0]
		c.Net.IP = daemon
		if c.Net.AgentDataPort == 0 {
			c.Net.AgentDataPort = 1902
		}
		if c.Net.AgentMgntPort == 0 {
			c.Net.AgentMgntPort = 1901
		}
		if c.Net.DataPort == 0 {
			c.Net.DataPort = 1902
		}
		if c.Net.MgntPort == 0 {
			c.Net.MgntPort = 1901
		}
	}
	if c.Control.Speed == 0 {
		c.Control.Speed = 20
	}
	if c.Control.MaxActive == 0 {
		c.Control.MaxActive = 10
	}
	if c.Control.CacheSize == 0 {
		c.Control.CacheSize = 50
	}
}

func (c *Config) validate() error {
	if c.Server {
		if c.Net.AgentMgntPort == 0 {
			return errors.New("Not set Net.AgentMgntPort in server config file")
		}
		if c.Net.AgentDataPort == 0 {
			return errors.New("Not set Net.AgentDataPort in server config file")
		}
	}

	if !c.Server {
		if c.DownDir == "" {
			return errors.New("Not set DownDir in client config file")
		}
	}

	if c.Auth.Username == "" || c.Auth.Password == "" {
		return errors.New("Not set auth in config file")
	}
	return nil
}

func ParserConfig(c *Config) (*Config, error) {
	c.defaultValue()
	return c, nil
}

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// 随机字符串
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

func GenPasswd(pwd string, pwdlen int) (string, string) {
	solt := string(Krand(pwdlen, KC_RAND_KIND_LOWER))
	h := md5.New()
	h.Write([]byte(pwd + solt))
	sum := h.Sum(nil)
	sendpwd := hex.EncodeToString(sum)

	return sendpwd, solt
}
func CmpPasswd(pwd string, solt string, sendpwd string) bool {
	h := md5.New()
	h.Write([]byte(pwd + solt))
	sum := h.Sum(nil)
	sendpwd1 := hex.EncodeToString(sum)
	if sendpwd1 == sendpwd {
		return true
	} else {
		return false
	}
}

/**
 * 根据path读取文件中的内容，返回字符串
 * 建议使用绝对路径，例如："./schema/search/appoint.json"
 */
func ReadFile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

func ReadJson(path string) Config {
	jsonStr := ReadFile(path)
	ret := Config{}
	err := json.Unmarshal([]byte(jsonStr), &ret)
	if err != nil {
		panic("文件[" + path + "]的内容不是指定格式")
	}
	return ret
}
