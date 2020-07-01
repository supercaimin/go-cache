package cache

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type MemDriver struct {
	// 使用sync.Map解决并发问题
	dataMap sync.Map
	// 当前缓存个数
	count             int64
	memorySize        int64
	maxMemory         int64
	defaultExpiration time.Duration
}

// 设置⼀个缓存项，并且在expire时间之后过期
func (m *MemDriver) Set(key string, val interface{}, expire time.Duration) {
	m.MemoryCheck(int64(unsafe.Sizeof(val)))
	m.dataMap.Store(key, val)
	atomic.AddInt64(&m.count, 1)
	atomic.AddInt64(&m.memorySize, int64(unsafe.Sizeof(val)))
	if expire != NoExpiration {
		if expire == DefaultExpiration {
			expire = m.defaultExpiration
		}
		time.AfterFunc(time.Second*time.Duration(expire), func() {
			atomic.AddInt64(&m.count, -1)
			m.dataMap.Delete(key)
		})
	}

}

// 获取⼀个值
func (m *MemDriver) Get(key string) (interface{}, bool) {
	return m.dataMap.Load(key)
}

// 删除⼀个值
func (m *MemDriver) Del(key string) bool {
	v, ok := m.Get(key)
	if ok {
		atomic.AddInt64(&m.count, -1)
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

	atomic.StoreInt64(&m.count, 0)
	atomic.StoreInt64(&m.memorySize, 0)

	return true
}

// 返回所有的key 多少
func (m *MemDriver) Keys() int64 {
	return atomic.LoadInt64(&m.count)
}

// 设置最大内存
func (m *MemDriver) SetMaxMemory(size int64) {
	atomic.StoreInt64(&m.maxMemory, size)
}

func (m *MemDriver) SetDefaultExpiration(expire time.Duration) {
	m.defaultExpiration = expire
}

// 检查内存是否够用，若不够用可用采用以下算法淘汰缓存
// 1.LRU缓存淘汰算法
// 2.最近过期时间淘汰算法
// 3.随机淘汰
// 这里使用随机淘汰
func (m *MemDriver) MemoryCheck(size int64) {
	memorySize := atomic.LoadInt64(&m.memorySize)
	maxMemory := atomic.LoadInt64(&m.maxMemory)

	if memorySize+size > maxMemory {
		m.dataMap.Range(func(k, v interface{}) bool {
			m.dataMap.Delete(k)
			atomic.AddInt64(&m.count, -1)
			atomic.AddInt64(&m.memorySize, -int64(unsafe.Sizeof(v)))
			memorySize -= int64(unsafe.Sizeof(v))
			if memorySize+size > maxMemory {
				return true
			}
			return false
		})
	}
}
