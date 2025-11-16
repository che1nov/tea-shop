# 📦 E-commerce Tea Shop Platform - Repository Description

**Documentation Date:** November 16, 2025  
**Project Status:** In Development  
**Architecture Type:** Microservices Architecture

---

## 🎯 Project Purpose

**E-commerce Tea Shop Platform** is a full-featured e-commerce platform built with microservices architecture principles. The project demonstrates modern approaches to developing distributed systems using Go for backend and React for frontend.

### Key Features:
- 🔐 User registration and authentication (JWT)
- 📦 Product catalog with search and filtering
- 🛒 Shopping cart
- 💳 Order processing and management
- 💰 Payment system integration
- 🚚 Delivery management
- 📧 Notification system
- 👨‍💼 Administrative panel

---

## 🏗️ System Architecture

### Overview Diagram
```
┌─────────────┐
│   Frontend  │ (React + TypeScript)
│   :5173     │
└──────┬──────┘
       │ HTTP/REST
       ▼
┌─────────────────┐
│  API Gateway    │ (Gin + gRPC Client)
│     :8080       │
└────────┬────────┘
         │ gRPC
    ┌────┴─────┬──────┬──────┬──────┐
    ▼          ▼      ▼      ▼      ▼
┌────────┐ ┌────┐ ┌────┐ ┌────┐ ┌────┐
│ Users  │ │Good│ │Ordr│ │Pay │ │Delv│
│ :8001  │ │8002│ │8003│ │8004│ │8005│
└────┬───┘ └─┬──┘ └─┬──┘ └─┬──┘ └─┬──┘
     │       │      │      │      │
     └───────┴──────┴──────┴──────┘
                    │
              ┌─────┴─────┐
              ▼           ▼
         ┌────────┐  ┌─────────┐
         │ Kafka  │  │Postgres │
         │  :9092 │  │  :5432  │
         └────┬───┘  └─────────┘
              │
              ▼
       ┌─────────────┐
       │ Notify Svc  │
       └─────────────┘
```

### Microservices

#### 1. **API Gateway** (:8080)
- **Role:** Single entry point for all client requests
- **Technologies:** Gin (HTTP Router), gRPC Client
- **Functions:**
  - HTTP → gRPC routing
  - JWT validation
  - CORS handling
  - Swagger documentation
  - Rate limiting (planned)

#### 2. **Users Service** (:8001)
- **Role:** User management and authentication
- **Database:** PostgreSQL (users_db)
- **Functions:**
  - User registration
  - Authentication (JWT)
  - Profile management
  - Roles (user, admin)

#### 3. **Goods Service** (:8002)
- **Role:** Product catalog management
- **Database:** PostgreSQL (goods_db)
- **Functions:**
  - CRUD operations for products
  - Stock management
  - Stock reservation
  - Search and filtering

#### 4. **Order Service** (:8003)
- **Role:** Order processing
- **Database:** PostgreSQL (orders_db)
- **Integrations:** Goods, Payment, Delivery services
- **Functions:**
  - Order creation
  - Status management
  - Stock reservation
  - Kafka event publishing

#### 5. **Payment Service** (:8004)
- **Role:** Payment processing
- **Database:** PostgreSQL (payments_db)
- **Functions:**
  - Payment creation
  - Transaction processing
  - Payment status management
  - Payment gateway integration (mock)

#### 6. **Delivery Service** (:8005)
- **Role:** Delivery management
- **Database:** PostgreSQL (deliveries_db)
- **Functions:**
  - Delivery creation
  - Status tracking
  - Delivery cost calculation
  - Courier management (planned)

#### 7. **Notify Service**
- **Role:** Notification delivery
- **Technologies:** Kafka Consumer
- **Functions:**
  - Email notifications
  - SMS notifications (planned)
  - Push notifications (planned)

---

## 🛠️ Technology Stack

### Backend

#### Programming Language
- **Go 1.25+**

#### Frameworks and Libraries
- **gRPC** - inter-service communication
- **Protocol Buffers** - data serialization
- **Gin** - HTTP router for API Gateway
- **database/sql** + **lib/pq** - PostgreSQL interaction
- **golang-jwt/jwt** - JWT tokens
- **bcrypt** - password hashing
- **Kafka Go Client** - Kafka integration

#### Databases
- **PostgreSQL 15** - primary storage (5 separate databases)
  - users_db
  - goods_db
  - orders_db
  - payments_db
  - deliveries_db

#### Message Queue
- **Apache Kafka** - asynchronous event processing
- **Zookeeper** - Kafka coordination

