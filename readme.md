# Work Tracker Backend API Documentation

A secure, Redis-powered work time tracking backend with authentication, automated timesheet generation, and comprehensive reporting.

## Overview

This backend provides a RESTful API for tracking work hours, managing user settings, generating automated timesheets, and creating monthly work recaps. Built with Go, Redis (Upstash), and JWT authentication.

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
