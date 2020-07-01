package utils
//
//import (
//	"github.com/kataras/golog"
//	"runtime"
//	"path"
//	"fmt"
//	"sync"
//)
//
////var Log1 = golog.New()
//const DebugLevel = golog.DebugLevel
//
//
//type Singleton struct {
//
//}
//
////var instance *Singleton
//var Logger *golog.Logger
//var once sync.Once
//
////func GetInstance() *Singleton {
////	once.Do(func() {
////		instance = &Singleton{}
////	})
////	return instance
////}
//
//
//func NEW() *golog.Logger {
//	once.Do(func() {
//		Logger = golog.New()
//		Logger.TimeFormat = "2006/01/02 15:04"
//		Logger.Level = golog.WarnLevel
//
//		Logger.Handle(func(l *golog.Log) bool {
//			prefix := golog.GetTextForLevel(l.Level, false)
//			_, fn, line, _ := runtime.Caller(5)
//			fn = path.Base(fn)
//			message := fmt.Sprintf("%s %s %s:%d %s", prefix, "03/01/2006 15:04", fn, line, l.Message)
//
//			if l.NewLine {
//				message += "\n"
//			}
//
//			fmt.Print(message)
//			return true
//		})
//		//	golog.Default.Level = golog.WarnLevel
//
//		//func New() *Logger {
//		//	return &Logger{
//		//	Level:      golog.InfoLevel,
//		//	TimeFormat: "2006/01/02 15:04",
//		//	Printer:    pio.NewPrinter("", os.Stdout).EnableDirectOutput().Hijack(logHijacker),
//		//	children:   newLoggerMap(),
//		//}
//		//}
//	})
//	return Logger
//}
//
//
//
//func init() {
//	golog.Default.Level = golog.WarnLevel
//	golog.Default.SetTimeFormat("03/01/2006 15:04")
//	//golog.SetTimeFormat("03/01/2006 15:04")
//
//
//	fmt.Printf("%p\n", golog.Default)
//
//	golog.Handle(func(l *golog.Log) bool {
//		prefix := golog.GetTextForLevel(l.Level, false)
//		_, fn, line, _ := runtime.Caller(6)
//		fn = path.Base(fn)
//		message := fmt.Sprintf("%s %s %s:%d %s", prefix, "03/01/2006 15:04", fn, line, l.Message)
//
//		if l.NewLine {
//			message += "\n"
//		}
//
//		fmt.Print(message)
//		return true
//	})
//
//
//	//Log1.Handle(func(l *golog.Log) bool {
//	//	//prefix := golog.GetTextForLevel(l.Level., false)
//	//	//_, fn, line, _ := runtime.Caller(7)
//	//	message := fmt.Sprintf("%s", l.Message)
//	//
//	//	if l.NewLine {
//	//		message += "\n"
//	//	}
//	//
//	//	fmt.Print(message)
//	//	return true
//	//})
//}
//
