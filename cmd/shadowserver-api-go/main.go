package main

import (
	"encoding/json"
	"flag"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
	"os"
	"path"
	"runtime"
	"shadowserver"
	"shadowserver/model"
	"strconv"
	"time"
)

var (
	method      string
	param       string
	pretty      bool
	reports     bool
	reportsCron bool
)

func init() {
	var exampleFormatter = &zt_formatter.ZtFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		Formatter: nested.Formatter{
			//HideKeys:    true,
			FieldsOrder: []string{"component", "category"},
		},
	}

	log.SetFormatter(exampleFormatter)
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	flag.StringVar(&method, "method", "test/ping", "Request URI")
	flag.StringVar(&param, "param", "{}", "JSON parameter")
	flag.BoolVar(&pretty, "pretty", true, "Pretty print JSON")
	flag.BoolVar(&reports, "reports", false, "Download reports")
	flag.BoolVar(&reportsCron, "reportsCron", false, "Schedule report download")
}

func main() {
	log.WithFields(log.Fields{
		"app": os.Args[0],
	}).Info("starting")

	// load .env environment variables
	err := godotenv.Load()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("failed to load .env file")
	}

	if os.Getenv("LOG_LEVEL") == "info" {
		log.SetLevel(log.InfoLevel)
	}

	flag.Parse()

	if reports {
		_, err = shadowserver.DownloadReports()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("failed to make API call")
		}
	} else if reportsCron {
		// Schedule report download
		scheduleEveryHours, err := strconv.Atoi(os.Getenv("REPORTS_DOWNLOAD_EVERY"))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("failed to get integer every hours value for scheduling")
		}
		s := gocron.NewScheduler(time.UTC)
		_, err = s.Every(scheduleEveryHours).Hours().Do(shadowserver.DownloadReports)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("failed to schedule report download")
		}
		s.StartBlocking()
	} else {
		p := make(model.ShadowserverParam)
		err = json.Unmarshal([]byte(param), &p)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("failed to unmarshal json param")
		}
		data, err := shadowserver.CallApi(method, p)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("failed to make API call")
		}
		shadowserver.PrintJson(data, pretty)
	}
}
