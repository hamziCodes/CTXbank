/**
 * VERTEX Universal Design — Tailwind v3 preset
 *
 * For projects still on Tailwind v3. Tailwind v4 projects use
 * tailwind-v4.css instead; never load both.
 *
 * Install:
 *   1. Import tokens/tokens.css in the global stylesheet.
 *   2. tailwind.config.js:
 *
 *      module.exports = {
 *        presets: [require("./vertex-universal-design/tokens/tailwind-v3.cjs")],
 *        content: ["./src/**\/*.{js,ts,jsx,tsx}"],
 *      };
 *
 * Every value below resolves to a CSS variable, so light and dark mode work
 * through `html.dark` with no duplicated colour config.
 */

/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: ["class", ".dark"],
  theme: {
    extend: {
      colors: {
        "brand-primary": "var(--brand-primary)",
        "brand-hover": "var(--brand-hover)",
        cta: "var(--cta)",
        "cta-hover": "var(--cta-hover)",
        "cta-foreground": "var(--cta-foreground)",
        "accent-text": "var(--accent-text)",
        "accent-text-hover": "var(--accent-text-hover)",
        "focus-ring": "var(--focus-ring)",
        "product-signal": "var(--product-signal)",

        navy: {
          50: "var(--navy-50)",
          100: "var(--navy-100)",
          200: "var(--navy-200)",
          300: "var(--navy-300)",
          400: "var(--navy-400)",
          500: "var(--navy-500)",
          600: "var(--navy-600)",
          700: "var(--navy-700)",
          800: "var(--navy-800)",
          900: "var(--navy-900)",
          950: "var(--navy-950)",
        },

        background: "var(--background)",
        foreground: "var(--foreground)",
        card: "var(--card)",
        muted: "var(--muted)",
        border: "var(--border)",
        "button-secondary": "var(--button-secondary)",
        "surface-dark": "var(--surface-dark)",
        "surface-hover": "var(--surface-hover)",

        "app-canvas": "var(--app-canvas)",
        "app-surface": "var(--app-surface)",
        "app-surface-sunken": "var(--app-surface-sunken)",
        "app-surface-inverse": "var(--app-surface-inverse)",
        "app-surface-hover": "var(--app-surface-hover)",
        "app-on-inverse": "var(--app-on-inverse)",
        "app-line": "var(--app-line)",
        "app-line-strong": "var(--app-line-strong)",
        "app-line-inverse": "var(--app-line-inverse)",
        "app-ink": "var(--app-ink)",
        "app-ink-soft": "var(--app-ink-soft)",
        "app-ink-muted": "var(--app-ink-muted)",

        footer: "var(--footer-bg)",
        "footer-elevated": "var(--footer-elevated)",
        "footer-foreground": "var(--footer-foreground)",
        "footer-muted": "var(--footer-muted)",
        "footer-border": "var(--footer-border)",

        positive: "var(--positive)",
        "positive-soft": "var(--positive-soft)",
        "positive-line": "var(--positive-line)",
        negative: "var(--negative)",
        "negative-soft": "var(--negative-soft)",
        "negative-line": "var(--negative-line)",
        attention: "var(--attention)",
        "attention-soft": "var(--attention-soft)",
        "attention-line": "var(--attention-line)",
        info: "var(--info)",
        "info-soft": "var(--info-soft)",
        "info-line": "var(--info-line)",
        scrim: "var(--scrim)",
      },

      fontFamily: {
        sans: "var(--font-sans)",
        mono: "var(--font-mono)",
      },

      fontSize: {
        display: "var(--fs-display)",
        hero: "var(--fs-hero)",
        title: "var(--fs-title)",
        lead: "var(--fs-lead)",
        menu: "var(--fs-menu)",
      },

      borderRadius: {
        control: "var(--radius-control)",
        slab: "var(--radius-slab)",
        panel: "var(--radius-panel)",
        sheet: "var(--radius-sheet)",
        card: "var(--radius-marketing)",
      },

      boxShadow: {
        slab: "var(--shadow-slab)",
        "slab-strong": "var(--shadow-slab-strong)",
        "slab-pressed": "var(--shadow-slab-pressed)",
        sheet: "var(--shadow-sheet)",
      },

      spacing: {
        page: "var(--page-pad)",
        gutter: "var(--page-gutter)",
        inset: "var(--page-inset)",
        panel: "var(--page-panel)",
        nav: "var(--nav-height)",
        "bottom-nav": "var(--bottom-nav-height)",
        control: "var(--control-h)",
        "slab-pad": "var(--slab-pad)",
        stack: "var(--stack-gap)",
      },

      maxWidth: {
        page: "var(--page-max)",
        app: "var(--app-max)",
        "prose-vx": "var(--prose-max)",
      },

      width: {
        sidebar: "var(--sidebar-width)",
        "sidebar-collapsed": "var(--sidebar-width-collapsed)",
      },

      height: {
        nav: "var(--nav-height)",
        "bottom-nav": "var(--bottom-nav-height)",
        control: "var(--control-h)",
      },

      transitionTimingFunction: {
        vertex: "var(--ease-vertex)",
      },

      transitionDuration: {
        micro: "200ms",
        swap: "320ms",
        panel: "450ms",
        reveal: "800ms",
        hero: "950ms",
      },

      zIndex: {
        1: "1",
        2: "2",
        nav: "50",
        scrim: "40",
        sheet: "50",
      },
    },
  },
  plugins: [],
};
