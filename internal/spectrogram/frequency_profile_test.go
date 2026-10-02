package spectrogram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatProfileWithRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		sampleRate   int
		wantResample int
		wantSuffix   string
	}{
		{"48 kHz capture (24 kHz Nyquist)", 48000, 48000, "bat-24k"},
		{"96 kHz capture (48 kHz Nyquist)", 96000, 96000, "bat-48k"},
		{"192 kHz capture (96 kHz Nyquist)", 192000, 192000, "bat-96k"},
		{"256 kHz capture (128 kHz Nyquist)", 256000, 256000, "bat-128k"},
		{"384 kHz capture (192 kHz Nyquist)", 384000, 384000, "bat-192k"},
		{"0 rate falls back to 48000 (24 kHz Nyquist)", 0, 48000, "bat-24k"},
		{"negative rate falls back to 48000", -1, 48000, "bat-24k"},
		{"non-standard rate 44100", 44100, 44100, "bat-22.1k"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := BatProfileWithRate(tt.sampleRate)
			assert.Equal(t, tt.wantResample, p.ResampleRate)
			assert.Equal(t, tt.wantSuffix, ProfileSuffix(p))
		})
	}
}

func TestProfileForModelType_WithSampleRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		modelType    string
		sampleRate   int
		wantResample int
		wantSuffix   string
	}{
		{"bat model with 96 kHz mic", "bat", 96000, 96000, "bat-48k"},
		{"bat model with 192 kHz mic", "bat", 192000, 192000, "bat-96k"},
		{"bat model with 256 kHz mic", "bat", 256000, 256000, "bat-128k"},
		{"bat model with 0 rate uses default BatProfile", "bat", 0, 256000, "bat-v2"},
		{"bird model ignores sample rate", "bird", 96000, 24000, ""},
		{"empty model ignores sample rate", "", 192000, 24000, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := ProfileForModelType(tt.modelType, tt.sampleRate)
			assert.Equal(t, tt.wantResample, p.ResampleRate)
			assert.Equal(t, tt.wantSuffix, ProfileSuffix(p))

			// Also verify ProfileForModelTypeWithRate helper
			p2 := ProfileForModelTypeWithRate(tt.modelType, tt.sampleRate)
			assert.Equal(t, p, p2)
		})
	}
}
