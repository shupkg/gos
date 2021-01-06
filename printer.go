package gos

type Printer interface {
	Printf(format string, args ...interface{})
}
