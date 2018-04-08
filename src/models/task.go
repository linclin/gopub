package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Task struct {
	Id             int       `orm:"column(id);auto"`
	UserId         uint      `orm:"column(user_id)"`
	ProjectId      int       `orm:"column(project_id)"`
	Action         int16     `orm:"column(action)"`
	Status         int16     `orm:"column(status)"`
	Title          string    `orm:"column(title);size(100);null"`
	LinkId         string    `orm:"column(link_id);size(20);null"`
	ExLinkId       string    `orm:"column(ex_link_id);size(20);null"`
	CommitId       string    `orm:"column(commit_id);size(800);null"`
	CreatedAt      time.Time `orm:"column(created_at);type(datetime);null"`
	UpdatedAt      time.Time `orm:"column(updated_at);type(datetime);null"`
	Branch         string    `orm:"column(branch);size(100);null"`
	FileList       string    `orm:"column(file_list);null"`
	EnableRollback int       `orm:"column(enable_rollback)"`
	PmsBatchId     int       `orm:"column(pms_batch_id);null"`
	PmsUworkId     int       `orm:"column(pms_uwork_id);null"`
	IsRun          int       `orm:"column(is_run);null"`
	FileMd5        string    `orm:"column(file_md5);size(200);null"`
	Hosts          string    `orm:"column(hosts);null"`
}

func (t *Task) TableName() string {
	return "task"
}

func init() {
	orm.RegisterModel(new(Task))
}

// AddTask insert a new Task into database and returns
// last inserted Id on success.
func AddTask(m *Task) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTaskById retrieves Task by Id. Returns error if
// Id doesn't exist
func GetTaskById(id int) (v *Task, err error) {
	o := orm.NewOrm()
	v = &Task{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTask retrieves all Task matches certain condition. Returns empty list if
// no records exist
func GetAllTask(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Task))
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

	var l []Task
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

// UpdateTask updates Task by Id and returns error if
// the record to be updated doesn't exist
func UpdateTaskById(m *Task) (err error) {
	o := orm.NewOrm()
	v := Task{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTask deletes Task by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTask(id int) (err error) {
	o := orm.NewOrm()
	v := Task{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Task{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
func GetAllTaskAndPro(pro_name string, startTime string, endTime string) ([]orm.Params, error) {
	o := orm.NewOrm()
	var events []orm.Params
	sql := "SELECT title,task.id,commit_id,branch,project_id,project.`name` as project_name,task.updated_at FROM task LEFT JOIN project on task.project_id=project.id WHERE action=0 AND task.`status`=3 AND project.level=3 "
	if pro_name != "" {
		sql += `AND project.repo_url like"%` + pro_name + `%" `
	}
	if startTime != "" {
		sql += `and task.updated_at>"` + startTime + `" `
	}
	if endTime != "" {
		sql += `and task.updated_at<"` + endTime + `" `
	}
	sql += `order by task.id DESC `
	_, err := o.Raw(sql).Values(&events)
	return events, err
}
