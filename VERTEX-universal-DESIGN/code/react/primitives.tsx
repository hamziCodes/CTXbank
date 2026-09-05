"use client";

/**
 * VERTEX Universal Design — React primitives
 *
 * Reference implementation of every component in spec/DESIGN.md §7.
 * Density-aware and archetype-neutral: the same components serve a marketing
 * page and an admin console, because they read the density CSS variables
 * instead of hardcoding heights and padding.
 *
 * REQUIREMENTS
 *   tokens/tokens.css imported once at the app root, plus one Tailwind adapter
 *   (tailwind-v4.css or tailwind-v3.cjs).
 *   Packages: react, motion, lucide-react, clsx, tailwind-merge.
 *
 * NO TAILWIND?
 *   Port these components against the .vx-* utility classes in tokens.css.
 *   Keep the same states, semantics and token bindings — those are the
 *   contract; the class syntax is not.
 *
 * DROP-IN LOCATION
 *   src/components/ui/primitives.tsx (or split per component as it grows)
 */

import {
  useCallback,
  useEffect,
  useId,
  useRef,
  useState,
  useSyncExternalStore,
  type ButtonHTMLAttributes,
  type InputHTMLAttributes,
  type ReactNode,
  type TextareaHTMLAttributes,
} from "react";
import {
  AnimatePresence,
  cubicBezier,
  motion,
  useReducedMotion,
  useScroll,
  useTransform,
} from "motion/react";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import {
  ChevronRight,
  ChevronsUpDown,
  Inbox,
  Loader2,
  Moon,
  Plus,
  Search,
  Sun,
  TriangleAlert,
  X,
  type LucideIcon,
} from "lucide-react";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/* ==========================================================================
   Motion — one curve for the entire system
   ========================================================================== */

export const easeVertex = [0.16, 1, 0.3, 1] as const;
export const easeVertexFn = cubicBezier(0.16, 1, 0.3, 1);

export const duration = {
  micro: 0.2,
  swap: 0.32,
  panel: 0.45,
  reveal: 0.8,
  hero: 0.95,
} as const;

/* ==========================================================================
   Types
   ========================================================================== */

export type Theme = "light" | "dark";
export type Density = "comfortable" | "compact" | "dense";
export type SurfaceArchetype =
  | "marketing"
  | "dashboard"
  | "mobile-app"
  | "form-flow"
  | "console";

export type SlabTone = "surface" | "brand" | "inverse" | "attention";
export type SemanticTone = "neutral" | "positive" | "negative" | "attention" | "info";
export type MoneySign = "none" | "in" | "out";

const slabToneClass: Record<SlabTone, string> = {
  surface:
    "bg-app-surface text-app-ink border-[1.5px] border-app-line-strong shadow-slab",
  brand: "bg-brand-primary text-white border-0",
  inverse: "bg-app-surface-inverse text-app-on-inverse border-0",
  attention:
    "bg-attention-soft text-app-ink border-[1.5px] border-attention-line shadow-slab",
};

const semanticTextClass: Record<SemanticTone, string> = {
  neutral: "text-app-ink",
  positive: "text-positive",
  negative: "text-negative",
  attention: "text-attention",
  info: "text-info",
};

const semanticFillClass: Record<SemanticTone, string> = {
  neutral: "bg-app-surface-sunken text-app-ink-soft",
  positive: "bg-positive-soft text-positive",
  negative: "bg-negative-soft text-negative",
  attention: "bg-attention-soft text-attention",
  info: "bg-info-soft text-info",
};

const focusRing =
  "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-focus-ring";

/* ==========================================================================
   Theme
   --------------------------------------------------------------------------
   The pre-paint bootstrap lives in tokens/theme-init.ts and must run in
   <head> before first paint. This hook is only the React binding: it reads
   the class the bootstrap applied, so state can never disagree with what is
   actually painted. Keep THEME_KEY identical in both files.
   ========================================================================== */

export const THEME_KEY = "vertex-theme";

function subscribeTheme(onChange: () => void) {
  const observer = new MutationObserver(onChange);
  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["class"],
  });
  window.addEventListener("storage", onChange);
  return () => {
    observer.disconnect();
    window.removeEventListener("storage", onChange);
  };
}

function readTheme(): Theme {
  return document.documentElement.classList.contains("dark") ? "dark" : "light";
}

export function useVertexTheme() {
  const theme = useSyncExternalStore(
    subscribeTheme,
    readTheme,
    () => "light" as const,
  );

  const setTheme = useCallback((next: Theme) => {
    document.documentElement.classList.toggle("dark", next === "dark");
    document.documentElement.style.colorScheme = next;
    try {
      window.localStorage.setItem(THEME_KEY, next);
    } catch {
      // Storage unavailable; the class still applies for this session.
    }
  }, []);

  return { theme, setTheme };
}

export function ThemeToggle({ className }: { className?: string }) {
  const reduce = useReducedMotion();
  const { theme, setTheme } = useVertexTheme();

  // Renders both icons until mounted so the server and client markup agree.
  const mounted = useSyncExternalStore(
    (onChange) => {
      queueMicrotask(onChange);
      return () => {};
    },
    () => true,
    () => false,
  );

  return (
    <button
      type="button"
      aria-label={theme === "dark" ? "Switch to light theme" : "Switch to dark theme"}
      aria-pressed={theme === "dark"}
      onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
      className={cn(
        "inline-flex size-11 cursor-pointer items-center justify-center rounded-control",
        "text-foreground transition-colors hover:text-accent-text-hover",
        focusRing,
        className,
      )}
    >
      <span className="relative flex size-4.5 items-center justify-center">
        {mounted ? (
          <AnimatePresence mode="wait" initial={false}>
            <motion.span
              key={theme}
              initial={reduce ? false : { opacity: 0, y: 5 }}
              animate={{ opacity: 1, y: 0 }}
              exit={reduce ? undefined : { opacity: 0, y: -5 }}
              transition={{ duration: duration.swap, ease: easeVertex }}
              className="absolute inset-0 flex items-center justify-center"
            >
              {theme === "dark" ? (
                <Sun size={18} strokeWidth={1.75} />
              ) : (
                <Moon size={18} strokeWidth={1.75} />
              )}
            </motion.span>
          </AnimatePresence>
        ) : (
          <>
            <Moon size={18} strokeWidth={1.75} className="dark:hidden" aria-hidden />
            <Sun size={18} strokeWidth={1.75} className="hidden dark:block" aria-hidden />
          </>
        )}
      </span>
    </button>
  );
}

