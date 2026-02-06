import { describe, it, expect } from 'vitest';

/**
 * WCAG AA Contrast Validation Tests
 *
 * Tests verify that all Kabuki Adaptatif design system colors
 * meet WCAG AA accessibility standards (4.5:1 for normal text, 3:1 for large text)
 */

// WCAG relative luminance calculation
// Based on https://www.w3.org/TR/WCAG20/#relativeluminancedef
function getLuminance(r: number, g: number, b: number): number {
  const [rs, gs, bs] = [r, g, b].map(c => {
    c = c / 255;
    return c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4);
  });
  return 0.2126 * rs + 0.7152 * gs + 0.0722 * bs;
}

// Contrast ratio calculation
// Based on https://www.w3.org/TR/WCAG20/#contrast-ratiodef
function getContrastRatio(rgb1: string, rgb2: string): number {
  // Parse hex color
  const hexToRgb = (hex: string): [number, number, number] => {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    if (!result) throw new Error(`Invalid hex color: ${hex}`);
    return [parseInt(result[1], 16), parseInt(result[2], 16), parseInt(result[3], 16)];
  };

  // Parse rgba color
  const rgbaToRgb = (rgba: string): [number, number, number] => {
    const match = /rgba?\((\d+),\s*(\d+),\s*(\d+)/i.exec(rgba);
    if (!match) throw new Error(`Invalid rgba color: ${rgba}`);
    return [parseInt(match[1]), parseInt(match[2]), parseInt(match[3])];
  };

  const rgb1Parsed = rgb1.startsWith('#') ? hexToRgb(rgb1) : rgbaToRgb(rgb1);
  const rgb2Parsed = rgb2.startsWith('#') ? hexToRgb(rgb2) : rgbaToRgb(rgb2);

  const l1 = getLuminance(...rgb1Parsed);
  const l2 = getLuminance(...rgb2Parsed);

  const lighter = Math.max(l1, l2);
  const darker = Math.min(l1, l2);

  return (lighter + 0.05) / (darker + 0.05);
}

describe('WCAG AA Contrast Validation', () => {
  // Test data: [foreground, background, description, minRatio]
  const testCases: Array<[string, string, string, number]> = [
    // Primary text combinations
    ['#F5F7FF', '#0A0E1A', 'Primary text on base background', 4.5],
    ['#F5F7FF', '#1A1F2E', 'Primary text on surface background', 4.5],
    ['#F5F7FF', '#252A3B', 'Primary text on elevated background', 4.5],

    // Secondary text combinations
    ['#A8B3D1', '#1A1F2E', 'Secondary text on surface background', 3],
    ['#A8B3D1', '#0A0E1A', 'Secondary text on base background', 3],
    ['#A8B3D1', '#252A3B', 'Secondary text on elevated background', 3],

    // Muted text combinations
    ['#6B7694', '#1A1F2E', 'Muted text on surface background', 3],
    ['#6B7694', '#0A0E1A', 'Muted text on base background', 3],

    // Accent text combinations (primary CTA on base, lighter variant on surface)
    ['#E63946', '#0A0E1A', 'Red accent (CTA) on base background', 4.5],
    ['#FB6F8A', '#1A1F2E', 'Red accent (lighter) on surface background', 4.5],

    ['#F59E0B', '#0A0E1A', 'Gold accent on base background', 4.5],
    ['#06B6D4', '#0A0E1A', 'Cyan accent on base background', 4.5],

    // Semantic alerts
    ['#4ADE80', '#0A0E1A', 'Success text on base background', 4.5],
    ['#F87171', '#0A0E1A', 'Error text on base background', 4.5],
    ['#22D3EE', '#0A0E1A', 'Info text on base background', 4.5],
  ];

  testCases.forEach(([fg, bg, description, minRatio]) => {
    it(`${description} - should have contrast ratio >= ${minRatio}:1`, () => {
      const ratio = getContrastRatio(fg, bg);
      console.log(`  ${description}: ${ratio.toFixed(2)}:1 ✓`);
      expect(ratio).toBeGreaterThanOrEqual(minRatio);
    });
  });

  // Print summary
  it('Summary: All color combinations meet WCAG AA standards', () => {
    console.log('\n📊 Contrast Ratio Summary:\n');
    let allPass = true;

    testCases.forEach(([fg, bg, description, minRatio]) => {
      const ratio = getContrastRatio(fg, bg);
      const meets = ratio >= minRatio ? '✅' : '❌';
      console.log(`  ${meets} ${description}`);
      console.log(`     Foreground: ${fg} | Background: ${bg}`);
      console.log(`     Ratio: ${ratio.toFixed(2)}:1 (min: ${minRatio}:1)\n`);

      if (ratio < minRatio) {
        allPass = false;
      }
    });

    expect(allPass).toBe(true);
  });
});
