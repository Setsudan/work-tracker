# Work Tracker Backend

A secure, Redis-powered work time tracking backend with authentication, automated timesheet generation, and comprehensive reporting. Built with Go, Redis (Upstash), and JWT authentication.

## 🚀 Features

- **Secure Authentication** - JWT-based auth with Argon2id password hashing
- **Time Tracking** - Start/stop time logs with smart toggle logic
- **Automated Reports** - Weekly timesheets and monthly recaps generated automatically
- **Days Off Management** - Declare vacation days and holidays
- **Overtime Tracking** - Automatic calculation of worked vs expected hours
- **Lunch Break Logic** - Smart deduction based on log patterns
- **Data Retention** - Automatic cleanup of old data
- **Production Ready** - Optimized for Coolify deployment with health checks

## 🛠️ Tech Stack

- **Backend**: Go 1.22+
- **Database**: Redis (Upstash)
- **Authentication**: JWT with session invalidation
- **Password Hashing**: Argon2id
- **HTTP Router**: Chi
- **Scheduling**: Cron jobs
- **Containerization**: Docker with multi-stage builds
- **Deployment**: Coolify-ready with health checks

## 📋 Prerequisites

- Go 1.22 or higher
- Redis instance (Upstash recommended for production)
- Docker (for containerized deployment)

## 🚀 Quick Start

### Local Development

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/work-tracker.git
   cd work-tracker
   ```

2. **Set up environment variables**

   ```bash
   export REDIS_URL=redis://localhost:6379
   export JWT_SECRET=your-dev-secret-key
   export CORS_ALLOWED_ORIGINS=http://localhost:3000
   ```

3. **Run with Docker Compose (recommended)**

   ```bash
   docker-compose up -d
   ```

4. **Or run directly with Go**

   ```bash
   go mod download
   go run ./cmd/server
   ```

5. **Test the API**

   ```bash
   curl http://localhost:8080/v1/health
   ```

### Production Deployment

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed deployment instructions.

**Quick Coolify Deploy:**

1. Connect your GitHub repository to Coolify
2. Set environment variables: `REDIS_URL`, `JWT_SECRET`, `CORS_ALLOWED_ORIGINS`
3. Deploy with the provided `coolify.yml` configuration

## 📖 API Documentation

The backend provides a comprehensive RESTful API with the following endpoints:

### Authentication

- `POST /v1/auth/register` - User registration
- `POST /v1/auth/login` - User login
- `GET /v1/auth/me` - Get current user profile
- `POST /v1/auth/logout` - Logout and invalidate session

### Time Tracking

- `POST /v1/time-logs/toggle` - Create start/stop time log
- `GET /v1/time-logs` - Get time logs within date range

### User Management

- `PUT /v1/users/me` - Update user settings

### Days Off

- `GET /v1/days-off` - List days off
- `POST /v1/days-off` - Add day off
- `DELETE /v1/days-off/{dateISO}` - Remove day off

### Reports

- `GET /v1/timesheets` - Get weekly timesheets
- `GET /v1/month-recaps` - Get monthly recaps

### Health

- `GET /v1/health` - Application health check

For complete API documentation with request/response examples, see the [API Documentation](#api-documentation) section below.

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `REDIS_URL` | Redis connection URL | ✅ | - |
| `JWT_SECRET` | JWT signing secret | ✅ | - |
| `PORT` | Application port | ❌ | `8080` |
| `JWT_TTL_HOURS` | JWT token expiration (hours) | ❌ | `336` (14 days) |
| `TIMEZONE` | Application timezone | ❌ | `Europe/Paris` |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | ❌ | `*` |

### Example Configuration

```bash
# Development
REDIS_URL=redis://localhost:6379
JWT_SECRET=dev-secret-key
CORS_ALLOWED_ORIGINS=http://localhost:3000

