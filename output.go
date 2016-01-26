package airstrike

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fatih/color"
)

// A Reporter manages logging and console ouput.
type Reporter struct {
	CountGoroutines bool
	Gauge           bool

	// The Glyph is the character that will make up the horizontal bar gauge
	// if gauge output is enabled.
	Glyph byte

	// The Logger can be anything API-compatible with logrus.Logger.
	Logger *logrus.Logger

	// This channel receives types that fulfill the logrus.Fields interface.
	LogFields         chan map[string]interface{}
	MaxColumns        int
	URLInvariant      string
	WarningThreshold  time.Duration
	ThresholdReceiver chan time.Duration
}

// Run should be invoked in a goroutine. Log data fulfilling the logrus.Fields
// interface should be sent down its channel.
func (r *Reporter) Run(ch chan map[string]interface{}) {
	r.LogFields = ch
	for {
		select {
		case fields := <-r.LogFields:
			responseTime, _ := fields["response_time"].(time.Duration)
			r.writeLog(responseTime, fields)
			if r.Gauge {
				if r.MaxColumns == 0 {
					r.MaxColumns = 80
				}
				r.writeConsoleGauge(responseTime)
			}
		case r.WarningThreshold = <-r.ThresholdReceiver:
		default:
		}
	}
}

func (r Reporter) writeLog(responseTime time.Duration, fields map[string]interface{}) {
	desc := "airstrike.log"

	if fields["response_time"] != nil {
		if responseTime > r.WarningThreshold {
			r.Logger.WithFields(fields).Warn(desc)
		} else {
			r.Logger.WithFields(fields).Info(desc)
		}
	} else {
		switch fields["severity"] {
		case "INFO", "info":
			r.Logger.WithFields(fields).Info(desc)
		case "WARN", "warn":
			r.Logger.WithFields(fields).Warn(desc)
		case "ERROR", "error":
			r.Logger.WithFields(fields).Error(desc)
		default:
			if r.CountGoroutines {
				desc = fmt.Sprintf("%s (gr: %v)", desc, runtime.NumGoroutine())
			}
			r.Logger.WithFields(fields).Debug(desc)
		}
	}
}

func (r Reporter) makeBar(numBlocks int, responseTime time.Duration) string {
	var a []byte

	// allow chars for "nnnnms" text
	for i := 0; i < numBlocks-charsToSave(responseTime); i++ {
		if i < r.MaxColumns-charsToSave(responseTime) {
			a = append(a, r.Glyph)
		}
	}
	return string(a[:])
}

func (r Reporter) writeConsoleGauge(responseTime time.Duration) error {
	// 80 chars * 1 block per 10 ms = 800 ms max resolution
	maxRes := r.MaxColumns * 10
	blockMS := time.Duration(maxRes) / time.Duration(r.MaxColumns)
	numBlocks := int(responseTime / blockMS / time.Millisecond)
	if responseTime != time.Duration(0) {
		fmt.Println()
		r.writeConsoleLegend()
		bar := r.makeBar(numBlocks, responseTime)

		color.Set(color.FgBlue, color.Bold)
		defer color.Unset()

		if responseTime > r.WarningThreshold {
			color.Set(color.FgHiYellow, color.Bold)
		}
		if responseTime > time.Duration(maxRes)*time.Millisecond {
			color.Set(color.FgHiRed, color.Bold)
		}
		fmt.Printf("%s[%dms]", bar, responseTime/time.Millisecond)
	}
	return nil
}

func (r Reporter) writeConsoleLegend() {
	for tens := 0; tens <= r.MaxColumns/10; tens++ {
		for ones := 0; ones < 10; ones++ {
			warningMsec := r.WarningThreshold / time.Nanosecond
			gaugeMsec := time.Duration(((10*tens)+ones)*10) / 1 * time.Millisecond
			// fmt.Printf("warningMsec: %v, gaugeMsec: %v\n", warningMsec, gaugeMsec)
			if ones == 0 {
				// print the multiples-of-ten grid line
				num := fmt.Sprintf("%d", tens)
				if tens >= 10 {
					num = string(num[len(num)-1])
				}
				fmt.Print(num)
			} else if gaugeMsec > warningMsec && (10*tens)+ones <= r.MaxColumns {
				color.Set(color.FgRed, color.Faint)
				fmt.Print(`░`)
				color.Unset()
			} else {
				fmt.Print(" ")
			}
		}
	}
	fmt.Printf("\r")
}

func charsToSave(responseTime time.Duration) int {
	if responseTime/time.Millisecond >= 1000 {
		return 7
	}
	return 6
}
