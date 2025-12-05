// File: pkg/utils/logger.go
package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// Nivel de log
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// Logger es un logger simple y concurrente
type Logger struct {
	mu     sync.Mutex
	out    io.Writer
	level  Level
	logger *log.Logger
}

// paquete-global (fácil uso)
var std = New(os.Stdout, INFO)

// New crea un nuevo Logger apuntando a out con un nivel inicial
func New(out io.Writer, level Level) *Logger {
	l := &Logger{
		out:   out,
		level: level,
	}
	l.logger = log.New(out, "", 0) // manejamos formato en este paquete
	return l
}

// Init inicializa el logger por defecto del paquete. Llamar early en main.
func Init(level Level, out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = out
	std.level = level
	std.logger = log.New(out, "", 0)
}

// SetOutput permite cambiar el destino (por ejemplo, un archivo)
func SetOutput(out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = out
	std.logger = log.New(out, "", 0)
}

// SetLevel cambia el nivel del logger (DEBUG, INFO, ...)
func SetLevel(level Level) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.level = level
}

// helper para formatear prefijo
func prefix(level Level) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	var lvl string
	switch level {
	case DEBUG:
		lvl = "DEBUG"
	case INFO:
		lvl = "INFO "
	case WARN:
		lvl = "WARN "
	case ERROR:
		lvl = "ERROR"
	default:
		lvl = "INFO "
	}
	return fmt.Sprintf("%s [%s] ", t, lvl)
}

// Métodos de instancia (útiles si quieres instancias separadas)
func (l *Logger) logf(level Level, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if level < l.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	l.logger.Output(3, prefix(level)+msg) // depth 3 para que apunte al llamador
}

func (l *Logger) Debugf(format string, v ...interface{}) { l.logf(DEBUG, format, v...) }
func (l *Logger) Infof(format string, v ...interface{})  { l.logf(INFO, format, v...) }
func (l *Logger) Warnf(format string, v ...interface{})  { l.logf(WARN, format, v...) }
func (l *Logger) Errorf(format string, v ...interface{}) { l.logf(ERROR, format, v...) }

// Métodos del logger por defecto (paquete-global)
func Debugf(format string, v ...interface{}) { std.Debugf(format, v...) }
func Infof(format string, v ...interface{})  { std.Infof(format, v...) }
func Warnf(format string, v ...interface{})  { std.Warnf(format, v...) }
func Errorf(format string, v ...interface{}) { std.Errorf(format, v...) }

// Utility: nivel desde string
func LevelFromString(s string) Level {
	switch s {
	case "debug", "DEBUG":
		return DEBUG
	case "info", "INFO":
		return INFO
	case "warn", "WARN", "warning", "WARNING":
		return WARN
	case "error", "ERROR":
		return ERROR
	default:
		return INFO
	}
}