/* ==========================================================================
   Surface — declares the archetype and density for everything inside
   ========================================================================== */

export function Surface({
  archetype,
  density = "comfortable",
  className,
  children,
}: {
  archetype: SurfaceArchetype;
  density?: Density;
  className?: string;
  children: ReactNode;
}) {
  return (
    <div
      data-vertex-surface={archetype}
      data-density={density}
      className={cn("min-h-svh", className)}
    >
      {children}
    </div>
  );
}

/** Nests a different density inside a surface, e.g. a dense table on a compact page. */
export function DensityScope({
  density,
  className,
  children,
}: {
  density: Density;
  className?: string;
  children: ReactNode;
}) {
  return (
    <div data-density={density} className={className}>
      {children}
    </div>
  );
}

/* ==========================================================================
   Money
   ========================================================================== */

/**
 * The single money formatter. Pass a positive magnitude and express direction
 * through `sign`, so the prefix and the colour can never disagree.
 * Locale and currency come from the project profile.
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

export function Money({
  amount,
  sign = "none",
  currency = "PKR",
  precision = 0,
  locale = "en-PK",
  className,
}: {
  amount: number;
  sign?: MoneySign;
  currency?: string;
  precision?: 0 | 2;
  locale?: string;
  className?: string;
}) {
  const tone: SemanticTone =
    sign === "in" ? "positive" : sign === "out" ? "negative" : "neutral";

  return (
    <span
      className={cn(
        "font-bold tabular-nums tracking-[-0.01em]",
        semanticTextClass[tone],
        className,
      )}
    >
      {formatMoney(amount, sign, currency, precision, locale)}
    </span>
  );
}

/* ==========================================================================
   Type primitives
   ========================================================================== */

export function Eyebrow({
  children,
  className,
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <p
      className={cn(
        "text-[11px] font-bold uppercase leading-[1.2] tracking-[0.12em] text-app-ink-muted",
        className,
      )}
    >
      {children}
    </p>
  );
}

/** Parenthetical section kicker — the marketing signature, e.g. (HOW WE WORK) */
export function Kicker({
  children,
  className,
}: {
  children: string;
  className?: string;
}) {
  return (
    <p
      className={cn(
        "text-[14px] font-bold leading-4 text-foreground lg:text-[16px]",
        className,
      )}
    >
      ({children.toUpperCase()})
    </p>
  );
}

export function BrandWordmark({
  brand = "VERTEX",
  suffix = "studio",
  className,
  slashClassName = "text-brand-primary",
}: {
  brand?: string;
  suffix?: string;
  className?: string;
  slashClassName?: string;
}) {
  return (
    <span className={cn("whitespace-nowrap", className)}>
      <span className="font-black">{brand}</span>
      <span className={cn("font-black", slashClassName)}>/</span>
      <span className="font-light">{suffix}</span>
    </span>
  );
}

/* ==========================================================================
   Reveal — fires on first mount only, never on state change
   ========================================================================== */

export function Reveal({
  children,
  delay = 0,
  travel = 24,
  amount = 0.6,
  disabled = false,
  className,
}: {
  children: ReactNode;
  delay?: number;
  travel?: number;
  amount?: number;
  /** Set true when the project's motion budget is `minimal`. */
  disabled?: boolean;
  className?: string;
}) {
  const reduce = useReducedMotion();
  const off = reduce || disabled;

  return (
    <motion.div
      className={className}
      initial={off ? false : { opacity: 0, y: travel }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, amount }}
      transition={{ duration: duration.reveal, delay, ease: easeVertex }}
    >
      {children}
    </motion.div>
  );
}

/* ==========================================================================
   Buttons
   ========================================================================== */

export type ButtonVariant = "primary" | "inverse" | "ghost" | "danger";

const buttonVariantClass: Record<ButtonVariant, string> = {
  primary: "bg-brand-primary text-white hover:bg-brand-hover",
  inverse:
    "bg-app-surface-inverse text-app-on-inverse shadow-slab-strong active:translate-x-0.5 active:translate-y-0.5 active:shadow-slab-pressed",
  ghost:
    "border border-app-line bg-transparent text-app-ink hover:border-app-line-strong hover:text-accent-text-hover",
  danger: "bg-negative text-white hover:opacity-90",
};

export function Button({
  variant = "primary",
  icon: Icon,
  loading = false,
  className,
  children,
  disabled,
  ...rest
}: ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: ButtonVariant;
  icon?: LucideIcon;
  loading?: boolean;
}) {
  return (
    <button
      type="button"
      aria-busy={loading || undefined}
      aria-disabled={disabled || undefined}
      disabled={disabled || loading}
      className={cn(
        "inline-flex items-center justify-center rounded-control",
        "min-h-[var(--control-h)] gap-[var(--control-gap)] px-[var(--control-pad-x)]",
        "text-[var(--text-body)] font-bold leading-none",
        "transition-[background-color,color,border-color,transform,box-shadow] duration-200 ease-vertex",
        focusRing,
        buttonVariantClass[variant],
        (disabled || loading) && "pointer-events-none opacity-40",
        className,
      )}
      {...rest}
    >
      {loading ? (
        <Loader2 size={18} strokeWidth={1.75} className="animate-spin" />
      ) : Icon ? (
        <Icon size={18} strokeWidth={1.75} />
      ) : null}
      {children}
    </button>
  );
}

