package underarock

import (
	"sync"
)

// ConcurrentSlice type that can be safely shared between goroutines
type ConcurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

// ConcurrentSliceItem contains the index/value pair of an item in a
// concurrent slice
type ConcurrentSliceItem struct {
	Index int
	Value interface{}
}

// NewConcurrentSlice creates a new concurrent slice
func NewConcurrentSlice() *ConcurrentSlice {
	cs := &ConcurrentSlice{
		items: make([]interface{}, 0),
	}

	return cs
}

// Append adds an item to the concurrent slice
func (cs *ConcurrentSlice) Append(item interface{}) {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, item)
}

// Append adds an item to the concurrent slice
func (cs *ConcurrentSlice) Size() int {
	cs.Lock()
	defer cs.Unlock()

	return len(cs.items)
}

// Iter iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builin range keyword
func (cs *ConcurrentSlice) Iter(limit int) <-chan ConcurrentSliceItem {
	ch := make(chan ConcurrentSliceItem)

	f := func() {
		cs.Lock()
		defer cs.Unlock()
		// for index, value := range cs.items {
		// 	c <- ConcurrentSliceItem{index, value}
		// }
		c := 0
		for n := len(cs.items) - 1; n >= 0; n-- {
			ch <- ConcurrentSliceItem{Index: n, Value: cs.items[n]}
			if limit > 0 {
				if c > (limit - 1) {
					break
				}
				c++
			}
		}
		close(ch)
	}
	go f()

	return ch
}

var MessageSlice = &ConcurrentSlice{}

func Top20() []Message {
	var top20 []Message
	queue := MessageSlice.Iter(20)
	for elem := range queue {
		//log.Println("from channel", elem.Value.(*Message))

		original, ok := elem.Value.(*Message)
		if ok {
		}
		top20 = append(top20, *original)
	}
	return top20
}

func Top20For(id string) []Message {
	var top20 []Message
	queue := MessageSlice.Iter(-1)
	for elem := range queue {
		//log.Println("from channel", elem.Value.(*Message))

		original, ok := elem.Value.(*Message)
		if ok {
		}
		if original.ToID == id {
			top20 = append(top20, *original)
		}
	}
	return top20
}
func Top20From(myid string, fid string) []Message {
	var top20 []Message
	queue := MessageSlice.Iter(-1)
	for elem := range queue {
		//log.Println("from channel", elem.Value.(*Message))

		original, ok := elem.Value.(*Message)
		if ok {
		}
		if original.ToID == myid && original.FromID == fid {
			top20 = append(top20, *original)
		}
		if original.ToID == fid && original.FromID == myid {
			top20 = append(top20, *original)
		}
	}
	return top20
}

func AddMessage(m *Message) {
	//log.Println("to Queue", m)
	MessageSlice.Append(m)
}
