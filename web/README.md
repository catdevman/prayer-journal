# Prayer Journal Frontend

A Single Page Application (SPA) built with **Vue 3**, **TypeScript**, **Vite**, and **Pinia**.

## ⚡ Quick Start

```bash
# Install dependencies
npm install

# Run development server
npm run dev
```

## 🔧 Configuration

Create a `.env` file in the root of this directory. These variables are exposed to the client via Vite.

```ini
# Auth0 Configuration
VITE_AUTH0_DOMAIN=your-tenant.us.auth0.com
VITE_AUTH0_CLIENT_ID=your-client-id
VITE_AUTH0_AUDIENCE=https://prayerapi.faithforge.academy

# Backend API
# Local: http://localhost:8080
# Prod: https://prayerapi.faithforge.academy
VITE_API_URL=http://localhost:8080
```

## 🏗 Architecture

### Auth0 Integration
We use the **Closed Door** pattern.
* `LoginView.vue`: Public landing page.
* `DashboardView.vue`: Protected by `authGuard`.
* **Router:** Configured in `src/router/index.ts` to automatically redirect unauthenticated users to the Universal Login page if they attempt to access protected routes.

### Type Synchronization (`src/types`)
**⚠️ DO NOT EDIT FILES IN `src/types/` MANUALLY.**

TypeScript interfaces are generated automatically from the Go backend structs using `tygo`.
To update types:
1. Modify the Go struct in `../internal/models`.
2. Run `make gen` from the repository root.

## 📦 Scripts

| Script | Description |
| :--- | :--- |
| `npm run dev` | Starts the dev server (usually localhost:5173). |
| `npm run build` | Compiles assets to `dist/` for S3 deployment. |
| `npm run type-check` | Runs `vue-tsc` to verify types without building. |
| `npm run lint` | Runs ESLint. |

## 🚀 Deployment

Deployment is handled by the root **AWS CDK** stack.
The `make deploy` command in the root directory will:
1. Run `npm run build` in this directory.
2. Upload the contents of `dist/` to the S3 Web Bucket.
3. Invalidate the CloudFront distribution.
