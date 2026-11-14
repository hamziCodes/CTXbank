/**
 * VERTEX Universal Design — typed token mirror
 *
 * Use this for anything CSS cannot reach: chart libraries, canvas and SVG
 * fills, theme-colour meta tags, native shells, Figma sync scripts, PDF and
 * e-mail rendering.
 *
 * For anything rendered in the DOM, prefer the CSS variable or the framework
 * class. This is the second copy of the truth: if it ever disagrees with
 * tokens.css, tokens.css wins and this file is the bug.
 */

export type Mode = "light" | "dark";
export type Density = "comfortable" | "compact" | "dense";
export type SurfaceArchetype =
  | "marketing"
  | "dashboard"
  | "mobile-app"
  | "form-flow"
  | "console";

/* -------------------------------------------------------------------------- */
/* Navy ramp — the single accent family, hue locked 203-205                    */
/* -------------------------------------------------------------------------- */

export const navy = {
  50: "#eff6fb",
  100: "#dbebf5",
  200: "#bad6e8",
  300: "#90bbd5",
  400: "#609abe",
  500: "#236798",
  600: "#1a537c",
  700: "#0b3b5b",
  800: "#042940",
  900: "#021c2c",
  950: "#01111b",
} as const;

/* -------------------------------------------------------------------------- */
/* Colour, per mode                                                           */
/* -------------------------------------------------------------------------- */

export const colors = {
  light: {
    brandPrimary: "#042940",
    brandHover: "#0b3b5b",
    cta: "#042940",
    ctaHover: "#0b3b5b",
    ctaForeground: "#ffffff",
    accentText: "#042940",
    accentTextHover: "#0b3b5b",
    focusRing: "#042940",

    background: "#ffffff",
    foreground: "#000000",
    card: "#ffffff",
    muted: "#000000",
    border: "rgb(0 0 0 / 0.08)",
    buttonSecondary: "#000000",
    surfaceDark: "#000000",
    surfaceHover: "rgb(0 0 0 / 0.04)",

    appCanvas: "#f2f6f9",
    appSurface: "#ffffff",
    appSurfaceSunken: "#e7eef4",
    appSurfaceInverse: "#042940",
    appSurfaceHover: "rgb(4 41 64 / 0.04)",
    appOnInverse: "#ffffff",
    appLine: "#dbe4ec",
    appLineStrong: "#c2d0dd",
    appLineInverse: "rgb(255 255 255 / 0.18)",
    appInk: "#000000",
    appInkSoft: "rgb(0 0 0 / 0.64)",
    appInkMuted: "rgb(0 0 0 / 0.52)",

    positive: "#0f6b4f",
    positiveSoft: "#e2f2ec",
    positiveLine: "#bfdfd2",
    negative: "#b02a25",
    negativeSoft: "#fbe9e7",
    negativeLine: "#f0cdc8",
    attention: "#8a5b00",
    attentionSoft: "#fdf1dc",
    attentionLine: "#eed9b0",
    info: "#042940",
    infoSoft: "#eff6fb",
    infoLine: "#cbdeeb",

    shadowColor: "#dce5ec",
    shadowColorStrong: "#b9c8d5",
    scrim: "rgb(4 41 64 / 0.38)",
  },
  dark: {
    brandPrimary: "#1a537c",
    brandHover: "#236798",
    cta: "#1a537c",
    ctaHover: "#236798",
    ctaForeground: "#ffffff",
    accentText: "#609abe",
    accentTextHover: "#90bbd5",
    focusRing: "#609abe",

    background: "#161719",
    foreground: "#f3f4f6",
    card: "#212326",
    muted: "#9ca3af",
    border: "#2e3238",
    buttonSecondary: "#2d3035",
    surfaceDark: "#212326",
    surfaceHover: "#2d3035",

    appCanvas: "#161719",
    appSurface: "#212326",
    appSurfaceSunken: "#101113",
    appSurfaceInverse: "#f3f4f6",
    appSurfaceHover: "#2d3035",
    appOnInverse: "#161719",
    appLine: "#2e3238",
    appLineStrong: "#3d434b",
    appLineInverse: "rgb(22 23 25 / 0.18)",
    appInk: "#f3f4f6",
    appInkSoft: "rgb(243 244 246 / 0.72)",
    appInkMuted: "#9ca3af",

    positive: "#3fae86",
    positiveSoft: "rgb(63 174 134 / 0.16)",
    positiveLine: "rgb(63 174 134 / 0.34)",
    negative: "#f0766a",
    negativeSoft: "rgb(240 118 106 / 0.16)",
    negativeLine: "rgb(240 118 106 / 0.34)",
    attention: "#e0a94a",
    attentionSoft: "rgb(224 169 74 / 0.16)",
    attentionLine: "rgb(224 169 74 / 0.34)",
    info: "#609abe",
    infoSoft: "rgb(96 154 190 / 0.16)",
    infoLine: "rgb(96 154 190 / 0.34)",

    shadowColor: "#0d0e10",
    shadowColorStrong: "#08090a",
    scrim: "rgb(0 0 0 / 0.58)",
  },
} as const satisfies Record<Mode, Record<string, string>>;

