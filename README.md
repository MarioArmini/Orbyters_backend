
# Orbyters Backend

**Orbyters Backend** is a server-side application written in Go. This project handles APIs for a web or mobile application, supporting authentication, user management, and other services.

## Project Structure

### `/cmd/server`
Contains the main entry point to start the backend server. It includes the logic needed to configure and initialize the application.

### `/config`
Configuration files and global settings. For example:
- **API Keys**
- **Database credentials**
- **External service configurations**

### `/models`
Defines the data models used in the application, generally representing database tables.  
For example:
- **User**: Model for users.
- **Tokens**: Manages access or refresh tokens.

### `/routes`
Contains files for handling HTTP routes.  
Organized into modules such as:
- **`/auth`**: Endpoints for login, registration, and authentication.
- **`/user`**: Routes for managing user details and operations.
- **`/huggingFace`**: Interface for calling Hugging Face AI models.

### `/services`
Encapsulates the business logic of the backend.  
For example:
- **`jwt`**: Service for validating and generating JSON Web Tokens (JWT).
- **`huggingFace`**: Service to send requests to AI models hosted on Hugging Face.

### `/docs`
Includes additional documentation for developers and users of the project.

## Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/MarioArmini/Orbyters_backend.git
   cd Orbyters_backend
   ```

2. **Configure environment variables**:
   - Create a `.env` file.
   - Add keys such as:
     ```
      DB_CONNECTION_STRING
      JWT_SECRET
      HUGGING_FACE_KEY
      HUGGING_FACE_URL="https://api-inference.huggingface.co/models/mistralai/Mistral-7B-Instruct-v0.2/v1/chat/completions"
      MODEL_NAME="mistralai/Mistral-7B-Instruct-v0.2"
     ```

3. **Install dependencies**:
   ```bash
   go mod tidy
   ```

4. **Update swagger docs**:
   ```bash
   make swag
   ```

5. **Start the server**:
   ```bash
   make run
   ```

## Key APIs

| Method | Endpoint           | Description                        |
|--------|--------------------|------------------------------------|
| GET    | `/user/details`    | Fetches details of the authenticated user. |
| POST   | `/auth/login`      | Logs in and returns a JWT token.    |
| POST   | `/mistral/generate` | Generates a chatbot response starting from a prompt.    |

## Tools

- **Golang**: Main backend framework.
- **Gin Framework**: Lightweight web framework for Go.
- **Gorm**: ORM for database management.
- **Hugging Face**: API for artificial intelligence models.
- **JWT**: Token-based authentication.

## Contributions

Feel free to contribute to the project. Open a pull request or report issues in the [Issues section](https://github.com/MarioArmini/Orbyters_backend/issues).

