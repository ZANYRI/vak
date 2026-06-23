# Recipe Haven

A production-ready bilingual recipe website built with Next.js 16, React 19, TypeScript, Tailwind CSS 4, and `next-intl`.

## Features

- Bilingual routing and UI: English and Russian
- Locale switcher that preserves the current page
- Recipe list with pagination, sorting, and client-side filtering
- Recipe detail pages with SEO metadata and JSON-LD `Recipe` schema
- Category and cuisine pages
- Favorites stored in `localStorage`
- Responsive, accessible design with warm food-inspired styling
- Smooth animations that respect `prefers-reduced-motion`
- Remote recipe images from Unsplash with fallback UI
- Docker Compose one-command deployment

## Tech Stack

- Next.js 16 (App Router)
- React 19.2
- TypeScript 6
- Tailwind CSS 4.3
- next-intl 4
- Framer Motion
- Zod
- ESLint 9
- Prettier
- pnpm via Corepack
- Docker / Docker Compose

## Local Development

Clone the repository and install dependencies:

```bash
pnpm install
```

Run the development server:

```bash
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000).

For a production build locally:

```bash
pnpm build
pnpm start
```

## Docker Deployment

Start the application with Docker Compose:

```bash
docker compose up --build
```

The app is then available at [http://localhost:3000](http://localhost:3000).

To stop:

```bash
docker compose down
```

The Dockerfile is multi-stage and runs the final Next.js standalone server as a non-root user.

## Environment Variables

Copy `.env.example` to `.env` and adjust if needed:

```text
NEXT_PUBLIC_SITE_URL=http://localhost:3000
```

## Available Routes

- `/en`, `/ru` — Home page
- `/en/recipes`, `/ru/recipes` — Recipe list
- `/en/recipes/page/2`, `/ru/recipes/page/2` — Paginated recipe list
- `/en/recipes/<slug>`, `/ru/recipes/<slug>` — Recipe detail
- `/en/categories`, `/ru/categories` — Category index
- `/en/categories/<category>`, `/ru/categories/<category>` — Category detail
- `/en/cuisines`, `/ru/cuisines` — Cuisine index
- `/en/cuisines/<cuisine>`, `/ru/cuisines/<cuisine>` — Cuisine detail
- `/en/favorites`, `/ru/favorites` — Favorites
- `/en/about`, `/ru/about` — About
- `/api/health` — Health check

## How to Add a Recipe

1. Open `data/recipes.ts`.
2. Duplicate an existing recipe object and update all fields.
3. Make sure the `slug` is unique and URL-friendly.
4. Provide both `en` and `ru` localizations for title, description, ingredients, steps, tips, and image alt text.
5. Run `pnpm typecheck` and `pnpm build` to verify.

## How to Add a New Language

1. Add the new locale to `i18n/routing.ts`.
2. Create a new messages file under `i18n/messages/<locale>.json`.
3. Copy the structure from `en.json` and translate all keys.
4. Update the locale label in the `LanguageSwitcher` if needed.
5. Update `generateStaticParams` usages if you want to pre-render the new locale.

## Image Attribution

Recipe images are sourced from Unsplash using public photo URLs. Individual photographers retain their rights; please respect Unsplash license terms when reusing images. If a remote image fails to load, a styled fallback UI is shown.

## Production Deployment Notes

- Set `NEXT_PUBLIC_SITE_URL` to the public domain.
- Keep `NODE_ENV=production` in the container.
- The container exposes port `3000` and includes a health check against `/api/health`.

## Scripts

- `pnpm dev` — Start development server
- `pnpm build` — Build production app
- `pnpm start` — Start production server
- `pnpm lint` — Run ESLint
- `pnpm lint:fix` — Fix ESLint issues
- `pnpm typecheck` — Run TypeScript without emitting
- `pnpm format` — Format code with Prettier
- `pnpm format:check` — Check formatting