export function TextLink({
  href,
  children,
  className,
}: {
  href: string;
  children: ReactNode;
  className?: string;
}) {
  return (
    <a
      href={href}
      className={cn(
        "inline-flex items-center gap-1 rounded-control text-[13px] font-bold text-accent-text",
        "transition-colors hover:text-accent-text-hover",
        focusRing,
        className,
      )}
    >
      {children}
      <ChevronRight size={14} strokeWidth={2} />
    </a>
  );
}

export function FilterChip({
  active,
  onClick,
  children,
}: {
  active: boolean;
  onClick: () => void;
  children: ReactNode;
}) {
  return (
    <button
      type="button"
      aria-pressed={active}
      onClick={onClick}
      className={cn(
        "min-h-[var(--control-h)] rounded-full px-4 text-[13px] font-bold",
        "transition-colors duration-200 ease-vertex",
        focusRing,
        active
          ? "border border-transparent bg-app-surface-inverse text-app-on-inverse"
          : "border border-app-line bg-app-surface text-app-ink-soft hover:border-app-line-strong",
      )}
    >
      {children}
    </button>
  );
}

/**
 * Sticker-press floating action button. `mobile-app` archetype.
 * Sits above the bottom nav; the press offset is the only feedback a touch
 * user gets, so it matters more here than anywhere else.
 */
export function FabButton({
  label,
  onClick,
  icon: Icon = Plus,
}: {
  label: string;
  onClick: () => void;
  icon?: LucideIcon;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "fixed bottom-[86px] right-[max(20px,calc((100vw-740px)/2))] z-20",
        "inline-flex h-13 items-center gap-2 rounded-[17px] px-[18px]",
        "bg-app-surface-inverse text-[13px] font-bold text-app-on-inverse",
        "shadow-slab-strong transition-[transform,box-shadow] duration-200 ease-vertex",
        "active:translate-x-0.5 active:translate-y-0.5 active:shadow-slab-pressed",
        focusRing,
      )}
    >
      <Icon size={18} strokeWidth={2} />
      {label}
    </button>
  );
}

/* ==========================================================================
   Badge
   ========================================================================== */

export function Badge({
  tone = "neutral",
  children,
  className,
}: {
  tone?: SemanticTone;
  children: ReactNode;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[12px] font-bold",
        semanticFillClass[tone],
        className,
      )}
    >
      {children}
    </span>
  );
}

/* ==========================================================================
   Slabs
   ========================================================================== */

export function Slab({
  tone = "surface",
  radius = "panel",
  className,
  children,
}: {
  tone?: SlabTone;
  radius?: "slab" | "panel";
  className?: string;
  children: ReactNode;
}) {
  return (
    <div
      className={cn(
        "p-[var(--slab-pad)]",
        radius === "panel" ? "rounded-panel" : "rounded-slab",
        slabToneClass[tone],
        className,
      )}
    >
      {children}
    </div>
  );
}

export function MetricSlab({
  eyebrow,
  amount,
  sign = "none",
  currency = "PKR",
  locale = "en-PK",
  delta,
  caption,
  action,
  icon: Icon,
  tone = "brand",
}: {
  eyebrow: string;
  amount: number;
  sign?: MoneySign;
  currency?: string;
  locale?: string;
  /** A delta without a stated period is not shippable. */
  delta?: { percent: number; direction: "up" | "down" | "flat"; period: string };
  caption?: string;
  action?: { label: string; href: string };
  icon?: LucideIcon;
  tone?: SlabTone;
}) {
  const onFilled = tone === "brand" || tone === "inverse";

  return (
    <Slab tone={tone}>
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0">
          <p
            className={cn(
              "text-[11px] font-bold uppercase leading-[1.2] tracking-[0.12em]",
              onFilled ? "text-white/70" : "text-app-ink-muted",
            )}
          >
            {eyebrow}
          </p>

          <div className="mt-2 flex flex-wrap items-baseline gap-x-3 gap-y-1">
            <span className="text-[30px] font-bold tabular-nums leading-none tracking-[-0.02em] lg:text-[40px]">
              {formatMoney(amount, sign, currency, 0, locale)}
            </span>

            {delta ? (
              <span
                className={cn(
                  "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[12px] font-bold",
                  onFilled
                    ? "bg-white/14 text-navy-200"
                    : delta.direction === "down"
                      ? "bg-negative-soft text-negative"
                      : "bg-positive-soft text-positive",
                )}
              >
                <span aria-hidden>
                  {delta.direction === "down" ? "▼" : delta.direction === "up" ? "▲" : "—"}
                </span>
                {delta.direction === "down" ? "-" : "+"}
                {Math.abs(delta.percent)}%
              </span>
            ) : null}
          </div>

          {delta ? (
            <p
              className={cn(
                "mt-1.5 text-[12px] font-medium",
                onFilled ? "text-white/60" : "text-app-ink-muted",
              )}
            >
              {delta.period}
            </p>
          ) : null}
        </div>

        {Icon ? (
          <span
            className={cn(
              "grid size-11 shrink-0 place-items-center rounded-[14px]",
              onFilled ? "bg-white/14" : "bg-app-surface-sunken",
            )}
          >
            <Icon size={20} strokeWidth={1.75} />
          </span>
        ) : null}
      </div>

      {caption || action ? (
        <div
          className={cn(
            "mt-[22px] flex items-center justify-between gap-4 border-t pt-4 text-[12px] font-medium",
            onFilled
              ? "border-white/18 text-white/72"
              : "border-app-line text-app-ink-muted",
          )}
        >
          <span className="min-w-0 truncate">{caption}</span>
          {action ? (
            <a
              href={action.href}
              className={cn(
                "inline-flex shrink-0 items-center gap-1 rounded-control font-bold transition-colors",
                onFilled
                  ? "text-white hover:text-navy-200"
                  : "text-accent-text hover:text-accent-text-hover",
                focusRing,
              )}
            >
              {action.label}
              <ChevronRight size={14} strokeWidth={2} />
            </a>
          ) : null}
        </div>
      ) : null}
    </Slab>
  );
}

