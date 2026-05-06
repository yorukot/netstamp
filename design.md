# Netstamp Design System

> Category: Network Observability / Developer Infrastructure Current frontend: dark engineering console, orange-accented probe fleet dashboard, cut-corner industrial panels, monospace operational UI.

This document describes the visual language already used by the Netstamp frontend. Keep new UI aligned with `packages/ui/src/styles/tokens.css`, `web/src/index.css`, `web/src/layouts/AppShell.module.css`, `web/src/features/dashboard/components/LandingPage.module.css`, and the reusable components in `packages/ui/src/components`.

## 1. Design Direction

Netstamp should feel like a network operations console rather than a generic SaaS landing page. The interface is precise, dark, gridded, technical, and built around distributed measurement probes.

Core atmosphere:

- Near-black operational canvas.
- Orange as the primary brand and interaction accent.
- Thin engineering grid lines and diagnostic overlays.
- Cut-corner panels, buttons, badges, fields, and data surfaces.
- Display typography for large claims and screen titles.
- Monospace labels for navigation, metadata, buttons, tables, and telemetry.
- Sparse layouts with strong hierarchy and high-contrast data blocks.
- Network/probe/route language instead of broad AI startup wording.

The product voice is infrastructure-grade: "See the network before it fails you", "Probe Fleet", "Network Operations Console", "Measurement origin", "Recent system events".

## 2. Design Tokens

Use the existing `--ns-*` tokens. Do not introduce parallel token names unless a new primitive is genuinely required.

### Fonts

- `--ns-font-sans`: `TASAOrbiter`, `TASA Orbiter`, sans-serif.
- `--ns-font-display`: `TASAExplorer`, `TASA Explorer`, sans-serif.
- `--ns-font-mono`: `JetBrainsMono`, `JetBrains Mono`, monospace.

### Backgrounds And Surfaces

- `--ns-bg`: `#030406` main app canvas.
- `--ns-bg-section`: `#07090d` section background.
- `--ns-bg-subtle`: `#0b0e13` subtle elevated area.
- `--ns-surface`: `#111318` base panel surface.
- `--ns-surface-raised`: `#181b20` stronger raised surface.
- `--ns-surface-deep`: `#07080b` deep panel surface.
- `--ns-glass-dark`: `rgba(15, 18, 23, 0.94)` sticky/nav dark glass.
- `--ns-glass-light`: `rgba(255, 255, 255, 0.025)` faint overlay.

### Text

- `--ns-text`: `#fff7ec` primary text.
- `--ns-text-muted`: `#ddd4c8` body and secondary text.
- `--ns-text-subtle`: `#b8b3aa` supporting copy.
- `--ns-text-low`: `#77736b` metadata and low-priority labels.
- `--ns-text-on-accent`: `#0b0704` text on orange controls.

### Accent And State

- `--ns-accent`: `#ff7a1a` primary CTA and active brand color.
- `--ns-accent-hover`: `#ff9a3d` hover/highlight orange.
- `--ns-accent-active`: `#ff5f00` pressed/active orange.
- `--ns-accent-muted`: `rgba(255, 122, 26, 0.16)` muted selected background.
- `--ns-accent-subtle`: `rgba(255, 122, 26, 0.22)` stronger selected background.
- `--ns-accent-border`: `rgba(255, 122, 26, 0.66)` active frame color.
- `--ns-accent-glow`: `rgba(255, 122, 26, 0.3)` restrained glow.
- `--ns-critical`: `#ff453a` destructive/error state.
- `--ns-warning`: `#ff9f0a` warning state.
- `--ns-success`: `#30d158` healthy/online state.
- `--ns-metal`: `#a8b4c7` neutral technical accent.

### Borders, Cuts, Shadows

- `--ns-border`: `rgba(255, 255, 255, 0.22)` default frame.
- `--ns-border-strong`: `rgba(255, 255, 255, 0.38)` hover/active frame.
- `--ns-border-faint`: `rgba(255, 255, 255, 0.13)` quiet dividers and nested frames.
- `--ns-cut-xs`: `0.375rem` badges and small tags.
- `--ns-cut-sm`: `0.5rem` buttons, nav items, small frames.
- `--ns-cut-md`: `0.75rem` tables and medium cards.
- `--ns-cut-lg`: `1rem` panels, footer, landing sections.
- `--ns-shadow-sm`: dark technical elevation.
- `--ns-shadow-md`: stronger dark elevation.
- `--ns-shadow-glow`: orange-focused glow for active brand marks.
- `--ns-transition`: `180ms cubic-bezier(0.2, 0.8, 0.2, 1)`.

