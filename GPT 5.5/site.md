# Prompt: Bilingual Recipe Website in Docker

## Task

Create a complete production-ready bilingual recipe website application.

The application must support **Russian and English**, have animated UI, recipe images loaded from the internet, recipe pages, recipe categories, search, filters, responsive design, SEO metadata, and full containerized deployment.

Use current stable tools and best practices as of **June 18, 2026**. Prefer official stable releases, avoid deprecated APIs, and verify package versions from official documentation before implementation.

---

## Main Requirements

Build a full-stack web application for browsing recipes.

The website must include:

* Russian and English language versions
* A visible language switcher: `RU / EN`
* Recipe list pages
* Individual recipe detail pages
* Recipe category pages
* Pagination or page-based recipe separation
* Search by recipe title and ingredients
* Filters by cuisine, cooking time, difficulty, diet type, and meal type
* Animated page transitions and UI interactions
* Recipe cards with images from the internet
* Responsive layout for mobile, tablet, and desktop
* SEO-friendly pages
* Accessible UI
* Dockerized deployment
* One-command startup using Docker Compose

---

## Recommended Technology Stack

Use the latest stable versions available on the request date.

Preferred stack:

* **Next.js 16** with App Router
* **React 19.2**
* **TypeScript**
* **Node.js 24 LTS**
* **Tailwind CSS 4.3**
* **Motion / Framer Motion or native CSS animations**
* **next-intl** or an equivalent modern i18n solution
* **Zod** for schema validation
* **ESLint**
* **Prettier**
* **Docker**
* **Docker Compose Specification**
* **pnpm** via Corepack

Use `node:24-alpine` or another official Node.js 24 LTS image for the production container.

Do not use outdated routing patterns, deprecated Next.js APIs, old Tailwind configuration patterns, or obsolete Docker Compose syntax.

---

## Application Structure

Use a clean modular structure similar to:

```txt
recipe-app/
├── app/
│   ├── [locale]/
│   │   ├── page.tsx
│   │   ├── recipes/
│   │   │   ├── page.tsx
│   │   │   ├── page/
│   │   │   │   └── [pageNumber]/
│   │   │   │       └── page.tsx
│   │   │   └── [slug]/
│   │   │       └── page.tsx
│   │   ├── categories/
│   │   │   ├── page.tsx
│   │   │   └── [category]/
│   │   │       └── page.tsx
│   │   ├── cuisines/
│   │   │   ├── page.tsx
│   │   │   └── [cuisine]/
│   │   │       └── page.tsx
│   │   ├── favorites/
│   │   │   └── page.tsx
│   │   └── about/
│   │       └── page.tsx
│   ├── api/
│   │   └── health/
│   │       └── route.ts
│   ├── globals.css
│   └── layout.tsx
├── components/
│   ├── layout/
│   ├── recipes/
│   ├── filters/
│   ├── animations/
│   └── ui/
├── data/
│   └── recipes.ts
├── i18n/
│   ├── routing.ts
│   ├── request.ts
│   └── messages/
│       ├── en.json
│       └── ru.json
├── lib/
│   ├── recipes.ts
│   ├── seo.ts
│   ├── pagination.ts
│   └── validation.ts
├── public/
├── Dockerfile
├── docker-compose.yml
├── next.config.ts
├── package.json
├── tsconfig.json
├── .env.example
└── README.md
```

---

## Pages

Create the following pages.

### Home Page

Route examples:

```txt
/en
/ru
```

The home page must include:

* Hero section with animated food imagery
* Short intro text
* Featured recipes
* Popular cuisines
* Search input
* CTA button to browse all recipes
* Animated recipe cards
* Smooth hover states

English hero example:

```txt
Discover recipes from around the world.
Cook smarter, eat better, and explore new flavors every day.
```

Russian hero example:

```txt
Откройте рецепты со всего мира.
Готовьте проще, ешьте вкуснее и находите новые идеи каждый день.
```

---

### Recipes Listing Page

Route examples:

```txt
/en/recipes
/ru/recipes
/en/recipes/page/2
/ru/recipes/page/2
```

The recipes must be divided into pages.

Requirements:

