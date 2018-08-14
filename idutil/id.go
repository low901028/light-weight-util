package idutil

import (
	"sync"
	"time"
	"math"
)

const (
	tsLen  = 5 * 8  			// 时间戳长度
	cntLen = 8      			// 计数长度
	suffixLen = tsLen + cntLen	// 后缀长度
)

// 全局唯一ID生成器：基于计数、时间戳、节点ID
// 初始化id格式： 高位2bytes存放节点ID; 接下来5bytes存放时间戳; 接着低位1byte存放计数
// | prefix   | suffix              |
// | 2 bytes  | 5 bytes   | 1 byte  |
// | memberID | timestamp | cnt     |
// 说明：
// 1、时间戳：当机器重启在1ms后，35年前该机器的时间戳是不同的
// 2、全局id在上一个id的基础上 在后缀累加生成
// 3、计数字段有可能会溢出到时间戳字段，这是被允许的：主要是为了拓展事件窗口到2^56.同样这并不能破坏机器重启生成ID的唯一性
//     由于etcd的吞吐量远远小于1毫秒256个(1秒2.5个)请求
type Generator struct {
	mu sync.Mutex
	prefix uint64  // 高 2 bytes
	suffix uint64  // 低 6 bytes
}

func NewGenerator(memberID uint64, now time.Time) *Generator{
	// 将节点ID左移suffixLen，低位补0 填充
	prefix := uint64(memberID) << suffixLen
	// 获取当前节点的时间戳转为毫秒单位
	unixMilli := uint64(now.UnixNano()) / uint64(time.Millisecond / time.Nanosecond)
	// 将当前节点毫秒时间戳右移tsLen；高位置为0 填充时间戳内容
	// 再进行左移 cntLen 填充计数内容
	suffix := lowbit(unixMilli, tsLen) << cntLen
	return &Generator{
		prefix: prefix,
		suffix: suffix,
	}
}

// 生成id
func (g *Generator) Next() uint64{
	g.mu.Lock()
	defer  g.mu.Unlock()
	g.suffix++
	id := g.prefix | lowbit(g.suffix, suffixLen)
	return id
}

// 获取低位
func lowbit(x uint64, n uint) uint64{
	return x & (math.MaxUint64 >> (64 - n))
}