#### Monitoring and Observability
- **Prometheus** - metrics collection
- **Grafana** - metrics visualization
- **Jaeger** - distributed tracing
- **Zap Logger** - structured logging

### Frontend

#### Framework
- **React 18** - UI library
- **TypeScript** - type safety

#### Build Tools
- **Vite** - fast build tool
- **ESLint** - code linting

#### UI and Styling
- **Tailwind CSS** - utility-first CSS
- **Headless UI** - accessible components

#### State Management
- **Zustand** - state management
- **React Query** - server state

#### Routing
- **React Router v6** - client-side routing

#### HTTP Client
- **Axios** - HTTP requests

### DevOps

#### Containerization
- **Docker** - service containerization
- **Docker Compose** - local development orchestration

#### CI/CD (planned)
- GitHub Actions
- Docker Registry

---

## 📁 Repository Structure

```
ecommerce/
│
├── 📂 api-gateway/               # API Gateway service
│   ├── cmd/
│   │   └── main.go              # Entry point
│   ├── config/
│   │   └── config.go            # Configuration
│   ├── internal/
│   │   ├── handler/             # HTTP handlers
│   │   └── middleware/          # Middleware (auth, admin)
│   ├── docs/                    # Swagger documentation
│   ├── go.mod
│   └── README.md
│
├── 📂 users-service/            # Users service
│   ├── cmd/
│   │   └── main.go
│   ├── config/
│   ├── internal/
│   │   ├── handler/            # gRPC handlers
│   │   ├── service/            # Business logic
│   │   ├── repository/         # Database access
│   │   └── model/              # Data models
│   ├── migrations/             # SQL migrations
│   ├── go.mod
│   └── README.md
│
├── 📂 goods-service/            # Goods service
│   ├── cmd/
│   ├── config/
│   ├── internal/
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   └── model/
│   ├── migrations/
│   ├── go.mod
│   └── README.md
│
├── 📂 order-service/            # Order service
│   ├── cmd/
│   ├── config/
│   ├── internal/
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── model/
│   │   └── kafka/              # Kafka producer/consumer
│   ├── migrations/
│   ├── go.mod
│   └── README.md
│
├── 📂 payment-service/          # Payment service
│   ├── cmd/
│   ├── config/
│   ├── internal/
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   └── model/
│   ├── go.mod
│   └── README.md
│
├── 📂 delivery-service/         # Delivery service
│   ├── cmd/
│   ├── config/
│   ├── internal/
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   └── model/
│   ├── go.mod
│   └── README.md
│
├── 📂 notify-service/           # Notification service
│   ├── cmd/
│   ├── config/
│   ├── internal/
│   │   ├── kafka/              # Kafka consumer
│   │   └── service/            # Email sender
│   ├── go.mod
│   └── README.md
│
├── 📂 shared/                   # Shared components
│   ├── pb/                     # Protocol Buffers
│   │   ├── users.proto
│   │   ├── goods.proto
│   │   ├── orders.proto
│   │   ├── payments.proto
│   │   ├── delivery.proto
│   │   └── *.pb.go             # Generated files
│   ├── pkg/
│   │   ├── logger/             # Shared logger
│   │   ├── errors/             # Error handling
│   │   └── config/             # Config utilities
│   ├── models/                 # Shared models
│   ├── go.mod
│   └── README.md
│
├── 📂 frontend/                 # React frontend
│   ├── public/
│   ├── src/
│   │   ├── api/                # API clients
│   │   ├── components/         # React components
│   │   ├── pages/              # Pages
│   │   ├── store/              # Zustand stores
│   │   ├── types/              # TypeScript types
│   │   ├── utils/              # Utilities
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── package.json
│   ├── vite.config.ts
│   └── README.md
│
├── 📂 monitoring/               # Monitoring configuration
│   ├── prometheus.yml          # Prometheus config
│   ├── grafana-datasource.yml  # Grafana datasources
│   ├── dashboard-*.json        # Grafana dashboards
│   └── *.md                    # Documentation
│
├── 📂 scripts/                  # Utility scripts
│   ├── start_all_services.sh   # Start all services
│   ├── stop_all_services.sh    # Stop services
│   ├── check_services.sh       # Health checks
│   └── test_complete_workflow.sh
│
├── docker-compose.yml           # Infrastructure orchestration
├── go.mod                       # Root go.mod
├── README.md                    # Main documentation
├── ARCHITECTURE_REVIEW.md       # Architecture review
├── CONTRIBUTING.md              # Contributing guide
└── ORDER_FLOW.md                # Order flow description

```

