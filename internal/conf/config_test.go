package conf

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAudioSettings_NeedsFfprobeWorkaround(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		audio          AudioSettings
		wantWorkaround bool
	}{
		{
			name: "FFmpeg 5.x needs workaround",
			audio: AudioSettings{
				FfmpegVersion: "5.1.7-0+deb12u1+rpt1",
				FfmpegMajor:   5,
				FfmpegMinor:   1,
			},
			wantWorkaround: true,
		},
		{
			name: "FFmpeg 7.x does not need workaround",
			audio: AudioSettings{
				FfmpegVersion: "7.1.2-0+deb13u1",
				FfmpegMajor:   7,
				FfmpegMinor:   1,
			},
			wantWorkaround: false,
		},
		{
			name: "FFmpeg 6.x does not need workaround",
			audio: AudioSettings{
				FfmpegVersion: "6.0",
				FfmpegMajor:   6,
				FfmpegMinor:   0,
			},
			wantWorkaround: false,
		},
		{
			name: "FFmpeg 4.x does not need workaround",
			audio: AudioSettings{
				FfmpegVersion: "4.4.2",
				FfmpegMajor:   4,
				FfmpegMinor:   4,
			},
			wantWorkaround: false,
		},
		{
			name: "Unknown version does not need workaround",
			audio: AudioSettings{
				FfmpegVersion: "",
				FfmpegMajor:   0,
				FfmpegMinor:   0,
			},
			wantWorkaround: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.audio.NeedsFfprobeWorkaround()
			assert.Equal(t, tt.wantWorkaround, got,
				"AudioSettings.NeedsFfprobeWorkaround() mismatch")
		})
	}
}

func TestAudioSettings_HasFfmpegVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		audio   AudioSettings
		wantHas bool
	}{
		{
			name: "Valid version detected",
			audio: AudioSettings{
				FfmpegVersion: "7.1.2",
				FfmpegMajor:   7,
				FfmpegMinor:   1,
			},
			wantHas: true,
		},
		{
			name: "No version detected",
			audio: AudioSettings{
				FfmpegVersion: "",
				FfmpegMajor:   0,
				FfmpegMinor:   0,
			},
			wantHas: false,
		},
		{
			name: "Version string but no major version",
			audio: AudioSettings{
				FfmpegVersion: "unknown",
				FfmpegMajor:   0,
				FfmpegMinor:   0,
			},
			wantHas: false,
		},
		{
			name: "Major version but no version string",
			audio: AudioSettings{
				FfmpegVersion: "",
				FfmpegMajor:   7,
				FfmpegMinor:   1,
			},
			wantHas: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.audio.HasFfmpegVersion()
			assert.Equal(t, tt.wantHas, got,
				"AudioSettings.HasFfmpegVersion() mismatch")
		})
	}
}

func TestBackwardCompatibility_RemovedBatFilterFields(t *testing.T) {
	t.Parallel()

	// YAML content containing the removed bat high-pass filter fields.
	yamlContent := `
bat:
  threshold: 0.75
  nighttimeonly: true
  filterenabled: true
  filtercutoffhz: 5000.0
  filterpasscount: 2
`

	v := viper.New()
	v.SetConfigType("yaml")

	err := v.ReadConfig(strings.NewReader(yamlContent))
	require.NoError(t, err)

	var settings Settings
	err = v.Unmarshal(&settings, viper.DecodeHook(DurationDecodeHook()))
	require.NoError(t, err)

	// Verify that the parser successfully ignored the removed filter fields,
	// but correctly loaded the other active BatConfig fields.
	assert.InDelta(t, 0.75, settings.Bat.Threshold, 1e-9)
	assert.True(t, settings.Bat.NighttimeOnly)
}

func TestAudioSettings_SampleRateHelpers(t *testing.T) {
	t.Parallel()

	t.Run("nil AudioSettings returns default 48000", func(t *testing.T) {
		var a *AudioSettings
		assert.Equal(t, SampleRate, a.GetPrimarySampleRate())
		assert.Equal(t, SampleRate, a.GetSourceSampleRate("any"))
		assert.Equal(t, SampleRate, a.DefaultBatSampleRate())
	})

	t.Run("empty sources returns default 48000", func(t *testing.T) {
		a := &AudioSettings{Sources: []AudioSourceConfig{}}
		assert.Equal(t, SampleRate, a.GetPrimarySampleRate())
		assert.Equal(t, SampleRate, a.GetSourceSampleRate("any"))
		assert.Equal(t, SampleRate, a.DefaultBatSampleRate())
	})

	t.Run("source with 0 sample rate returns default 48000", func(t *testing.T) {
		a := &AudioSettings{
			Sources: []AudioSourceConfig{
				{Name: "mic1", SampleRate: 0},
			},
		}
		assert.Equal(t, SampleRate, a.GetPrimarySampleRate())
		assert.Equal(t, SampleRate, a.GetSourceSampleRate("mic1"))
		assert.Equal(t, SampleRate, a.DefaultBatSampleRate())
	})

	t.Run("primary source with configured rate", func(t *testing.T) {
		a := &AudioSettings{
			Sources: []AudioSourceConfig{
				{Name: "mic1", SampleRate: 96000},
			},
		}
		assert.Equal(t, 96000, a.GetPrimarySampleRate())
		assert.Equal(t, 96000, a.GetSourceSampleRate("mic1"))
		assert.Equal(t, 96000, a.DefaultBatSampleRate())
	})

	t.Run("multiple sources with specific bat source", func(t *testing.T) {
		a := &AudioSettings{
			Sources: []AudioSourceConfig{
				{Name: "bird_mic", SampleRate: 48000, Model: "birdnet"},
				{Name: "bat_mic", SampleRate: 192000, Models: []string{"bat"}},
			},
		}
		assert.Equal(t, 48000, a.GetPrimarySampleRate())
		assert.Equal(t, 48000, a.GetSourceSampleRate("bird_mic"))
		assert.Equal(t, 192000, a.GetSourceSampleRate("bat_mic"))
		assert.Equal(t, 192000, a.DefaultBatSampleRate())
	})

	t.Run("multiple sources with bat model field", func(t *testing.T) {
		a := &AudioSettings{
			Sources: []AudioSourceConfig{
				{Name: "bird_mic", SampleRate: 48000, Model: "birdnet"},
				{Name: "bat_mic", SampleRate: 256000, Model: "bat"},
			},
		}
		assert.Equal(t, 256000, a.DefaultBatSampleRate())
	})
}

func TestSettings_SampleRateHelpers(t *testing.T) {
	t.Parallel()

	t.Run("nil Settings returns default 48000", func(t *testing.T) {
		var s *Settings
		assert.Equal(t, SampleRate, s.GetPrimaryAudioSourceSampleRate())
		assert.Equal(t, SampleRate, s.GetAudioSourceSampleRate("any"))
		assert.Equal(t, SampleRate, s.DefaultBatSampleRate())
	})

	t.Run("Settings delegates to Realtime.Audio", func(t *testing.T) {
		s := &Settings{
			Realtime: RealtimeSettings{
				Audio: AudioSettings{
					Sources: []AudioSourceConfig{
						{Name: "ultrasonic", SampleRate: 384000, Model: "bat"},
					},
				},
			},
		}
		assert.Equal(t, 384000, s.GetPrimaryAudioSourceSampleRate())
		assert.Equal(t, 384000, s.GetAudioSourceSampleRate("ultrasonic"))
		assert.Equal(t, 384000, s.DefaultBatSampleRate())
	})
}
