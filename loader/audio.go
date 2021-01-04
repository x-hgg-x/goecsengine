package loader

import (
	"fmt"
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
	f, err := os.Open(audioFilePath)
	utils.LogError(err)

	var d io.ReadSeeker
	switch filepath.Ext(audioFilePath) {
	case ".mp3":
		d, err = mp3.Decode(audioContext, f)
	case ".ogg":
		d, err = vorbis.Decode(audioContext, f)
	case ".wav":
		d, err = wav.Decode(audioContext, f)
	default:
		utils.LogError(fmt.Errorf("unknown audio file extension: '%s'", filepath.Ext(audioFilePath)))
	}
	utils.LogError(err)

	audioPlayer, err := audio.NewPlayer(audioContext, d)
	utils.LogError(err)
	return audioPlayer
}
