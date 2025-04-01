/**
 * @Author: zjj
 * @Date: 2025/4/16
 * @Desc:
**/

package sqllite

import (
	"database/sql"
	"github.com/oldbai555/lbtool/log"

	_ "github.com/mattn/go-sqlite3"
)

func Test() error {
	// 打开数据库
	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer db.Close()

	// 创建表
	sqlStmt := `CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 插入数据
	stmt, err := db.Prepare("INSERT INTO users(name) VALUES(?)")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec("Alice")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	_, err = stmt.Exec("Bob")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 查询数据
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		log.Infof("ID: %d, Name: %s", id, name)
	}

	// 检查查询错误
	err = rows.Err()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
