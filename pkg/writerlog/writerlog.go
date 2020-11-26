package writerlog

import "io"

type WriterLog struct {
	Writer io.Writer
}

func (e WriterLog) Log(str string) {
	_, _ = e.Writer.Write([]byte(str + "\n"))
}
