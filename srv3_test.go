package astisub_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/asticode/go-astisub"
	"github.com/stretchr/testify/assert"
)

func TestSRV3(t *testing.T) {
	// Open
	s, err := astisub.OpenFile("./testdata/example-in.srv3")
	assert.NoError(t, err)
	assertSubtitleItems(t, s)

	// No subtitles to write
	w := &bytes.Buffer{}
	err = astisub.Subtitles{}.WriteToSRV3(w)
	assert.EqualError(t, err, astisub.ErrNoSubtitlesToWrite.Error())

	// Write
	c, err := ioutil.ReadFile("./testdata/example-out.srv3")
	assert.NoError(t, err)
	err = s.WriteToSRV3(w)
	assert.NoError(t, err)
	assert.Equal(t, string(c), w.String())
}
