package cache

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	driver := new(MemDriver)
	c := NewCache(driver, 11)

	c.Set("test", "343435", 1)
	v, _ := c.Get("test")
	t.Log(v)

	time.Sleep(time.Duration(2) * time.Second)
	vv, _ := c.Get("test")
	t.Log(vv)

}
