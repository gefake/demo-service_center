package logger

import (
	"io"
	"log"
	logging "log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var (
	Log *logrus.Logger // share will all packages
)

func init() {
	err := os.Mkdir("logs/", os.ModeDir)

	f, err := os.OpenFile("logs/"+time.Now().Format("02-01-2006-15-04-05")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		// use go's logger, while we configure Logrus
		logging.Fatalf("error opening file: %v", err)
	}

	// configure Logrus
	Log = logrus.New()
	Log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%\n",
	})
	Log.SetReportCaller(true)

	mw := io.MultiWriter(os.Stdout, f)
	Log.SetOutput(mw)
}
