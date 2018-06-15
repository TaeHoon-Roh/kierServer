package counter

import (
	"time"
)

type Counter struct {
	MsecPerDay int
	Start int
}

func (c *Counter) InitCounter()  {
	c.MsecPerDay = 86400000
	c.Start = c.counter()
}

func (c *Counter) InitCounterInParameter(msec int) {
	c.Start = msec
}

func (c *Counter) counter() int {
	temp := time.Now().UnixNano()
	return int((temp/(1000000)) % int64(c.MsecPerDay))
}

func (c *Counter) diff(mark int, start int) int {
	if(mark >= start){
		return mark - start
	} else {
		return mark - start + c.MsecPerDay
	}
}

func (c *Counter) reset() int {
	now := c.counter()
	elapsed := c.diff(now, c.Start)
	c.Start = now
	return elapsed
}

func (c *Counter) Elapsed() int  {
	return c.diff(c.counter(), c.Start)
}

func (c *Counter) setStart(msec int) {
	c.Start = msec
}