# Production
REDIS_URL=rediss://username:password@host:port
JWT_SECRET=your-super-secret-jwt-key
CORS_ALLOWED_ORIGINS=https://yourdomain.com
```

## 🏗️ Project Structure

```
work-tracker/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── auth/           # JWT authentication
│   ├── config/         # Configuration management
│   ├── httpserver/     # HTTP server and routes
│   ├── model/          # Data models
│   ├── scheduler/      # Cron jobs and automation
│   ├── service/        # Business logic
│   ├── store/          # Redis data access
│   └── timeutil/       # Time utilities
├── Dockerfile          # Production container
├── docker-compose.yml  # Local development
├── coolify.yml         # Coolify deployment config
└── DEPLOYMENT.md       # Deployment guide
```

## 🔄 Automated Tasks

The application runs several automated tasks:

- **Weekly Timesheet Generation** - Every Saturday at 01:00 (Europe/Paris)
- **Monthly Recap Generation** - 1st day of each month at 02:00 (Europe/Paris)
- **Data Cleanup** - Daily at 03:30, removes time logs older than 14 days

## 🛡️ Security Features

- **Password Security**: Argon2id hashing with secure parameters
- **JWT Tokens**: Secure token-based authentication with expiration
- **Session Management**: Redis-backed session invalidation
- **CORS Protection**: Configurable cross-origin resource sharing
- **Container Security**: Non-root user, read-only filesystem
- **Data Retention**: Automatic cleanup of sensitive data

## 📊 Business Logic

### Time Log Toggle Logic

- No previous log → Create START log
- Last log was START → Create STOP log
- Last log was STOP → Create START log

### Lunch Break Calculation

- **2 logs per day**: Total work time minus lunch break
- **4 logs per day**: (First start to first stop) + (Second start to second stop) - lunch break
- **Other patterns**: Total work time minus lunch break

### Overtime Calculation

- Daily overtime = worked minutes - expected minutes
- Weekly overtime = sum of daily overtimes
- Expected daily hours = weekly hours ÷ 5

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Build application
go build ./cmd/server
```

## 📈 Monitoring

The application includes built-in monitoring:

- **Health Check**: `GET /v1/health` returns application status
- **Docker Health**: Container health checks for deployment platforms
- **Logging**: Structured logging for debugging and monitoring

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Check [DEPLOYMENT.md](DEPLOYMENT.md) for deployment issues
- **API Issues**: Review the API documentation below
- **General Questions**: Open an issue on GitHub

---

## 📚 API Documentation

This section provides detailed API documentation for frontend developers.

## Base URL

```
http://localhost:8080/v1
```

For production, replace with your deployed domain.

## Authentication

All endpoints except `/auth/*` require a valid JWT token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

## API Endpoints

### Authentication

#### POST /auth/register

