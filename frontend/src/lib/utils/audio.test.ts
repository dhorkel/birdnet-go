import { describe, it, expect } from 'vitest';
import { generateBatTicks, DEFAULT_BAT_NYQUIST_KHZ, dbToGain } from './audio';

describe('audio utilities', () => {
  describe('generateBatTicks', () => {
    it('generates ticks for 48 kHz Nyquist (96 kHz capture)', () => {
      const ticks = generateBatTicks(48);
      expect(ticks).toEqual([45, 40, 35, 30, 25, 20, 15, 10, 5]);
    });

    it('generates ticks for 96 kHz Nyquist (192 kHz capture)', () => {
      const ticks = generateBatTicks(96);
      expect(ticks).toEqual([90, 80, 70, 60, 50, 40, 30, 20, 10]);
    });

    it('generates ticks for 128 kHz Nyquist (256 kHz capture)', () => {
      const ticks = generateBatTicks(128);
      expect(ticks).toEqual([120, 100, 80, 60, 40, 20]);
    });

    it('generates ticks for 24 kHz Nyquist (48 kHz capture)', () => {
      const ticks = generateBatTicks(24);
      expect(ticks).toEqual([20, 15, 10, 5]);
    });

    it('generates ticks for 192 kHz Nyquist (384 kHz capture)', () => {
      const ticks = generateBatTicks(192);
      expect(ticks).toEqual([180, 160, 140, 120, 100, 80, 60, 40, 20]);
    });

    it('handles DEFAULT_BAT_NYQUIST_KHZ', () => {
      expect(DEFAULT_BAT_NYQUIST_KHZ).toBe(48);
      const ticks = generateBatTicks(DEFAULT_BAT_NYQUIST_KHZ);
      expect(ticks).toEqual([45, 40, 35, 30, 25, 20, 15, 10, 5]);
    });
  });

  describe('dbToGain', () => {
    it('converts 0 dB to 1.0 gain', () => {
      expect(dbToGain(0)).toBeCloseTo(1.0);
    });

    it('converts +6 dB to ~2.0 gain', () => {
      expect(dbToGain(6)).toBeCloseTo(1.995, 2);
    });

    it('converts -6 dB to ~0.5 gain', () => {
      expect(dbToGain(-6)).toBeCloseTo(0.501, 2);
    });
  });
});
