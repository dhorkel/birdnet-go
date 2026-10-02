package spectrogram

import "fmt"

// FrequencyProfile controls spectrogram frequency range and resampling per
// detection. The gate is the detection's model type: bat models scale their
// Nyquist frequency to match the microphone capture sample rate (e.g. 96 kHz
// mic -> Nyquist = 48 kHz), eliminating wasted high-frequency space; everything
// else gets bird defaults (resample to birdResampleHz).
type FrequencyProfile struct {
	ResampleRate int    // Target sample rate in Hz; 0 means keep native rate
	suffix       string // Cache-filename token identifying the profile; "" for the default bird render
}

const (
	birdResampleHz  = 24000
	batResampleHz   = 256000 // Nyquist = 128 kHz; legacy default bat capture rate
	modelTypeBatStr = "bat"  // ai_models.model_type value that selects the bat profile

	// batCacheSuffix is the legacy cache-filename token for bat spectrograms
	// rendered at the fixed 256 kHz rate ("bat-v2").
	batCacheSuffix = "bat-v2"
)

// BirdProfile returns the default frequency profile for bird detections.
func BirdProfile() FrequencyProfile {
	return FrequencyProfile{
		ResampleRate: birdResampleHz,
	}
}

// BatProfile returns the default frequency profile for bat detections (fixed 256 kHz resample).
func BatProfile() FrequencyProfile {
	return FrequencyProfile{
		ResampleRate: batResampleHz,
		suffix:       batCacheSuffix,
	}
}

// BatProfileWithRate returns a frequency profile for bat detections scaled to the
// specified capture sample rate. The spectrogram will span from 0 Hz up to Nyquist
// (rate / 2), and the cache suffix will include the Nyquist frequency in kHz
// (e.g., "bat-48k" for a 96 kHz mic) so renders with different rates do not collide.
func BatProfileWithRate(rate int) FrequencyProfile {
	if rate <= 0 {
		rate = 48000
	}
	nyquistKhz := rate / 2000
	var suffix string
	if rate%2000 == 0 {
		suffix = fmt.Sprintf("bat-%dk", nyquistKhz)
	} else {
		suffix = fmt.Sprintf("bat-%.1fk", float64(rate)/2000.0)
	}
	return FrequencyProfile{
		ResampleRate: rate,
		suffix:       suffix,
	}
}

// ProfileForModelType selects the appropriate frequency profile based on the
// AI model's type string (as stored in ai_models.model_type).
// Bat models use the bat profile scaled to the optional sampleRate (or default
// BatProfile if omitted or non-positive); everything else uses bird defaults.
func ProfileForModelType(modelType string, sampleRate ...int) FrequencyProfile {
	if modelType == modelTypeBatStr {
		if len(sampleRate) > 0 && sampleRate[0] > 0 {
			return BatProfileWithRate(sampleRate[0])
		}
		return BatProfile()
	}
	return BirdProfile()
}

// ProfileForModelTypeWithRate selects the frequency profile for a model type with an explicit sample rate.
func ProfileForModelTypeWithRate(modelType string, sampleRate int) FrequencyProfile {
	return ProfileForModelType(modelType, sampleRate)
}

// ProfileSuffix returns a short, stable token identifying the frequency profile
// for use in spectrogram cache filenames and queue keys, so renders made with
// different profiles do not collide on disk. The default bird profile returns ""
// for backward compatibility with existing cached spectrograms.
func ProfileSuffix(p FrequencyProfile) string {
	return p.suffix
}