All radii tokens are currently `0`; use cut corners instead of rounded corners.

## 3. Typography

Use typography to separate marketing, console structure, and machine-readable data.

- Display type uses `--ns-font-display` with large scale, tight line-height, and strong weight around `650-800`.
- Body type uses `--ns-font-sans`, `1rem-1.25rem`, muted text, and short technical sentences.
- Operational UI uses `--ns-font-mono`, uppercase labels, `0.6875rem-0.875rem`, and letter-spacing between `0.07em` and `0.14em`.
- Navigation and buttons are uppercase monospace with heavy weight.
- Tables use monospace for headers and cells to preserve console/data-table tone.

Recommended hierarchy:

- Landing hero `h1`: `clamp(2.4rem, 5.2vw, 6rem)`, line-height `0.9`.
- Large story headline: `clamp(3.25rem, 7vw, 8rem)`, line-height `0.9`.
- App screen title: `clamp(3rem, 6vw, 5.75rem)`, line-height `0.9`.
- Drawer title: `clamp(2rem, 8vw, 3.75rem)`, line-height `0.92`.
- Panel title: `1rem`, display font, weight `650`.
- Eyebrow/label: `0.6875rem`, mono, uppercase, orange or muted.

## 4. Color Usage Rules

Orange is the only strong brand accent. It should carry CTAs, selected states, probe highlights, network activity, scanlines, and important telemetry.

Use limited state colors only for operational meaning:

- Green only means healthy, online, connected, or success.
- Yellow only means warning, waiting, degraded, or pending.
- Red only means critical, destructive, invalid, or failed.
- Slate/metal only supports neutral chart baselines and inactive technical lines.

Avoid adding blue, purple, rainbow gradients, pastel colors, or glossy SaaS-style gradients. If a non-orange color is necessary, it must represent data state or a small technical signal.

## 5. Layout System

The layout is modular, grid-first, and asymmetric when useful.

### Global Grid Background

Most full-page surfaces use layered backgrounds:

- Orange large grid: `8rem 8rem`, `rgba(255, 122, 26, 0.08-0.09)`.
- White fine grid: `2rem 2rem`, `rgba(255, 255, 255, 0.045-0.06)`.
- Optional diagonal micro-pattern: `0.75rem-0.875rem` for panels, sidebar, maps, and footer.
- Optional radial orange glow for landing/CTA depth.

Use the grid as engineering structure, not decoration. Major panels should align to visible grid rhythm where possible.

### Landing Page

The current landing page uses:

- Sticky top nav with brand mark, GitHub link, and orange Deploy CTA.
- Split hero with copy on the left and CSS-generated global network animation on the right.
- Floating telemetry chip in the hero.
- Story section with a single oversized framed card.
- Three staggered feature cards on desktop; stacked/2-column on smaller screens.
- Final CTA panel with trust signals and orange glow.
- Full footer using framed grid cells.

Keep landing copy short, technical, and direct. Do not center everything into a generic marketing template.

### App Shell

The dashboard app uses:

- Two-column shell: `17rem` fixed sidebar and fluid content.
- Sticky full-height sidebar with brand, team selector, nav links, and user card.
- Content area with `3rem 1.5rem 1.5rem` padding and compact footer.
- Sidebar collapses to `6rem` around `58rem`, hiding secondary labels.
- Main app backgrounds retain orange/white grid layers.

### Product Pages

Product screens share:

- `ScreenHeader` for eyebrow, large title, copy, and actions.
- 1rem grid gaps.
- Two-column layouts for dashboards, results, insight, alerts, probe detail, settings, and team views.
- Single-column layout below `58rem-78rem` depending on content density.
- Panels as primary containers; nested cells for key-value data, event feeds, steps, timelines, and route diffs.

## 6. Component Language

### Buttons

Buttons live in `@netstamp/ui` and should remain rectangular with cut corners.

- Use mono uppercase text, heavy weight, and `0.06em` letter-spacing.
- Primary uses orange fill and dark text.
- Secondary uses dark surface and light text.
- Outline uses dark surface, muted text, and orange hover.
- Ghost is for low-priority sidebar/user actions.
- Danger is reserved for destructive actions and uses red technical patterning.
- Hover may translate upward by `-0.0625rem`; avoid bounce.

Sizes:

