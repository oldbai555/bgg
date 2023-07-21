package game

import (
	"container/list"
	"github.com/name5566/leaf/gate"
	"sync"
)

// todo 后续使用 redis
var matchQueue = NewMatchQueue()

type MatchQueue struct {
	// 匹配队列
	queue *list.List

	survivalMap map[uint64]struct{}

	sync.Mutex
}

func NewMatchQueue() *MatchQueue {
	return &MatchQueue{
		queue:       list.New(),
		survivalMap: make(map[uint64]struct{}),
	}
}

func (w *MatchQueue) Pop() (*MatchItem, bool) {
	w.Lock()
	defer w.Unlock()

	pItem := w.queue.Front()

	w.queue.Remove(pItem)

	p, ok := pItem.Value.(*MatchItem)

	delete(w.survivalMap, p.PId)

	return p, ok
}

func (w *MatchQueue) Push(vList ...*MatchItem) {
	w.Lock()
	defer w.Unlock()

	for _, v := range vList {

		_, ok := w.survivalMap[v.PId]

		if !ok {
			w.queue.PushBack(v)
			w.survivalMap[v.PId] = struct{}{}
		}
	}
}

type MatchItem struct {
	PId uint64
	A   gate.Agent
	Seq uint32
}

func NewMatchItem(pid uint64, a gate.Agent) *MatchItem {
	return &MatchItem{PId: pid, A: a}
}
