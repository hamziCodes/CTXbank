/**
 * VERTEX Universal Design — theme bootstrap
 *
 * Dependency-free. No React, no framework. Works in Next.js, Vite, Astro,
 * plain HTML or a web view.
 *
 * A dark-mode flash on first paint is a shipping defect, not a rough edge, so
 * the resolved theme must be applied by a synchronous script in `<head>`
 * BEFORE the first paint. Everything else here reads that decision.
 */

export type Theme = "light" | "dark";
export type ThemePreference = Theme | "system";

export const THEME_STORAGE_KEY = "vertex-theme";

/* -------------------------------------------------------------------------- */
/* 1. Pre-paint script — the part that prevents the flash                     */
/* -------------------------------------------------------------------------- */

/**
 * Minified synchronous script. Must run in `<head>`, before any stylesheet
 * that paints a background, and must not be deferred or made async.
 *
 * Next.js App Router — app/layout.tsx:
 *
 *   import { themeInitScript } from "@/lib/theme-init";
 *
 *   <head>
 *     <script dangerouslySetInnerHTML={{ __html: themeInitScript }} />
 *   </head>
 *
 * Plain HTML — inline it directly in <head>:
 *
 *   <script>...contents of themeInitScript...</script>
 */
export const themeInitScript = `(function(){try{var k="${THEME_STORAGE_KEY}",s=localStorage.getItem(k),m=window.matchMedia("(prefers-color-scheme: dark)").matches,d=s==="dark"||(s!=="light"&&m);document.documentElement.classList.toggle("dark",d);document.documentElement.style.colorScheme=d?"dark":"light";}catch(e){}})();`;

/* -------------------------------------------------------------------------- */
/* 2. Reading and writing the theme                                           */
/* -------------------------------------------------------------------------- */

/** The user's stored preference, or `"system"` when they have not chosen. */
export function getStoredPreference(): ThemePreference {
  if (typeof window === "undefined") return "system";
  try {
    const stored = window.localStorage.getItem(THEME_STORAGE_KEY);
    if (stored === "light" || stored === "dark") return stored;
  } catch {
    // Private mode, or storage disabled. Fall through to system.
  }
  return "system";
}

/** What the OS is currently asking for. */
export function getSystemTheme(): Theme {
  if (typeof window === "undefined") return "light";
  return window.matchMedia("(prefers-color-scheme: dark)").matches
    ? "dark"
    : "light";
}

/** The theme that should actually be on the page right now. */
export function resolveTheme(preference: ThemePreference = getStoredPreference()): Theme {
  return preference === "system" ? getSystemTheme() : preference;
}

/**
 * Applies a theme to the document and persists it.
 *
 * Pass `"system"` to clear the stored preference and follow the OS again —
 * a three-state toggle is friendlier than a two-state one, because it lets
 * someone undo an explicit choice.
 */
export function applyTheme(preference: ThemePreference): Theme {
  const theme = resolveTheme(preference);

  if (typeof document !== "undefined") {
    document.documentElement.classList.toggle("dark", theme === "dark");
    document.documentElement.style.colorScheme = theme;
  }

  try {
    if (preference === "system") {
      window.localStorage.removeItem(THEME_STORAGE_KEY);
    } else {
      window.localStorage.setItem(THEME_STORAGE_KEY, preference);
    }
  } catch {
    // Storage unavailable. The class is still applied for this session.
  }

  return theme;
}

/* -------------------------------------------------------------------------- */
/* 3. Subscribing — for useSyncExternalStore or any observer                   */
/* -------------------------------------------------------------------------- */

/**
 * Reads the theme from the DOM rather than from state, so it can never
 * disagree with what is actually painted.
 */
export function getThemeSnapshot(): Theme {
  if (typeof document === "undefined") return "light";
  return document.documentElement.classList.contains("dark") ? "dark" : "light";
}

/** Stable server snapshot. Must be a constant to avoid a hydration mismatch. */
export function getServerThemeSnapshot(): Theme {
  return "light";
}

/**
 * Fires whenever the theme changes: in this tab, in another tab, or because
 * the OS switched while the user was on `"system"`.
 *
 * React usage:
 *
 *   const theme = useSyncExternalStore(
 *     subscribeToTheme,
 *     getThemeSnapshot,
 *     getServerThemeSnapshot,
 *   );
 */
export function subscribeToTheme(onChange: () => void): () => void {
  if (typeof window === "undefined") return () => {};

  const observer = new MutationObserver(onChange);
  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["class"],
  });

  const media = window.matchMedia("(prefers-color-scheme: dark)");
  const onMediaChange = () => {
    // Only follow the OS while the user has made no explicit choice.
    if (getStoredPreference() === "system") applyTheme("system");
    onChange();
  };

  const onStorage = (event: StorageEvent) => {
    if (event.key === THEME_STORAGE_KEY) {
      applyTheme(getStoredPreference());
      onChange();
    }
  };

  media.addEventListener("change", onMediaChange);
  window.addEventListener("storage", onStorage);

  return () => {
    observer.disconnect();
    media.removeEventListener("change", onMediaChange);
    window.removeEventListener("storage", onStorage);
  };
}

/* -------------------------------------------------------------------------- */
/* 4. Density                                                                 */
/* -------------------------------------------------------------------------- */

export type Density = "comfortable" | "compact" | "dense";

/**
 * Sets density on an element, or on the document when no element is given.
 * Nests freely — a dense table inside a comfortable page is valid.
 */
export function applyDensity(density: Density, element?: HTMLElement): void {
  if (typeof document === "undefined") return;
  (element ?? document.documentElement).setAttribute("data-density", density);
}

/* -------------------------------------------------------------------------- */
/* 5. Surface archetype                                                       */
/* -------------------------------------------------------------------------- */

export type SurfaceArchetype =
  | "marketing"
  | "dashboard"
  | "mobile-app"
  | "form-flow"
  | "console";

/** Selects the canvas. Set once on the outermost element of a surface. */
export function applySurface(
  archetype: SurfaceArchetype,
  element?: HTMLElement,
): void {
  if (typeof document === "undefined") return;
  (element ?? document.body).setAttribute("data-vertex-surface", archetype);
}
