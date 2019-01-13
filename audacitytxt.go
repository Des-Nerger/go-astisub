package astisub

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func ReadFromAudacityTXT(r io.Reader) (subs *Subtitles, err error) {
	subs = NewSubtitles()
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), "\t", 3)

		item := &Item{}
		if item.StartAt, err = time.ParseDuration(fields[0] + "s"); err != nil {
			err = errors.Wrapf(err, "astisub: parsing audacitytxt duration %s failed", fields[0])
			return
		}
		if item.EndAt, err = time.ParseDuration(fields[1] + "s"); err != nil {
			err = errors.Wrapf(err, "astisub: parsing audacitytxt duration %s failed", fields[1])
			return
		}
		for _, dialogLine := range strings.Split(fields[2], `\n`) {
			item.Lines = append(item.Lines, Line{Items: []LineItem{{Text: dialogLine}}})
		}
		subs.Items = append(subs.Items, item)
	}
	return
}


func (subs Subtitles) WriteToAudacityTXT(w io.Writer) (err error) {
	// Do not write anything if no subtitles
	if len(subs.Items) == 0 {
		err = ErrNoSubtitlesToWrite
		return
	}

	var bs []byte
	for _, item := range subs.Items {
		bs = strconv.AppendFloat(bs, item.StartAt.Seconds(), 'f', -1, 64)
		bs = append(bs, '\t')
		bs = strconv.AppendFloat(bs, item.EndAt.Seconds(), 'f', -1, 64)
		bs = append(bs, '\t')
		var linesStrings []string
		for _, line := range item.Lines {
			linesStrings = append(linesStrings, line.String())
		}
		bs = append(bs, strings.Join(linesStrings, `\n`)...)
		bs = append(bs, bytesLineSeparator...)
	}
	// Remove last new line
	bs = bs[:len(bs)-len(bytesLineSeparator)]

	// Write
	if _, err = w.Write(bs); err != nil {
		err = errors.Wrap(err, "astisub: writing failed")
		return
	}
	return
}