Register a new user account.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "weeklyHours": 40,
  "defaultLunchBreakMinutes": 60,
  "timezone": "Europe/Paris"
}
```

**Response (201):**

```json
{
  "user": {
    "id": "01HXYZ1234567890ABCDEF",
    "email": "user@example.com",
    "weeklyHours": 40,
    "defaultLunchBreakMinutes": 60,
    "timezone": "Europe/Paris",
    "createdAt": "2024-01-15T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### POST /auth/login

Authenticate existing user.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200):**

```json
{
  "user": {
    "id": "01HXYZ1234567890ABCDEF",
    "email": "user@example.com",
    "weeklyHours": 40,
    "defaultLunchBreakMinutes": 60,
    "timezone": "Europe/Paris",
    "createdAt": "2024-01-15T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### GET /auth/me

Get current user profile (requires authentication).

**Response (200):**

```json
{
  "id": "01HXYZ1234567890ABCDEF",
  "email": "user@example.com",
  "weeklyHours": 40,
  "defaultLunchBreakMinutes": 60,
  "timezone": "Europe/Paris",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

#### POST /auth/logout

Invalidate current session (requires authentication).

**Response (200):**

```json
{
  "message": "Logged out successfully"
}
```

### User Management

#### PUT /users/me

Update current user settings (requires authentication).

**Request Body:**

```json
{
  "weeklyHours": 35,
  "defaultLunchBreakMinutes": 45,
  "timezone": "America/New_York"
}
```

**Response (200):**

```json
{
  "id": "01HXYZ1234567890ABCDEF",
  "email": "user@example.com",
  "weeklyHours": 35,
  "defaultLunchBreakMinutes": 45,
  "timezone": "America/New_York",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

### Time Logs

#### POST /time-logs/toggle

Create a start or stop time log. The system automatically determines if this should be a start or stop based on the last log entry.

**Request Body:**

```json
{
  "timestamp": "2024-01-15T09:00:00Z"
}
```

**Response (201):**

```json
{
  "id": "01HXYZ1234567890ABCDEF",
  "userId": "01HXYZ1234567890ABCDEF",
  "type": "start",
  "timestamp": "2024-01-15T09:00:00Z"
}
```

#### GET /time-logs

Get time logs within a date range (requires authentication).

**Query Parameters:**

- `from` (optional): Start date in ISO format (default: 7 days ago)
- `to` (optional): End date in ISO format (default: today)

**Response (200):**

```json
{
  "logs": [
    {
      "id": "01HXYZ1234567890ABCDEF",
      "userId": "01HXYZ1234567890ABCDEF",
      "type": "start",
      "timestamp": "2024-01-15T09:00:00Z"
    },
    {
      "id": "01HXYZ1234567890ABCDEF",
      "userId": "01HXYZ1234567890ABCDEF",
      "type": "stop",
      "timestamp": "2024-01-15T17:30:00Z"
    }
  ]
}
```

### Days Off

#### GET /days-off

Get all days off for the current user (requires authentication).

**Response (200):**

```json
{
  "daysOff": [
    {
      "id": "01HXYZ1234567890ABCDEF",
      "userId": "01HXYZ1234567890ABCDEF",
      "dateISO": "2024-01-20",
      "reason": "Vacation",
      "createdAt": "2024-01-15T10:30:00Z"
    }
  ]
}
```

#### POST /days-off

Add a day off (requires authentication).

**Request Body:**

```json
{
  "dateISO": "2024-01-20",
  "reason": "Vacation"
}
```

**Response (201):**

```json
{
  "id": "01HXYZ1234567890ABCDEF",
  "userId": "01HXYZ1234567890ABCDEF",
  "dateISO": "2024-01-20",
  "reason": "Vacation",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /days-off/{dateISO}

Remove a day off (requires authentication).

**Response (200):**

```json
{
  "message": "Day off removed successfully"
}
```

### Timesheets

#### GET /timesheets

Get weekly timesheets (requires authentication).

**Query Parameters:**

- `from` (optional): Start date in ISO format (default: 4 weeks ago)
- `to` (optional): End date in ISO format (default: today)

**Response (200):**

```json
{
  "timesheets": [
    {
      "id": "01HXYZ1234567890ABCDEF",
      "userId": "01HXYZ1234567890ABCDEF",
      "weekStart": "2024-01-15T00:00:00Z",
      "weekEnd": "2024-01-21T23:59:59Z",
      "totalWorkedMinutes": 2400,
      "totalExpectedMinutes": 2400,
      "totalOvertimeMinutes": 0,
      "averageWorkedMinutesPerDay": 480,
      "averageStartTime": "09:15",
      "averageEndTime": "17:15",
      "days": [
        {
          "dateISO": "2024-01-15",
          "workedMinutes": 480,
          "expectedMinutes": 480,
          "overtimeMinutes": 0,
          "firstStartTime": "09:00",
          "lastEndTime": "17:00",
          "averageStartTime": "09:00",
          "averageEndTime": "17:00"
        }
      ],
      "createdAt": "2024-01-20T01:00:00Z"
    }
  ]
}
```

### Month Recaps

#### GET /month-recaps

Get monthly recaps (requires authentication).

**Query Parameters:**

- `from` (optional): Start date in ISO format (default: 6 months ago)
- `to` (optional): End date in ISO format (default: today)

**Response (200):**

```json
{
  "recaps": [
    {
      "id": "01HXYZ1234567890ABCDEF",
      "userId": "01HXYZ1234567890ABCDEF",
      "monthStart": "2024-01-01T00:00:00Z",
      "monthEnd": "2024-01-31T23:59:59Z",
      "totalWorkedMinutes": 9600,
      "totalExpectedMinutes": 9600,
      "totalOvertimeMinutes": 0,
      "averageWorkedMinutesPerDay": 480,
      "averageStartTime": "09:15",
      "averageEndTime": "17:15",
      "mostProductiveDay": "Monday",
      "leastProductiveDay": "Friday",
      "daysWorked": 20,
      "daysOff": 11,
      "createdAt": "2024-02-01T02:00:00Z"
    }
  ]
}
```

## Data Models

### User

```typescript
interface User {
  id: string;
  email: string;
  weeklyHours: number;
  defaultLunchBreakMinutes: number;
  timezone: string;
  createdAt: string;
}
```

### TimeLog

```typescript
interface TimeLog {
  id: string;
  userId: string;
  type: 'start' | 'stop';
  timestamp: string;
}
```

### DayOff

```typescript
interface DayOff {
  id: string;
  userId: string;
  dateISO: string;
  reason: string;
  createdAt: string;
}
```

### DaySummary

```typescript
interface DaySummary {
  dateISO: string;
  workedMinutes: number;
  expectedMinutes: number;
  overtimeMinutes: number;
  firstStartTime: string;
  lastEndTime: string;
  averageStartTime: string;
  averageEndTime: string;
}
```

### TimeSheet

```typescript
interface TimeSheet {
  id: string;
  userId: string;
  weekStart: string;
  weekEnd: string;
  totalWorkedMinutes: number;
  totalExpectedMinutes: number;
  totalOvertimeMinutes: number;
  averageWorkedMinutesPerDay: number;
  averageStartTime: string;
  averageEndTime: string;
  days: DaySummary[];
  createdAt: string;
}
```

### MonthRecap

```typescript
interface MonthRecap {
  id: string;
  userId: string;
  monthStart: string;
  monthEnd: string;
  totalWorkedMinutes: number;
  totalExpectedMinutes: number;
  totalOvertimeMinutes: number;
  averageWorkedMinutesPerDay: number;
  averageStartTime: string;
  averageEndTime: string;
  mostProductiveDay: string;
  leastProductiveDay: string;
  daysWorked: number;
  daysOff: number;
  createdAt: string;
}
```

## Error Responses

All endpoints return consistent error responses:

**400 Bad Request:**

```json
{
  "error": "Invalid request data",
  "details": "Field 'email' is required"
}
```

**401 Unauthorized:**

```json
{
  "error": "Authentication required"
}
```

**403 Forbidden:**

```json
{
  "error": "Access denied"
}
```

**404 Not Found:**

```json
{
  "error": "Resource not found"
}
```

**500 Internal Server Error:**

```json
{
  "error": "Internal server error"
}
```

## Business Logic

### Time Log Toggle Logic

- If no previous log exists → create START log
- If last log was START → create STOP log
- If last log was STOP → create START log

### Lunch Break Calculation

- **2 logs per day**: Total work time minus lunch break
- **4 logs per day**: (First start to first stop) + (Second start to second stop) - lunch break
- **Other patterns**: Total work time minus lunch break

### Overtime Calculation

- Daily overtime = worked minutes - expected minutes
- Weekly overtime = sum of daily overtimes
- Expected daily hours = weekly hours ÷ 5

### Days Off

- Days marked as "off" count as fulfilled expected hours
- They don't affect overtime calculations
- Useful for vacations, holidays, sick days

## Automated Tasks

### Weekly Timesheet Generation

- **Schedule**: Every Saturday at 01:00 (Europe/Paris timezone)
- **Process**: Generates timesheet for the previous week
- **Retention**: 90 days

### Monthly Recap Generation

- **Schedule**: 1st day of each month at 02:00 (Europe/Paris timezone)
- **Process**: Compiles data from weekly timesheets
- **Retention**: 365 days

### Data Cleanup

- **Schedule**: Daily at 03:30 (Europe/Paris timezone)
- **Process**: Removes time logs older than 14 days

## Environment Variables

```bash
# Required
REDIS_URL=rediss://username:password@host:port
JWT_SECRET=your-super-secret-jwt-key

# Optional
PORT=8080
JWT_TTL_HOURS=336
TIMEZONE=Europe/Paris
CORS_ALLOWED_ORIGINS=*
```

## Frontend Considerations

### Authentication Flow

1. User registers/logs in
2. Store JWT token securely (httpOnly cookie recommended)
3. Include token in all API requests
4. Handle token expiration (401 responses)
5. Implement logout to clear token

### Real-time Updates

- Poll for new timesheets/recaps or implement WebSocket
- Update UI when time logs are added
- Show current work status (started/stopped)

### Timezone Handling

- Display times in user's timezone
- Send timestamps in ISO format
- Handle DST transitions

### Data Visualization

- Charts for daily/weekly/monthly trends
- Overtime tracking
- Productivity patterns
- Export functionality

### Mobile Considerations

- Touch-friendly time log toggle
- Offline capability for logging
- Push notifications for reminders

## Security Notes

- JWT tokens expire after 14 days
- Passwords hashed with Argon2id
- CORS configured for cross-origin requests
- All sensitive data stored in Redis with TTL
- Rate limiting recommended for production

## Deployment

### Coolify Deployment

This application is optimized for deployment on Coolify. Use the provided configuration files:

1. **Dockerfile** - Multi-stage build with security best practices
2. **coolify.yml** - Coolify-specific configuration
3. **docker-compose.yml** - Local development setup

#### Quick Deploy on Coolify

1. Connect your GitHub repository to Coolify
2. Select the repository and branch
3. Configure environment variables:
   - `REDIS_URL` - Your Upstash Redis URL
   - `JWT_SECRET` - Secure random string for JWT signing
   - `CORS_ALLOWED_ORIGINS` - Your frontend domain(s)

4. Deploy with the provided `coolify.yml` configuration

#### Environment Variables for Production

```bash
# Required
REDIS_URL=rediss://username:password@host:port
JWT_SECRET=your-super-secret-jwt-key

# Optional
PORT=8080
JWT_TTL_HOURS=336
TIMEZONE=Europe/Paris
CORS_ALLOWED_ORIGINS=https://yourdomain.com
```

#### Local Development

```bash
# Using Docker Compose
docker-compose up -d

# Or build and run manually
docker build -t work-tracker .
docker run -p 8080:8080 \
  -e REDIS_URL=redis://localhost:6379 \
  -e JWT_SECRET=dev-secret-key \
  work-tracker
```

### Docker Example

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
```

### Environment Setup

```bash
# Production
REDIS_URL=rediss://your-upstash-url
JWT_SECRET=your-production-secret
CORS_ALLOWED_ORIGINS=https://yourdomain.com

# Development
REDIS_URL=redis://localhost:6379
JWT_SECRET=dev-secret-key
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

## API Rate Limits

Consider implementing rate limiting for production:

- Authentication endpoints: 5 requests/minute
- Time log endpoints: 60 requests/minute
- Other endpoints: 100 requests/minute

## Testing

Test all endpoints with tools like Postman or curl:

```bash
# Register
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","weeklyHours":40}'

# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Toggle time log (with token)
curl -X POST http://localhost:8080/v1/time-logs/toggle \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"timestamp":"2024-01-15T09:00:00Z"}'
```
