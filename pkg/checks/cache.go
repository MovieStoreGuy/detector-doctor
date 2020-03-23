package checks

import (
	"sync"
	"time"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

var (
	global *cache
	once   sync.Once
)

func getGlobalCache() *cache {
	once.Do(func() {
		global = &cache{
			detectors: make(map[string]*types.Detector),
		}
	})
	return global
}

type cache struct {
	rw        sync.RWMutex
	detectors map[string]*types.Detector
}

func (c *cache) getDetector(ID string) *types.Detector {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.detectors[ID]
}

func (c *cache) setDetector(ID string, det *types.Detector) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.detectors[ID] = det
}

func (c *cache) explictWaitGetDetector(ID string, timeout time.Duration) *types.Detector {
	// Quickly check to see if we have the detector already
	if res := c.getDetector(ID); res != nil {
		return res
	}
	background := func(r chan *types.Detector) {
		timer := time.NewTicker(100 * time.Millisecond)
		explictStop := time.NewTicker(timeout)
		for {
			select {
			case <-timer.C:
				if res := c.getDetector(ID); res != nil {
					r <- res
					return
				}
			case <-explictStop.C:
				r <- nil
				return
			}
		}
	}
	result := make(chan *types.Detector)
	defer close(result)
	go background(result)
	return <-result
}
