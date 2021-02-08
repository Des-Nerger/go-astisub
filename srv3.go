package astisub

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func ReadFromSRV3(r io.Reader) (*Subtitles, error) {
	subs, d, it := NewSubtitles(), xml.NewDecoder(r), (*Item)(nil)
	for {
		t, err := d.RawToken(); switch err{case nil:; case io.EOF:return subs,nil; default:return subs,err}
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "p":
				it=&Item{}; subs.Items=append(subs.Items, it)
				for _, a := range t.Attr {
					switch a.Name.Local {case "t", "d": default: continue}
					i, err := strconv.ParseInt(a.Value, 10, 64); if err!=nil{return subs,err}
					*(func() *time.Duration {
						if a.Name.Local=="t" {return &it.StartAt}
						return &it.EndAt
					} ()) = time.Duration(i) * time.Millisecond
				}
				it.EndAt += it.StartAt
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "p": it=nil
			}
		case xml.CharData:
			if it != nil {
				for _, lineString := range strings.Split(string(t), "\n") {
					it.Lines = append(it.Lines, Line{Items: []LineItem{{Text: lineString}}})
				}
			}
		}
	}
}

func (subs Subtitles) WriteToSRV3(w io.Writer) (err error) {{
	if len(subs.Items) == 0 {return ErrNoSubtitlesToWrite}
	w := bufio.NewWriter(w); defer w.Flush()
	w.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\" ?><timedtext format=\"3\">\n<body>\n")
	for _, it := range subs.Items {
		startMs := it.StartAt.Milliseconds()
		fmt.Fprintf(w, `<p t="%v" d="%v">`, startMs, it.EndAt.Milliseconds()-startMs)
		for i, line := range it.Lines {
			w.WriteString(line.String())
			if i != len(it.Lines)-1 {w.WriteByte('\n')}
		}
		w.WriteString("</p>\n")
	}
	w.WriteString("</body>\n</timedtext>\n")
	return nil
}}
