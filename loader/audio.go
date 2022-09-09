package loader

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/x-hgg-x/goecsengine/utils"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// InitAudio creates a new audio context
func InitAudio(sampleRate int) *audio.Context {
	audioContext := audio.NewContext(sampleRate)
	return audioContext
}

// LoadAudio loads an audio file and returns an audio player
func LoadAudio(audioContext *audio.Context, audioFilePath string) *audio.Player {
	f := bytes.NewReader(utils.Try(os.ReadFile(audioFilePath)))

	var d io.ReadSeeker
	switch filepath.Ext(audioFilePath) {
	case ".mp3":
		d = utils.Try(mp3.DecodeWithSampleRate(audioContext.SampleRate(), f))
	case ".ogg":
		d = utils.Try(vorbis.DecodeWithSampleRate(audioContext.SampleRate(), f))
	case ".wav":
		d = utils.Try(wav.DecodeWithSampleRate(audioContext.SampleRate(), f))
	default:
		utils.LogFatalf("unknown audio file extension: '%s'", filepath.Ext(audioFilePath))
	}

	return utils.Try(audioContext.NewPlayer(d))
}
