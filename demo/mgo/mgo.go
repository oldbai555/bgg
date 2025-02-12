/**
 * @Author: zjj
 * @Date: 2025/2/12
 * @Desc:
**/

package mgo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/oldbai555/lbtool/pkg/lberr"
)

type MongoSt struct {
	Port     uint32
	Host     string
	Username string
	Password string
	Database string
	session  *mgo.Session
}

func NewMongoSt(port uint32, host string, username, password string) *MongoSt {
	return &MongoSt{Port: port, Host: host, Username: username, Password: password}
}

func (m *MongoSt) Dial() error {
	session, err := mgo.Dial(fmt.Sprintf("%s:%d", m.Host, m.Port))
	if err != nil {
		return err
	}
	m.session = session
	return nil
}

func (m *MongoSt) Close() {
	if m.session != nil {
		m.session.Close()
	}
}

func (m *MongoSt) DoLogic(f func(s *mgo.Session) error) error {
	if m.session == nil {
		return lberr.NewCustomErr("mongo db session is nil")
	}
	cloneSession := m.session.Clone()
	defer cloneSession.Close()
	err := f(cloneSession)
	return lberr.Wrap(err)
}
