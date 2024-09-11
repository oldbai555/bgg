/**
 * @Author: zjj
 * @Date: 2024/8/22
 * @Desc:
**/

package roommgr

import (
	"github.com/oldbai555/bgg/singlesrv/server/wsmgr"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"google.golang.org/protobuf/proto"
	"sync"
)

const globalRoomId = 20240822

var onceChatRoom sync.Once
var singleChatRoom *ChatRoomSt

type ChatRoomSt struct {
	roomId  uint64
	userMgr map[uint64]struct{}
}

func GetSingleChatRoom() *ChatRoomSt {
	if singleChatRoom == nil {
		onceChatRoom.Do(func() {
			singleChatRoom = &ChatRoomSt{
				roomId:  globalRoomId,
				userMgr: make(map[uint64]struct{}),
			}
		})

	}
	return singleChatRoom
}

func (s *ChatRoomSt) AddUser(uid uint64) {
	_, ok := s.userMgr[uid]
	if ok {
		return
	}
	s.userMgr[uid] = struct{}{}
}

func (s *ChatRoomSt) DelUser(uid uint64) {
	_, ok := s.userMgr[uid]
	if !ok {
		return
	}
	delete(s.userMgr, uid)
}

func (s *ChatRoomSt) Broadcast(msg proto.Message) error {
	bytes, err := jsonpb.Marshal(msg)
	if err != nil {
		log.Errorf("unmarshal err:%v", err)
		return err
	}
	for uid := range s.userMgr {
		wsmgr.WriteBytes(uid, bytes)
	}
	return nil
}
