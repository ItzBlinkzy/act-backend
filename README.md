# ACT – AI-Driven Agentic Trading Application  

## Overview  

As part of our third-year group project, we developed an **AI-driven agentic trading application**. The system consists of a **React (Vite) frontend** and a **Go (Echo) backend**, working together to provide portfolio insights, AI-driven recommendations, and client/asset management.  

🌐 **Live Frontend:** [https://act-frontend.netlify.app](https://act-frontend.netlify.app)  
🚀 **Frontend Github:** [https://github.com/ItzBlinkzy/act-frontend](https://github.com/ItzBlinkzy/act-frontend)
---

## Architecture  

Our system is split into two main components:  

- **Frontend:** React + Vite (Netlify deployment)  
- **Backend:** Go + Echo (REST API, deployed separately)  

![Architecture Diagram](https://github.com/user-attachments/assets/bb904a95-9d83-45f7-ba39-9e156c5d0fca)  

---

## Getting Started  

### Prerequisites  

- [Go](https://go.dev/dl/) (1.21+)  
- [Git](https://git-scm.com/)  
- Optional: [Docker](https://www.docker.com/) for containerized setup  


### Running the Backend (Go + Echo)

1. Clone the repository and navigate to the backend folder:

```
git clone https://github.com/your-org/act-backend.git
cd act-backend
```

2. Download dependencies:
```
go mod tidy
```


3. Run the application:

```
go run main.go
```
---