---

## 🔄 Data Flows

### 1. User Registration
```
Frontend → API Gateway → Users Service → PostgreSQL
                                ↓
                             Kafka → Notify Service → Email
```

### 2. Order Creation
```
Frontend → API Gateway → Order Service
                            ↓
                    ┌───────┼───────┐
                    ▼       ▼       ▼
            Goods Service Payment Delivery
                 (gRPC)    (gRPC)   (gRPC)
                    ↓       ↓       ▼
                PostgreSQL PostgreSQL PostgreSQL
                    ↓
                  Kafka → Notify Service → Email
```

### 3. Administrative Operations
```
Frontend → API Gateway (JWT + Admin check)
              ↓
        Goods Service → PostgreSQL
              ↓
            Kafka → Notify Service
```

---

## 🔌 API Endpoints

### Public Endpoints

#### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

#### Goods
- `GET /api/v1/goods` - List goods (with pagination)
- `GET /api/v1/goods/:id` - Get good details

### Protected Endpoints (JWT required)

#### Users
- `GET /api/v1/users/me` - Get current user info

#### Orders
- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders/:id` - Get order details
- `GET /api/v1/orders` - Get order history

#### Payments
- `GET /api/v1/payments/:id` - Get payment info

#### Delivery
- `POST /api/v1/deliveries` - Create delivery
- `GET /api/v1/deliveries/:id` - Get delivery status

### Admin Endpoints (Admin role required)

#### Goods Management
- `POST /api/v1/admin/goods` - Create good
- `PUT /api/v1/admin/goods/:id` - Update good
- `DELETE /api/v1/admin/goods/:id` - Delete good

---

## 📊 Database Schemas

### Schema Overview

#### users_db
```sql
users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)
```

#### goods_db
```sql
goods (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)

stock_reservations (
    id SERIAL PRIMARY KEY,
    good_id INT REFERENCES goods(id),
    order_id INT NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP NOT NULL
)
```

#### orders_db
```sql
orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    items JSONB NOT NULL,
    status VARCHAR(50) NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    address TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)
```

#### payments_db
```sql
payments (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    method VARCHAR(50),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)
```

#### deliveries_db
```sql
deliveries (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    address TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    tracking_number VARCHAR(100),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)
```

---

## 🔐 Security

### Authentication and Authorization
- **JWT tokens** for user identification
- **Bcrypt** for password hashing
- **Role-based access control** (user, admin)

### Data Protection
- **Prepared statements** to prevent SQL injection
- **CORS** settings for access control
- **Rate limiting** (planned)

### Secrets Management
- Environment variables for sensitive data
- `.env` files (not committed to git)

---

## 📈 Monitoring and Metrics

### Prometheus Metrics
Each service exports metrics on port 900X:
- Users Service: `:9001/metrics`
- Goods Service: `:9002/metrics`
- Order Service: `:9003/metrics`
- Payment Service: `:9004/metrics`
- Delivery Service: `:9005/metrics`

### Available Metrics
- Request count
- Request duration
- Error rate
- Database connection pool stats
- gRPC call latency

### Grafana Dashboards
- **Go Services Dashboard** - overview of all Go services
- **Node Exporter Dashboard** - system metrics
- **Services Health** - health checks for all services

---

## 🧪 Testing

### Unit Tests
Each service contains unit tests for:
- Domain models
- Service logic
- Repository operations

### Integration Tests
Tests for interaction with:
- PostgreSQL
- gRPC clients
- Kafka

### Running Tests
```bash
# All tests in service
cd users-service
go test ./...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

---

## 🚀 Deployment

### Local Development

1. **Start infrastructure:**
```bash
docker-compose up -d
```

2. **Start microservices:**
```bash
./start_all_services.sh
```

3. **Start frontend:**
```bash
cd frontend
npm install
npm run dev
```

### Production (planned)
- **Kubernetes** for orchestration
- **Helm charts** for deployment
- **CI/CD** via GitHub Actions
- **Docker Registry** for images

---

## 📝 Documentation

### Main Documents
- `README.md` - project overview and quick start
- `ARCHITECTURE_REVIEW.md` - architecture analysis
- `CONTRIBUTING.md` - contribution guidelines
- `ORDER_FLOW.md` - order flow description
- `ADMIN_ROLES.md` - role management

### Service Documentation
Each service has its own `README.md` with:
- API endpoints
- Data structures
- Usage examples
- Troubleshooting

