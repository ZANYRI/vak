# Misen — bilingual recipe app

Misen is a production-oriented bilingual (English/Russian) recipe catalogue built with Next.js App Router. It contains 24 localized recipes, client-side filtering, URL pagination, accessible interactions, structured SEO metadata, remote image fallbacks, and a production Docker image.

## Features

- Locale-first routing: `/en/...` and `/ru/...`, with a route-preserving language switcher.
- 24 recipes with explicit English and Russian titles, descriptions, ingredients, steps and cook's notes.
- Search across localized titles, descriptions, ingredients and tags; composable cuisine, category, difficulty, diet, meal and time filters.
- URL-backed filtering and pagination (8 recipes a page), category pages and cuisine pages.
- Detail pages with nutrition, print, share, localStorage favourites, related recipes and Recipe JSON-LD.
- Responsive editorial design, keyboard focus states, semantic landmarks and `prefers-reduced-motion` support.
- Remote Unsplash images rendered with `next/image`; image errors fall back to an in-app decorative placeholder.
- SEO metadata: canonical URLs, alternate languages, Open Graph and Twitter cards.
- Health endpoint at `/api/health`.

## Technology

- Next.js 16, React 19, TypeScript and Tailwind CSS 4 (with a deliberately small custom CSS layer)
- Zod for query-string validation
- ESLint, Prettier, pnpm/Corepack
- Docker multi-stage build on the official Node.js 24 Alpine image

## Local development

Prerequisite: Node.js 24 LTS with Corepack enabled.

```bash
corepack enable
pnpm install
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000). The root path redirects to `/en`.

Quality commands:

```bash
pnpm lint
pnpm build
pnpm start
pnpm format:check
```

## Docker

Create a local environment file only if you need a public site URL other than the default:

```bash
cp .env.example .env
docker compose up --build
```

The application is available at [http://localhost:3000](http://localhost:3000). Stop the stack with:

```bash
docker compose down
```

The runner image only contains standalone Next output, runs as the unprivileged `nextjs` user, and uses `/api/health` for its health check.

## Environment variables

| Variable | Default | Purpose |
| --- | --- | --- |
| `NEXT_PUBLIC_SITE_URL` | `http://localhost:3000` | Canonical base URL used for SEO metadata. Set this to the final HTTPS site URL in production. |

## Routes

| Route | Purpose |
| --- | --- |
| `/{en,ru}` | Home page |
| `/{en,ru}/recipes` | Recipe list, search, filters and sorting |
| `/{en,ru}/recipes/page/2` | Paginated recipe list |
| `/{en,ru}/recipes/classic-carbonara` | Localized recipe detail page |
| `/{en,ru}/categories` | All recipe categories |
| `/{en,ru}/categories/desserts` | Filtered category |
| `/{en,ru}/cuisines` | All cuisines |
| `/{en,ru}/cuisines/italian` | Filtered cuisine |
| `/{en,ru}/favorites` | Browser-local saved recipes |
| `/{en,ru}/about` | Project, image and accessibility notes |
| `/api/health` | JSON liveness response |

## Add a recipe

1. Add a fully localized seed to `data/recipes.ts`. Include explicit English and Russian recipe text, a permissively licensed remote image URL, and searchable tags in both languages.
2. Keep `slug` stable because it is used in URLs and localStorage favourites.
3. Choose an existing `cuisine` and `category`, or add the key to `data/recipes.ts` and both translation files.
4. Run `pnpm lint` and `pnpm build` before shipping.

## Add a language

1. Add the locale to `i18n/config.ts`.
2. Copy `i18n/messages/en.json` into a new translation file with exactly the same keys.
3. Register it in `i18n/dictionary.ts`.
4. Add each recipe's explicit localized fields to its `title`, `description`, `ingredients`, `steps`, `tips` and `imageAlt` records.
5. Confirm metadata and language-switching URLs for the new locale.

## Image attribution and production notes

The demonstration uses remote images from Unsplash, selected as permissive demonstration imagery. Retain a creator/source record for every image before publishing a real catalogue, and replace any image that does not match your licensing policy.

For production, set `NEXT_PUBLIC_SITE_URL` to the public HTTPS origin, serve the image hosts over HTTPS, run the image through your licence review process, and place a CDN/reverse proxy in front of the container as appropriate.
