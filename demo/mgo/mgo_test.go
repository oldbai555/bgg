/**
 * @Author: zjj
 * @Date: 2025/2/12
 * @Desc:
**/

package mgo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"testing"
)

// 定义数据结构
type User struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Name  string        `bson:"name"`
	Age   int           `bson:"age"`
	Email string        `bson:"email"`
}

// 开启认证
// sudo vim /etc/mongod.conf
// 追加内容
// security:
//     authorization: enabled
// 重启
// sudo systemctl restart mongod

// 创建管理员
// use admin

// db.createUser({
//  	user: "admin",
//  	pwd: "admin",  // 替换为你的密码
//  	roles: [{ role: "userAdminAnyDatabase", db: "admin" }]
// })

// db.getUsers()

// mongo -u admin -p yourpassword --authenticationDatabase admin

// 创建普通用户

// use testDB

//db.createUser({
//	user: "oldbai",
//	pwd: "oldbai",
//	roles: [{ role: "readWrite", db: "testDB" }]
//})

// db.auth("oldbai","oldbai")

func TestNewMongoSt(t *testing.T) {
	st := NewMongoSt(17017, "192.168.226.4", "oldbai", "oldbai")
	err := st.Dial()
	if err != nil {
		t.Log(err)
		return
	}
	err = st.DoLogic(func(s *mgo.Session) error {
		db := s.DB("testDB")
		collection := db.C("user")
		// 插入文档
		user := User{
			Name:  "李四",
			Age:   25,
			Email: "lisi@example.com",
		}
		err = collection.Insert(&user)
		if err != nil {
			return lberr.Wrap(err)
		}
		t.Log("插入文档成功，ID:", user.ID)

		// 查询单个文档
		var result User
		err = collection.Find(bson.M{"name": "李四"}).One(&result)
		if err != nil {
			return lberr.Wrap(err)
		}
		t.Logf("查询结果: %+v\n", result)

		// 更新文档
		err = collection.Update(bson.M{"name": "李四"}, bson.M{"$set": bson.M{"age": 26}})
		if err != nil {
			return lberr.Wrap(err)
		}
		t.Log("更新文档成功")

		// 查询多个文档
		var users []User
		err = collection.Find(bson.M{}).All(&users)
		if err != nil {
			return lberr.Wrap(err)
		}
		t.Log("查询多个文档:")
		for _, user := range users {
			fmt.Printf("%+v\n", user)
		}

		// 删除文档
		err = collection.Remove(bson.M{"name": "李四"})
		if err != nil {
			return lberr.Wrap(err)
		}
		t.Log("删除文档成功")

		return nil
	})
	if err != nil {
		t.Log(err)
		return
	}
}
