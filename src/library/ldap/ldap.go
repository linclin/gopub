package ldap

import (
	"fmt"
	"github.com/astaxie/beego"
	ldap "gopkg.in/ldap.v3"
	"strings"
)

type Ldap struct {
	link *ldap.Conn
}
type Ldap_user struct {
	Uid       string
	Cn        string
	Sn        string
	Email     string
	UidNumber string
}

func (l *Ldap) Connect() (e error) {
	ldapHost := beego.AppConfig.String("ldapHost")
	ldapPort, _ := beego.AppConfig.Int("ldapPort")
	beego.Debug(fmt.Sprintf("try to connect ldap: %s:%d", ldapHost, ldapPort))

	link, e := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))
	if e != nil {
		beego.Info(e)
		return e
	}

	l.link = link
	return e
}

func (l *Ldap) Search(baseDn string, query string) (rs *ldap.SearchResult, e error) {
	beego.Debug("search:" + query)
	searchRequest := ldap.NewSearchRequest(
		baseDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		query,         // The filter to apply
		[]string{"*"}, // A list attributes to retrieve
		nil,
	)
	rs, e = l.link.Search(searchRequest)
	if e != nil {
		beego.Info(e)
	}
	return rs, e
}

func (l *Ldap) SearchUser(uid string) (users []Ldap_user, e error) {
	baseDn := beego.AppConfig.String("ldapPeopleDn")
	query := fmt.Sprintf("(uid=%s)", uid)
	beego.Debug("search user" + baseDn + query)
	rs, e := l.Search(baseDn, query)
	if e == nil {
		for _, entry := range rs.Entries {
			user := l.formatUser(entry)
			users = append(users, user)
		}
		return users, nil
	} else {
		return nil, e
	}

}
func (l *Ldap) formatUser(entry *ldap.Entry) (user Ldap_user) {
	user.Cn = entry.GetAttributeValue("cn")
	user.Sn = entry.GetAttributeValue("sn")
	user.Uid = entry.GetAttributeValue("uid")
	user.UidNumber = entry.GetAttributeValue("uidNumber")
	user.Email = entry.GetAttributeValue("mail")
	return user
}
func (l *Ldap) AuthByUidAndPassword(uid string, password string) (user Ldap_user, e error) {
	userDn := strings.Replace(beego.AppConfig.String("ldapPeopleDnTpl"), "{uid}", uid, -1)
	beego.Debug("auth user: " + userDn)
	e = l.link.Bind(userDn, password)

	if e == nil {
		beego.Debug("ldap auth succ")
		baseDn := beego.AppConfig.String("ldapPeopleDn")
		userEntry, e := l.Search(baseDn, fmt.Sprintf("(uid=%s)", uid))
		if e != nil {
			return user, e
		} else {
			user := l.formatUser(userEntry.Entries[0])
			return user, nil
		}
	} else {
		return user, e
	}
}
func (l *Ldap) SearchGroupCn(query string) (cn string, e error) {

	baseDn := beego.AppConfig.String("ldapGroupDn")
	rs, e := l.Search(baseDn, query)
	if e == nil {
		if len(rs.Entries) > 0 {
			return rs.Entries[0].GetAttributeValue("cn"), nil
		} else {
			beego.Info("user is not in cronsun group " + query)
			return "", fmt.Errorf("user is not in cronsun group")
		}
	} else {
		return "", e
	}
}
