package airstrike

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fatih/color"
)

// A Reporter manages logging and console ouput.
type Reporter struct {
	// Display an ANSI-colorized graph of response times on STDOUT.
	Gauge bool

	// The Glyph is the character that will make up the horizontal bar gauge
	// if gauge output is enabled.
	Glyph byte

	// A reference to a logger.
	Logger *logrus.Logger

	// This channel receives types that fulfill the logrus.Fields interface.
	LogCh chan map[string]interface{}

	// Number of columns the gauge will occupy.
	GaugeWidth int

	// A string to omit from URLs in order to shorten log messages, i.e.,
	// the API's base URL.
	URLInvariant string

	// Response times over this threshold will be logged at the WARN level.
	WarningThreshold time.Duration

	// Values received on this channel will become the new WarningThreshold.
	ThresholdCh chan time.Duration

	// This channel holds the response time of the last result.
	LastResponseTimeCh chan time.Duration
}

func NewReporter(logger *logrus.Logger) *Reporter {
	return &Reporter{
		Logger:             logger,
		LogCh:              make(chan map[string]interface{}),
		ThresholdCh:        make(chan time.Duration, 1),
		LastResponseTimeCh: make(chan time.Duration, 1),
	}
}

// Run should be invoked in a goroutine. Log data fulfilling the
// logrus.Fields interface should be sent down its channel.
func (r *Reporter) Listen() {
	for {
		select {
		case fields := <-r.LogCh:
			responseTime, ok := fields["response_time"].(time.Duration)
			if ok {
				// fmt.Println("reporting the response time")
				select {
				case r.LastResponseTimeCh <- responseTime:
					fmt.Println("successfully reported the response time")
				default:
				}
				r.writeLog(responseTime, fields)
			}
			if r.Gauge {
				if r.GaugeWidth == 0 {
					r.GaugeWidth = 80
				}
				r.writeConsoleGauge(responseTime)
			}
		case r.WarningThreshold = <-r.ThresholdCh:
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
		switch strings.ToUpper(fields["severity"].(string)) {
		case "INFO":
			r.Logger.WithFields(fields).Info(desc)
		case "WARN":
			r.Logger.WithFields(fields).Warn(desc)
		case "ERROR":
			r.Logger.WithFields(fields).Error(desc)
		default:
			r.Logger.WithFields(fields).Debug(desc)
		}
	}
}

func (r Reporter) makeBar(numBlocks int, responseTime time.Duration) string {
	var a []byte

	// allow chars for "nnnnms" text
	for i := 0; i < numBlocks-charsToSave(responseTime); i++ {
		if i < r.GaugeWidth-charsToSave(responseTime) {
			a = append(a, r.Glyph)
		}
	}
	return string(a[:])
}

func (r Reporter) writeConsoleGauge(responseTime time.Duration) error {
	// 80 chars * 1 block per 10 ms = 800 ms max resolution
	maxRes := r.GaugeWidth * 10
	blockMS := time.Duration(maxRes) / time.Duration(r.GaugeWidth)
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
	for tens := 0; tens <= r.GaugeWidth/10; tens++ {
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
			} else if gaugeMsec > warningMsec && (10*tens)+ones <= r.GaugeWidth {
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
