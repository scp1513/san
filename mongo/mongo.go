package mongo

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// mongodb数据库封装

type Proxy struct {
	session  *mgo.Session // 默认最大连接数 4096
	database *mgo.Database
	daname   string
}

func NewProxy() *Proxy {
	return &Proxy{}
}

// Dial 连接
func (p *Proxy) Dial(host, usr, pwd, db string) error {
	p.daname = db
	var err error
	p.session, err = mgo.Dial(host)
	if err != nil {
		return err
	}
	if usr != "" && pwd != "" && db != "" {
		p.session.DB(p.daname).Login(usr, pwd)
	}
	if err != nil {
		return err
	}
	return nil
}

// M 执行mongodb shell
func (p *Proxy) M(collection string, f func(*mgo.Collection) error) error {
	session := p.session.Clone()
	defer func() {
		session.Close()
		if err := recover(); err != nil {
			log.Println("M", err)
		}
	}()

	c := p.session.DB(p.daname).C(collection)
	return f(c)
}

// Select 执行查询，此方法可拆分做为公共方法
// [Select description]
// @param {[type]} collectionName string [description]
// @param {[type]} query          bson.M [description]
// @param {[type]} sort           bson.M [description]
// @param {[type]} fields         bson.M [description]
// @param {[type]} skip           int    [description]
// @param {[type]} limit          int)   (results      []interface{}, err error [description]
func (p *Proxy) Select(collectionName string, query bson.M, sort string, fields bson.M, skip, limit int) (results []interface{}, err error) {
	exop := func(c *mgo.Collection) error {
		if len(sort) == 0 {
			return c.Find(query).Select(fields).Skip(skip).Limit(limit).All(&results)
		}
		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	err = p.M(collectionName, exop)
	return
}

// SelectAllWithParam 查询所有结果
// 类似Select方法,需要传入比较详细的参数列表，注意：result 必须传递一个slice指针
func (p *Proxy) SelectAllWithParam(collectionName string, query bson.M, sort string, fields bson.M, skip, limit int, results interface{}) error {
	exop := func(c *mgo.Collection) error {
		if len(sort) == 0 {
			return c.Find(query).Select(fields).Skip(skip).Limit(limit).All(results)
		}
		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(results)
	}
	return p.M(collectionName, exop)
}

// SelectAll 查询所有结果
// 比较常用的查询，只带field筛选条件
func (p *Proxy) SelectAll(collectionName string, query, fields bson.M, results interface{}) error {
	exop := func(c *mgo.Collection) error {
		if fields == nil {
			return c.Find(query).Select(fields).All(results)
		}
		return c.Find(query).All(results)
	}
	return p.M(collectionName, exop)
}

// SelectByID 根据id查询结果
// 类似Select方法，注意：result 必须传递一个slice指针
func (p *Proxy) SelectByID(collectionName string, id interface{}, fields bson.M, results interface{}) error {
	exop := func(c *mgo.Collection) error {
		if fields == nil {
			return c.FindId(id).One(results)
		}
		return c.FindId(id).Select(fields).One(results)
	}
	return p.M(collectionName, exop)
}

// SelectOne 查询一个结果，如果查询结果超过一个，会报错
// 类似Select方法，注意：result 必须传递一个slice指针
func (p *Proxy) SelectOne(collectionName string, query, fields bson.M, results interface{}) error {
	exop := func(c *mgo.Collection) error {
		if fields == nil {
			return c.Find(query).Select(fields).One(results)
		}
		return c.Find(query).One(results)
	}
	return p.M(collectionName, exop)
}

// SelectMin 查询最小值
func (p *Proxy) SelectMin(collectionName string, query, fields bson.M, column string, results interface{}) error {
	exop := func(c *mgo.Collection) error {
		if fields == nil {
			return c.Find(query).Sort(column).Select(fields).One(results)
		}
		return c.Find(query).Sort(column).One(results)
	}
	return p.M(collectionName, exop)
}

// SelectMax 查询最大值
func (p *Proxy) SelectMax(collectionName string, query, fields bson.M, column string, results interface{}) error {
	exop := func(c *mgo.Collection) error {
		if fields == nil {
			return c.Find(query).Sort("-" + column).Select(fields).One(results)
		}
		return c.Find(query).Sort("-" + column).One(results)
	}
	return p.M(collectionName, exop)
}

// Exists 查询是否存在
func (p *Proxy) Exists(collectionName string, query bson.M) bool {
	var count int
	exop := func(c *mgo.Collection) error {
		var err error
		count, err = c.Find(query).Count()
		return err
	}
	err := p.M(collectionName, exop)
	if err != nil {
		return false
	}

	return count > 0
}

// Insert 插入文档
func (p *Proxy) Insert(collectionName string, docs interface{}) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(docs)
	}

	return p.M(collectionName, query)
}

// Update 更新文档
func (p *Proxy) Update(collectionName string, selector interface{}, update interface{}) error {
	query := func(c *mgo.Collection) error {
		return c.Update(selector, update)
	}

	return p.M(collectionName, query)
}

// UpdateByID 通过Id更新文档
func (p *Proxy) UpdateByID(collectionName string, id interface{}, update interface{}) error {
	query := func(c *mgo.Collection) error {
		return c.UpdateId(id, update)
	}

	return p.M(collectionName, query)
}

// Delete 删除文档
func (p *Proxy) Delete(collectionName string, selector interface{}) error {
	query := func(c *mgo.Collection) error {
		return c.RemoveId(selector)
	}

	return p.M(collectionName, query)
}
