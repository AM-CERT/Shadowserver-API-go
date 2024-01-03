package shadowserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AM-CERT/Shadowserver-API-go/model"
	"github.com/go-resty/resty/v2"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*

The API is mostly RESTful with the following stipulations:

API calls must be made with HTTP POST.
All requests are JSON objects.
Each request requires an API key and HMAC header.
Any non-200 HTTP response is an error. Error details will be included in the response object.

https://github.com/The-Shadowserver-Foundation/api_utils/wiki/call_api-Documentation

*/

var (
	shadowserverSecret string
	shadowserverApiKey string
	reportsBaseDir     string
	reportsMinDiskFree string
	client             *resty.Client
)

func init() {
	shadowserverSecret = os.Getenv("SHADOWSERVER_SECRET")
	shadowserverApiKey = os.Getenv("SHADOWSERVER_API_KEY")
	reportsBaseDir = os.Getenv("REPORTS_BASEDIR")
	reportsMinDiskFree = os.Getenv("REPORT_MIN_DISK_FREE")
	client = resty.New()
	client.SetHeader("Accept", "application/json")
	client.SetBaseURL(os.Getenv("SHADOWSERVER_URI"))
}

func DownloadReports() ([]*model.ShadowserverReport, error) {
	// check disk usage
	reportsMinDiskFreeInt, err := strconv.ParseUint(reportsMinDiskFree, 10, 64)
	if err != nil {
		return nil, err
	}
	disk, err := DiskUsage(reportsBaseDir)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"Total MB": float64(disk.All) / float64(MB),
		"Used  MB": float64(disk.Used) / float64(MB),
		"Free  MB": float64(disk.Free) / float64(MB),
	}).Info("Disk space check")

	if disk.Free < reportsMinDiskFreeInt {
		return nil, errors.New("insufficient disk space")
	}

	reportDays, err := strconv.Atoi(os.Getenv("REPORT_DAYS"))
	if err != nil {
		return nil, err
	}
	// store actual downloaded new reports in this list
	reportListFinal := make([]*model.ShadowserverReport, 0)

	// get reports for today minus reportDays
	end := time.Now()
	start := end.AddDate(0, 0, -reportDays)
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		dayStr := d.Format("2006-01-02")

		// make tree folders as needed
		dir := reportsBaseDir
		for _, d := range strings.Split(dayStr, "-") {
			dir = filepath.Join(dir, d)
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}

		param := make(model.ShadowserverParam)
		param["date"] = dayStr
		reportList, err := GetReportList(param)
		if err != nil {
			return nil, err
		}

		log.WithFields(log.Fields{
			"day": dayStr,
			"dir": dir,
		}).Info("downloading reports")

		for _, report := range reportList {
			reportFilePath := filepath.Join(dir, report.File)

			if FileExists(reportFilePath) {
				log.WithFields(log.Fields{
					"reportFile": reportFilePath,
				}).Info("report already downloaded")
				continue
			}

			err = DownloadReport(report.Id, reportFilePath)
			if err != nil {
				log.WithFields(log.Fields{
					"err":    err,
					"report": report,
				}).Error("failed to download report")
				continue
			}

			report.FilePath = reportFilePath
			reportListFinal = append(reportListFinal, report)

			log.WithFields(log.Fields{
				"reportFile": report.File,
			}).Info("successfully downloaded report")
		}
	}

	return reportListFinal, nil
}

func CallApi(method string, param model.ShadowserverParam) ([]byte, error) {
	// add API key to the request
	param["apikey"] = shadowserverApiKey

	// convert to json string
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	hmacData := ComputeHmac(shadowserverSecret, paramBytes)

	resp, err := client.R().
		SetHeader("HMAC2", hmacData).
		SetBody(paramBytes).
		Post(method)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return resp.Body(), errors.New(resp.Status())
	}

	return resp.Body(), nil
}

func GetReportList(param model.ShadowserverParam) ([]*model.ShadowserverReport, error) {
	data, err := CallApi("reports/list", param)
	if err != nil {
		return nil, err
	}

	var reportList []*model.ShadowserverReport
	if err := json.Unmarshal(data, &reportList); err != nil {
		panic(err)
	}

	return reportList, nil
}

func DownloadReport(id string, path string) error {
	resp, err := client.R().
		SetOutput(path).
		Get(fmt.Sprintf("https://dl.shadowserver.org/%s", id))
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New(resp.Status())
	}

	return nil
}
