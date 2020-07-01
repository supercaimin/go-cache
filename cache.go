package cache

import (
	"strconv"
	"strings"
	"time"
)

const (
	// 过期时间设置为-1时，永不过期
	NoExpiration time.Duration = -1
	// 过期时间设置为0，使用默认过期时间
	DefaultExpiration time.Duration = 0
)

/**
⽀持过期时间和最⼤内存⼤⼩的的内存缓存库。
*/
type Cache struct {
	driver Driver

	// 解析后的最大内存
	maxMemory int64
	// 默认最大内存
	defaultmaxMemory int64
	// 默认过期时间
	defaultExpiration time.Duration
}

func New(driver Driver, defaultExpiration time.Duration, defaultMaxMemory string) *Cache {
	c := &Cache{
		driver:            driver,
		defaultExpiration: defaultExpiration,
	}
	c.SetMaxMemory(defaultMaxMemory)
	c.SetDefaultExpiration(defaultExpiration)
	return c
}

//size 是⼀个字符串。⽀持以下参数: 1KB，100KB，1MB，2MB，1GB 等
func (c *Cache) SetMaxMemory(size string) bool {
	if len(size) <= 2 {
		return false
	}
	size = strings.ToUpper(size)
	unit := size[len(size)-2:]
	x, err := strconv.ParseInt(size[:len(size)-2], 10, 64)
	if err != nil {
		return false
	}
	switch unit {
	case "KB":
		c.maxMemory = x * 1024
	case "MB":
		c.maxMemory = x * 1024 * 1024
	case "GB":
		c.maxMemory = x * 1024 * 1024 * 1024
	default:
		return false
	}
	c.driver.SetMaxMemory(c.maxMemory)
	return true
}

func (c *Cache) SetDefaultExpiration(expire time.Duration) {
	c.driver.SetDefaultExpiration(expire)
}

// 设置⼀个缓存项，并且在expire时间之后过期
func (c *Cache) Set(key string, val interface{}, expire time.Duration) {
	c.driver.Set(key, val, expire)
}

// 获取⼀个值
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.driver.Get(key)
}

// 删除⼀个值
func (c *Cache) Del(key string) bool {
	return c.driver.Del(key)
}

// 检测⼀个值 是否存在
func (c *Cache) Exists(key string) bool {
	return c.driver.Exists(key)
}

// 情况所有值
func (c *Cache) Flush() bool {
	return c.driver.Flush()
}

// 返回所有的key 多少
func (c *Cache) Keys() int64 {
	return c.driver.Keys()
}
