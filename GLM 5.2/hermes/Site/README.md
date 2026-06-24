# Saveur — Bilingual Recipe Website

A production-ready bilingual (Russian / English) recipe website built with Next.js 16, React 19, TypeScript, and Tailwind CSS 4. It ships fully containerized with a one-command Docker Compose startup.

---

## Features

- **Bilingual** — full Russian and English support with explicit localized content (no runtime machine translation).
- **Animated UI** — page transitions, card entrances, hero animations, favorite-button interactions, and loading skeletons, all respecting `prefers-reduced-motion`.
- **Recipe browsing** — listing with pagination (6–12 per page), individual detail pages, categories, and cuisines.
- **Search & filters** — search by title, ingredients, and tags; filter by cuisine, category, difficulty, diet, meal type, and max cooking time. Filters combine and sync to the URL.
- **Sorting** — newest, fastest, easiest, highest rated.
- **Favorites** — saved in `localStorage`, with a localized empty state.
- **Internet images** — recipe images from Unsplash / Pexels / Wikimedia via `next/image`, with a graceful fallback when a remote image fails.
- **SEO** — localized titles, descriptions, Open Graph & Twitter cards, canonical URLs, alternate language links, and `Recipe` JSON-LD structured data on detail pages.
- **Accessible** — semantic HTML, keyboard navigation, visible focus, ARIA labels, screen-reader-friendly pagination and language switcher, reduced-motion support.
- **Responsive** — mobile-first layout for phone, tablet, and desktop.
- **Dockerized** — multi-stage `node:24-alpine` image, non-root user, healthcheck, standalone output.

---

## Tech stack

| Layer            | Tool                                   |
| ---------------- | -------------------------------------- |
| Framework        | Next.js 16 (App Router, standalone)    |
| UI               | React 19.2                             |
| Language         | TypeScript                             |
| Styling          | Tailwind CSS 4.3                       |
| Animation        | Motion (Framer Motion)                 |
| i18n             | next-intl                              |
| Validation       | Zod                                    |
| Tooling          | ESLint, Prettier, pnpm (Corepack)      |
| Containerization | Docker, Docker Compose                 |
| Runtime          | Node.js 24 LTS (Alpine)                |

---

## Quick start (Docker)

```bash
docker compose up --build
```

The app is then available at **http://localhost:3000**.

To stop it:

```bash
docker compose down
```

---

## Local development

Prerequisites: Node.js 24+ and pnpm (enabled via Corepack).

```bash
pnpm install
pnpm dev
```

Open http://localhost:3000.

Other scripts:

```bash
pnpm build      # production build
pnpm start      # run the production build
pnpm lint       # ESLint
pnpm typecheck  # TypeScript check
pnpm format     # Prettier write
```

---

## Environment variables

Copy `.env.example` to `.env.local` and adjust if needed:

| Variable               | Default                  | Description                                            |
| ---------------------- | ------------------------ | ----------------------------------------------------- |
| `NEXT_PUBLIC_SITE_URL` | `http://localhost:3000`  | Public site URL, used for canonical / OG absolute URLs |
| `NEXT_PUBLIC_IMAGE_HOSTS` | (informational)       | Comma-separated list of allowed image hosts           |

Remote image hosts are configured statically via `remotePatterns` in `next.config.ts`.

---

## Available routes

| Route                              | Description                          |
| ---------------------------------- | ------------------------------------ |
| `/en`, `/ru`                       | Home page (hero, featured, cuisines) |
| `/en/recipes`, `/ru/recipes`       | Paginated recipe listing + filters   |
| `/en/recipes/page/2`               | Listing, page 2                      |
| `/en/recipes/[slug]`               | Recipe detail page                   |
| `/en/categories`, `/ru/categories` | All categories                       |
| `/en/categories/[category]`        | Recipes in a category                |
| `/en/cuisines`, `/ru/cuisines`     | All cuisines                         |
| `/en/cuisines/[cuisine]`           | Recipes for a cuisine                |
| `/en/favorites`, `/ru/favorites`   | Saved favorites                      |
| `/en/about`, `/ru/about`           | About the project                    |
| `/api/health`                      | Health endpoint (`{status,service}`) |

The `RU / EN` switcher in the header preserves the current route when changing locale.

---

## How to add a recipe

1. Open `data/recipes.ts`.
2. Append a new object to the `recipes` array conforming to the `Recipe` type from `lib/types.ts`.
3. Provide both `en` and `ru` strings for every localized field (title, description, imageAlt, ingredients, steps, tips). Array lengths must match between languages.
4. Use taxonomy slugs defined in `lib/taxonomy.ts` for `cuisine`, `category`, `mealType`, and `diet`.
5. Use a remote image URL from an allowed host (Unsplash, Pexels, Wikimedia) and add descriptive alt text.
6. Run `pnpm build` to validate types and generate the static pages.

---

## How to add a new language

1. Add the locale code to `locales` in `i18n/routing.ts` (e.g. `"de"`).
2. Create `i18n/messages/de.json` by translating `en.json`.
3. Add the locale to the `messages` import resolver in `i18n/request.ts` (it already loads `${locale}.json`, so only the file is needed).
4. Add localized names for any new taxonomy entries in `lib/taxonomy.ts`.
5. The language switcher and middleware will pick up the new locale automatically.

---

## Image attribution

Recipe photographs are loaded from open and permissive sources — primarily **Unsplash**, **Pexels**, and **Wikimedia Commons**. When a remote image fails to load, a tasteful in-app fallback is rendered. Please respect the license of each source when reusing images outside this project.

---

## Production deployment notes

- The Dockerfile produces a minimal standalone image (no `node_modules` shipped in the final layer) running as a non-root `nextjs` user.
- `next.config.ts` sets `output: "standalone"` so the container runs `node server.js` directly.
- A `HEALTHCHECK` polls `/api/health`; the compose service mirrors it.
- Set `NEXT_PUBLIC_SITE_URL` to your production domain so canonical and Open Graph URLs are absolute.
- For platforms without Docker, `pnpm build && pnpm start` runs the same production server.
- All recipe pages are statically generated (`generateStaticParams`), giving fast cold starts and easy CDN caching.
