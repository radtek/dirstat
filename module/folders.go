package module

import (
	"dirstat/module/internal/sys"
	"errors"
	"github.com/aegoroff/godatastruct/rbtree"
)

// Folder interface
type folderI interface {
	Path() string
	Size() int64
	Count() int64
}

// folder represents file system container that described by path
// and has size and the number of elements in it (count field).
type folder struct {
	path  string
	size  int64
	count int64
	pd    *pathDecorator
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
func (f *folder) Path() string {
	if f.pd == nil {
		return f.path
	}
	return f.pd.decorate(f.path)
}

// Count sortable folder methods

func (fc *folderC) LessThan(y interface{}) bool { return fc.count < y.(*folderC).count }
func (fc *folderC) EqualTo(y interface{}) bool  { return fc.count == y.(*folderC).count }

// Size sortable folder methods

func (fs *folderS) LessThan(y interface{}) bool { return fs.size < y.(*folderS).size }
func (fs *folderS) EqualTo(y interface{}) bool  { return fs.size == y.(*folderS).size }

type foldersWorker struct {
	voidInit
	voidFinalize
	total   *totalInfo
	bySize  *fixedTree
	byCount *fixedTree
	pd      *pathDecorator
}

type foldersRenderer struct {
	*foldersWorker
}

func newFoldersWorker(ctx *Context) *foldersWorker {
	return &foldersWorker{
		total:   ctx.total,
		bySize:  newFixedTree(ctx.top),
		byCount: newFixedTree(ctx.top),
		pd:      ctx.pd,
	}
}

func newFoldersRenderer(work *foldersWorker) renderer {
	return &foldersRenderer{foldersWorker: work}
}

// Worker methods

func (m *foldersWorker) handler(evt *sys.ScanEvent) {
	if evt.Folder == nil {
		return
	}
	fe := evt.Folder

	fn := folder{
		path:  fe.Path,
		count: fe.Count,
		size:  fe.Size,
		pd:    m.pd,
	}

	fs := folderS{fn}
	m.bySize.insert(&fs)

	fc := folderC{fn}
	m.byCount.insert(&fc)
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

func (f *foldersRenderer) print(p printer) {
	const format = "%v\t%v\t%v\t%v\t%v\t%v\n"

	p.cprint("\n<gray>TOP %d folders by size:</>\n\n", f.bySize.size)

	f.printTop(f.bySize, p, format, castSize)

	p.cprint("\n<gray>TOP %d folders by count:</>\n\n", f.byCount.size)

	f.printTop(f.byCount, p, format, castCount)
}

func (f *foldersRenderer) printTop(ft *fixedTree, p printer, format string, cast folderCast) {
	p.tprint(format, " #", "Folder", "Files", "%", "Size", "%")
	p.tprint(format, "--", "------", "-----", "------", "----", "------")

	i := 1

	ft.tree.Descend(func(n rbtree.Node) bool {
		fi, err := cast(n.Key())
		if err != nil {
			p.cprint("<red>%v</>", err)
			return false
		}
		f.printTableRow(&i, fi, p)

		return true
	})

	p.flush()
}

func (f *foldersRenderer) printTableRow(i *int, fi folderI, p printer) {
	count := fi.Count()
	sz := uint64(fi.Size())

	f.total.printCountAndSizeStatLine(p, *i, count, sz, fi.Path())

	*i++
}