/** Footer chrome stays charcoal in both modes and does not flip. */
export const footer = {
  bg: "#161719",
  elevated: "#212326",
  foreground: "#f3f4f6",
  muted: "#9ca3af",
  border: "#2e3238",
} as const;

/**
 * Chart series order. Fixed so the same data reads the same way in every
 * VERTEX product. Never reorder per chart.
 */
export const chartSeries = {
  light: [navy[800], navy[500], navy[300], navy[600], navy[200]],
  dark: [navy[300], navy[400], navy[500], navy[200], navy[600]],
} as const;

/* -------------------------------------------------------------------------- */
/* Shape                                                                      */
/* -------------------------------------------------------------------------- */

export const radius = {
  control: "12px",
  slab: "18px",
  panel: "24px",
  sheet: "28px",
  pill: "9999px",
  /** Marketing archetype only. Never inside an application surface. */
  marketing: "clamp(40px, 4.167vw, 60px)",
} as const;

export const borderWidth = {
  hairline: "1px",
  slab: "1.5px",
  emphasis: "2px",
} as const;

export const shadow = {
  slab: (mode: Mode) => `3px 4px 0 ${colors[mode].shadowColor}`,
  slabStrong: (mode: Mode) => `3px 4px 0 ${colors[mode].shadowColorStrong}`,
  slabPressed: (mode: Mode) => `1px 2px 0 ${colors[mode].shadowColorStrong}`,
  /** The only blurred shadow in the system. */
  sheet: (mode: Mode) =>
    mode === "light"
      ? "0 -5px 20px rgb(4 41 64 / 0.12)"
      : "0 -5px 20px rgb(0 0 0 / 0.45)",
} as const;

/* -------------------------------------------------------------------------- */
/* Type                                                                       */
/* -------------------------------------------------------------------------- */

export const font = {
  sans: '"Inter Tight", "Inter", system-ui, -apple-system, "Segoe UI", sans-serif',
  /** Permitted only for code, logs, IDs and diffs. Never UI chrome. */
  mono: 'ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace',
  weight: {
    display: 900,
    bold: 700,
    body: 500,
    light: 300,
  },
} as const;

export const fontSize = {
  display: "clamp(4.25rem, 18vw, 17.525rem)",
  hero: "clamp(3rem, 5.56vw, 5rem)",
  title: "clamp(2rem, 3.61vw, 3.25rem)",
  lead: "clamp(1.125rem, 1.67vw, 1.5rem)",
  menu: "3rem",
  metricXl: "40px",
  metric: "30px",
  metricSm: "22px",
  heading: "22px",
  subheading: "18px",
  body: "15px",
  bodySm: "14px",
  label: "13px",
  eyebrow: "11px",
  nav: "10px",
} as const;

/* -------------------------------------------------------------------------- */
/* Density                                                                    */
/* -------------------------------------------------------------------------- */

/**
 * Mirrors the `[data-density]` blocks in tokens.css. Components should read
 * the CSS variables; this is here for canvas, native and layout maths.
 *
 * `compact` and `dense` are pointer-precise. On a coarse pointer, controlHeight
 * clamps back to 44px — see the `@media (pointer: coarse)` block in tokens.css.
 */
export const density = {
  comfortable: {
    controlHeight: 48,
    controlPadX: 20,
    controlGap: 8,
    slabPad: 24,
    stackGap: 18,
    rowPadY: 12,
    rowGap: 12,
    fieldGap: 20,
    iconTile: 36,
    textBody: 15,
  },
  compact: {
    controlHeight: 40,
    controlPadX: 16,
    controlGap: 6,
    slabPad: 18,
    stackGap: 14,
    rowPadY: 10,
    rowGap: 8,
    fieldGap: 16,
    iconTile: 32,
    textBody: 14,
  },
  dense: {
    controlHeight: 34,
    controlPadX: 12,
    controlGap: 4,
    slabPad: 14,
    stackGap: 10,
    rowPadY: 7,
    rowGap: 4,
    fieldGap: 12,
    iconTile: 28,
    textBody: 13,
  },
} as const satisfies Record<Density, Record<string, number>>;

export const touchFloor = 44;

/* -------------------------------------------------------------------------- */
/* Layout                                                                     */
/* -------------------------------------------------------------------------- */

