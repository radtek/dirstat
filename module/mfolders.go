package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
)

type folderNode struct {
	container
}

type folderCount struct {
	container
}

func (f *folderNode) LessThan(y interface{}) bool {
	return f.name < y.(*folderNode).name
}

func (f *folderNode) EqualTo(y interface{}) bool {
	return f.name == y.(*folderNode).name
}

func (f *folderCount) LessThan(y interface{}) bool {
	return f.count < y.(*folderCount).count
}

func (f *folderCount) EqualTo(y interface{}) bool {
	return f.count == y.(*folderCount).count
}

// NewFoldersModule creates new folders module
func NewFoldersModule(ctx *Context) Module {
	work := newFoldersWorker(ctx)
	rend := foldersRenderer{work}
	m := moduleFolders{
		work,
		rend,
	}
	return &m
}

// NewFoldersHiddenModule creates new folders module
// that has disabled output
func NewFoldersHiddenModule(ctx *Context) Module {
	work := newFoldersWorker(ctx)
	m := moduleFoldersNoOut{
		work,
		emptyRenderer{},
	}
	return &m
}

type foldersWorker struct {
	total    *totalInfo
	folders  *rbtree.RbTree
	topSize  *rbtree.RbTree
	topCount *rbtree.RbTree
	top      int
}

type foldersRenderer struct {
	foldersWorker
}

type moduleFolders struct {
	foldersWorker
	foldersRenderer
}

type moduleFoldersNoOut struct {
	foldersWorker
	emptyRenderer
}

func newFoldersWorker(ctx *Context) foldersWorker {
	return foldersWorker{
		total:    ctx.total,
		folders:  rbtree.NewRbTree(),
		topSize:  rbtree.NewRbTree(),
		topCount: rbtree.NewRbTree(),
		top:      ctx.top,
	}
}

func (m *foldersWorker) init() {
}

func (m *foldersWorker) postScan() {
	m.folders.WalkInorder(func(node *rbtree.Node) {
		fn := node.Key.(*folderNode)

		insertTo(m.topSize, m.top, &fn.container)

		fcn := folderCount{
			fn.container,
		}

		insertTo(m.topCount, m.top, &fcn)
	})

	m.total.CountFolders = m.folders.Root.Size
}

func (m *foldersWorker) folderHandler(fe *sys.FolderEntry) {
	fn := folderNode{
		container{name: fe.Path, count: fe.Count, size: fe.Size},
	}
	m.folders.Insert(rbtree.NewNode(&fn))
}

func (m *foldersWorker) fileHandler(*sys.FileEntry) {

}

func (f *foldersRenderer) print(p printer) {
	const format = "%v\t%v\t%v\t%v\t%v\n"

	p.print("\nTOP %d folders by size:\n\n", f.top)

	f.outputTableHead(p, format)

	i := 1

	f.topSize.Descend(func(c rbtree.Comparable) bool {

		folder := c.(*container)
		f.outputTableRow(&i, folder, p)

		return true
	})

	p.flush()

	p.print("\nTOP %d folders by count:\n\n", f.top)

	f.outputTableHead(p, format)

	i = 1

	f.topCount.Descend(func(c rbtree.Comparable) bool {

		folder := c.(*folderCount)
		f.outputTableRow(&i, &folder.container, p)

		return true
	})

	p.flush()
}

func (f *foldersRenderer) outputTableRow(i *int, folder *container, p printer) {
	h := fmt.Sprintf("%d. %s", *i, folder.name)

	*i++

	count := folder.count
	sz := uint64(folder.size)

	f.total.printCountAndSizeStatLine(p, count, sz, h)
}

func (f *foldersRenderer) outputTableHead(rc printer, format string) {
	rc.printtab(format, "Folder", "Files", "%", "Size", "%")
	rc.printtab(format, "------", "-----", "------", "----", "------")
}
