package examples

import (
	"io"
	"os"
	"time"
)

//go:generate goor -type=csv -setter -output=csv_gen.go
type csv struct {
	name string
	date time.Time
	file os.File
	w    io.Writer `goor:"constructor:-"`
}
