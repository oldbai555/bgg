package mysql

import (
	"github.com/oldbai555/bgg/service/lbddzserver/impl/dao"
	"github.com/oldbai555/lbtool/log"
)

var (
	PlayerOrm dao.IPlayerDao

	RoomOrm dao.IRoomDao

	GameOrm dao.IGameDao

	GamePlayerOrm dao.IGamePlayerDao
)

func RegisterOrm(dsn string) (err error) {
	log.Infof("start init db orm......")

	err = InitMasterOrm(dsn)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	AutoMigrate()

	PlayerOrm, err = NewPlayerImpl()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	RoomOrm, err = NewRoomImpl()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	GameOrm, err = NewGameImpl()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	GamePlayerOrm, err = NewGamePlayerImpl()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	log.Infof("end init db orm......")
	return
}
