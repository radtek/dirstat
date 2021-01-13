package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"sort"
)

type extRenderer struct {
	*baseRenderer
	total *totalInfo
	top   int
}

func newExtRenderer(ctx *Context, order int) renderer {
	return &extRenderer{
		baseRenderer: newBaseRenderer(order),
		total:        ctx.total,
		top:          ctx.top,
	}
}

// Renderer method

func (e *extRenderer) render(p out.Printer) {
	extBySize := e.evolventMap(func(agr countSizeAggregate) int64 {
		return int64(agr.Size)
	})

	extByCount := e.evolventMap(func(agr countSizeAggregate) int64 {
		return agr.Count
	})

	sizePrint := fileExtPrint{
		count:   func(f *file) int64 { return e.total.extensions[f.path].Count },
		size:    func(f *file) uint64 { return uint64(f.size) },
		p:       p,
		headfmt: "<gray>TOP %d file extensions by size:</>",
		total:   e.total,
	}

	sizePrint.print(extBySize, e.top)

	countPrint := fileExtPrint{
		count:   func(f *file) int64 { return f.size },
		size:    func(f *file) uint64 { return e.total.extensions[f.path].Size },
		p:       p,
		headfmt: "<gray>TOP %d file extensions by count:</>",
		total:   e.total,
	}

	countPrint.print(extByCount, e.top)
}

type fileExtPrint struct {
	count   func(f *file) int64
	size    func(f *file) uint64
	p       out.Printer
	headfmt string
	total   *totalInfo
}

func (fp *fileExtPrint) print(data files, top int) {
	fp.p.Println()
	fp.p.Cprint(fp.headfmt, top)
	fp.p.Println()
	fp.p.Println()

	tw := newTableWriter(fp.p)

	tw.tab.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 2, Align: text.AlignLeft, AlignHeader: text.AlignLeft, WidthMax: 100},
		{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignLeft},
		{Number: 4, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.percentTransformer},
		{Number: 5, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.sizeTransformer},
		{Number: 6, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.percentTransformer},
	})

	tw.appendHeaders([]string{"#", "Extension", "Count", "%", "Size", "%"})

	sort.Sort(sort.Reverse(data))

	for i := 0; i < top && i < len(data); i++ {
		h := data[i].path

		count := fp.count(data[i])
		sz := fp.size(data[i])

		percentOfCount := fp.total.countPercent(count)
		percentOfSize := fp.total.sizePercent(sz)

		tw.tab.AppendRow([]interface{}{
			i + 1,
			h,
			count,
			percentOfCount,
			sz,
			percentOfSize,
		})
	}

	tw.tab.Render()
}

func (e *extRenderer) evolventMap(mapper func(countSizeAggregate) int64) files {
	var result = make(files, len(e.total.extensions))
	i := 0
	for k, v := range e.total.extensions {
		result[i] = &file{size: mapper(v), path: k}
		i++
	}
	return result
}
