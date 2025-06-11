# Go Gin Backend Response Standardization

## 🚀 Overview

This project provides a standardized approach to handling API responses in Go applications built with the Gin web framework. It aims to solve common inconsistencies in API design by offering a consistent structure for success, error, and validation payloads across all your backend services.

By using this standardization, you'll:

- ✅ **Improve API consistency**: Make your APIs predictable and easier for frontend developers to consume.  
- 🔧 **Enhance maintainability**: Centralize response logic, reducing boilerplate and simplifying updates.  
- ⚠️ **Streamline error handling**: Provide clear, structured error messages that are easy to parse and act upon.  
- 🚀 **Boost developer experience**: Create a more efficient and less error-prone development workflow.

---

## ✨ Features

- **Consistent Success Responses**  
  Define a clear structure for `200 OK` and other success status codes.

- **Structured Error Handling**  
  Standardize error messages, codes, and details for various scenarios (e.g., internal server errors, bad requests).

- **Validation Error Uniformity**  
  Provide a unified way to return validation failures, often with specific field errors.

- **Customizable Payloads**  
  Easily extend and customize response structures to fit your specific needs.

- **Gin-compatible Middleware/Helpers**  
  Integrate seamlessly with your existing Gin routes and handlers.

---

## 🛠️ Getting Started

To get started with this project:

1. Clone the repository:
   ```bash
   mkdir your_app
   cd your_app
   git clone https://github.com/Palguna1121/response-std .
   ```

that command is for make a directory and clone the repository to it

2. Initialize Go modules:
   ```bash
   go mod init your_app
   ```

that command is for initialize go modules

3. Install dependencies:
   ```bash
   go mod tidy
   ```

that command is for install dependencies

4. Setup your environment variables:
   ```bash  
   cp .env.example .env
   ```

that command is for setup your environment variables

5. Run the server:
   ```bash
   go run main.go
   ```
   or if use makefile
   ```bash
   make run
   ```

that command is for run the server