# Antimog

<div align="center">

![Antimg Logo](static/logo.svg)

**Advanced Watermark Attack Tool**

A powerful web-based image processing platform that removes watermarks using advanced algorithms.

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Supported-blue.svg)](Dockerfile)

[🌐 Live Demo](https://antimg.neurocoda.com) | [📖 Documentation](#documentation) | [🚀 Quick Start](#quick-start)

</div>

## 📖 Overview

Antimg is a sophisticated image processing tool designed to remove watermarks from images using multiple attack algorithms. It provides both a user-friendly web interface and a powerful REST API for automated processing.

"""
api-core - 基于Go的服务器端图像处理API版本
web-core - 纯浏览器本地计算的客户端版本
"""

### ✨ Key Features

- **🎯 Advanced Watermark Removal**: Multi-layered attack algorithms including geometric, noise, frequency, compression, and color attacks
- **🌐 Dual Interface**: Beautiful web UI and comprehensive REST API
- **📱 Responsive Design**: Mobile-friendly interface with dark/light themes
- **🔒 Secure Authentication**: JWT-based authentication with API token management
- **🌍 Multi-language Support**: English and Chinese interface
- **📊 No Size Limits**: Process images of any size and resolution
- **🐳 Docker Ready**: Easy deployment with Docker and Docker Compose
- **⚡ High Performance**: Optimized Go backend with efficient image processing

### 🎨 Supported Formats

- **Input**: JPEG, PNG, BMP, WebP
- **Output**: Maintains original format or converts to JPEG

## 🚀 Quick Start

### 🐳 Docker Deployment (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/neurocoda/Antimg.git
   cd Antimg
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Deploy with Docker Compose**
   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - Web Interface: http://localhost:8080
   - Default credentials: `admin` / `admin123`

### 🔧 Manual Deployment

#### Prerequisites
- Go 1.22 or higher
- Git

#### Installation Steps

1. **Clone and build**
   ```bash
   git clone https://github.com/neurocoda/Antimg.git
   cd Antimg
   go mod download
   go build -o marknullifier .
   ```

2. **Configure environment**
   ```bash
   export PORT=8080
   export JWT_SECRET=your-super-secret-key
   export ADMIN_USERNAME=admin
   export ADMIN_PASSWORD=your-secure-password
   ```

3. **Run the application**
   ```bash
   ./marknullifier
   ```

## 🎮 Web Interface Usage

### 🔐 Login Process

1. Navigate to the application URL
2. Use your admin credentials to log in
3. You'll be redirected to the Image Processing Workbench

### 🖼️ Image Processing

1. **Upload Image**
   - Drag and drop an image or click to select
   - Supports JPEG, PNG, BMP, WebP formats
   - No size or resolution limits

2. **Configure Attack Strength**
   - Adjust the slider from 0.0 (weak) to 1.0 (strong)
   - Recommended range: 0.5 - 0.8
   - Higher values provide stronger watermark removal

3. **Process Image**
   - Click "Process Image" to start
   - Download the processed result automatically

### 🔑 API Token Management

- View your permanent API token in the interface
- Copy token for API usage
- Reset token if needed for security

## 🔌 API Usage

### 🔐 Authentication

All API requests require authentication using Bearer tokens:

```bash
Authorization: Bearer YOUR_API_TOKEN
```

### 📡 Endpoints

#### POST `/api/attack`

Remove watermarks from uploaded images.

**Request:**
```bash
curl -X POST "https://your-domain.com/api/attack" \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -F "image=@your_image.jpg" \
  -F "attackLevel=0.65" \
  --output processed_image.jpg
```

**Parameters:**
- `image` (required): Image file (multipart/form-data)
- `attackLevel` (optional): Attack strength (0.00-1.00, default: 0.50)

**Response:**
- Success: Processed image file
- Error: JSON error message

#### POST `/api/login`

Authenticate and receive access token.

**Request:**
```bash
curl -X POST "https://your-domain.com/api/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your-password"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Login successful"
}
```

### 🐍 Python Example

```python
import requests

# Login and get token
login_response = requests.post('https://your-domain.com/api/login', json={
    'username': 'admin',
    'password': 'your-password'
})
token = login_response.json()['token']

# Process image
with open('watermarked_image.jpg', 'rb') as f:
    response = requests.post(
        'https://your-domain.com/api/attack',
        headers={'Authorization': f'Bearer {token}'},
        files={'image': f},
        data={'attackLevel': '0.7'}
    )

# Save processed image
with open('processed_image.jpg', 'wb') as f:
    f.write(response.content)
```

### 🟢 Node.js Example

```javascript
const axios = require('axios');
const FormData = require('form-data');
const fs = require('fs');

async function processImage() {
  // Login
  const loginResponse = await axios.post('https://your-domain.com/api/login', {
    username: 'admin',
    password: 'your-password'
  });
  
  const token = loginResponse.data.token;
  
  // Process image
  const form = new FormData();
  form.append('image', fs.createReadStream('watermarked_image.jpg'));
  form.append('attackLevel', '0.7');
  
  const response = await axios.post('https://your-domain.com/api/attack', form, {
    headers: {
      'Authorization': `Bearer ${token}`,
      ...form.getHeaders()
    },
    responseType: 'stream'
  });
  
  // Save processed image
  response.data.pipe(fs.createWriteStream('processed_image.jpg'));
}
```

## ⚙️ Configuration

### 🌍 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key-change-in-production` |
| `ADMIN_USERNAME` | Admin username | `admin` |
| `ADMIN_PASSWORD` | Admin password | `password` |

### 🐳 Docker Configuration

**docker-compose.yml:**
```yaml
version: '3.8'
services:
  marknullifier:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - JWT_SECRET=your-super-secret-key
      - ADMIN_USERNAME=admin
      - ADMIN_PASSWORD=secure-password
    volumes:
      - uploads:/app/uploads
    restart: unless-stopped
```

## 🔧 Development

### 🏗️ Project Structure

```
Antimg/
├── config/          # Configuration management
├── handlers/        # HTTP request handlers
├── middleware/      # Authentication middleware
├── models/          # Data models
├── routes/          # Route definitions
├── services/        # Business logic
├── static/          # Static assets (CSS, JS, images)
├── templates/       # HTML templates
├── utils/           # Utility functions
├── main.go          # Application entry point
├── Dockerfile       # Docker configuration
└── docker-compose.yml
```

### 🔨 Building from Source

```bash
# Clone repository
git clone https://github.com/neurocoda/Antimg.git
cd Antimg

# Install dependencies
go mod download

# Build application
go build -o marknullifier .

# Run tests (if available)
go test ./...
```

### 🐛 Development Mode

```bash
# Install air for hot reloading
go install github.com/cosmtrek/air@latest

# Run in development mode
air
```

## 🔒 Security Considerations

- **Change default credentials** before production deployment
- **Use strong JWT secrets** (recommended: 32+ random characters)
- **Enable HTTPS** in production environments
- **Regularly rotate API tokens** for enhanced security
- **Monitor API usage** to detect unusual activity

## 🚀 Production Deployment

### 🌐 Reverse Proxy Setup (Nginx)

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Increase upload size limit
        client_max_body_size 100M;
    }
}
```

### 🔐 SSL/TLS Configuration

```bash
# Using Certbot for Let's Encrypt
sudo certbot --nginx -d your-domain.com
```

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Neurocoda**
- Website: [https://neurocoda.com](https://neurocoda.com)
- Email: [public@neurocoda.com](mailto:public@neurocoda.com)
- Demo: [https://antimg.neurocoda.com](https://antimg.neurocoda.com)

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/) and [Gin Framework](https://gin-gonic.com/)
- Image processing powered by [imaging](https://github.com/disintegration/imaging)
- UI components inspired by modern design principles
- Special thanks to the open-source community

## 📊 Algorithm Details

Antimg employs a sophisticated multi-layered approach to watermark removal:

1. **Geometric Attack**: Rotation, scaling, and perspective transformations
2. **Noise Attack**: Strategic noise injection to disrupt watermark patterns
3. **Frequency Attack**: Fourier domain manipulation for frequency-based watermarks
4. **Compression Attack**: JPEG compression artifacts to degrade watermark quality
5. **Color Attack**: Color space manipulation and channel mixing
6. **Mixed Attack**: Combination of multiple techniques for maximum effectiveness

The attack strength parameter (0.0-1.0) controls the intensity of these algorithms, allowing fine-tuned control over the removal process.

---

<div align="center">

**⭐ Star this repository if you find it useful!**

[🐛 Report Bug](https://github.com/neurocoda/Antimg/issues) | [✨ Request Feature](https://github.com/neurocoda/Antimg/issues) | [💬 Discussions](https://github.com/neurocoda/Antimg/discussions)

</div>