* Show recipe cards in a responsive grid
* Add pagination
* Use URL-based page numbers
* Display 6–12 recipes per page
* Include sorting:

  * newest
  * fastest
  * easiest
  * highest rated
* Include filters:

  * cuisine
  * meal type
  * difficulty
  * diet
  * cooking time
* Include search by:

  * recipe title
  * ingredients
  * tags

Pagination UI must include:

* Previous page
* Next page
* Page numbers
* Disabled state
* Active page state

---

### Recipe Detail Page

Route examples:

```txt
/en/recipes/classic-carbonara
/ru/recipes/classic-carbonara
```

Each recipe detail page must include:

* Recipe title
* Localized title
* Large hero image
* Short description
* Cooking time
* Preparation time
* Total time
* Difficulty
* Servings
* Cuisine
* Meal type
* Diet tags
* Ingredients list
* Step-by-step instructions
* Nutrition block
* Tips section
* Related recipes
* Print recipe button
* Favorite button stored in localStorage
* Share button
* SEO metadata
* JSON-LD structured data using `Recipe` schema

The recipe page must be fully localized.

Do not machine-translate in the UI at runtime. Store Russian and English text explicitly in the recipe data.

---

### Category Pages

Route examples:

```txt
/en/categories
/ru/categories
/en/categories/desserts
/ru/categories/desserts
```

Categories should include:

* Breakfast
* Lunch
* Dinner
* Desserts
* Soups
* Salads
* Vegetarian
* Quick meals
* Baking
* Drinks

Each category page must show filtered recipes and localized category names.

---

### Cuisine Pages

Route examples:

```txt
/en/cuisines
/ru/cuisines
/en/cuisines/italian
/ru/cuisines/italian
```

Cuisine pages should include:

* Italian
* French
* Georgian
* Japanese
* Mexican
* Indian
* Mediterranean
* Ukrainian
* Latvian
* American

Each cuisine page must have a short localized intro and matching recipes.

---

### Favorites Page

Route examples:

```txt
/en/favorites
/ru/favorites
```

Requirements:

* Store favorites in `localStorage`
* Allow adding/removing recipes from favorites
* Show empty state if no favorites exist
* Keep UI localized

---

### About Page

Route examples:

```txt
/en/about
/ru/about
```

Include:

* About the project
* Cooking philosophy
* Data/source note
* Image attribution note
* Accessibility and responsive design note

---

## Internationalization

Implement proper bilingual routing.

Required locales:

```ts
const locales = ["en", "ru"] as const;
```

Default locale:

```ts
"en"
```

The language switcher must:

* Preserve the current page where possible
* Switch `/en/recipes/classic-carbonara` to `/ru/recipes/classic-carbonara`
* Switch `/ru/categories/desserts` to `/en/categories/desserts`
* Be accessible by keyboard
* Clearly show active language

All UI strings must be stored in translation files:

```txt
i18n/messages/en.json
i18n/messages/ru.json
```

Recipe data must contain localized fields:

```ts
type LocalizedText = {
  en: string;
  ru: string;
};

type Recipe = {
  id: string;
  slug: string;
  title: LocalizedText;
  description: LocalizedText;
  ingredients: {
    en: string[];
    ru: string[];
  };
  steps: {
    en: string[];
    ru: string[];
  };
};
```

---

## Recipe Data

Create at least **24 recipes** so pagination is meaningful.

Each recipe must include:

```ts
type Recipe = {
  id: string;
  slug: string;
  title: {
    en: string;
    ru: string;
  };
  description: {
    en: string;
    ru: string;
  };
  image: string;
  imageAlt: {
    en: string;
    ru: string;
  };
  cuisine: string;
  category: string;
  mealType: string[];
  diet: string[];
  difficulty: "easy" | "medium" | "hard";
  prepTimeMinutes: number;
  cookTimeMinutes: number;
  totalTimeMinutes: number;
  servings: number;
  rating: number;
  ingredients: {
    en: string[];
    ru: string[];
  };
  steps: {
    en: string[];
    ru: string[];
  };
  nutrition: {
    calories: number;
    protein: number;
    fat: number;
    carbs: number;
  };
  tips: {
    en: string[];
    ru: string[];
  };
  tags: string[];
};
```

Use varied recipes, for example:

