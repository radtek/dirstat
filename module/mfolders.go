package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"io"
	"text/tabwriter"
)

type moduleFolders struct {
	total   *totalInfo
	folders map[string]*container
	tree    *rbtree.RbTree
}

func (m *moduleFolders) init() {
}

func (m *moduleFolders) postScan() {
	for _, cont := range m.folders {
		cont.insertTo(m.tree)
	}
}

func (m *moduleFolders) handler() sys.FileHandler {
	return func(f *sys.FileEntry) {
	}
}

func (m *moduleFolders) output(tw *tabwriter.Writer, w io.Writer) {
	const format = "%v\t%v\t%v\t%v\t%v\n"

	_, _ = fmt.Fprintf(w, "\nTOP %d folders by size:\n\n", top)
	_, _ = fmt.Fprintf(tw, format, "Folder", "Files", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "------", "-----", "------", "----", "------")

	i := 1

	m.tree.Descend(func(c *rbtree.Comparable) bool {

		folder := (*c).(*container)
		h := fmt.Sprintf("%d. %s", i, folder.name)

		i++

		count := folder.count
		sz := uint64(folder.size)

		outputTopStatLine(tw, count, m.total, sz, h)

		return true
	})

	_ = tw.Flush()
}
