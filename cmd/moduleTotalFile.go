package cmd

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type moduleTotalFile struct {
	aggregate map[Range]fileStat
	total     *totalInfo
}

func (m *moduleTotalFile) postScan() {

}

func (m *moduleTotalFile) handler() fileHandler {
	return func(f *fileEntry) {
	}
}

func (m *moduleTotalFile) output(tw *tabwriter.Writer, w io.Writer) {
	_, _ = fmt.Fprintf(w, "Total files stat:\n\n")

	const format = "%v\t%v\t%v\t%v\t%v\n"

	_, _ = fmt.Fprintf(tw, format, "File size", "Amount", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "---------", "------", "------", "----", "------")

	heads := createRangesHeads()
	for i, r := range fileSizeRanges {
		count := m.aggregate[r].TotalFilesCount
		sz := m.aggregate[r].TotalFilesSize

		outputTopStatLine(tw, count, m.total, sz, heads[i])
	}
	_ = tw.Flush()
}
