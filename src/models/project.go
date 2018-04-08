package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Project struct {
	Id             int       `orm:"column(id);auto"`
	UserId         uint      `orm:"column(user_id)"`
	Name           string    `orm:"column(name);size(100);null"`
	Level          int16     `orm:"column(level)"`
	Status         int16     `orm:"column(status)"`
	Version        string    `orm:"column(version);size(32);null"`
	RepoUrl        string    `orm:"column(repo_url);size(800);null"`
	RepoUsername   string    `orm:"column(repo_username);size(50);null"`
	RepoPassword   string    `orm:"column(repo_password);size(100);null"`
	RepoMode       string    `orm:"column(repo_mode);size(50);null"`
	RepoType       string    `orm:"column(repo_type);size(10);null"`
	DeployFrom     string    `orm:"column(deploy_from);size(200)"`
	Excludes       string    `orm:"column(excludes);null"`
	ReleaseUser    string    `orm:"column(release_user);size(50)"`
	ReleaseTo      string    `orm:"column(release_to);size(200)"`
	ReleaseLibrary string    `orm:"column(release_library);size(200)"`
	Hosts          string    `orm:"column(hosts);null"`
	PreDeploy      string    `orm:"column(pre_deploy);null"`
	PostDeploy     string    `orm:"column(post_deploy);null"`
	PreRelease     string    `orm:"column(pre_release);null"`
	PostRelease    string    `orm:"column(post_release);null"`
	LastDeploy     string    `orm:"column(last_deploy);null"`
	Audit          int16     `orm:"column(audit);null"`
	KeepVersionNum int       `orm:"column(keep_version_num)"`
	CreatedAt      time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(datetime);null"`
	P2p            int16     `orm:"column(p2p)"`
	HostGroup      string    `orm:"column(host_group)"`
	Gzip           int16     `orm:"column(gzip)"`
	PmsProName     string    `orm:"column(pms_pro_name);size(200)"`
}

func (t *Project) TableName() string {
	return "project"
}

func init() {
	orm.RegisterModel(new(Project))
}

// AddProject insert a new Project into database and returns
// last inserted Id on success.
func AddProject(m *Project) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetProjectById retrieves Project by Id. Returns error if
// Id doesn't exist
func GetProjectById(id int) (v *Project, err error) {
	o := orm.NewOrm()
	v = &Project{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllProject retrieves all Project matches certain condition. Returns empty list if
// no records exist
func GetAllProject(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Project))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Project
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateProject updates Project by Id and returns error if
// the record to be updated doesn't exist
func UpdateProjectById(m *Project) (err error) {
	o := orm.NewOrm()
	v := Project{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteProject deletes Project by Id and returns error if
// the record to be deleted doesn't exist
func DeleteProject(id int) (err error) {
	o := orm.NewOrm()
	v := Project{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Project{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
