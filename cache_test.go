package cache

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	driver := new(MemDriver)
	c := New(driver, 11, "3MB")

	c.Set("test", "343435", 1)
	c.Set("test1", "3434235", 11)
	c.Set("test4", "3434235", 11)
	c.Set("test3", "3434235", 133)
	c.Set("test2", "3434235", 122)

	v, _ := c.Get("test")
	t.Log(v)

	time.Sleep(time.Duration(2) * time.Second)
	vv, _ := c.Get("test")
	t.Log(vv)

	c.SetMaxMemory("1GB")
	t.Log(c.Exists("test1"))
	t.Log(c.Exists("testzz"))
	c.Del("test1")
	t.Log(c.Exists("test1"))
	t.Log(c.Keys())
	c.Flush()
	t.Log(c.Exists("test3"))
	t.Log(c.Keys())
}
