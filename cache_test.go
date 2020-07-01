package cache

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	driver := new(MemDriver)
	c := New(driver, 11, "3MB")

	c.Set("test", "343435", 1)
	v, _ := c.Get("test")
	t.Log(v)

	time.Sleep(time.Duration(2) * time.Second)
	vv, _ := c.Get("test")
	t.Log(vv)

	c.SetMaxMemory("1GB")

}
