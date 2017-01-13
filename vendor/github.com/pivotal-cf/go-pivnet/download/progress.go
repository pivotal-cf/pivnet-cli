package download

import (
	"io"

	pb "gopkg.in/cheggaaa/pb.v1"
)

type Bar struct {
	*pb.ProgressBar
}

func NewBar() Bar {
	b := pb.New(0)
	b.SetUnits(pb.U_BYTES)
	b.SetWidth(80)
	return Bar{b}
}

func (b Bar) SetTotal(contentLength int64) {
	b.Total = contentLength
}

func (b Bar) Kickoff() {
	b.Start()
}

func (b Bar) SetOutput(output io.Writer) {
	b.Output = output
}
