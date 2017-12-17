package logger

import (
	"log"
	"os"
)

var levels = map[string]int{
	"fatal": 0,
	"error": 1,
	"info":  2,
	"debug": 3,
}

type Logger struct {
	Level string
}

func New() *Logger {
	level := os.Getenv("LOGLEVEL")
	if _, ok := levels[level]; !ok {
		level = "info"
	}
	return &Logger{
		Level: level,
	}
}

func (logger *Logger) Info(args ...interface{}) {
	if levels[logger.Level] >= levels["info"] {
		newArgs := appendToArray([]interface{}{"info:"}, args)
		log.Println(newArgs...)
	}
}

func (logger *Logger) Debug(args ...interface{}) {
	if levels[logger.Level] >= levels["debug"] {
		newArgs := appendToArray([]interface{}{"debug:"}, args)
		log.Println(newArgs...)
	}
}

func (logger *Logger) Error(args ...interface{}) {
	if levels[logger.Level] >= levels["error"] {
		newArgs := appendToArray([]interface{}{"error:"}, args)
		log.Println(newArgs...)
	}
}

func (logger *Logger) Fatalln(args ...interface{}) {
	newArgs := appendToArray([]interface{}{"fatal:"}, args)
	log.Fatalln(newArgs...)
}

func appendToArray(a1 []interface{}, a2 []interface{}) []interface{} {
	for _, item := range a2 {
		a1 = append(a1, item)
	}
	return a1
}
