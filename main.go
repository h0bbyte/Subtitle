package main

import (
	"flag"
	"log"
	"os"
	"io/ioutil"
	"strings"
	"strconv"
	"fmt"
	"math"
)

func main() {
	initFlags()
	output(parse(read()))
}

func read() string {
	//flagPath := "../data/1.srt"
	f, err := os.Open(flagPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err.Error())
	}

	return string(b)
}

type Subtitle struct {
	Index uint64
	Begin int32
	End int32
	Content string
}

func parse(str string) ([]Subtitle) {
	result := []Subtitle{Subtitle{}}
	idx := 0
	slice := strings.Split(str, "\r\n")
	for _, v := range slice {
		if len(v) == 0 {
			result = append(result, Subtitle{})
			idx += 1
			continue
		} else {
			_, err := strconv.Atoi(v)
			if err == nil {
					// Index
					result[idx].Index, _ = strconv.ParseUint(v, 10, 0)
			} else {
				if strings.Index(v, ":") != -1 {
					// Time
					var bh, bm, bs, bms, eh, em, es, ems int32
					_, err = fmt.Sscanf(v, "%02d:%02d:%02d,%03d --> %02d:%02d:%02d,%03d", &bh, &bm, &bs, &bms, &eh, &em, &es, &ems)
					if err != nil {
						log.Fatal(err)
					}
					result[idx].Begin = bh * 3600000 + bm * 60000 + bs * 1000 + bms
					result[idx].End = eh * 3600000 + em * 60000 + es * 1000 + ems
					if result[idx].Begin >= begin && result[idx].End <= end {
						result[idx].Begin += flagDelta
						result[idx].End += flagDelta
					}
				} else {
					// Content
					result[idx].Content = v
				}
			}
		}
	}
	return result
}

func output(slice []Subtitle) {
	var format = func(val int32) string {
		h := val / 3600000
		m := (val - (h * 3600000)) / 60000
		s := (val - h * 3600000 - m * 60000) / 1000
		ms := val - h * 3600000 - m * 60000 - s * 1000
		return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, s, ms)
	}

	for _, v := range slice {
		fmt.Printf("%d\n%s --> %s\n%s\n\n", v.Index, format(v.Begin), format(v.End), v.Content)
	}
}

// Flags
var flagPath string
var flagBegin string
var flagEnd string
var flagDelta int32
var begin, end int32
func initFlags() {
	flag.StringVar(&flagPath, "path", "", "The subtitle file path.")
	flag.StringVar(&flagBegin, "begin", "", "The begin of duration, format hh:mm:ss,ms.")
	flag.StringVar(&flagEnd, "end", "", "The end of duration, format hh:mm:ss,ms.")
	var delta int
	flag.IntVar(&delta, "delta", 0, "The shift value in milliseconds.")
	flag.Parse()

	flagDelta = int32(delta)
	var h, m, s, ms int32
	_, err := fmt.Sscanf(flagBegin, "%02d:%02d:%02d,%03d", &h, &m, &s, &ms)
	if err != nil {
		log.Fatal(err)
	}
	begin = int32(h * 3600000 + m * 60000 + s * 1000 + ms)

	// 99 * 3600000 + 59 * 60000 + 59 * 1000 + 999 = 359941999, max end second.
	end = math.MaxInt32
	if len(flagEnd) > 0 {
		h = 0
		m = 0
		s = 0
		ms = 0
		_, err = fmt.Sscanf(flagEnd, "%02d:%02d:%02d,%03d", &h, &m, &s, &ms)
		if err != nil {
			log.Fatal(err)
		}
		end = int32(h * 3600000 + m * 60000 + s * 1000 + ms)
	}

	fmt.Printf("begin: %d, end: %d, delta: %d", begin, end, flagDelta)
}
