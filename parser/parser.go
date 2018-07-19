package parser

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/adwpc/logmetrics/conf"
	"github.com/adwpc/logmetrics/metrics"
	"github.com/adwpc/logmetrics/zlog"
	"github.com/buger/jsonparser"
	"github.com/hpcloud/tail"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	METRIC_GAUGE     = "gauge"
	METRIC_COUNTER   = "counter"
	METRIC_HISTOGRAM = "histogram"
)

var (
	log = zlog.Log
)

func Monitor(c *conf.Config) {
	for _, v := range c.Logs {
		go Run(v)
	}
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal().Msg(http.ListenAndServe(c.Listen, nil).Error())
}

type LogJson struct {
	Type     string  //METRIC_GAUGE METRIC_COUNTER METRIC_HISTOGRAM
	ValKey   string  //"model_val":"floatval"  "webserver1_httpok":"1"  ValKey=webserver1_httpok
	ValValue float64 //"model_val":"floatval"  "webserver1_httpok":"1"  ValValue=1
	Alert    string  //"alert":"xxx"
}

func (j *LogJson) GetKV(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
	fmt.Printf("---------Key: '%s'------- Value: '%s'--------- Type: %s\n", string(key), string(value), dataType)
	//{"interface_alart":"1", "type":"counter", "alert":"l-cltvprocess-work1.vps.dev.ten.dm offline"}
	switch string(key) {
	case "type":
		j.Type = string(value)
	case "alert":
		j.Alert = string(value)
	default:
		if f, err := strconv.ParseFloat(string(value), 64); err != nil {
			log.Error().Msg(err.Error())
		} else {
			j.ValKey = string(key)
			j.ValValue = f
		}
	}
	return nil

}

func Run(l conf.Log) error {

	if l.Path == "" {
		log.Error().Msg("path == \"\"")
		return errors.New("path == \"\"")
	}

	tails, err := tail.TailFile(l.Path, tail.Config{
		ReOpen: true, Follow: true,
		MustExist: false,
		Poll:      true,
	})

	if err != nil {
		log.Error().Msgf("tail file err:", err)
		return err
	}

	for true {
		msg, ok := <-tails.Lines
		if !ok {
			log.Error().Msgf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		r, _ := regexp.Compile("{.*?}")
		jsons := r.FindAllString(msg.Text, -1)
		for i := 0; i < len(jsons); i++ {
			// log.Info().Msg(jsons[i])
			var val string
			var err error
			if val, err = jsonparser.GetString([]byte(jsons[i]), "type"); err != nil {
				log.Error().Msg("get type error  json=" + jsons[i])
				continue
			}
			if val != METRIC_GAUGE && val != METRIC_COUNTER && val != METRIC_HISTOGRAM {
				log.Error().Msg("type is invalid" + jsons[i])
				continue
			}
			var j LogJson
			if err = jsonparser.ObjectEach([]byte(jsons[i]), j.GetKV); err != nil {
				log.Error().Msg("jsonparser.ObjectEach failed : " + err.Error() + "   " + jsons[i])
			}
			switch j.Type {
			case METRIC_COUNTER:
				log.Debug().Msg(j.ValKey)
				metrics.Get(j.ValKey, j.Alert).Deal(j.ValValue)
			}
		}
	}

	return nil
}