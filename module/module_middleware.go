package module

import "dirstat/module/internal/sys"

type onlyFilesWorker struct {
	w worker
}

type onlyFoldersWorker struct {
	w worker
}

func newOnlyFilesWorker(w worker) worker {
	return &onlyFilesWorker{
		w: w,
	}
}

func newOnlyFoldersWorker(w worker) worker {
	return &onlyFoldersWorker{
		w: w,
	}
}

func (f *onlyFilesWorker) handler(evt *sys.ScanEvent) {
	if evt.File != nil {
		f.w.handler(evt)
	}
}

func (f *onlyFoldersWorker) handler(evt *sys.ScanEvent) {
	if evt.Folder != nil {
		f.w.handler(evt)
	}
}
