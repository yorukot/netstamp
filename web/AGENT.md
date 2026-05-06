# Frontend Agent Notes

## Scope
- The frontend app is `@netstamp/web` in `web`; shared React UI is `@netstamp/ui` in `packages/ui` and is consumed through package exports, not TS path aliases.
- Use pnpm only; root `package.json` has `preinstall: npx only-allow pnpm`.
- Workspace config lists `web`, `docs`, and `packages/*`, but `pnpm-lock.yaml` still has old `apps/web` and `apps/docs` importers; avoid incidental `pnpm install`/lockfile churn unless changing dependencies.

## Commands
- Dev server: `pnpm dev:web` from the repo root, or `pnpm --filter @netstamp/web dev`.
- Lint frontend: `pnpm --filter @netstamp/web lint`.
- Typecheck frontend: `pnpm --filter @netstamp/web exec tsc -p tsconfig.json --noEmit`; there is no `typecheck` script.
- Build frontend: `pnpm build:web` from the repo root, or `pnpm --filter @netstamp/web build`.
- Typecheck shared UI after editing `packages/ui`: `pnpm --filter @netstamp/ui exec tsc -p tsconfig.json --noEmit`; that package has no lint/build scripts.
- No frontend test runner or test files are configured; use lint/typecheck/build for verification unless you add tests.
- Root `pnpm format` formats the whole repo with Prettier, import organization, tabs, semicolons, and double quotes; for small changes prefer `pnpm prettier <touched-files> --write`.
- The repo hook installed by root `prepare` is `.githooks/pre-push`; it only formats staged files and re-adds them, so run lint/typecheck yourself.

## App Wiring
- `web/src/main.tsx` mounts `App` inside `HelmetProvider`; route/page metadata uses `react-helmet-async`.
- Browser routing lives in `web/src/routes/AppRouter.tsx` and `web/src/routes/routePaths.ts`; dashboard routes render inside `AppShell`.
- Adding a route usually requires updating the `Route` unions in `web/src/shared/utils/mockData.ts`, `routePaths`, `AppRouter`, and `sidebarItems` if it belongs in the shell sidebar.
- `web/src/shared/hooks/useHashNavigation.ts` is currently unused legacy hash-routing code; the app uses `createBrowserRouter`.
- Dashboard data and route types are mock-driven from `web/src/shared/utils/mockData.ts`; auth is mocked in `web/src/features/auth/services/authService.ts` and returns `controller: "waiting-for-api"`.

## UI And Styles
- Global CSS imports shared tokens once via `@import "@netstamp/ui/styles";` in `web/src/index.css`; do not duplicate token imports in component CSS.
- Shared UI components live under `packages/ui/src/components/*` with colocated CSS modules and must be exported from `packages/ui/src/index.ts` for `@netstamp/ui` consumers.
- Preserve the existing visual system: `--ns-*` tokens, dark surfaces, orange accent, mono labels, square/cut-corner frames, and CSS modules per component/page.
- ECharts setup is centralized in `web/src/shared/components/ChartPanel.tsx`; option factories live in `web/src/shared/utils/chartOptions.ts`.
