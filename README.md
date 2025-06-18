# 🎯 Antimg - Advanced Image Watermark Attack Tool

[![Main Branch Docker Build](https://github.com/Neurocoda/Antimg/actions/workflows/docker-image.yml/badge.svg)](https://github.com/Neurocoda/Antimg/actions/workflows/docker-image.yml)[![Docker Pulls](https://img.shields.io/docker/pulls/neurocoda/antimg)](https://hub.docker.com/r/neurocoda/antimg)
[![API-CORE Branch Docker Build](https://github.com/Neurocoda/Antimg/actions/workflows/docker-image-api.yml/badge.svg)](https://github.com/Neurocoda/Antimg/actions/workflows/docker-image-api.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Neurocoda/Antimg)](https://goreportcard.com/report/github.com/Neurocoda/Antimg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Advanced image watermark attack processing tool with multi-platform support and comprehensive security features.

## ✨ Features

- 🎯 **Advanced Watermark Attack**: Multi-round attack algorithms
- 🔒 **Security First**: JWT authentication, rate limiting, input validation
- 🌐 **Multi-Platform**: Supports AMD64, ARM64, ARMv7
- 🐳 **Docker Ready**: Multi-arch Docker images available
- 🚀 **High Performance**: Optimized Go implementation
- 📱 **Web Interface**: User-friendly web UI
- 🔌 **REST API**: Complete API for integration

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
# Pull and run the latest image
docker run -d \
  --name antimg \
  -p 8080:8080 \
  -e JWT_SECRET="your-secure-jwt-key-at-least-32-characters-long" \
  -e ADMIN_PASSWORD="your-secure-password" \
  ghcr.io/neurocoda/antimg:latest
```

### Using Docker Compose

```bash
# Clone the repository
git clone https://github.com/Neurocoda/Antimg.git
cd antimg

# Copy environment file
cp .env.example .env

# Edit .env with your secure values
nano .env

# Start the service
docker-compose up -d
```

### Binary Installation

Download the latest binary from [Releases](https://github.com/Neurocoda/Antimg/releases):

```bash
# Linux AMD64
wget https://github.com/Neurocoda/Antimg/releases/latest/download/antimg_Linux_x86_64.tar.gz
tar -xzf antimg_Linux_x86_64.tar.gz
./antimg

# macOS ARM64
wget https://github.com/Neurocoda/Antimg/releases/latest/download/antimg_Darwin_arm64.tar.gz
tar -xzf antimg_Darwin_arm64.tar.gz
./antimg
```

## 🏗️ Supported Platforms

### Docker Images
- `linux/amd64` - Intel/AMD 64-bit
- `linux/arm64` - ARM 64-bit (Apple Silicon, AWS Graviton)
- `linux/arm/v7` - ARM 32-bit (Raspberry Pi)

### Binary Releases
- Linux: AMD64, ARM64, ARMv6, ARMv7
- macOS: AMD64, ARM64 (Apple Silicon)
- Windows: AMD64

## 📖 Usage

### Web Interface

1. Open http://localhost:8080 in your browser
2. Login with admin credentials
3. Upload an image and select attack level
4. Download the processed image

### API Usage

```bash
# Login to get token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'

# Process image
curl -X POST http://localhost:8080/api/attack \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@your-image.jpg" \
  -F "attackLevel=0.8" \
  --output processed-image.jpg
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `JWT_SECRET` | JWT signing key (min 32 chars) | - | **Yes** |
| `ADMIN_USERNAME` | Admin username | `admin` | No |
| `ADMIN_PASSWORD` | Admin password | - | **Yes** |

### Security Features

- ✅ JWT authentication with configurable expiration
- ✅ Rate limiting (API: 30/min, Processing: 10/min)
- ✅ File size validation (max 100MB)
- ✅ Processing timeout (30 seconds)
- ✅ Input sanitization and validation
- ✅ Non-root container execution

## 🔧 Development

### Prerequisites

- Go 1.21+
- Docker (for containerization)
- Make (optional)

### Local Development

```bash
# Clone repository
git clone https://github.com/Neurocoda/Antimg.git
cd antimg

# Install dependencies
go mod download

# Set environment variables
export JWT_SECRET="your-development-jwt-secret-key-32-chars"
export ADMIN_PASSWORD="dev123456"

# Run application
go run main.go
```

### Building

```bash
# Build for current platform
go build -o antimg .

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o antimg-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o antimg-darwin-arm64 .
```

## 🐳 Docker

### Multi-Platform Build

```bash
# Build for multiple platforms
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  --tag neurocoda/antimg:latest \
  --push .
```

### Production Deployment

```bash
# Use production compose file
docker-compose -f docker-compose.prod.yml up -d
```

## 📊 Monitoring

### Health Check

```bash
curl http://localhost:8080/
```

### Metrics (if enabled)

```bash
curl http://localhost:8080/metrics
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Imaging Library](https://github.com/disintegration/imaging)
- [JWT-Go](https://github.com/golang-jwt/jwt)

## 📞 Support

- 🌐 Website: [https://neurocoda.com](https://neurocoda.com)
- 🎯 Demo: [https://antimg.neurocoda.com](https://antimg.neurocoda.com)
- 🐛 Issues: [GitHub Issues](https://github.com/Neurocoda/Antimg/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/Neurocoda/Antimg/discussions)

## 👨‍💻 Author

**Neurocoda**
- Website: [https://neurocoda.com](https://neurocoda.com)
- GitHub: [@Neurocoda](https://github.com/Neurocoda)
