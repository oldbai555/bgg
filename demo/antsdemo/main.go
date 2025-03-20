package main

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/panjf2000/ants/v2"
	"math"
	"runtime"
	"time"
)

func main() {
	defer func() {
		log.GetLogger().Flush()
	}()
	numCPU := runtime.NumCPU()
	pool, err := ants.NewPoolWithFunc(numCPU, func(data interface{}) {
		log.Infof("data:%v", data)
		time.Sleep(time.Second)
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	defer func() {
		pool.Release()
	}()
	for i := 0; i < math.MaxInt32; i++ {
		err := pool.Invoke(i)
		if err != nil {
			log.Errorf("err:%v", err)
		}
	}
	for pool.Waiting() != 0 {
	}
	log.Infof("All tasks are done ")
}