### API Documentation
- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- **Proto files:** `shared/pb/*.proto`

---

## 📊 Project Statistics

### Codebase Size
```
Backend (Go):     ~15,000 lines
Frontend (React): ~5,000 lines
Config/Scripts:   ~2,000 lines
Tests:            ~3,000 lines
────────────────────────────
Total:            ~25,000 lines
```

### Service Count
- **7 microservices** (including API Gateway)
- **5 PostgreSQL databases**
- **1 Kafka cluster**

### Technologies
- **2 programming languages** (Go, TypeScript)
- **8+ external dependencies** (PostgreSQL, Kafka, Prometheus, etc.)
- **15+ Go libraries**
- **20+ npm packages**

---

## 🎯 Roadmap

### In Development
- [ ] Improved error handling
- [ ] Caching layer (Redis)
- [ ] Advanced product filtering
- [ ] User wishlist

### Planned
- [ ] Kubernetes manifests
- [ ] CI/CD pipeline
- [ ] E2E tests
- [ ] GraphQL API (in addition to REST)
- [ ] WebSocket for real-time notifications
- [ ] Full-featured admin panel
- [ ] Internationalization (i18n)
- [ ] Mobile app (React Native)

### Optimization
- [ ] Refactoring to Clean Architecture
- [ ] Redis cache integration
- [ ] Database query optimization
- [ ] Index optimization
- [ ] Connection pooling improvements

---

## 🏛️ Architectural Principles

### Current Architecture
The project currently uses a **layered architecture** pattern:
```
Handler → Service → Repository → Database
```

### Target Architecture (Clean Architecture)
Migration to **Clean Architecture** is planned:
```
Controllers → Use Cases ← Interfaces ← Adapters
     ↓           ↓
    DTO        Domain
```

**Benefits:**
- ✅ Better testability
- ✅ Clear separation of concerns
- ✅ Business logic independence
- ✅ Industry best practices compliance

**See:** `ARCHITECTURE_REVIEW.md` for detailed analysis

---

## 🔧 Development Guidelines

### Adding a New Service

1. Create service structure based on existing services
2. Add proto definitions to `shared/pb/`
3. Rebuild proto files
4. Add service to `start_all_services.sh`
5. Update API Gateway for request proxying
6. Add Swagger annotations to new endpoints

### Rebuilding Proto Files

```bash
cd shared/pb
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       *.proto
```

### Updating Swagger Documentation

After changing API endpoints:

```bash
cd api-gateway
swag init -g cmd/main.go
```

Swagger UI will automatically update after restarting API Gateway.

---

## 🤝 Contributing

The project is open for improvements. See `CONTRIBUTING.md` for details.

### How to Help
1. Report bugs
2. Propose new features
3. Improve documentation
4. Write tests
5. Optimize code

### Contribution Process
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write/update tests
5. Submit a pull request

---

## 🐛 Known Issues

### Current Limitations
- No caching layer (all requests go to database)
- Limited error handling in some services
- No rate limiting on API Gateway
- Database migrations embedded in main.go

### Architectural Issues
See `ARCHITECTURE_REVIEW.md` for detailed analysis of architectural concerns and recommendations.

---

## 📞 Contact

**Author:** [Author Name]  
**Email:** [email]  
**GitHub:** [github username]  
**LinkedIn:** [linkedin profile]

---

## 📄 License

MIT License

Copyright (c) 2025 [Author Name]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

---

## 🙏 Acknowledgments

### Technologies Used
- [Go](https://golang.org/) - Programming language
- [gRPC](https://grpc.io/) - RPC framework
- [React](https://reactjs.org/) - UI library
- [PostgreSQL](https://www.postgresql.org/) - Database
- [Kafka](https://kafka.apache.org/) - Message broker
- [Docker](https://www.docker.com/) - Containerization

### Inspired By
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Uncle Bob
- [Microservices Patterns](https://microservices.io/) by Chris Richardson
- [12 Factor App](https://12factor.net/) methodology

---

## 📚 Additional Resources

### Documentation
- [Go Standard Project Layout](https://github.com/golang-standards/project-layout)
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/quickstart/)
- [React Documentation](https://react.dev/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)

### Related Projects
- [go-micro](https://github.com/micro/go-micro) - Go microservices framework
- [go-kit](https://github.com/go-kit/kit) - Programming toolkit for microservices
- [kratos](https://github.com/go-kratos/kratos) - Microservice framework

---

**Last Updated:** November 16, 2025  
**Version:** 1.0.0-alpha  
**Documentation Version:** 1.0

