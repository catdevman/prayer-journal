# Technical Spec: Prayer Journal

## 1. Architecture Overview
* **Pattern:** Serverless Monorepo (Backend + Frontend + Infra).
* **Region:** `us-east-1`.
* **Repository Strategy:** Monorepo with strict separation of concerns.

## 2. Technology Stack

### Backend (Lambda)
* **Language:** Go 1.23+
* **Router:** `github.com/go-chi/chi/v5`
* **Lambda Proxy:** `github.com/awslabs/aws-lambda-go-api-proxy` (Chi adapter)
* **Authentication:** Auth0 (JWT validation via `github.com/auth0/go-jwt-middleware/v2`)
* **Logging:** `log/slog` (JSON structured)

### Persistence (DynamoDB)
* **Table Name:** `prayer-journal-data`
* **Partition Key (PK):** `pk` (String)
* **Sort Key (SK):** `sk` (String)
* **Design Pattern:** Single Table Design.

### Frontend (SPA)
* **Framework:** Vue 3 + TypeScript
* **Build Tool:** Vite
* **Auth:** `@auth0/auth0-vue`
* **Type Sync:** `tygo` (Generates TS interfaces from Go structs)

### Infrastructure (IaC)
* **Tool:** AWS CDK (Go)
* **Resources:**
    * **S3:** Web hosting bucket (block public access, OAC).
    * **CloudFront:** HTTPS distribution pointing to S3.
    * **Lambda:** Go binary (Arm64).
    * **API Gateway:** HTTP API (v2) integrated with Lambda.
    * **Route53:** Custom domain + ACM Certificate.

---

## 3. Data Model (DynamoDB)

### Entity: Prayer Item
* **PK:** `USER#{Auth0_Subject_ID}`
* **SK:** `PRAYER#{ULID}` (Sortable by time via ULID)
* **Attributes:**
    * `id` (String, ULID)
    * `title` (String)
    * `content` (String)
    * `is_answered` (Boolean)
    * `created_at` (ISO8601 String)

### Access Patterns
| Query | Key Condition | Filter |
| :--- | :--- | :--- |
| **Get User's Prayers** | `PK = USER#123` | None (Sorted by SK desc) |
| **Get Specific Prayer** | `PK = USER#123 AND SK = PRAYER#abc` | None |

---

## 4. API Contract (Go Interfaces)

**Type Sync Strategy:**
We will run `tygo generate` against `internal/models`.

```go
// internal/models/prayer.go

type Prayer struct {
    ID         string    `json:"id" dynamodbav:"id"`
    Title      string    `json:"title" dynamodbav:"title"`
    Content    string    `json:"content" dynamodbav:"content"`
    IsAnswered bool      `json:"is_answered" dynamodbav:"is_answered"`
    CreatedAt  time.Time `json:"created_at" dynamodbav:"created_at"`
}
