package media

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tphakala/birdnet-go/internal/api/v2/apitest"
	"github.com/tphakala/birdnet-go/internal/conf"
	"github.com/tphakala/birdnet-go/internal/datastore/mocks"
	"github.com/tphakala/birdnet-go/internal/spectrogram"
)

func TestResolveDetectionFrequencyProfile(t *testing.T) {
	withRestoredGlobalSettings(t)

	tests := []struct {
		name         string
		modelType    string
		modelTypeErr error
		micRate      int
		wantResample int
		wantSuffix   string
	}{
		{
			name:         "bird model defaults to 24000 Hz with empty suffix",
			modelType:    "bird",
			micRate:      96000,
			wantResample: 24000,
			wantSuffix:   "",
		},
		{
			name:         "GetNoteModelType error falls back to bird profile",
			modelTypeErr: assert.AnError,
			micRate:      96000,
			wantResample: 24000,
			wantSuffix:   "",
		},
		{
			name:         "bat model with 96 kHz mic scales to 48 kHz Nyquist",
			modelType:    "bat",
			micRate:      96000,
			wantResample: 96000,
			wantSuffix:   "bat-48k",
		},
		{
			name:         "bat model with 192 kHz mic scales to 96 kHz Nyquist",
			modelType:    "bat",
			micRate:      192000,
			wantResample: 192000,
			wantSuffix:   "bat-96k",
		},
		{
			name:         "bat model with 256 kHz mic scales to 128 kHz Nyquist",
			modelType:    "bat",
			micRate:      256000,
			wantResample: 256000,
			wantSuffix:   "bat-128k",
		},
		{
			name:         "bat model with unconfigured mic defaults to 48 kHz (24 kHz Nyquist)",
			modelType:    "bat",
			micRate:      0,
			wantResample: 48000,
			wantSuffix:   "bat-24k",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := apitest.NewValidTestSettings()
			if tt.micRate > 0 {
				settings.Realtime.Audio.Sources = []conf.AudioSourceConfig{
					{Name: "mic", SampleRate: tt.micRate},
				}
			} else {
				settings.Realtime.Audio.Sources = nil
			}

			mockDS := new(mocks.MockInterface)
			if tt.modelTypeErr != nil {
				mockDS.On("GetNoteModelType", "42").Return("", tt.modelTypeErr)
			} else {
				mockDS.On("GetNoteModelType", "42").Return(tt.modelType, nil)
			}

			core := apitest.NewCore(t,
				apitest.WithDatastore(mockDS),
				apitest.WithSettings(settings),
			)
			handler := New(core)

			profile := handler.resolveDetectionFrequencyProfile("42")
			assert.Equal(t, tt.wantResample, profile.ResampleRate)
			assert.Equal(t, tt.wantSuffix, spectrogram.ProfileSuffix(profile))
			assert.Equal(t, tt.wantSuffix, handler.spectrogramProfileSuffix("42"))
		})
	}
}
