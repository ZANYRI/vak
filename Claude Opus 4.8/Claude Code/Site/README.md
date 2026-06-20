# 🍅 Savora — Bilingual Recipe Website

A production-ready, fully internationalized (English / Russian) recipe website
built with **Next.js 16 (App Router)**, **React 19**, **TypeScript**,
**Tailwind CSS 4**, **next-intl** and **Motion** — containerized for one-command
deployment with Docker Compose.

> Discover recipes from around the world. Cook smarter, eat better, and explore
> new flavors every day.

---

## ✨ Features

- **Bilingual (RU / EN)** with a fully accessible `RU / EN` language switcher that
  preserves the current page and query parameters.
- **24 hand-written bilingual recipes** — every title, description, ingredient,
  step and tip is authored in both languages (no runtime machine translation).
- **Recipe listing** with URL-based pagination (`/recipes/page/2`), sorting
  (newest / fastest / easiest / highest rated) and combinable filters
  (cuisine, category, meal type, difficulty, diet, max total time).
- **Search** across localized titles, descriptions, ingredients, tags, cuisine
  and category — works in both languages.
- **Recipe detail pages** with hero image, full metadata, ingredients,
  step-by-step instructions, nutrition, tips, related recipes, print / share /
  favorite actions and **`Recipe` JSON-LD** structured data.
- **Category & cuisine pages** with localized names and intros.
- **Favorites** stored in `localStorage` (with empty state), synced across tabs.
- **Tasteful animations** with Motion — page transitions, card entrances, hero
  text, filter panels, favorite button — all respecting `prefers-reduced-motion`.
- **Internet images** via `next/image` with skeleton loading states and a
  graceful fallback when a remote URL fails.
- **SEO**: localized titles/descriptions, Open Graph + Twitter cards, canonical
  URLs, `hreflang` alternates, `sitemap.xml`, `robots.txt`.
- **Accessible**: semantic HTML, keyboard navigation, visible focus, skip link,
  ARIA where needed, screen-reader-friendly pagination & language switcher.
- **Dockerized**: multi-stage build, non-root user, standalone output, health
  check — `docker compose up --build` and you're live on port 3000.

---

## 🧱 Tech stack

| Area            | Choice                                            |
| --------------- | ------------------------------------------------- |
| Framework       | Next.js 16 (App Router, Turbopack)                |
| UI library      | React 19.2                                        |
| Language        | TypeScript (strict)                               |
| Styling         | Tailwind CSS 4 (CSS-first `@theme` tokens)        |
| i18n            | next-intl 4                                       |
| Animation       | Motion 12 (`motion/react`)                        |
| Validation      | Zod 4                                             |
| Runtime         | Node.js 24 LTS                                     |
| Package manager | pnpm (via Corepack)                               |
| Tooling         | ESLint 9, Prettier 3                              |
| Container       | Docker + Docker Compose (Compose Specification)   |

---

## 🚀 Quick start (local)

> Requires **Node.js 24+**. pnpm is provided through Corepack.

```bash
# Enable pnpm (one-time)
corepack enable

# Install dependencies
pnpm install

# Start the dev server (http://localhost:3000)
pnpm dev

# Production build & run
pnpm build
pnpm start
```

Then open <http://localhost:3000> — you'll be redirected to `/en`.

### Other scripts

```bash
pnpm lint          # ESLint
pnpm typecheck     # tsc --noEmit
pnpm format        # Prettier write
pnpm format:check  # Prettier check
```

---

## 🐳 Run with Docker

The entire app runs in a container. One command builds and starts it:

```bash
docker compose up --build
```

The app is then available at <http://localhost:3000>.

To stop and remove the container:

```bash
docker compose down
```

The image:

- uses the official **`node:24-alpine`** base,
- enables **Corepack/pnpm** and installs from the lockfile,
- builds Next.js in **standalone** mode,
- runs as a **non-root** `nextjs` user,
- exposes port **3000** and ships a **healthcheck** hitting `/api/health`.

---

## 🔧 Environment variables

Copy `.env.example` to `.env` (or set them in your environment / compose file):

| Variable               | Default                 | Description                                        |
| ---------------------- | ----------------------- | -------------------------------------------------- |
| `NEXT_PUBLIC_SITE_URL` | `http://localhost:3000` | Public base URL for canonical & Open Graph links.  |
| `NODE_ENV`             | `development`           | `production` in the container.                     |
| `PORT`                 | `3000`                  | Port the server listens on.                        |

---

## 🗺️ Available routes

