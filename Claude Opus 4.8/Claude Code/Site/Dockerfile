# syntax=docker/dockerfile:1

# ---------- Base ----------
FROM node:24-alpine AS base
WORKDIR /app
# Enable pnpm via Corepack (pinned by package.json "packageManager").
RUN corepack enable
ENV NEXT_TELEMETRY_DISABLED=1

# ---------- Dependencies ----------
FROM base AS deps
# pnpm-workspace.yaml carries the build-script allowlist (sharp, swc, …).
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

# ---------- Builder ----------
FROM base AS builder
ENV NODE_ENV=production
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN pnpm build

# ---------- Runner ----------
FROM node:24-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME=0.0.0.0
ENV NEXT_TELEMETRY_DISABLED=1

# Run as a non-root user.
RUN addgroup -S nextjs && adduser -S nextjs -G nextjs

# Copy the standalone server output and static assets only.
COPY --from=builder --chown=nextjs:nextjs /app/public ./public
COPY --from=builder --chown=nextjs:nextjs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nextjs /app/.next/static ./.next/static

USER nextjs
EXPOSE 3000

# Container-level healthcheck against the /api/health route.
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD node -e "fetch('http://127.0.0.1:3000/api/health').then(r=>{if(r.status!==200)process.exit(1)}).catch(()=>process.exit(1))"

CMD ["node", "server.js"]