- `sm`: `0.75rem`, min-height `2rem`.
- `md`: `0.8125rem`, min-height `2.5rem`.
- `lg`: `0.875rem`, min-height `3rem`.
- `xl`: `0.9375rem`, min-height `3.5rem`.

### Panels

Panels are the default page container primitive.

- Use `tone="glass"` for default raised sections.
- Use `tone="matte"` for lower-contrast nested sections.
- Use `tone="deep"` for maps, terminals, and high-depth diagnostic blocks.
- Keep panel header structure: eyebrow, title, optional actions, separator.
- Use `--ns-cut-lg` and frame color from `--ns-border` unless nested.

### Badges

Badges are small operational tags.

- Use mono uppercase text and optional dot.
- Use `neutral`, `accent`, `success`, `warning`, `critical`, or `muted` only.
- Do not use badge tones for decorative color variety.

### Fields

Fields use a clipped `controlFrame` around input/select/textarea.

- Labels are mono uppercase.
- Control background is near black.
- Focus uses `--ns-accent-border` and subtle inset orange glow.
- Invalid state uses red border and low red glow.
- Select chevron is CSS-generated with orange triangles.

### Tables

Data tables should feel like controller output.

- Use monospace throughout.
- Minimum width may exceed mobile viewport and scroll horizontally.
- Sticky table header uses dark raised surface.
- Hover and selected rows use faint orange backgrounds.
- Keep row borders thin and low-contrast.

### Metric Cards

Metric cards are compact status summaries.

- Use large display numeric value.
- Use mono orange label.
- Optional badge shows category or state.
- Use bottom-right corner bracket as diagnostic ornament.

### Terminal

Terminal blocks are command surfaces.

- Use clipped shell, 1rem orange micro-grid, dark command background.
- Top bar contains status dots, title, and right-aligned meta.
- Body uses mono text, `#ffd9bd`, and generous line-height.

### Network Map And Telemetry Widgets

Network visuals should be abstract, map-like, and diagnostic.

- Use grid background and skewed frame overlays.
- Nodes are square/dot-based, not playful pins.
- Active nodes may use green only for online state.
- Scanline animation is acceptable when slow and subtle.
- FleetMatrix uses small square cells, orange for online, low-opacity for offline.

## 7. Cut-Corner Frames

Cut corners are a core identity element. Use one of the approved frame techniques whenever `clip-path` is used.

### Rule

Do not rely on `border: 1px` alone when a clipped rectangle has diagonal cut corners. The straight borders will not draw the diagonal edges and the frame will look broken. Either use the mask frame technique or add diagonal patch lines with `::before` and `::after`.

### Standard Polygon

Use this shape for the current top-left and bottom-right cuts:

```css
.frame {
	--cut: var(--ns-cut-sm);
	clip-path: polygon(var(--cut) 0, 100% 0, 100% calc(100% - var(--cut)), calc(100% - var(--cut)) 100%, 0 100%, 0 var(--cut));
}
```

### Diagonal Border Patch

Use this when the element already has `border: 1px solid var(--frame-color)` and needs the missing diagonal strokes filled. The diagonal length must be `sqrt(2) * cut`, implemented as `1.41421356`.

```css
.frame {
	--cut: var(--ns-cut-sm);
	--frame-color: var(--ns-border);
	position: relative;
	border: 1px solid var(--frame-color);
	clip-path: polygon(var(--cut) 0, 100% 0, 100% calc(100% - var(--cut)), calc(100% - var(--cut)) 100%, 0 100%, 0 var(--cut));
}

.frame::before,
.frame::after {
	content: "";
	position: absolute;
	width: calc(var(--cut) * 1.41421356);
	height: 1px;
	background: var(--frame-color);
	pointer-events: none;
}

.frame::before {
	top: var(--cut);
	left: 0;
	transform: rotate(-45deg);
	transform-origin: left center;
}

.frame::after {
	right: 0;
	bottom: var(--cut);
	transform: rotate(-45deg);
	transform-origin: right center;
}
```

This pattern is already used by controls such as `Button`, `Panel`, and `Field`. Reuse it for new clipped controls when a real CSS border is used.

### Mask Frame Alternative

Use the mask frame technique when a full clipped outline is easier than patching individual diagonal strokes. This works well for nested cards, badges, nav items, app shell frames, map labels, and other elements whose actual border is transparent.

