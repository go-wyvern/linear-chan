package models

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/go-wyvern/linear-chan/dbs"
)

type ModelControl struct {
	Models       []interface{}
	SelectParams map[reflect.Type][]string
	Attribute    map[reflect.Type]*Attribute
}

type Attribute struct {
	Select   []string `json:"select" `
	Limit    int      `json:"limit" `
	Unscope  bool     `json:"unscope" `
	Order    string   `json:"order" `
	Orders   []string `json:"orders" `
	Group    string   `json:"group" `
	Transfer bool     `json:"transfer" `
	Join     []string `json:"join" `
}

var Models = new(ModelControl)

func (modelscontrol *ModelControl) RegistModel(model interface{}, params []string) {
	v := reflect.TypeOf(model)
	modelscontrol.Models = append(modelscontrol.Models, model)
	modelscontrol.SelectParams = make(map[reflect.Type][]string)
	modelscontrol.SelectParams[v] = params
	modelscontrol.Attribute = make(map[reflect.Type]*Attribute, len(modelscontrol.Models))
}

func (modelscontrol *ModelControl) SetAttribute(model interface{}, attribute *Attribute) {
	v := reflect.TypeOf(model)
	modelscontrol.Attribute[v] = attribute

}

func (modelscontrol *ModelControl) HasRegist(model interface{}) bool {
	_, ok := modelscontrol.SelectParams[reflect.TypeOf(model)]
	return ok
}

func (modelscontrol *ModelControl) Select(model interface{}, params map[string]interface{}, p *Pagination) *gorm.DB {
	var queryDb *gorm.DB = dbs.Db
	var querystring, querystring2,querystring3 string
	newparams := make(map[string]interface{})
	if Models.HasRegist(model) {
		for _, param := range Models.SelectParams[reflect.TypeOf(model)] {
			if d, ok := params[param]; ok {
				newparams[param] = d
			}
		}
	} else {
		newparams = params
	}
	attr := modelscontrol.Attribute[reflect.TypeOf(model)]
	if p == nil {
		if attr != nil {
			if len(attr.Select) != 0 {
				queryDb = queryDb.Select(attr.Select)
			}
			if len(attr.Group) != 0 {
				queryDb = queryDb.Group(attr.Group)
			}
			if attr.Unscope {
				queryDb = queryDb.Unscoped()
			}
			if attr.Limit != 0 {
				queryDb = queryDb.Limit(attr.Limit)
			}
			if attr.Order != "" {
				queryDb = queryDb.Order(attr.Order)
			} else {
				queryDb = queryDb.Order("id desc")
			}
			for _, order := range attr.Orders {
				if order != "" {
					queryDb = queryDb.Order(order)
				}
			}
		} else {
			queryDb = queryDb.Order("id desc")
		}
		queryDb = queryDb.Model(model)
		for k, v := range newparams {
			querystring = fmt.Sprintf(" %v = ?", k)
			querystring3 = fmt.Sprintf(" %v like ?", k)
			if reflect.ValueOf(v).Kind() == reflect.String {
				if strings.Contains(v.(string), "or") && len(strings.Split(v.(string), " ")) >= 3 {
					querystring = fmt.Sprintf(" %v in (?)", k)
					var newparam []string
					for _, vl := range strings.Split(v.(string), " ") {
						if vl == "or" {
							continue
						}
						newparam = append(newparam, vl)
					}
					queryDb = queryDb.Where(querystring, newparam)
				} else if strings.Contains(v.(string), "and") && len(strings.Split(v.(string), " ")) >= 3 {
					for _, vl := range strings.Split(v.(string), " ") {
						if vl == "and" {
							continue
						}
						queryDb = queryDb.Where(querystring, vl)
					}
				} else if strings.Contains(v.(string), "not") && len(strings.Split(v.(string), " ")) >= 2 {
					for _, vl := range strings.Split(v.(string), " ") {
						if vl == "not" {
							continue
						}
						vl = strings.Replace(vl, " ", "", -1)
						queryDb = queryDb.Not(querystring, vl)
					}
				} else {
					if strings.Contains(v.(string),"like"){
						v="%"+strings.Split(v.(string)," ")[1]+"%"
						fmt.Println(querystring3,v)
						queryDb = queryDb.Where(querystring3, v)
					}else{
						queryDb = queryDb.Where(querystring, v)
					}
				}
			} else if reflect.ValueOf(v).Kind() == reflect.Slice {
				var whereparam []interface{}
				for i := 0; i < reflect.ValueOf(v).Len(); i++ {
					if i == 0 {
						querystring2 = querystring
					} else {
						querystring2 = querystring2 + " or " + querystring
					}
					// if i == 0 {
					// 	queryDb = queryDb.Where(querystring, reflect.ValueOf(v).Index(i).Interface())
					// } else {
					// 	queryDb = queryDb.Or(querystring, reflect.ValueOf(v).Index(i).Interface())
					// }
				}
				for i := 0; i < reflect.ValueOf(v).Len(); i++ {
					whereparam = append(whereparam, reflect.ValueOf(v).Index(i).Interface())
				}
				queryDb = queryDb.Where(querystring2, whereparam...)
			} else {
				queryDb = queryDb.Where(querystring, v)

			}
		}
	} else {
		ipage, _ := strconv.ParseInt(p.UrlQuery.Get("page"), 10, 0)
		page := int(ipage)
		if page <= 0 {
			page = 1
		}
		p.Page = page
		p.Query = dbs.Db.Model(model)
		if attr != nil {
			if len(attr.Select) != 0 {
				p.Query = p.Query.Select(attr.Select)
			}
			if len(attr.Group) != 0 {
				p.Query = p.Query.Group(attr.Group)
			}
			if attr.Unscope {
				p.Query = p.Query.Unscoped()
			}
			if attr.Order != "" {
				p.Query = p.Query.Order(attr.Order)
			} else {
				p.Query = p.Query.Order("id desc")
			}
			for _, order := range attr.Orders {
				if order != "" {
					queryDb = queryDb.Order(order)
				}
			}
		} else {
			p.Query = p.Query.Order("id desc")
		}
		for k, v := range newparams {
			querystring = fmt.Sprintf(" %v = ?", k)
			querystring3 = fmt.Sprintf(" %v like ?", k)
			if reflect.ValueOf(v).Kind() == reflect.String {
				if strings.Contains(v.(string), "or") && len(strings.Split(v.(string), " ")) >= 3 {
					querystring = fmt.Sprintf(" %v in (?)", k)
					var newparam []string
					for _, vl := range strings.Split(v.(string), " ") {
						if vl == "or" {
							continue
						}
						newparam = append(newparam, vl)

					}
					p.Query = p.Query.Where(querystring, newparam)
				} else if strings.Contains(v.(string), "and") && len(strings.Split(v.(string), " ")) >= 3 {
					for _, vl := range strings.Split(v.(string), " ") {
						if vl == "and" {
							continue
						}
						p.Query = p.Query.Where(querystring, vl)
					}
				} else if strings.Contains(v.(string), "not") && len(strings.Split(v.(string), " ")) >= 2 {
					for _, vl := range strings.Split(v.(string), " ") {
						if vl == "not" {
							continue
						}
						vl = strings.Replace(vl, " ", "", -1)
						p.Query = p.Query.Not(querystring, vl)
					}
				} else {
					if strings.Contains(v.(string),"like"){
						v=strings.Split(v.(string)," ")[1]+"%"
						fmt.Println(querystring3,v)
						p.Query = p.Query.Where(querystring3, v)
					}else{
						p.Query = p.Query.Where(querystring, v)
					}
				}
			} else if reflect.ValueOf(v).Kind() == reflect.Slice {
				var whereparam []interface{}
				for i := 0; i < reflect.ValueOf(v).Len(); i++ {
					if i == 0 {
						querystring2 = querystring
					} else {
						querystring2 = querystring2 + " or " + querystring
					}
					// if i == 0 {
					// 	p.Query = p.Query.Where(querystring, reflect.ValueOf(v).Index(i).Interface())
					// } else {
					// 	p.Query = p.Query.Or(querystring, reflect.ValueOf(v).Index(i).Interface())
					// }
				}
				for i := 0; i < reflect.ValueOf(v).Len(); i++ {
					whereparam = append(whereparam, reflect.ValueOf(v).Index(i).Interface())
				}
				p.Query = p.Query.Where(querystring2, whereparam...)
			} else {
				p.Query = p.Query.Where(querystring, v)
			}
		}
		queryDb = p.Paginate(page)
	}
	return queryDb
}

