package module

import (
	"errors"
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/aegoroff/godatastruct/rbtree/special"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// folder represents file system container that described by path
// and has size and the number of elements in it (count field).
type folder struct {
	path  string
	size  int64
	count int64
	pd    decorator
}

// Count sortable folder
type folderC struct {
	folder
}

// Size sortable folder
type folderS struct {
	folder
}

// Path sortable folder methods

func (f *folder) String() string { return f.Path() }
func (f *folder) Size() int64    { return f.size }
func (f *folder) Count() int64   { return f.count }
func (f *folder) Path() string   { return f.pd.decorate(f.path) }

// Count sortable folder methods

func (fc *folderC) Less(y rbtree.Comparable) bool  { return fc.count < y.(*folderC).count }
func (fc *folderC) Equal(y rbtree.Comparable) bool { return fc.count == y.(*folderC).count }

// Size sortable folder methods

func (fs *folderS) Less(y rbtree.Comparable) bool  { return fs.size < y.(*folderS).size }
func (fs *folderS) Equal(y rbtree.Comparable) bool { return fs.size == y.(*folderS).size }

type folders struct {
	bySize  rbtree.RbTree
	byCount rbtree.RbTree
}

type foldersHandler struct {
	*folders
	pd decorator
}

type foldersRenderer struct {
	*folders
	*baseRenderer
	total *totalInfo
}

func newFolders(top int) *folders {
	return &folders{
		bySize:  special.NewMaxTree(int64(top)),
		byCount: special.NewMaxTree(int64(top)),
	}
}

func newFoldersHandler(fc *folders, pd decorator) scan.Handler {
	h := &foldersHandler{
		folders: fc,
		pd:      pd,
	}
	return newOnlyFoldersHandler(h)
}

func newFoldersRenderer(f *folders, ctx *Context, order int) renderer {
	return &foldersRenderer{
		folders:      f,
		total:        ctx.total,
		baseRenderer: newBaseRenderer(order),
	}
}

// Worker method

func (m *foldersHandler) Handle(evt *scan.Event) {
	fe := evt.Folder

	fn := folder{
		path:  fe.Path,
		count: fe.Count,
		size:  fe.Size,
		pd:    m.pd,
	}

	fs := folderS{fn}
	m.bySize.Insert(&fs)

	fc := folderC{fn}
	m.byCount.Insert(&fc)
}

// Renderer method

type folderCast func(c rbtree.Comparable) (folderI, error)

func castSize(c rbtree.Comparable) (folderI, error) {
	f, ok := c.(*folderS)

	if !ok {
		return nil, errors.New("invalid casting: expected *folderS key type but it wasn`t")
	}
	return f, nil
}

func castCount(c rbtree.Comparable) (folderI, error) {
	f, ok := c.(*folderC)

	if !ok {
		return nil, errors.New("invalid casting: expected *folderC key type but it wasn`t")
	}
	return f, nil
}

func (f *foldersRenderer) render(p out.Printer) {
	p.Cprint("\n<gray>TOP %d folders by size:</>\n\n", f.bySize.Len())

	f.printTop(f.bySize, p, castSize)

	p.Cprint("\n<gray>TOP %d folders by count:</>\n\n", f.byCount.Len())

	f.printTop(f.byCount, p, castCount)
}

func (f *foldersRenderer) printTop(ft rbtree.RbTree, p out.Printer, cast folderCast) {
	tw := newTableWriter(p)

	tw.tab.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 2, Align: text.AlignLeft, AlignHeader: text.AlignLeft, WidthMax: 100},
		{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignLeft},
		{Number: 4, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.percentTransformer},
		{Number: 5, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.sizeTransformer},
		{Number: 6, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.percentTransformer},
	})

	tw.appendHeaders([]string{"#", "Folder", "Files", "%", "Size", "%"})

	i := 1

	it := rbtree.NewDescend(ft).Iterator()

	for it.Next() {
		fi, err := cast(it.Current())
		if err != nil {
			p.Cprint("<red>%v</>", err)
			return
		}

		count := fi.Count()
		sz := uint64(fi.Size())
		percentOfCount := f.total.countPercent(count)
		percentOfSize := f.total.sizePercent(sz)

		tw.tab.AppendRow([]interface{}{
			i,
			fi.Path(),
			count,
			percentOfCount,
			sz,
			percentOfSize,
		})

		i++
	}

	tw.tab.Render()
}