export function StatCard({
  eyebrow,
  amount,
  sign = "none",
  currency = "PKR",
  locale = "en-PK",
  caption,
  icon: Icon,
  tone = "neutral",
}: {
  eyebrow: string;
  amount: number;
  sign?: MoneySign;
  currency?: string;
  locale?: string;
  caption?: string;
  icon: LucideIcon;
  tone?: SemanticTone;
}) {
  return (
    <div className="flex flex-col justify-between rounded-slab border-[1.5px] border-app-line-strong bg-app-surface p-[var(--slab-pad)] shadow-slab">
      <span
        className={cn(
          "grid size-[var(--icon-tile)] place-items-center rounded-[10px]",
          semanticFillClass[tone],
        )}
      >
        <Icon size={18} strokeWidth={1.75} />
      </span>
      <div className="mt-4">
        <Eyebrow>{eyebrow}</Eyebrow>
        <p className="mt-1.5 text-[22px] font-bold tabular-nums leading-none tracking-[-0.015em] text-app-ink">
          {formatMoney(amount, sign, currency, 0, locale)}
        </p>
        {caption ? (
          <p className="mt-1.5 text-[12px] font-medium text-app-ink-muted">{caption}</p>
        ) : null}
      </div>
    </div>
  );
}

/** Ring plus a text readout. Never a bare ring. */
export function ProgressRing({ percent, label }: { percent: number; label: string }) {
  const clamped = Math.max(0, Math.min(100, Math.round(percent)));

  return (
    <div
      role="img"
      aria-label={`${clamped}% ${label}`}
      className="flex size-17 shrink-0 flex-col items-center justify-center rounded-full border-[7px] border-attention-line"
    >
      <span className="text-[13px] font-bold tabular-nums leading-none text-app-ink">
        {clamped}%
      </span>
      <span className="mt-0.5 text-[10px] font-medium leading-none text-app-ink-muted">
        {label}
      </span>
    </div>
  );
}

export function ListRow({
  title,
  detail,
  initials,
  amount,
  sign = "none",
  currency = "PKR",
  locale = "en-PK",
  onClick,
  trailing,
}: {
  title: string;
  detail?: string;
  initials?: string;
  amount?: number;
  sign?: MoneySign;
  currency?: string;
  locale?: string;
  onClick?: () => void;
  trailing?: ReactNode;
}) {
  const Wrapper = onClick ? "button" : "div";

  return (
    <Wrapper
      {...(onClick ? { type: "button" as const, onClick } : {})}
      className={cn(
        "flex w-full items-center gap-3 rounded-slab border border-app-line bg-app-surface text-left",
        "px-4 py-[var(--row-pad-y)]",
        // Hover changes the border, not the fill — a fill change across
        // 40 rows flickers.
        "transition-colors duration-200 ease-vertex hover:border-app-line-strong",
        onClick && focusRing,
      )}
    >
      {initials ? (
        <span className="grid size-10 shrink-0 place-items-center rounded-full bg-app-surface-sunken text-[13px] font-bold text-app-ink-soft">
          {initials}
        </span>
      ) : null}

      <span className="min-w-0 flex-1">
        <span className="block truncate text-[var(--text-row)] font-bold text-app-ink">
          {title}
        </span>
        {detail ? (
          <span className="mt-0.5 block truncate text-[var(--text-detail)] font-medium text-app-ink-muted">
            {detail}
          </span>
        ) : null}
      </span>

      {typeof amount === "number" ? (
        <Money
          amount={amount}
          sign={sign}
          currency={currency}
          locale={locale}
          className="shrink-0 text-[var(--text-row)]"
        />
      ) : null}
      {trailing}
    </Wrapper>
  );
}

/* ==========================================================================
   Segmented switch
   ========================================================================== */

export function SegmentedSwitch<T extends string>({
  options,
  value,
  onChange,
  label,
}: {
  options: { id: T; label: string; icon?: LucideIcon }[];
  value: T;
  onChange: (next: T) => void;
  label: string;
}) {
  const reduce = useReducedMotion();
  const groupId = useId();

  return (
    <div
      role="group"
      aria-label={label}
      className="inline-flex w-max gap-1 rounded-[14px] bg-app-surface-sunken p-1"
    >
      {options.map((option) => {
        const active = option.id === value;
        const Icon = option.icon;

        return (
          <button
            key={option.id}
            type="button"
            aria-pressed={active}
            onClick={() => onChange(option.id)}
            className={cn(
              "relative inline-flex min-h-[var(--control-h)] items-center gap-2 rounded-[11px] px-4 text-[13px] font-bold",
              "transition-colors duration-200 ease-vertex",
              focusRing,
              active ? "text-accent-text" : "text-app-ink-soft hover:text-app-ink",
            )}
          >
            {active ? (
              <motion.span
                aria-hidden
                layoutId={reduce ? undefined : `vx-segment-${groupId}`}
                transition={{ duration: duration.swap, ease: easeVertex }}
                className="absolute inset-0 rounded-[11px] bg-app-surface shadow-slab-pressed"
              />
            ) : null}
            {Icon ? <Icon size={16} strokeWidth={1.75} className="relative z-1" /> : null}
            <span className="relative z-1">{option.label}</span>
          </button>
        );
      })}
    </div>
  );
}

/* ==========================================================================
   Form controls
   ========================================================================== */