export const layout = {
  pageMax: "1440px",
  appMax: "740px",
  proseMax: "720px",
  navHeight: "72px",
  bottomNavHeight: "70px",
  sidebarWidth: "264px",
  sidebarWidthCollapsed: "72px",
  desktopBottomNavWidth: "600px",
} as const;

/** Canvas token per archetype. See spec/SURFACES.md for the full rules. */
export const archetypeCanvas = {
  marketing: "--background",
  dashboard: "--app-canvas",
  "mobile-app": "--app-canvas",
  "form-flow": "--app-canvas",
  console: "--app-canvas",
} as const satisfies Record<SurfaceArchetype, string>;

/* -------------------------------------------------------------------------- */
/* Motion — one curve, never mixed                                            */
/* -------------------------------------------------------------------------- */

export const easeVertex = [0.16, 1, 0.3, 1] as const;
export const easeVertexCss = "cubic-bezier(0.16, 1, 0.3, 1)";

export const duration = {
  micro: 0.2,
  swap: 0.32,
  panel: 0.45,
  reveal: 0.8,
  hero: 0.95,
} as const;

export const revealTravel = {
  min: 16,
  default: 24,
  hero: 28,
  form: 32,
} as const;

export type MotionBudget = "full" | "reduced" | "minimal";

/**
 * What each budget permits. Recorded per project in project-profile.md.
 * `minimal` is a legitimate choice for low-end devices, not a failure.
 */
export const motionBudget = {
  full: {
    reveals: true,
    scrollLinkedTransforms: true,
    layoutAnimations: true,
    loops: true,
  },
  reduced: {
    reveals: true,
    scrollLinkedTransforms: false,
    layoutAnimations: true,
    loops: false,
  },
  minimal: {
    reveals: false,
    scrollLinkedTransforms: false,
    layoutAnimations: false,
    loops: false,
  },
} as const satisfies Record<MotionBudget, Record<string, boolean>>;

/* -------------------------------------------------------------------------- */
/* Money — one formatter for every VERTEX surface                             */
/* -------------------------------------------------------------------------- */

export type MoneySign = "none" | "in" | "out";

/**
 * Formats a figure as `PKR 184,250`, `+ PKR 82,000` or `- PKR 46,780`.
 *
 * Pass a positive magnitude and express direction through `sign`, so the
 * renderer can pick the colour and prefix together and they can never
 * disagree. Locale and currency come from project-profile.md — `PKR` and
 * `en-PK` are the defaults, not a constraint.
 */
export function formatMoney(
  amount: number,
  sign: MoneySign = "none",
  currency = "PKR",
  precision: 0 | 2 = 0,
  locale = "en-PK",
): string {
  const body = `${currency} ${new Intl.NumberFormat(locale, {
    minimumFractionDigits: precision,
    maximumFractionDigits: precision,
  }).format(Math.abs(amount))}`;

  if (sign === "in") return `+ ${body}`;
  if (sign === "out") return `- ${body}`;
  return body;
}

/** Maps a direction to its semantic colour token name. */
export function moneyTone(sign: MoneySign): "positive" | "negative" | "ink" {
  if (sign === "in") return "positive";
  if (sign === "out") return "negative";
  return "ink";
}

/* -------------------------------------------------------------------------- */
/* Contrast reference — measured, not estimated                               */
/* -------------------------------------------------------------------------- */

/**
 * WCAG 2.1 ratios for every pair that carries text. Below 4.5 is decorative
 * only and must never hold body copy.
 */
export const contrast = {
  "light: #000000 on #ffffff": 21.0,
  "light: #042940 on #ffffff": 15.0,
  "light: #ffffff on #042940": 15.0,
  "light: #0f6b4f on #ffffff": 6.5,
  "light: #0f6b4f on #e2f2ec": 5.6,
  "light: #b02a25 on #ffffff": 6.6,
  "light: #8a5b00 on #ffffff": 5.9,
  "dark: #f3f4f6 on #161719": 15.4,
  "dark: #609abe on #161719": 5.8,
  "dark: #90bbd5 on #161719": 8.8,
  "dark: #3fae86 on #161719": 6.5,
  "dark: #f0766a on #161719": 6.4,
  "dark: #e0a94a on #161719": 8.5,
  /** FAILS. --brand-primary is a fill token in dark mode, never text. */
  "dark: #1a537c on #161719": 2.2,
  /** Decorative only on a light canvas. */
  "light: #609abe on #ffffff": 3.1,
} as const;

export const tokens = {
  navy,
  colors,
  footer,
  chartSeries,
  radius,
  borderWidth,
  shadow,
  font,
  fontSize,
  density,
  touchFloor,
  layout,
  archetypeCanvas,
  easeVertex,
  easeVertexCss,
  duration,
  revealTravel,
  motionBudget,
  contrast,
} as const;

export default tokens;
