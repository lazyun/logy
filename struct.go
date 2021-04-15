package logy

import "fmt"

type Logger interface {
	Format(LogFields, ...interface{}) string
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Formatf(LogFields, ...interface{}) string
}

type loggerY struct {
}

func (y loggerY) Format(fields LogFields, args ...interface{}) []interface{} {
	args[0] = fmt.Sprintf("uuid: %s, path: %s, func: %s, line: %d, %+v",
		fields.UUID, fields.FilePath, fields.FuncName, fields.Line, args)

	return args
}

func (y loggerY) Debug(args ...interface{}) {
	fmt.Println(append([]interface{}{"Debug "}, args...)...)
}

func (y loggerY) Info(args ...interface{}) {
	fmt.Println(append([]interface{}{"Info "}, args...)...)
}

func (y loggerY) Warn(args ...interface{}) {
	fmt.Println(append([]interface{}{"Warn "}, args...)...)
}

func (y loggerY) Error(args ...interface{}) {
	fmt.Println(append([]interface{}{"Error "}, args...)...)
}

func (y loggerY) Fatal(args ...interface{}) {
	fmt.Println(append([]interface{}{"Error "}, args...)...)
}
