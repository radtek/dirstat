package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/aegoroff/godatastruct/rbtree/special"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile_Equal(t *testing.T) {
	var tests = []struct {
		size1  int64
		size2  int64
		result bool
	}{
		{0, 0, true},
		{1, 1, true},
		{1, 0, false},
		{0, 1, false},
		{2, 1, false},
		{1, 2, false},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		f1 := &file{
			path: "/",
			size: test.size1,
		}
		f2 := &file{
			path: "/f",
			size: test.size2,
		}

		// Act
		result := f1.Equal(f2)

		// Assert
		ass.Equal(test.result, result)
	}
}

func TestFile_Less(t *testing.T) {
	var tests = []struct {
		size1  int64
		size2  int64
		result bool
	}{
		{0, 0, false},
		{1, 1, false},
		{1, 0, false},
		{0, 1, true},
		{2, 1, false},
		{1, 2, true},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		f1 := &file{
			path: "/",
			size: test.size1,
		}
		f2 := &file{
			path: "/f",
			size: test.size2,
		}

		// Act
		result := f1.Less(f2)

		// Assert
		ass.Equal(test.result, result)
	}
}

func TestFile_String(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	f := &file{
		path: "/",
		size: 0,
	}

	// Act
	result := f.String()

	// Assert
	ass.Equal("/", result)
}

func Test_printTopFile_invalidCastingError(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	ft := special.NewMaxTree(3)
	ft.Insert(rbtree.Int(1))
	ft.Insert(rbtree.Int(2))
	ft.Insert(rbtree.Int(3))

	fr := topFilesRenderer{topFiles: &topFiles{
		tree: ft,
	}}
	e := out.NewMemoryEnvironment()
	p, _ := e.NewPrinter()

	// Act
	fr.render(p)

	// Assert
	ass.Contains(e.String(), "Invalid casting: expected *file key type but it wasn`t")
}