* Classic Carbonara
* Margherita Pizza
* Georgian Khachapuri
* Chicken Ramen
* Vegetable Curry
* Caesar Salad
* Borscht
* Shakshuka
* Beef Tacos
* French Onion Soup
* Pancakes
* Apple Pie
* Greek Salad
* Pad Thai
* Mushroom Risotto
* Latvian Rye Bread Dessert
* Falafel Bowl
* Salmon Teriyaki
* Lentil Soup
* Tiramisu
* Ratatouille
* Hummus Plate
* Cheesecake
* Berry Smoothie

---

## Internet Images

Use recipe images from the internet.

Requirements:

* Use reliable remote image URLs
* Use `next/image`
* Configure `remotePatterns` in `next.config.ts`
* Add descriptive localized alt text
* Add loading states
* Add fallback image if a remote image fails
* Do not use copyrighted images without permission
* Prefer images from open or permissive sources such as Unsplash, Pexels, Wikimedia Commons, or a placeholder image service suitable for demos
* Include attribution notes where appropriate
* Make images visually consistent:

  * same aspect ratio
  * rounded corners
  * subtle hover zoom
  * skeleton loading state

Example image handling requirements:

```ts
images: {
  remotePatterns: [
    {
      protocol: "https",
      hostname: "images.unsplash.com"
    },
    {
      protocol: "https",
      hostname: "images.pexels.com"
    },
    {
      protocol: "https",
      hostname: "upload.wikimedia.org"
    }
  ]
}
```

---

## Design Requirements

Create a warm, modern, elegant food website.

Visual direction:

* Clean editorial layout
* Soft warm colors
* Cream / ivory background
* Accent colors inspired by herbs, tomato, olive oil, and spices
* Rounded cards
* Large appetizing images
* Smooth shadows
* Clear typography
* Strong whitespace
* Mobile-first layout

Suggested design tokens:

```txt
background: warm ivory
text: deep charcoal
primary: tomato red or paprika
secondary: olive green
accent: golden saffron
muted: warm beige
card: white or off-white
```

The site must look polished, not like a basic template.

---

## Animation Requirements

Add tasteful animations.

Use Motion, Framer Motion, or modern CSS animations.

Animate:

* Page transitions
* Recipe card entrance
* Hero section text
* Hero image
* Filter panel open/close
* Favorite button interaction
* Language switcher hover
* Pagination transitions
* Loading skeletons
* Card image hover zoom
* Button hover/tap states

Respect accessibility:

* Add `prefers-reduced-motion`
* Disable or reduce animations when reduced motion is enabled
* Do not create distracting or excessive animation
* Keep interactions smooth and subtle

---

## Components

Create reusable components.

Required components:

```txt
Header
Footer
LanguageSwitcher
ThemeToggle optional
SearchBar
RecipeCard
RecipeGrid
RecipeFilters
RecipePagination
RecipeHero
RecipeMeta
IngredientList
InstructionSteps
NutritionCard
RelatedRecipes
FavoriteButton
ShareButton
ImageWithFallback
AnimatedSection
EmptyState
Badge
Button
Input
Select
Skeleton
```

Components must be:

* Typed with TypeScript
* Accessible
* Responsive
* Reusable
* Cleanly separated from data logic

---

## Search and Filtering

Implement client-side search and filters for the local recipe dataset.

Search must support:

* Localized recipe title
* Localized description
* Ingredients
* Tags
* Cuisine
* Category

Filters must support:

* Cuisine
* Category
* Difficulty
* Diet
* Meal type
* Maximum total cooking time

Use URL query parameters for state where appropriate:

```txt
/en/recipes?search=salad&difficulty=easy&time=30
/ru/recipes?search=суп&diet=vegetarian
```

Search and filtering should work in both languages.

---

## SEO Requirements

Implement SEO for all pages.

Include:

* Localized page titles
* Localized descriptions
* Open Graph metadata
* Twitter card metadata
* Canonical URLs
* Alternate language links
* JSON-LD structured data for recipe pages
* Clean semantic headings
* Descriptive image alt text

Recipe detail pages must include valid `Recipe` JSON-LD with:

