# Prayer Journal Monorepo

A serverless full-stack application using **Go (Lambda)**, **Vue 3 (SPA)**, and **AWS CDK (Infrastructure as Code)**.

## 🏗 Architecture

* **Frontend:** Vue 3 + TypeScript (Vite Build) -> S3 + CloudFront.
* **Backend:** Go 1.23 -> AWS Lambda (Monolith pattern via `chi` adapter) -> API Gateway HTTP API.
* **Database:** DynamoDB (Single Table Design).
* **Auth:** Auth0 (JWT).

## 🚀 Prerequisites

Ensure you have the following installed before starting:

1.  **Go 1.23+**: [Download](https://go.dev/dl/)
2.  **Node.js 20+**: [Download](https://nodejs.org/)
3.  **AWS CLI**: Configured with `aws configure`.
4.  **AWS CDK**: `npm install -g aws-cdk`
5.  **Tygo** (For syncing Go structs to TS):
    ```bash
    go install [github.com/gzuidhof/tygo@latest](https://github.com/gzuidhof/tygo@latest)
    ```

## 📂 Project Structure

```text
.
├── cmd/api/            # Main Lambda entrypoint
├── internal/
│   ├── api/            # HTTP Handlers & Middleware
│   └── models/         # Shared Go Structs (Source of Truth)
├── web/                # Vue 3 Frontend
├── infra/              # AWS CDK (Go)
├── scaffold.sh         # Bootstrapping script
└── Makefile            # Build & Deploy commands
```

## 🛠 Setup & Installation

1.  **Bootstrap the Repo:**
    If you haven't run the scaffold script yet:
    ```bash
    ./scaffold.sh
    go mod tidy
    ```

2.  **Install Frontend Deps:**
    ```bash
    cd web
    npm install
    ```

3.  **Environment Variables:**
    Create a `.env` file in `web/` (Vite requires `VITE_` prefix):
    ```ini
    VITE_AUTH0_DOMAIN=your-tenant.us.auth0.com
    VITE_AUTH0_CLIENT_ID=your-client-id
    VITE_API_URL=[https://your-api-gateway-url.com](https://your-api-gateway-url.com) # You get this after first deploy
    ```

## 💻 Local Development

### 1. The "Hybrid" Workflow (Recommended)
Since AWS Lambda is hard to mock perfectly locally, we recommend deploying a **Dev Stack** for the backend and running the frontend locally.

1.  **Deploy Backend First:**
    ```bash
    make deploy
    ```
    *Capture the API Gateway URL from the CDK Output.*

2.  **Run Frontend Locally:**
    Update `web/.env` with the API URL from step 1.
    ```bash
    cd web
    npm run dev
    ```

### 2. Syncing Types (Backend -> Frontend)
Whenever you modify a struct in `internal/models`, sync the changes to TypeScript:

1.  Edit `internal/models/prayer.go`.
2.  Run the sync command:
    ```bash
    make gen
    ```
3.  The TypeScript interfaces in `web/src/types/` (or configured path) will be updated automatically.

## 📦 Deployment

We use a unified `Makefile` workflow to build the Go binary, build the Vite assets, and deploy the CDK stack.

```bash
make deploy
```

**What this does:**
1.  Compiles `cmd/api/main.go` to a Linux/ARM64 binary (`bootstrap`).
2.  Runs `npm run build` in `web/` to generate `dist/`.
3.  Runs `cdk deploy` in `infra/` to upload the binary and static assets.

## 🔧 Troubleshooting

* **Tygo not found:** Ensure `$(go env GOPATH)/bin` is in your `$PATH`.
* **CDK Bootstrap Error:** If this is your first time using CDK in this region, run:
    ```bash
    cd infra
    cdk bootstrap aws://<ACCOUNT_ID>/<REGION>
    ```
* **CORS Issues:** Ensure the `AllowedOrigins` in your Go-Chi CORS middleware matches your CloudFront URL (and `localhost:5173` for dev).