const controlClass = cn(
  "w-full rounded-control border border-app-line bg-app-surface",
  "min-h-[var(--control-h)] px-4 py-3",
  // 16px minimum or iOS Safari zooms on focus.
  "text-base font-medium text-app-ink outline-none transition-colors",
  "placeholder:text-app-ink-muted hover:border-app-line-strong",
  "focus:border-brand-primary",
  focusRing,
  "disabled:cursor-not-allowed disabled:bg-app-surface-sunken disabled:text-app-ink-muted",
);

export function Field({
  id,
  label,
  required,
  hint,
  error,
  children,
}: {
  id: string;
  label: string;
  required?: boolean;
  hint?: string;
  error?: string;
  children: ReactNode;
}) {
  return (
    <div className="flex flex-col gap-2">
      <label htmlFor={id} className="text-[13px] font-bold tracking-[0.04em] text-app-ink">
        {label}
        {required ? (
          <span className="text-accent-text"> *</span>
        ) : (
          <span className="font-medium text-app-ink-muted"> · Optional</span>
        )}
      </label>
      {children}
      {hint && !error ? (
        <p id={`${id}-hint`} className="text-[12px] font-medium text-app-ink-muted">
          {hint}
        </p>
      ) : null}
      {error ? (
        <p id={`${id}-error`} role="alert" className="text-[12px] font-bold text-negative">
          {error}
        </p>
      ) : null}
    </div>
  );
}

export function TextInput({
  invalid,
  className,
  ...rest
}: InputHTMLAttributes<HTMLInputElement> & { invalid?: boolean }) {
  return (
    <input
      aria-invalid={invalid || undefined}
      className={cn(controlClass, invalid && "border-negative", className)}
      {...rest}
    />
  );
}

export function TextArea({
  invalid,
  className,
  rows = 3,
  ...rest
}: TextareaHTMLAttributes<HTMLTextAreaElement> & { invalid?: boolean }) {
  return (
    <textarea
      rows={rows}
      aria-invalid={invalid || undefined}
      className={cn(controlClass, "resize-none", invalid && "border-negative", className)}
      {...rest}
    />
  );
}

export function ChoiceButton({
  selected,
  onClick,
  children,
}: {
  selected: boolean;
  onClick: () => void;
  children: ReactNode;
}) {
  return (
    <button
      type="button"
      aria-pressed={selected}
      onClick={onClick}
      className={cn(
        "w-full rounded-control border px-4 py-3 text-left",
        "min-h-[var(--control-h)] text-[var(--text-body)] font-bold leading-[1.3]",
        "transition-colors duration-300 ease-vertex",
        focusRing,
        selected
          ? "border-brand-primary bg-brand-primary text-white"
          : "border-app-line bg-app-surface text-app-ink hover:border-brand-primary hover:text-accent-text-hover",
      )}
    >
      {children}
    </button>
  );
}

export function SearchRow({
  value,
  onChange,
  placeholder,
  label,
}: {
  value: string;
  onChange: (next: string) => void;
  placeholder: string;
  label: string;
}) {
  const id = useId();

  return (
    <div className="flex min-h-[var(--control-h)] items-center gap-2.5 rounded-[14px] border border-app-line bg-app-surface px-3.5 transition-colors focus-within:border-brand-primary">
      <Search size={18} strokeWidth={1.75} className="shrink-0 text-app-ink-muted" />
      <label htmlFor={id} className="sr-only">
        {label}
      </label>
      <input
        id={id}
        type="search"
        value={value}
        placeholder={placeholder}
        onChange={(event) => onChange(event.target.value)}
        className="min-w-0 flex-1 bg-transparent text-base font-medium text-app-ink outline-none placeholder:text-app-ink-muted"
      />
      {value ? (
        <button
          type="button"
          onClick={() => onChange("")}
          aria-label="Clear search"
          className={cn(
            "grid size-8 shrink-0 place-items-center rounded-full text-app-ink-muted",
            "transition-colors hover:text-app-ink",
            focusRing,
          )}
        >
          <X size={16} strokeWidth={2} />
        </button>
      ) : null}
    </div>
  );
}

export function ToggleSwitch({
  checked,
  onChange,
  label,
}: {
  checked: boolean;
  onChange: (next: boolean) => void;
  label: string;
}) {
  const reduce = useReducedMotion();

  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      aria-label={label}
      onClick={() => onChange(!checked)}
      className={cn(
        "inline-flex h-6 w-11 shrink-0 items-center rounded-full p-0.5",
        "transition-colors duration-200 ease-vertex",
        focusRing,
        checked ? "justify-end bg-brand-primary" : "justify-start bg-app-surface-sunken",
      )}
    >
      <motion.span
        layout={!reduce}
        transition={{ duration: duration.micro, ease: easeVertex }}
        className="block size-5 rounded-full bg-app-surface shadow-slab-pressed"
      />
    </button>
  );
}

/* ==========================================================================
   Data table — the dashboard and console workhorse
   ========================================================================== */

export type SortDirection = "asc" | "desc";

export type Column<T> = {
  id: string;
  header: string;
  /** Right-aligns and applies tabular-nums. Never centre a data column. */
  numeric?: boolean;
  sortable?: boolean;
  width?: string;
  render: (row: T) => ReactNode;
};