```css
.frame {
	--frame-color: var(--ns-border);
	position: relative;
	border: 1px solid transparent;
	clip-path: polygon(var(--ns-cut-sm) 0, 100% 0, 100% calc(100% - var(--ns-cut-sm)), calc(100% - var(--ns-cut-sm)) 100%, 0 100%, 0 var(--ns-cut-sm));
}

.frame::before {
	content: "";
	position: absolute;
	inset: 0;
	z-index: 2;
	padding: 1px;
	background: var(--frame-color);
	clip-path: inherit;
	pointer-events: none;
	mask:
		linear-gradient(#000 0 0) content-box,
		linear-gradient(#000 0 0);
	mask-composite: exclude;
	-webkit-mask:
		linear-gradient(#000 0 0) content-box,
		linear-gradient(#000 0 0);
	-webkit-mask-composite: xor;
}
```

### Frame Checklist

- Use `position: relative` on clipped frames.
- Use `isolation: isolate` when pseudo-elements, shadows, and children overlap.
- Use `overflow: hidden` only when contents must be clipped; avoid clipping focus rings unless the ring is outside the frame.
- Use `--ns-frame-color` or `--frame-color`; do not hardcode repeated border colors unless local state requires it.
- Patch diagonal cuts with `::before` and `::after` or use the mask frame. Never leave clipped borders visually open.

## 8. Motion

Motion should be restrained and mechanical.

- Default transition: `--ns-transition` (`180ms`).
- Hover translation may be subtle: `translateY(-0.0625rem)`.
- Drawer entry uses `180ms ease-out`.
- Scanline can move slowly over `8s`.
- Hero network animation may use slow tilt, orbit, node blink, and packet route animations.
- Respect `prefers-reduced-motion: reduce`; turn off decorative looping animation.

Avoid springy, playful, elastic, or large parallax motion.

## 9. Data Visualization

Charts are transparent and embedded inside panels.

- Primary series: orange `#FF7A1A` / `#FF8F3D`.
- Secondary baseline: slate `#94A3B8` at reduced opacity.
- Tooltip background: `rgba(10,13,18,0.92)` with faint white border.
- Axis labels: muted mono, around `10px`.
- Grid/split lines must remain very low contrast.
- Area fills should fade to transparent; avoid multicolor chart palettes unless the data state requires it.

## 10. Copywriting

Netstamp copy should be direct, technical, and grounded in network measurement.

Use phrases like:

- "See the network before it fails you."
- "Open-source network observability from probes you control."
- "Measure latency, packet loss, DNS, and routes."
- "Grid and map views for distributed measurement agents."
- "Path hash changed from previous run."
- "Waiting for first heartbeat."

Avoid phrases like:

- "Unlock your potential."
- "Supercharge your workflow."
- "Beautifully simple."
- "AI-powered for everyone."
- "Seamless experience."

The interface should sound like a controller, not a lifestyle brand.

## 11. Accessibility And Responsiveness

- Keep `:focus-visible` outlines visible, usually `2px solid var(--ns-accent)` with `0.25rem` offset.
- Preserve semantic landmarks: nav, main, section, article, aside, table.
- Use `aria-hidden="true"` for decorative geometry and icons.
- Do not rely on color alone for important state; badges and labels should include text.
- Use `100svh` for full-height shells to behave correctly on mobile browsers.
- Collapse dense multi-column grids to one column at mobile breakpoints.
- Allow wide data tables to scroll horizontally instead of compressing columns until unreadable.
- Keep touch targets at least around `2rem-2.75rem` high for controls.

## 12. Anti-Patterns

Do not add:

- Large rounded cards or pill-shaped core controls.
- Pastel SaaS gradients, glassmorphism-heavy cards, or glossy blobs.
- Blue/purple/rainbow brand accents.
- Decorative emoji-style icons.
- Soft lifestyle imagery or stock photos.
- Centered generic landing sections that ignore the grid.
- Heavy shadows that make the UI feel like a consumer app.
- New global CSS systems outside CSS modules and existing `--ns-*` tokens.

## 13. Implementation Checklist

Before shipping new frontend UI:

- Does it use existing `--ns-*` tokens?
- Does it use the correct font family for display, body, and operational labels?
- Are cut-corner frames patched with diagonal pseudo-elements or a mask frame?
- Is orange the primary interactive accent?
- Are state colors used only for state meaning?
- Does the layout align with the 1rem gap/grid rhythm?
- Does it collapse cleanly on mobile?
- Are focus states visible?
- Are decorative animations disabled under reduced motion?
- Does the copy sound like network infrastructure rather than generic SaaS?
