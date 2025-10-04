# Deployment Guide

This guide covers deploying Minsix to production environments.

## Deployment Options

### Option 1: Railway (Recommended for Quick Deploy)

Railway provides simple deployment for both backend and database.

#### Backend Deployment

1. Create a new project on [Railway](https://railway.app)
2. Add PostgreSQL database:
   ```bash
   railway add postgresql
   ```
3. Deploy backend:
   ```bash
   cd backend
   railway up
   ```
4. Set environment variables in Railway dashboard:
   - `ALCHEMY_API_KEY`: Your Alchemy API key
   - `ALCHEMY_NETWORK`: eth-mainnet
   - `PORT`: 8080
   - `CORS_ORIGINS`: Your frontend URL

#### Frontend Deployment on Vercel

1. Install Vercel CLI:
   ```bash
   npm i -g vercel
   ```
2. Deploy frontend:
   ```bash
   cd frontend
   vercel
   ```
3. Set environment variables in Vercel dashboard:
   - `NEXT_PUBLIC_API_URL`: Your Railway backend URL
   - `NEXT_PUBLIC_WS_URL`: Your Railway WebSocket URL

### Option 2: Fly.io

Fly.io offers excellent WebSocket support and global edge deployment.

#### Backend on Fly.io

1. Install Fly CLI:
   ```bash
   curl -L https://fly.io/install.sh | sh
   ```

2. Create `fly.toml` in backend directory:
   ```toml
   app = "minsix-backend"
   
   [build]
     dockerfile = "Dockerfile"
   
   [env]
     PORT = "8080"
   
   [[services]]
     internal_port = 8080
     protocol = "tcp"
   
     [[services.ports]]
       handlers = ["http"]
       port = 80
   
     [[services.ports]]
       handlers = ["tls", "http"]
       port = 443
   ```

3. Create Dockerfile:
   ```dockerfile
   FROM golang:1.21-alpine AS builder
   WORKDIR /app
   COPY go.* ./
   RUN go mod download
   COPY . .
   RUN go build -o server cmd/server/main.go
   
   FROM alpine:latest
   RUN apk --no-cache add ca-certificates
   WORKDIR /root/
   COPY --from=builder /app/server .
   COPY --from=builder /app/migrations ./migrations
   EXPOSE 8080
   CMD ["./server"]
   ```

4. Deploy:
   ```bash
   fly launch
   fly secrets set ALCHEMY_API_KEY=your_key
   fly deploy
   ```

#### Database on Supabase

1. Create project on [Supabase](https://supabase.com)
2. Get connection string from project settings
3. Run migrations manually or via migration tool

### Option 3: AWS (Production Grade)

For production deployments with high availability.

#### Architecture

- **Backend**: AWS ECS (Fargate) or EC2
- **Database**: RDS PostgreSQL
- **Frontend**: S3 + CloudFront
- **Load Balancer**: Application Load Balancer

#### Backend on ECS

1. Build and push Docker image:
   ```bash
   docker build -t minsix-backend .
   aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin your-account.dkr.ecr.us-east-1.amazonaws.com
   docker tag minsix-backend:latest your-account.dkr.ecr.us-east-1.amazonaws.com/minsix:latest
   docker push your-account.dkr.ecr.us-east-1.amazonaws.com/minsix:latest
   ```

2. Create ECS task definition and service via AWS Console or Terraform

3. Set up RDS PostgreSQL instance

4. Configure ALB with WebSocket support

#### Frontend on S3 + CloudFront

1. Build frontend:
   ```bash
   cd frontend
   npm run build
   ```

2. Deploy to S3:
   ```bash
   aws s3 sync out/ s3://your-bucket-name
   ```

3. Configure CloudFront distribution

### Option 4: Docker Compose (Self-Hosted)

For self-hosted deployments on VPS (DigitalOcean, Linode, etc.)

1. Create production docker-compose.yml:
   ```yaml
   version: '3.8'
   services:
     postgres:
       image: postgres:15-alpine
       environment:
         POSTGRES_USER: postgres
         POSTGRES_PASSWORD: ${DB_PASSWORD}
         POSTGRES_DB: minsix
       volumes:
         - postgres-data:/var/lib/postgresql/data
       restart: always
     
     backend:
       build: ./backend
       environment:
         DATABASE_URL: postgres://postgres:${DB_PASSWORD}@postgres:5432/minsix?sslmode=disable
         ALCHEMY_API_KEY: ${ALCHEMY_API_KEY}
         ALCHEMY_NETWORK: eth-mainnet
         PORT: 8080
       ports:
         - "8080:8080"
       depends_on:
         - postgres
       restart: always
     
     frontend:
       build: ./frontend
       environment:
         NEXT_PUBLIC_API_URL: https://your-domain.com/api
         NEXT_PUBLIC_WS_URL: wss://your-domain.com/ws
       ports:
         - "3000:3000"
       depends_on:
         - backend
       restart: always
   
   volumes:
     postgres-data:
   ```

2. Deploy:
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. Set up Nginx reverse proxy with SSL (Let's Encrypt)

## Environment Variables

### Backend

| Variable | Description | Required |
|----------|-------------|----------|
| `ALCHEMY_API_KEY` | Alchemy API key for Ethereum access | Yes |
| `ALCHEMY_NETWORK` | Network to monitor (eth-mainnet, etc.) | Yes |
| `DATABASE_URL` | PostgreSQL connection string | Yes |
| `PORT` | Server port | No (default: 8080) |
| `CORS_ORIGINS` | Allowed CORS origins | No (default: localhost:3000) |

### Frontend

| Variable | Description | Required |
|----------|-------------|----------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | Yes |
| `NEXT_PUBLIC_WS_URL` | WebSocket URL | Yes |

## SSL/TLS Configuration

For production, always use HTTPS/WSS:

1. **Automatic**: Use Cloudflare, Vercel, or Railway (SSL included)
2. **Manual**: Use Let's Encrypt with certbot:
   ```bash
   certbot --nginx -d your-domain.com
   ```

## Monitoring & Logging

### Recommended Tools

- **Application Monitoring**: Datadog, New Relic, or Prometheus
- **Logs**: CloudWatch, Logstash, or Loki
- **Uptime**: UptimeRobot, Pingdom

### Health Checks

The backend exposes `/api/health` endpoint for monitoring:
```bash
curl https://your-backend.com/api/health
```

## Scaling Considerations

1. **Database**: Use connection pooling (already configured)
2. **Backend**: Deploy multiple instances behind load balancer
3. **WebSocket**: Use sticky sessions or Redis pub/sub for multi-instance setup
4. **Frontend**: Served via CDN (automatic with Vercel/CloudFront)

## Security Best Practices

1. Never commit `.env` files
2. Use secret management (AWS Secrets Manager, Railway secrets, etc.)
3. Enable CORS only for trusted origins
4. Use HTTPS/WSS in production
5. Regularly update dependencies
6. Implement rate limiting
7. Monitor for unusual API usage

## Cost Estimates (Monthly)

### Small Scale (< 1000 tx/day)
- Railway: ~$5
- Vercel: Free tier
- **Total**: ~$5/month

### Medium Scale (10k tx/day)
- Railway Pro: ~$20
- Vercel Pro: $20
- **Total**: ~$40/month

### Large Scale (100k+ tx/day)
- AWS ECS: ~$50
- RDS: ~$50
- CloudFront: ~$10
- **Total**: ~$110/month

## Troubleshooting

### WebSocket Connection Issues
- Ensure WebSocket protocol is enabled on proxy/load balancer
- Check CORS configuration
- Verify WSS (not WS) for HTTPS sites

### Database Connection Issues
- Check firewall rules
- Verify connection string format
- Ensure database is running and accessible

### High Memory Usage
- Review query patterns
- Implement pagination
- Add database indexes
- Increase server resources

## Support

For deployment issues, check:
1. GitHub Issues
2. Documentation
3. Platform-specific support (Railway, Vercel, etc.)