| Route                              | Description                              |
| ---------------------------------- | ---------------------------------------- |
| `/` → `/en`                        | Locale-detected redirect                 |
| `/en`, `/ru`                       | Home page                                |
| `/{locale}/recipes`                | Paginated, filterable recipe listing     |
| `/{locale}/recipes/page/{n}`       | Page _n_ of the listing                  |
| `/{locale}/recipes/{slug}`         | Recipe detail page (+ JSON-LD)           |
| `/{locale}/categories`             | All categories                           |
| `/{locale}/categories/{category}`  | Recipes in a category                    |
| `/{locale}/cuisines`               | All cuisines                             |
| `/{locale}/cuisines/{cuisine}`     | Recipes in a cuisine                     |
| `/{locale}/favorites`              | Saved recipes (localStorage)             |
| `/{locale}/about`                  | About the project                        |
| `/api/health`                      | `{"status":"ok","service":"recipe-app"}` |
| `/sitemap.xml`, `/robots.txt`      | SEO                                       |

Example query-param URLs (search + filters):

```text
/en/recipes?search=salad&difficulty=easy&maxTime=30
/ru/recipes?diet=vegetarian&cuisine=italian&sort=fastest
```

---

## 📁 Project structure

```text
recipe-app/
├── app/
│   ├── [locale]/
│   │   ├── layout.tsx              # <html>, providers, header/footer
│   │   ├── page.tsx                # Home
│   │   ├── recipes/
│   │   │   ├── page.tsx
│   │   │   ├── page/[pageNumber]/page.tsx
│   │   │   └── [slug]/page.tsx
│   │   ├── categories/[category]/…
│   │   ├── cuisines/[cuisine]/…
│   │   ├── favorites/page.tsx
│   │   ├── about/page.tsx
│   │   └── not-found.tsx
│   ├── api/health/route.ts
│   ├── globals.css                 # Tailwind 4 theme tokens
│   ├── sitemap.ts
│   └── robots.ts
├── components/                     # layout / recipes / ui / animations / providers
├── data/recipes.ts                 # 24 bilingual recipes
├── i18n/                           # routing, request, navigation, messages
│   └── messages/{en,ru}.json
├── lib/                            # recipes, seo, pagination, validation, types
├── public/
├── proxy.ts                        # next-intl request handling (Next 16 "proxy")
├── Dockerfile
├── docker-compose.yml
├── next.config.ts                  # standalone output + image remotePatterns
└── …
```

> Note: Next.js 16 renamed the `middleware` file convention to **`proxy`**, so the
> next-intl request handler lives in `proxy.ts`.

---

## ➕ How to add a recipe

1. Open `data/recipes.ts`.
2. Append a new object to the `recipes` array. Provide **both** `en` and `ru`
   text for `title`, `description`, `imageAlt`, `ingredients`, `steps`, `tips`.
3. Use existing `cuisine` / `category` / `mealType` / `diet` / `difficulty`
   values (see `lib/types.ts` & `lib/constants.ts`).
4. Use a unique kebab-case `slug` and an image URL whose host is allow-listed in
   `next.config.ts` (`images.unsplash.com`, `images.pexels.com`,
   `upload.wikimedia.org`).

The optional `recipeSchema` in `lib/validation.ts` documents the exact shape and
can be used to validate the dataset.

---

## 🌍 How to add a new language

1. Add the locale to `i18n/routing.ts` (e.g. `['en', 'ru', 'de']`).
2. Create `i18n/messages/{locale}.json` (copy `en.json` and translate).
3. Add the language's text to every recipe in `data/recipes.ts` — extend the
   `LocalizedText` / `LocalizedList` types in `lib/types.ts` accordingly.
4. That's it: routing, the language switcher, sitemap and `hreflang` alternates
   pick up the new locale automatically.

---

## 🖼️ Image attribution

Photographs are loaded from open / permissively licensed sources — **Unsplash**,
**Pexels** and **Wikimedia Commons**. All rights remain with their respective
photographers. A graceful fallback image is shown if a remote URL ever fails.
Do not use copyrighted images without permission.

---

## 📦 Production deployment notes

- The app builds to a **standalone** server (`output: 'standalone'`), so the
  runtime image only needs `.next/standalone`, `.next/static` and `public`.
- Set `NEXT_PUBLIC_SITE_URL` to your real domain so canonical URLs, Open Graph
  tags and the sitemap are correct.
- The container runs as a **non-root** user and exposes a `/api/health`
  endpoint suitable for liveness/readiness probes (used by the Docker
  `HEALTHCHECK` and the compose `healthcheck`).
- Recipe content is static, so most pages are pre-rendered at build time (SSG);
  the listing pages render on demand to honor search/filter query params.

---

## ✅ Quality

- Builds successfully (`pnpm build`) and runs in Docker.
- No TypeScript errors (`pnpm typecheck`) and no ESLint errors (`pnpm lint`).
- Responsive, accessible, reduced-motion aware.