type JoinAttr struct {
	SelfModel interface{}
	SelfTable string
	JoinModel interface{}
	SelfValue interface{}
	JoinTable string
	SelfId    string
	JoinId    string
}

func (j *JoinAttr) JoinString() string {
	if j.SelfTable == "" {
		j.SelfTable = dbs.Db.NewScope(j.SelfModel).TableName()
	}
	if j.JoinTable == "" {
		j.JoinTable = dbs.Db.NewScope(j.JoinModel).TableName()
	}
	return fmt.Sprintf("join %s on %s.%s=%s.%s", j.SelfTable, j.JoinTable, j.JoinId, j.SelfTable, j.SelfId)
}

func SetJoin(model interface{}, join_id string) *gorm.DB {
	joinattr := new(JoinAttr)
	joinattr.SelfModel = model
	v := reflect.ValueOf(model)
	t:=reflect.Indirect(v).Type()
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		json_tag := strings.Split(f.Tag.Get("json"),",")[0]
		if joinattr.JoinTable != ""&&joinattr.JoinId != "" {
			if json_tag == joinattr.JoinTable {
				fmt.Println(joinattr.JoinString())
				return dbs.Db.Joins(joinattr.JoinString()).Where(fmt.Sprintf("%s.%s = ?", joinattr.SelfTable, joinattr.SelfId), joinattr.SelfValue)
			}
		}
		if json_tag == join_id {
			joinattr.SelfId = json_tag
			join_tag := f.Tag.Get("join")
			if join_tag != "" {
				joinattr.SelfValue = v.Field(i).Interface()
				joinattr.JoinTable = strings.Split(join_tag, ".")[0]
				joinattr.JoinId = strings.Split(join_tag, ".")[1]
			}
		}else {
			continue
		}
	}
	return nil
}
