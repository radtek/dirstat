package module

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
	"time"
)

const (
	_ int64 = 1 << (10 * iota)
	kbyte
	mbyte
	gbyte
	tbyte
	pbyte
)

// Range defined integer value range
type Range struct {
	// Min value
	Min int64

	// Max value
	Max int64
}

type ranges []Range

// Contains defines whether the number specified within range
func (r *Range) Contains(num int64) bool {
	return num >= r.Min && num <= r.Max
}

type fileStat struct {
	TotalFilesSize  uint64
	TotalFilesCount int64
}

type totalInfo struct {
	ReadingTime   time.Duration
	FilesTotal    countSizeAggregate
	CountFolders  int64
	CountFileExts int
	extensions    map[string]countSizeAggregate
}

type countSizeAggregate struct {
	Count int64
	Size  uint64
}

type fixedTree struct {
	tree rbtree.RbTree
	size int
}

func (t *totalInfo) countPercent(count int64) float64 {
	return percent(float64(count), float64(t.FilesTotal.Count))
}

func (t *totalInfo) sizePercent(size uint64) float64 {
	return percent(float64(size), float64(t.FilesTotal.Size))
}

func percent(value float64, total float64) float64 {
	return (value / total) * 100
}

func (t *totalInfo) printCountAndSizeStatLine(p printer, ix int, count int64, sz uint64, title string) {
	percentOfCount := t.countPercent(count)
	percentOfSize := t.sizePercent(sz)

	p.tprint("%2d\t%v\t%v\t%.2f%%\t%v\t%.2f%%\n", ix, title, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
}

func newFixedTree(sz int) *fixedTree {
	return &fixedTree{
		tree: rbtree.NewRbTree(),
		size: sz,
	}
}

// insert inserts node into tree which size is limited
// Only <size> max nodes will be in the tree
func (t *fixedTree) insert(c rbtree.Comparable) {
	min := t.tree.Minimum()
	if t.tree.Len() < int64(t.size) || min.Key().LessThan(c) {
		if t.tree.Len() == int64(t.size) {
			t.tree.DeleteNode(min.Key())
		}

		t.tree.Insert(c)
	}
}

func (r ranges) heads(addNum bool) []string {
	var heads []string
	for i, r := range r {
		h := fmt.Sprintf("Between %s and %s", human(r.Min), human(r.Max))
		if addNum {
			h = fmt.Sprintf("%2d. %s", i+1, h)
		}
		heads = append(heads, h)
	}
	return heads
}
