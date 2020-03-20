package checks

import (
	"sync"

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
