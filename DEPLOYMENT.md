# Deployment Guide

This guide covers deploying the Work Tracker backend on Coolify and other platforms.

## Coolify Deployment

### Prerequisites

1. **Coolify Instance**: Set up and running
2. **GitHub Repository**: Code pushed to GitHub
3. **Upstash Redis**: Redis database for production
4. **Domain**: Custom domain for your application (optional)

### Step-by-Step Deployment

#### 1. Connect Repository to Coolify

1. Log into your Coolify dashboard
2. Click "New Application" → "Source: Git"
3. Connect your GitHub account
4. Select the `work-tracker` repository
5. Choose the `main` branch

#### 2. Configure Application Settings

1. **Application Name**: `work-tracker`
2. **Build Pack**: Docker
3. **Dockerfile Path**: `Dockerfile`
4. **Port**: `8080`

#### 3. Set Environment Variables

Configure these environment variables in Coolify:

| Variable | Description | Example | Required |
|----------|-------------|---------|----------|
| `REDIS_URL` | Upstash Redis connection URL | `rediss://username:password@host:port` | ✅ |
| `JWT_SECRET` | Secret for JWT token signing | `your-super-secret-key` | ✅ |
| `PORT` | Application port | `8080` | ❌ |
| `JWT_TTL_HOURS` | JWT token expiration (hours) | `336` | ❌ |
| `TIMEZONE` | Application timezone | `Europe/Paris` | ❌ |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | `https://yourdomain.com` | ❌ |

#### 4. Configure Domain (Optional)

1. Add your custom domain in Coolify
2. Configure SSL certificate
3. Set up DNS records pointing to Coolify

#### 5. Deploy

1. Click "Deploy" in Coolify
2. Monitor the build and deployment process
3. Verify the application is running at your domain

### Environment Variables Details

#### REDIS_URL

Get this from your Upstash dashboard:

```bash
# Format: rediss://username:password@host:port
REDIS_URL=rediss://default:password@host:port
```

#### JWT_SECRET

Generate a secure random string:

```bash
# Generate a secure secret
openssl rand -base64 32
```

#### CORS_ALLOWED_ORIGINS

For production, specify your frontend domain:

```bash
# Single domain
CORS_ALLOWED_ORIGINS=https://yourdomain.com

# Multiple domains
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Development (allow all)
CORS_ALLOWED_ORIGINS=*
```

## Local Development

### Using Docker Compose

```bash
# Clone the repository
git clone https://github.com/yourusername/work-tracker.git
cd work-tracker

# Start the application with Redis
docker-compose up -d

# View logs
docker-compose logs -f work-tracker

# Stop the application
docker-compose down
```

### Manual Docker Build

```bash
# Build the image
docker build -t work-tracker .

# Run with environment variables
docker run -p 8080:8080 \
  -e REDIS_URL=redis://localhost:6379 \
  -e JWT_SECRET=dev-secret-key \
  -e CORS_ALLOWED_ORIGINS=http://localhost:3000 \
  work-tracker
```

## Health Checks

The application includes a health check endpoint:

```bash
# Check application health
curl http://localhost:8080/v1/health

# Expected response
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

## Monitoring

### Logs

Monitor application logs in Coolify dashboard or via Docker:

```bash
# View application logs
docker logs work-tracker

# Follow logs in real-time
docker logs -f work-tracker
```

### Metrics

The application exposes basic health metrics at `/v1/health`.

## Troubleshooting

### Common Issues

#### 1. Redis Connection Failed

**Error**: `dial tcp: connection refused`

**Solution**:

- Verify `REDIS_URL` is correct
- Check Upstash dashboard for connection details
- Ensure Redis instance is running

#### 2. JWT Token Issues

**Error**: `invalid token`

**Solution**:

- Verify `JWT_SECRET` is set correctly
- Check token expiration settings
- Ensure consistent secret across deployments

#### 3. CORS Errors

**Error**: `CORS policy blocked`

**Solution**:

- Update `CORS_ALLOWED_ORIGINS` with your frontend domain
- Remove trailing slashes from URLs
- Use HTTPS for production domains

#### 4. Port Already in Use

**Error**: `bind: address already in use`

**Solution**:

- Change `PORT` environment variable
- Stop conflicting services
- Use different port in Coolify configuration

### Debug Mode

Enable debug logging by setting:

```bash
LOG_LEVEL=debug
```

## Security Considerations

### Production Checklist

- [ ] Use HTTPS in production
- [ ] Set strong `JWT_SECRET`
- [ ] Configure `CORS_ALLOWED_ORIGINS` properly
- [ ] Use Upstash Redis with SSL
- [ ] Enable Coolify security features
- [ ] Set up monitoring and alerts
- [ ] Configure backup strategy

### Environment Variables Security

- Never commit secrets to Git
- Use Coolify's secret management
- Rotate secrets regularly
- Use different secrets for each environment

## Backup and Recovery

### Data Backup

The application uses Redis for data storage. Backup strategies:

1. **Upstash Backups**: Configure automatic backups in Upstash
2. **Manual Export**: Use Redis dump commands
3. **Application Level**: Implement data export endpoints

### Recovery

1. Restore Redis data from backup
2. Redeploy application if needed
3. Verify data integrity
4. Test application functionality

## Scaling

### Horizontal Scaling

To scale the application:

1. Increase replicas in Coolify
2. Ensure Redis can handle increased load
3. Configure load balancing
4. Monitor resource usage

### Resource Limits

Default resource limits:

- CPU: 0.5 cores
- Memory: 512MB
- Storage: 1GB

Adjust based on your usage patterns.

## CI/CD Integration

### GitHub Actions

The repository includes a GitHub Actions workflow for automated deployment:

1. Push to `main` branch triggers deployment
2. Tests run automatically
3. Successful tests trigger Coolify deployment
4. Monitor deployment status in GitHub Actions

### Manual Deployment

For manual deployments:

```bash
# Build and push to registry
docker build -t yourregistry/work-tracker .
docker push yourregistry/work-tracker

# Deploy via Coolify CLI or dashboard
```

## Support

For deployment issues:

1. Check Coolify documentation
2. Review application logs
3. Verify environment variables
4. Test locally first
5. Contact support with specific error messages
