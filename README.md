
# Mockoon Control Panel

Mockoon Control Panel is a lightweight web-based interface that serves as an **intermediary** or **control panel** for managing the Mockoon CLI. This project does not function as a mock server itself but helps manage configurations and sync operations for Mockoon CLI in a user-friendly way.

## Backend Implementations

This project has two backend implementations:
1. **TypeScript/Node.js** - Located in the `backend/` directory
2. **Go with Gin Gonic** - Located in the `backend_new/` directory 

The Go implementation is a complete rewrite of the TypeScript backend, offering the same functionality with improved performance and memory usage.

---

## 🚀 Key Feature: Automatic Proxy Routing with Traefik

The main highlight of this project is the **auto proxy integration with Traefik**. For every mock service started via this panel, a unique route is automatically created and accessible like this:

```
http://HOST_NAME:PORT/SERVICE_PORT
```

✅ Example:
```
http://localhost:8080/4000
```

This enables seamless access to each mock service without manual routing configuration.

![Auto Proxy Demo](./images/demo.png)

---

## 🧪 Disclaimer

- This is an experimental and personal side project.
- Stability and security are not guaranteed.
- Feel free to fork and adapt it to suit your own use case.
- If you encounter issues or bugs while using this tool, you're welcome to contribute fixes or improvements. The author holds no responsibility for any risks or damages resulting from its use.

---

## 🖼️ UI Screenshots

### 🔐 Login Page
![Login Page](./images/login.png)

### 📊 Dashboard View
![Dashboard](./images/dashboard.png)

---

## 🎯 Features

- [x] **Auto Proxy Routing via Traefik**
- [x] **Login Page** for admin authentication.
- [x] **Dashboard** to view and manage Mockoon projects.
- [x] **Auto Git Sync via SSH** for automatic configuration synchronization.

![Git Sync](./images/git.png)

---

## 🛠️ Upcoming Features

- [ ] Realtime logs on the dashboard.
- [ ] Display and manage routes.
- [ ] Update responses directly via UI.
- [ ] Change admin password via DB or dashboard.
- [ ] Toggle proxy mode.
- [ ] Support for multiple proxy URLs.
- [ ] (More to come...) Full control features via the UI.

---

## 📦 Getting Started

Run the project using Docker with the following setup:

- **Port Mapping**:
    - Map host port `8080` to container port `80`.
- **Volume Mount**:
    - Mount volume `/app/configs/` to your local `./mockoon-control-panel/`.

You can then access the control panel via: [http://localhost:8080](http://localhost:8080)

```bash
docker run -d \
  -p 8080:80 \
  -v $(pwd)/mockoon-control-panel/:/app/configs/ \
  ghcr.io/yogasw/mockoon-control-panel:latest
```

---

## 🗄️ Database Configuration

By default, the application uses a **SQLite** database stored within the container. If you'd like to use a different database like PostgreSQL, you can override the database connection string by setting the `DATABASE_URL` environment variable.

### Example using PostgreSQL:

```bash
docker run -d \
  -p 8080:80 \
  -v $(pwd)/mockoon-control-panel/:/app/configs/ \
  -e DATABASE_URL=postgresql://root:root@172.17.0.2:5432/mockoon \
  ghcr.io/yogasw/mockoon-control-panel:latest
```

> Make sure your PostgreSQL instance is running and accessible from the container.

---

## 🐳 Docker Image

You can pull the latest image from GitHub Container Registry:

```bash
docker pull ghcr.io/yogasw/mockoon-control-panel:latest
```
