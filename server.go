package v2rayaguard

import (
	"fmt"
	"time"
)

type server struct {
	id          int
	_type       string
	sub         int
	pingLatency time.Duration
}

func (s server) String() string {
	return fmt.Sprintf(`{"id":%d, "type":%s, "sub":%d, "pingLatency":%v}`, s.id, s._type, s.sub, s.pingLatency)
}

type servers []*server

func (ss servers) Len() int {
	return len(ss)
}
func (ss servers) Less(i, j int) bool {
	return ss[i].pingLatency < ss[j].pingLatency
}
func (ss servers) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}
