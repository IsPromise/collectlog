package logs

import (
	"io"
	"log"
	"os"
)

// 这个包是系统的日志记录的相关配置

// 日志级别
var (
	Trace   *log.Logger // 几乎任何东西
	Info    *log.Logger //	重要信息
	Warning *log.Logger //	警告
	Error   *log.Logger //	错误
)

// 不同级别日志分别记录的位置
var (
	traceFile   *os.File
	infoFile    *os.File
	warningFile *os.File
	errorFile   *os.File
)

var err error

func init() {
	// 初始化文件
	traceFile, err = os.OpenFile("./logFile/trace_files.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("init trace file err:", err)
	}
	infoFile, err = os.OpenFile("./logFile/info_files.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("init info file err:", err)
	}
	warningFile, err = os.OpenFile("./logFile/warning_files.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("init warning file err:", err)
	}
	errorFile, err = os.OpenFile("./logFile/error_files.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("init error file err:", err)
	}

	// 日志记录的位置和格式设置
	Trace = log.New(traceFile,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(io.MultiWriter(infoFile, os.Stdout, traceFile),
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(io.MultiWriter(warningFile, os.Stdout, traceFile),
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(errorFile, os.Stdout, traceFile),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func Close() {
	_ = traceFile.Close()
	_ = infoFile.Close()
	_ = warningFile.Close()
	_ = errorFile.Close()
}
