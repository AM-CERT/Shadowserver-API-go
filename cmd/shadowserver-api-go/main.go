package main

import (
	"encoding/json"
	"flag"
	"github.com/AM-CERT/Shadowserver-API-go"
	"github.com/AM-CERT/Shadowserver-API-go/internal"
	"github.com/AM-CERT/Shadowserver-API-go/model"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"os"
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
	flag.StringVar(&method, "method", "test/ping", "Request URI")
	flag.StringVar(&param, "param", "{}", "JSON parameter")
	flag.BoolVar(&pretty, "pretty", true, "Pretty print JSON")
	flag.BoolVar(&reports, "reports", false, "Download reports")
	flag.BoolVar(&reportsCron, "reportsCron", false, "Schedule report download")
}

func main() {
	internal.InitLogger()

	internal.Logger.Info().
		Str("app", os.Args[0]).
		Msg("starting")

	// load .env environment variables
	err := godotenv.Load()
	if err != nil {
		internal.Logger.Fatal().
			Err(err).
			Msg("failed to load .env file")
	}

	if os.Getenv("DEBUG") == "true" {
		internal.SetDebug()
	}

	flag.Parse()

	if reports {
		_, err = shadowserver.DownloadReports()
		if err != nil {
			internal.Logger.Fatal().
				Err(err).
				Msg("failed to make API call")
		}
	} else if reportsCron {
		// Schedule report download
		scheduleEveryHours, err := strconv.Atoi(os.Getenv("REPORTS_DOWNLOAD_EVERY"))
		if err != nil {
			internal.Logger.Fatal().
				Err(err).
				Msg("failed to get integer every hours value for scheduling")
		}
		s := gocron.NewScheduler(time.UTC)
		_, err = s.Every(scheduleEveryHours).Hours().Do(shadowserver.DownloadReports)
		if err != nil {
			internal.Logger.Fatal().
				Err(err).
				Msg("failed to schedule report download")
		}
		s.StartBlocking()
	} else {
		p := make(model.ShadowserverParam)
		err = json.Unmarshal([]byte(param), &p)
		if err != nil {
			internal.Logger.Fatal().
				Err(err).
				Msg("failed to unmarshal json param")
		}
		data, err := shadowserver.CallApi(method, p)
		if err != nil {
			internal.Logger.Fatal().
				Err(err).
				Msg("failed to make API call")
		}
		shadowserver.PrintJson(data, pretty)
	}
}
