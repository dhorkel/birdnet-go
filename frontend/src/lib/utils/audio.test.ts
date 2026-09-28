import { describe, it, expect } from 'vitest';
import { SPEED_OPTIONS, DEFAULT_PLAYBACK_SPEED, applyPlaybackRate, dbToGain } from './audio';

describe('audio utilities', () => {
  describe('SPEED_OPTIONS', () => {
    it('includes 0.1x (10%) and 0.25x for ultrasonic and high-frequency call review', () => {
      expect(SPEED_OPTIONS).toContain(0.1);
      expect(SPEED_OPTIONS).toContain(0.25);
    });

    it('contains speed options in strictly ascending order', () => {
      const sorted = [...SPEED_OPTIONS].sort((a, b) => a - b);
      expect([...SPEED_OPTIONS]).toEqual(sorted);
    });

    it('includes the default playback speed', () => {
      expect(SPEED_OPTIONS).toContain(DEFAULT_PLAYBACK_SPEED);
      expect(DEFAULT_PLAYBACK_SPEED).toBe(1.0);
    });
  });

  describe('applyPlaybackRate', () => {
    it('sets playbackRate and disables pitch preservation for tape slow-down pitch shift', () => {
      const mockAudio = {
        playbackRate: 1.0,
        preservesPitch: true,
        mozPreservesPitch: true,
        webkitPreservesPitch: true,
      } as unknown as HTMLAudioElement & {
        mozPreservesPitch?: boolean;
        webkitPreservesPitch?: boolean;
      };

      applyPlaybackRate(mockAudio, 0.1);

      expect(mockAudio.playbackRate).toBe(0.1);
      expect(mockAudio.preservesPitch).toBe(false);
      expect(mockAudio.mozPreservesPitch).toBe(false);
      expect(mockAudio.webkitPreservesPitch).toBe(false);
    });
  });

  describe('dbToGain', () => {
    it('converts 0 dB to unity gain (1.0)', () => {
      expect(dbToGain(0)).toBe(1.0);
    });

    it('converts +20 dB to 10.0', () => {
      expect(dbToGain(20)).toBeCloseTo(10.0, 5);
    });

    it('converts -20 dB to 0.1', () => {
      expect(dbToGain(-20)).toBeCloseTo(0.1, 5);
    });
  });
});
