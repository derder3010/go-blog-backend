# Blog Backend API

A RESTful API backend for a blog system built with Go, using Gin framework, MongoDB, and Cloudflare R2 for image storage.

## Technologies

- **Go**: Programming language
- **Gin**: Web framework
- **MongoDB**: Database
- **Cloudflare R2**: Image storage
- **JWT**: Authentication

## Features

- User authentication (register, login)
- Blog post management (CRUD operations)
- Image upload to Cloudflare R2
- JWT-based authorization
- CORS support

## Prerequisites

- Go 1.21 or higher
- MongoDB
- Cloudflare R2 account

## Configuration

Create a `.env` file in the root directory:

```env
MONGO_URI="your_mongodb_connection_string"
DATABASE_NAME="blog"
JWT_SECRET="your_jwt_secret"
R2_ACCOUNT_ID="your_cloudflare_account_id"
R2_ACCESS_KEY="your_r2_access_key"
R2_SECRET_KEY="your_r2_secret_key"
R2_BUCKET="your_bucket_name"
R2_PUBLIC_URL="your_r2_public_url"
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/go-blog-backend.git
cd go-blog-backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run main.go
```

The server will start at `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/register`: Register a new user
  ```json
  {
    "username": "string",
    "email": "string",
    "password": "string"
  }
  ```
- `POST /api/login`: Login
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```

### Posts
- `GET /api/posts`: Get all posts
- `GET /api/posts/:id`: Get a specific post
- `POST /api/posts`: Create a new post (requires authentication)
  ```json
  {
    "title": "string",
    "content": "string",
    "image_url": "string"
  }
  ```
- `PUT /api/posts/:id`: Update a post (requires authentication)
- `DELETE /api/posts/:id`: Delete a post (requires authentication)

### User Management
- `PUT /api/user`: Update user profile (requires authentication)
- `DELETE /api/user`: Delete user account (requires authentication)

### Image Upload
- `POST /api/upload`: Upload an image (requires authentication)
  - Use form-data with key "image"
  - Supports jpeg, png, gif formats

## Authentication

All protected routes require a Bearer token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Project Structure

```
.
├── config/
│   └── config.go
├── handlers/
│   ├── handler_interfaces.go
│   ├── post_handler.go
│   ├── upload_handler.go
│   └── user_handler.go
├── middleware/
│   └── auth_middleware.go
├── models/
│   ├── post.go
│   └── user.go
├── pkg/
│   ├── cloudflare/
│   │   └── r2.go
│   └── utils/
│       ├── image.go
│       ├── jwt.go
│       └── password.go
├── repositories/
│   ├── post_repository.go
│   └── user_repository.go
├── services/
│   ├── post_service.go
│   ├── upload_service.go
│   └── user_service.go
├── .env
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details

## Contact

Your Name - duytranduc3010@gmail.com

Project Link: https://github.com/derder3010/go-blog-backend
