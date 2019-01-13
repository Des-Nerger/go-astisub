package astisub_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/asticode/go-astisub"
	"github.com/stretchr/testify/assert"
)

func TestAudacityTXT(t *testing.T) {
	// Open
	s, err := astisub.OpenFile("./testdata/example-in.txt")
	assert.NoError(t, err)
	assertSubtitleItems(t, s)

	// No subtitles to write
	w := &bytes.Buffer{}
	err = astisub.Subtitles{}.WriteToAudacityTXT(w)
	assert.EqualError(t, err, astisub.ErrNoSubtitlesToWrite.Error())

	// Write
	c, err := ioutil.ReadFile("./testdata/example-out.txt")
	assert.NoError(t, err)
	err = s.WriteToAudacityTXT(w)
	assert.NoError(t, err)
	assert.Equal(t, string(c), w.String())
}