export function DataTable<T>({
  columns,
  rows,
  getRowId,
  caption,
  sort,
  onSortChange,
  emptyState,
  loading = false,
  onRowClick,
}: {
  columns: Column<T>[];
  rows: T[];
  getRowId: (row: T) => string;
  /** Screen-reader description of the table's purpose. Required. */
  caption: string;
  sort?: { id: string; direction: SortDirection };
  onSortChange?: (id: string) => void;
  emptyState?: ReactNode;
  loading?: boolean;
  onRowClick?: (row: T) => void;
}) {
  if (loading) return <SkeletonRows rows={5} height={44} />;
  if (rows.length === 0) {
    return emptyState ?? <EmptyState title="Nothing here yet" body="No records match." />;
  }

  return (
    <div className="overflow-x-auto rounded-slab border border-app-line bg-app-surface">
      <table className="w-full border-collapse text-left">
        <caption className="sr-only">{caption}</caption>
        <thead className="sticky top-0 z-1 bg-app-surface-sunken">
          <tr>
            {columns.map((column) => {
              const isSorted = sort?.id === column.id;

              return (
                <th
                  key={column.id}
                  scope="col"
                  style={column.width ? { width: column.width } : undefined}
                  aria-sort={
                    !column.sortable
                      ? undefined
                      : isSorted
                        ? sort.direction === "asc"
                          ? "ascending"
                          : "descending"
                        : "none"
                  }
                  className={cn(
                    "border-b border-app-line px-4 py-[var(--row-pad-y)]",
                    "text-[11px] font-bold uppercase tracking-[0.12em] text-app-ink-muted",
                    column.numeric && "text-right",
                  )}
                >
                  {column.sortable && onSortChange ? (
                    <button
                      type="button"
                      onClick={() => onSortChange(column.id)}
                      className={cn(
                        "inline-flex items-center gap-1 rounded-control transition-colors hover:text-app-ink",
                        isSorted && "text-accent-text",
                        focusRing,
                      )}
                    >
                      {column.header}
                      <ChevronsUpDown size={12} strokeWidth={2} aria-hidden />
                    </button>
                  ) : (
                    column.header
                  )}
                </th>
              );
            })}
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr
              key={getRowId(row)}
              onClick={onRowClick ? () => onRowClick(row) : undefined}
              className={cn(
                "border-b border-app-line last:border-0",
                // No zebra striping — the rules already do that job.
                "transition-colors hover:bg-app-surface-hover",
                onRowClick && "cursor-pointer",
              )}
            >
              {columns.map((column) => (
                <td
                  key={column.id}
                  className={cn(
                    "px-4 py-[var(--row-pad-y)] text-[var(--text-body)] font-medium text-app-ink",
                    column.numeric && "text-right tabular-nums",
                  )}
                >
                  {column.render(row)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

/* ==========================================================================
   Navigation
   ========================================================================== */

export function BottomNav<T extends string>({
  tabs,
  active,
  onChange,
}: {
  /** 3-5 tabs. A sixth means the information architecture is wrong. */
  tabs: { id: T; label: string; icon: LucideIcon }[];
  active: T;
  onChange: (next: T) => void;
}) {
  return (
    <nav
      aria-label="Primary"
      className={cn(
        "fixed inset-x-0 bottom-0 z-10 flex h-[var(--bottom-nav-height)] justify-center",
        "gap-[clamp(25px,8vw,72px)] border-t border-app-line bg-app-surface/96 px-4 backdrop-blur-md",
        "min-[800px]:left-1/2 min-[800px]:right-auto min-[800px]:w-[600px]",
        "min-[800px]:-translate-x-1/2 min-[800px]:rounded-t-[20px]",
        "min-[800px]:border min-[800px]:border-b-0 min-[800px]:border-app-line",
      )}
    >
      {tabs.map((tab) => {
        const isActive = tab.id === active;
        const Icon = tab.icon;

        return (
          <button
            key={tab.id}
            type="button"
            onClick={() => onChange(tab.id)}
            aria-current={isActive ? "page" : undefined}
            className={cn(
              "flex min-w-11 flex-col items-center justify-center gap-1",
              "text-[10px] font-bold tracking-[0.02em]",
              "transition-colors duration-200 ease-vertex",
              focusRing,
              isActive ? "text-accent-text" : "text-app-ink-muted hover:text-app-ink",
            )}
          >
            <Icon size={20} strokeWidth={1.75} />
            {tab.label}
            <span
              aria-hidden
              className={cn(
                "size-1 rounded-full",
                isActive ? "bg-accent-text" : "bg-transparent",
              )}
            />
          </button>
        );
      })}
    </nav>
  );
}

export function SidebarNav<T extends string>({
  groups,
  active,
  onChange,
  collapsed = false,
}: {
  groups: {
    label?: string;
    items: { id: T; label: string; icon: LucideIcon }[];
  }[];
  active: T;
  onChange: (next: T) => void;
  collapsed?: boolean;
}) {
  return (
    <nav
      aria-label="Primary"
      className={cn(
        "flex h-full shrink-0 flex-col gap-6 border-r border-app-line bg-app-surface py-4",
        collapsed ? "w-[var(--sidebar-width-collapsed)] px-2" : "w-[var(--sidebar-width)] px-3",
      )}
    >
      {groups.map((group, groupIndex) => (
        <div key={group.label ?? groupIndex}>
          {group.label && !collapsed ? (
            <Eyebrow className="px-3 pb-2">{group.label}</Eyebrow>
          ) : null}
          <ul className="flex flex-col gap-1">
            {group.items.map((item) => {
              const isActive = item.id === active;
              const Icon = item.icon;

              return (
                <li key={item.id}>
                  <button
                    type="button"
                    onClick={() => onChange(item.id)}
                    aria-current={isActive ? "page" : undefined}
                    title={collapsed ? item.label : undefined}
                    className={cn(
                      "relative flex w-full items-center gap-3 rounded-control",
                      "min-h-[var(--control-h)] px-3 text-[var(--text-body)] font-bold",
                      "transition-colors duration-200 ease-vertex",
                      focusRing,
                      collapsed && "justify-center px-0",
                      isActive
                        ? "bg-app-surface-hover text-accent-text"
                        : "text-app-ink-soft hover:bg-app-surface-hover hover:text-app-ink",
                    )}
                  >
                    {isActive ? (
                      <span
                        aria-hidden
                        className="absolute inset-y-1.5 left-0 w-0.5 rounded-full bg-accent-text"
                      />
                    ) : null}
                    <Icon size={18} strokeWidth={1.75} className="shrink-0" />
                    {collapsed ? (
                      <span className="sr-only">{item.label}</span>
                    ) : (
                      <span className="truncate">{item.label}</span>
                    )}
                  </button>
                </li>
              );
            })}
          </ul>
        </div>
      ))}
    </nav>
  );
}

/* ==========================================================================
   Bottom sheet
   ========================================================================== */

export function Sheet({
  open,
  onClose,
  title,
  children,
}: {
  open: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
}) {
  const reduce = useReducedMotion();
  const panelRef = useRef<HTMLDivElement>(null);
  const titleId = useId();

  const trapFocus = useCallback((event: KeyboardEvent) => {
    const panel = panelRef.current;
    if (!panel || event.key !== "Tab") return;

    const focusable = panel.querySelectorAll<HTMLElement>(
      'a[href], button:not([disabled]), input:not([disabled]), textarea:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])',
    );
    if (focusable.length === 0) return;

    const first = focusable[0];
    const last = focusable[focusable.length - 1];

    if (event.shiftKey && document.activeElement === first) {
      event.preventDefault();
      last.focus();
    } else if (!event.shiftKey && document.activeElement === last) {
      event.preventDefault();
      first.focus();
    }
  }, []);

  useEffect(() => {
    if (!open) return;

    const previouslyFocused = document.activeElement as HTMLElement | null;
    document.body.style.overflow = "hidden";
    panelRef.current?.focus();

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        event.preventDefault();
        onClose();
        return;
      }
      trapFocus(event);
    }

    document.addEventListener("keydown", onKeyDown);

    return () => {
      document.removeEventListener("keydown", onKeyDown);
      document.body.style.overflow = "";
      previouslyFocused?.focus();
    };
  }, [open, onClose, trapFocus]);

  return (
    <AnimatePresence>
      {open ? (
        <>
          <motion.div
            aria-hidden
            onClick={onClose}
            initial={reduce ? false : { opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={reduce ? undefined : { opacity: 0 }}
            transition={{ duration: duration.swap, ease: easeVertex }}
            className="fixed inset-0 z-40 bg-scrim"
          />
          <motion.div
            ref={panelRef}
            role="dialog"
            aria-modal="true"
            aria-labelledby={titleId}
            tabIndex={-1}
            initial={reduce ? false : { y: "100%" }}
            animate={{ y: "0%" }}
            exit={reduce ? undefined : { y: "100%" }}
            transition={{ duration: duration.panel, ease: easeVertex }}
            className={cn(
              "fixed inset-x-0 bottom-0 z-50 mx-auto w-[min(100%,560px)]",
              "rounded-t-sheet bg-app-canvas px-[22px] pb-[34px] pt-3 shadow-sheet",
              "max-h-[92svh] overflow-y-auto outline-none",
            )}
          >
            <span
              aria-hidden
              className="mx-auto mb-6 block h-1 w-9.5 rounded-full bg-app-line-strong"
            />
            <h2
              id={titleId}
              className="text-[18px] font-bold tracking-tight text-app-ink md:text-[20px]"
            >
              {title}
            </h2>
            <div className="mt-5">{children}</div>
          </motion.div>
        </>
      ) : null}
    </AnimatePresence>
  );
}

/* ==========================================================================
   Feedback states
   ========================================================================== */

export function EmptyState({
  title,
  body,
  icon: Icon = Inbox,
  action,
}: {
  title: string;
  body: string;
  icon?: LucideIcon;
  /** Label names the action: "Add your first patient", not "Get started". */
  action?: { label: string; onClick: () => void };
}) {
  return (
    <div className="flex flex-col items-center gap-3 rounded-slab border border-dashed border-app-line-strong px-4 py-8 text-center">
      <span className="grid size-12 place-items-center rounded-full bg-app-surface-sunken text-app-ink-muted">
        <Icon size={22} strokeWidth={1.75} />
      </span>
      <p className="text-[var(--text-body)] font-bold text-app-ink">{title}</p>
      <p className="max-w-[40ch] text-[13px] font-medium leading-[1.45] text-app-ink-muted">
        {body}
      </p>
      {action ? (
        <Button className="mt-2" onClick={action.onClick}>
          {action.label}
        </Button>
      ) : null}
    </div>
  );
}

export function ErrorState({
  title,
  body,
  onRetry,
  retryLabel = "Try again",
}: {
  title: string;
  /** Names the fix in plain language. No stack traces, no bare codes. */
  body: string;
  onRetry: () => void;
  retryLabel?: string;
}) {
  return (
    <div
      role="alert"
      className="flex flex-col gap-2 rounded-slab border border-negative-line bg-negative-soft px-4 py-3.5"
    >
      <p className="inline-flex items-center gap-2 text-[14px] font-bold text-negative">
        <TriangleAlert size={16} strokeWidth={2} />
        {title}
      </p>
      <p className="text-[13px] font-medium leading-[1.45] text-app-ink-soft">{body}</p>
      <button
        type="button"
        onClick={onRetry}
        className={cn(
          "mt-1 self-start rounded-control text-[13px] font-bold text-accent-text",
          "transition-colors hover:text-accent-text-hover",
          focusRing,
        )}
      >
        {retryLabel}
      </button>
    </div>
  );
}

/** Skeleton at the exact height of the row it replaces, so nothing shifts. */
export function SkeletonRows({
  rows = 3,
  height = 66,
}: {
  rows?: number;
  height?: number;
}) {
  return (
    <div aria-busy="true" className="flex flex-col gap-[var(--row-gap)]">
      {Array.from({ length: rows }, (_, index) => (
        <div
          key={index}
          aria-hidden
          style={{ height }}
          className="animate-pulse rounded-slab bg-app-surface-sunken"
        />
      ))}
    </div>
  );
}

/* ==========================================================================
   Sticky stack — `marketing` archetype only
   ========================================================================== */

const stickyOffsets = ["top-24", "top-28", "top-32", "top-36"] as const;

export function StickyStackSlab({
  index,
  isLast,
  tone = "brand",
  children,
}: {
  index: number;
  isLast: boolean;
  tone?: "brand" | "dark";
  children: ReactNode;
}) {
  const ref = useRef<HTMLDivElement>(null);
  const reduce = useReducedMotion();
  const { scrollYProgress } = useScroll({
    target: ref,
    offset: ["start 96px", "end start"],
  });

  const scale = useTransform(scrollYProgress, [0, 0.85], [1, isLast ? 1 : 0.94], {
    ease: easeVertexFn,
  });
  const overlay = useTransform(scrollYProgress, [0, 0.85], [0, isLast ? 0 : 0.42], {
    ease: easeVertexFn,
  });

  return (
    <div
      ref={ref}
      style={{ zIndex: index + 1 }}
      className={cn(
        "sticky h-[calc(100svh-9rem)] min-w-0 px-page pb-6",
        stickyOffsets[index] ?? "top-36",
      )}
    >
      <motion.article
        style={reduce ? undefined : { scale }}
        className={cn(
          "relative isolate flex h-full min-h-0 origin-top overflow-hidden rounded-4xl lg:rounded-card",
          tone === "brand" ? "bg-brand-primary" : "bg-surface-dark",
        )}
      >
        <div className="relative z-1 flex h-full min-h-0 w-full min-w-0 flex-col justify-between gap-8 overflow-hidden px-6 py-8 md:px-10 md:py-12 lg:flex-row lg:items-end lg:gap-16 lg:px-panel lg:py-16">
          {children}
        </div>
        {reduce ? null : (
          <motion.div
            aria-hidden
            style={{ opacity: overlay }}
            className="pointer-events-none absolute inset-0 z-2 rounded-4xl bg-black lg:rounded-card"
          />
        )}
      </motion.article>
    </div>
  );
}

/* ==========================================================================
   App shell — single-column layout for `mobile-app` and `form-flow`
   ========================================================================== */

export function AppShell({
  archetype = "mobile-app",
  density = "comfortable",
  children,
  bottomNav,
  fab,
}: {
  archetype?: SurfaceArchetype;
  density?: Density;
  children: ReactNode;
  bottomNav?: ReactNode;
  fab?: ReactNode;
}) {
  return (
    <Surface archetype={archetype} density={density}>
      <div className="mx-auto w-full max-w-app px-5 pb-[86px] pt-6">{children}</div>
      {fab}
      {bottomNav}
    </Surface>
  );
}

/* ==========================================================================
   Example — proves the parts compose, and shows the intended slab order
   ========================================================================== */

export function DashboardExample() {
  const [scope, setScope] = useState<"today" | "month">("today");
  const [sheetOpen, setSheetOpen] = useState(false);
  const [note, setNote] = useState("");

  return (
    <AppShell fab={<FabButton label="Add entry" onClick={() => setSheetOpen(true)} />}>
      <SegmentedSwitch
        label="Choose period"
        value={scope}
        onChange={setScope}
        options={[
          { id: "today", label: "Today" },
          { id: "month", label: "This month" },
        ]}
      />

      <div className="mt-[var(--stack-gap)]">
        <Eyebrow>Monday, 24 August</Eyebrow>
        <h1 className="mt-2 text-[25px] font-bold tracking-tight text-app-ink md:text-[34px]">
          Good morning, Hamza
        </h1>
      </div>

      <div className="mt-[var(--stack-gap)] flex flex-col gap-[var(--stack-gap)]">
        <MetricSlab
          eyebrow="Booked revenue"
          amount={184250}
          delta={{ percent: 4.8, direction: "up", period: "vs last week" }}
          caption="Across 3 chairs"
          action={{ label: "View schedule", href: "#schedule" }}
        />

        <div className="grid grid-cols-2 gap-3">
          <StatCard
            eyebrow="Collected"
            amount={82000}
            sign="in"
            icon={Plus}
            tone="positive"
            caption="14 payments"
          />
          <StatCard
            eyebrow="Outstanding"
            amount={46780}
            sign="out"
            icon={TriangleAlert}
            tone="negative"
            caption="6 invoices"
          />
        </div>

        <Slab tone="attention" className="flex items-center justify-between gap-4">
          <div className="min-w-0">
            <Eyebrow>Chair utilisation</Eyebrow>
            <p className="mt-2 text-[18px] font-bold text-app-ink">3 slots left today</p>
            <p className="mt-1 text-[13px] font-medium text-app-ink-soft">
              Two are back-to-back at 4pm.
            </p>
          </div>
          <ProgressRing percent={71} label="booked" />
        </Slab>

        <section>
          <div className="flex items-baseline justify-between gap-4">
            <h2 className="text-[21px] font-bold tracking-tight text-app-ink">
              Recent activity
            </h2>
            <TextLink href="#activity">View all</TextLink>
          </div>
          <div className="mt-4 flex flex-col gap-[var(--row-gap)]">
            <ListRow
              initials="AR"
              title="Ahmed Raza"
              detail="Today, 10:42 AM · Scale and polish"
              amount={8200}
              sign="in"
            />
            <ListRow
              initials="SK"
              title="Saad Khan"
              detail="Today, 11:15 AM · Refund issued"
              amount={5450}
              sign="out"
            />
          </div>
        </section>
      </div>

      <Sheet open={sheetOpen} onClose={() => setSheetOpen(false)} title="Add entry">
        <div className="flex flex-col gap-[var(--field-gap)]">
          <Field id="entry-note" label="What was it for">
            <TextArea
              id="entry-note"
              value={note}
              onChange={(event) => setNote(event.target.value)}
              placeholder="Scale and polish, two visits"
            />
          </Field>
          <Button onClick={() => setSheetOpen(false)}>Save entry</Button>
        </div>
      </Sheet>
    </AppShell>
  );
}