* name
* description
* image
* author
* datePublished
* prepTime
* cookTime
* totalTime
* recipeYield
* recipeCuisine
* recipeCategory
* recipeIngredient
* recipeInstructions
* nutrition

---

## Accessibility Requirements

The app must be accessible.

Requirements:

* Semantic HTML
* Keyboard navigation
* Visible focus states
* Accessible buttons and links
* Proper labels for inputs and selects
* Correct heading hierarchy
* Alt text for all images
* ARIA only where necessary
* Good color contrast
* Reduced motion support
* Screen-reader-friendly pagination
* Screen-reader-friendly language switcher

---

## Containerization Requirements

The entire application must run inside a container.

Create:

```txt
Dockerfile
docker-compose.yml
.dockerignore
.env.example
README.md
```

The app must start with:

```bash
docker compose up --build
```

The app must then be available at:

```txt
http://localhost:3000
```

### Dockerfile Requirements

Use a multi-stage Dockerfile.

Requirements:

* Use official Node.js 24 LTS image
* Use Alpine or slim image where appropriate
* Enable Corepack
* Install dependencies with lockfile
* Build the Next.js app
* Use production mode
* Run as non-root user
* Expose port `3000`
* Add healthcheck or provide a `/api/health` route
* Keep the final image small
* Avoid copying unnecessary files
* Use `.dockerignore`

Example structure:

```dockerfile
FROM node:24-alpine AS base
WORKDIR /app
RUN corepack enable

FROM base AS deps
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

FROM base AS builder
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN pnpm build

FROM node:24-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000

RUN addgroup -S nextjs && adduser -S nextjs -G nextjs

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

USER nextjs
EXPOSE 3000

CMD ["node", "server.js"]
```

Configure Next.js for standalone output:

```ts
const nextConfig = {
  output: "standalone"
};

export default nextConfig;
```

### Docker Compose Requirements

Use the current Compose Specification.

Do not use obsolete `version: "3"` syntax unless required by the environment.

Example:

```yaml
services:
  recipe-app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: recipe-app
    ports:
      - "3000:3000"
    environment:
      NODE_ENV: production
      NEXT_PUBLIC_SITE_URL: http://localhost:3000
    restart: unless-stopped
```

---

## Health Check

Add a health endpoint:

```txt
/api/health
```

It should return:

```json
{
  "status": "ok",
  "service": "recipe-app"
}
```

---

## README Requirements

Create a detailed `README.md`.

It must include:

* Project description
* Features
* Tech stack
* Local development instructions
* Docker startup instructions
* Environment variables
* Available routes
* How to add a recipe
* How to add a new language
* Image attribution note
* Production deployment notes

Include these commands:

```bash
pnpm install
pnpm dev
pnpm build
pnpm start
docker compose up --build
docker compose down
```

---

## Quality Requirements

The final project must:

* Build successfully
* Run in Docker
* Have no TypeScript errors
* Have no ESLint errors
* Be responsive
* Be accessible
* Have working language switching
* Have working recipe pagination
* Have working recipe detail pages
* Have working search and filters
* Have working internet images with fallback
* Have smooth animations
* Use clean component architecture
* Include realistic bilingual recipe content
* Include proper SEO metadata
* Include a complete README

---

## Acceptance Criteria

The work is complete only if:

1. `docker compose up --build` starts the app successfully.
2. The app is available at `http://localhost:3000`.
3. `/en` and `/ru` both work.
4. The language switcher changes the active locale without breaking routes.
5. `/en/recipes` and `/ru/recipes` show paginated recipe lists.
6. Individual recipe pages work in both languages.
7. Category and cuisine pages work.
8. Search works in English and Russian.
9. Filters work and can be combined.
10. Recipe images load from remote URLs.
11. Broken image URLs show a fallback.
12. Animations are smooth and respect `prefers-reduced-motion`.
13. The Docker image builds successfully.
14. The app runs as a non-root user in production.
15. The README explains how to run, build, and deploy the project.

---

## Final Output Format

Return the complete project code.

For every file, include:

```txt
/path/to/file
```

followed by a fenced code block with the file contents.

Do not skip important files.

Do not provide pseudocode.

Do not leave TODO comments for core functionality.

Make the result ready to copy into a real project and run immediately.
