# Deploying Banner Generator to Fly.io

This guide walks you through deploying the Banner Generator API to Fly.io.

## Prerequisites

1. Install the Fly CLI:
   ```bash
   curl -L https://fly.io/install.sh | sh
   ```

2. Sign up and log in to Fly.io:
   ```bash
   fly auth signup
   # or if you already have an account
   fly auth login
   ```

## Deployment Steps

### 1. Initialize the Fly App

From the project root directory:

```bash
fly launch --no-deploy
```

This will:
- Detect the existing `fly.toml` configuration
- Create the app on Fly.io
- Set up the necessary resources

### 2. Configure Secrets (Optional)

If you want to use a GitHub token for better API rate limits:

```bash
fly secrets set GITHUB_TOKEN=your-github-personal-access-token
```

### 3. Deploy the Application

```bash
fly deploy
```

This will:
- Build the Docker image
- Push it to Fly's registry
- Deploy the application

### 4. Verify Deployment

Check the application status:

```bash
fly status
```

View logs:

```bash
fly logs
```

Open your app in the browser:

```bash
fly open
```

## Configuration

The `fly.toml` file includes the following configuration:

- **Region**: Set to `iad` (Ashburn, Virginia) by default
- **Auto-scaling**: Configured to scale down to 0 when idle and auto-start on requests
- **Resources**: 1 shared CPU and 256MB RAM (suitable for most banner generation workloads)
- **Concurrency**: Soft limit of 200 concurrent requests, hard limit of 250

### Environment Variables

You can modify environment variables in `fly.toml` or set them via the CLI:

```bash
fly secrets set VARIABLE_NAME=value
```

Common variables:
- `GITHUB_TOKEN`: For improved GitHub API rate limits
- `ACCESS_CONTROL_ENABLED`: Enable/disable access control
- `ALLOWED_ORGS`: Comma-separated list of allowed GitHub organizations
- `BANNER_CACHE_TTL`: Cache duration for generated banners

### Scaling

To scale your app:

```bash
# Scale to 2 instances
fly scale count 2

# Scale memory
fly scale memory 512

# Scale to different VM size
fly scale vm shared-cpu-2x
```

### Custom Domain

To add a custom domain:

```bash
fly certs add yourdomain.com
```

Then configure your DNS to point to the Fly.io app.

## Monitoring

View metrics:

```bash
fly dashboard
```

SSH into the running instance:

```bash
fly ssh console
```

## Troubleshooting

### Build Failures

If the build fails, check:
1. The Dockerfile syntax
2. Go module dependencies
3. Build logs with `fly logs`

### Runtime Issues

1. Check logs: `fly logs`
2. SSH into the container: `fly ssh console`
3. Verify environment variables: `fly secrets list`

### Performance Issues

1. Check metrics: `fly dashboard`
2. Scale up if needed: `fly scale vm shared-cpu-2x`
3. Add more instances: `fly scale count 2`

## Updating

To deploy updates:

```bash
git pull origin main
fly deploy
```

## Costs

The default configuration should run within Fly.io's free tier:
- Up to 3 shared-cpu-1x VMs
- 3GB persistent volume storage
- 160GB outbound data transfer

For production workloads, consider upgrading to dedicated CPU instances.