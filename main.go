package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/dustin/go-humanize"
	gf "github.com/jessevdk/go-flags"
	"github.com/periaate/common"
)

type Options struct {
	Bytes   bool `short:"b" long:"bytes" description:"Count number of bytes."`
	Runes   bool `short:"r" long:"runes" description:"Count number of UTF8 characters."`
	Lines   bool `short:"l" long:"lines" description:"Count number of lines."`
	Words   bool `short:"w" long:"words" description:"Count number of words."`
	Max     bool `short:"m" long:"max" description:"Max value found."`
	Human   bool `short:"q" long:"human" description:"Values as human readable."`
	Value   bool `short:"g" long:"get" description:"Get max value found."`
	Verbose bool `short:"v" long:"verbose" description:"Verbose."`
}

func main() {
	opts := &Options{}
	rest, err := gf.Parse(opts)
	if err != nil {
		if gf.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatalln("Error parsing flags:", err)
	}

	if !opts.Bytes && !opts.Runes && !opts.Lines && !opts.Words {
		opts.Lines = true
	}

	vals := append(rest, common.ReadPipe()...)
	if len(vals) == 0 {
		log.Fatalln("No input.")
	}
	switch {
	case opts.Max:
		var ft bool
		switch {
		case opts.Lines:
			cnt, v := max(vals, func(s string) int { return len([]rune(s)) })
			fmt.Print(val(v, "lines", cnt, opts))
		case opts.Words:
			inp := strings.Split(strings.Join(vals, " "), " ")
			cnt, v := max(inp, func(s string) int { return len([]rune(s)) })
			fmt.Print(val(v, "words", cnt, opts))
		default:
			ft = true
		}
		if !ft {
			break
		}
		fallthrough
	default:
		switch {
		case opts.Lines:
			cnt := len(vals)
			fmt.Print(val("", "lines", cnt, opts))
		case opts.Words:
			cnt := len(strings.Fields(strings.Join(vals, " ")))
			fmt.Print(val("", "words", cnt, opts))
		case opts.Bytes:
			cnt := len(strings.Join(vals, ""))
			fmt.Print(val("", "bytes", cnt, opts))
		case opts.Runes:
			cnt := utf8.RuneCountInString(strings.Join(vals, ""))
			fmt.Print(val("", "runes", cnt, opts))
		}
	}
}

func val(v, obj string, count int, opts *Options) (res string) {
	switch {
	case opts.Value && opts.Human && opts.Verbose:
		return fmt.Sprintf("Max value found:\t%s\nTotal %s count:\t%s\n", v, obj, humanize.Comma(int64(count)))
	case opts.Value && opts.Human:
		return fmt.Sprintf("Max value found:\t%s\n", v)
	case opts.Value && opts.Verbose:
		return fmt.Sprintf("%s\n%d", v, count)
	case opts.Human:
		return fmt.Sprintf("Total %s count:\t%s", obj, humanize.Comma(int64(count)))
	case opts.Verbose:
		fallthrough
	default:
		return fmt.Sprintf("%d", count)
	}
}

func max[T any](sar []T, fn func(T) int) (H int, t T) {
	for _, v := range sar {
		h := fn(v)
		if h > H {
			H = h
			t = v
		}
	}
	return
}
