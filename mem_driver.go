package cache

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type MemDriver struct {
	dataMap    sync.Map
	length     int64
	memorySize int64
	maxMemory  int64
}

// 设置⼀个缓存项，并且在expire时间之后过期
func (m *MemDriver) Set(key string, val interface{}, expire time.Duration) {
	m.dataMap.Store(key, val)
	atomic.AddInt64(&m.length, 1)
	atomic.AddInt64(&m.memorySize, int64(unsafe.Sizeof(val)))
	time.AfterFunc(time.Second*time.Duration(expire), func() {
		atomic.AddInt64(&m.length, -1)
		m.dataMap.Delete(key)
	})
}

// 获取⼀个值
func (m *MemDriver) Get(key string) (interface{}, bool) {
	return m.dataMap.Load(key)
}

// 删除⼀个值
func (m *MemDriver) Del(key string) bool {
	v, ok := m.Get(key)
	if ok {
		atomic.AddInt64(&m.length, -1)
		atomic.AddInt64(&m.memorySize, -int64(unsafe.Sizeof(v)))
		m.dataMap.Delete(key)
		return true
	}
	return false
}

// 检测⼀个值 是否存在
func (m *MemDriver) Exists(key string) bool {
	_, ok := m.dataMap.Load(key)
	return ok
}

// 情况所有值
func (m *MemDriver) Flush() bool {
	m.dataMap.Range(func(k, v interface{}) bool {
		m.dataMap.Delete(k)
		return true
	})

	atomic.StoreInt64(&m.length, 0)
	atomic.StoreInt64(&m.memorySize, 0)

	return true
}

// 返回所有的key 多少
func (m *MemDriver) Keys() int64 {
	return m.length
}
