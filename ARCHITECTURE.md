# Architecture de RustDesk Interface

Ce document décrit l'architecture complète du projet RustDesk Interface et comment ses différents composants s'intègrent.

## Table des Matières

- [Vue d'ensemble](#vue-densemble)
- [Composants du Projet](#composants-du-projet)
- [Architecture Technique](#architecture-technique)
- [Intégration du Frontend Web](#intégration-du-frontend-web)
- [Flux de Données](#flux-de-données)
- [Ports et Services](#ports-et-services)

---

## Vue d'ensemble

RustDesk Interface est une **solution complète de gestion** pour RustDesk, le logiciel de bureau à distance open source. Ce projet fournit :

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        RUSTDESK INTERFACE                                    │
│                                                                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────────┐ │
│  │   Backend API   │  │  Admin Panel    │  │     Web Client RustDesk     │ │
│  │      (Go)       │  │  (React/Vue)    │  │       (Flutter/WASM)        │ │
│  │                 │  │                 │  │                             │ │
│  │ • API REST      │  │ • Gestion users │  │ • Bureau à distance via    │ │
│  │ • Auth JWT      │  │ • Gestion peers │  │   navigateur web           │ │
│  │ • OAuth/LDAP    │  │ • Carnet adres. │  │ • Synchronisation auto     │ │
│  │ • Audit logs    │  │ • Logs/Stats    │  │ • Partage temporaire       │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────────────────┘ │
│           │                   │                        │                    │
│           └───────────────────┴────────────────────────┘                    │
│                               │                                             │
│                        Port 21114                                           │
└─────────────────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         RUSTDESK SERVER                                      │
│                                                                              │
│  ┌─────────────────┐              ┌─────────────────┐                       │
│  │  ID Server      │              │  Relay Server   │                       │
│  │   (hbbs)        │              │    (hbbr)       │                       │
│  │  Port 21116     │              │   Port 21117    │                       │
│  └─────────────────┘              └─────────────────┘                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       CLIENTS RUSTDESK                                       │
│                                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Windows   │  │    Linux    │  │    macOS    │  │   Android   │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Composants du Projet

### 1. Backend API (Go) - Ce Repository

**Localisation** : `/` (racine du projet)

Le cœur du système, écrit en Go, qui fournit :

| Fonctionnalité | Description |
|----------------|-------------|
| **API REST** | Endpoints pour clients RustDesk et admin |
| **Authentification** | JWT, OAuth2 (GitHub, Google, OIDC), LDAP |
| **Gestion utilisateurs** | CRUD utilisateurs, groupes, permissions |
| **Carnet d'adresses** | Gestion des peers et partage |
| **Logs d'audit** | Journalisation sécurisée JSON |
| **Rate limiting** | Protection contre les attaques brute-force |

**Structure du code Go** :
```
├── cmd/apimain.go        # Point d'entrée
├── http/
│   ├── controller/       # Handlers HTTP
│   │   ├── admin/       # API Administration
│   │   ├── api/         # API Client PC
│   │   └── web/         # Interface Web
│   ├── middleware/       # JWT, CORS, Rate Limit
│   └── router/          # Définition des routes
├── service/              # Logique métier
├── model/                # Modèles de données
├── lib/                  # Bibliothèques utilitaires
└── config/               # Configuration
```

### 2. Frontend Admin Panel - Intégré

**Localisation** : `frontend/` → compilé vers `resources/admin/`

Interface d'administration web moderne :

- **Technologies** : Vue.js 3.5, Vite 6.3, Element Plus, Pinia
- **Fonctionnalités** :
  - Tableau de bord avec statistiques
  - Gestion des utilisateurs et groupes
  - Gestion des appareils (peers)
  - Configuration OAuth (désactivable)
  - Visualisation des logs
  - Contrôle serveur RustDesk
  - Support multilingue (fr, en, es, ru, ko)

> **Source originale** : https://github.com/RobertLesgros/rustdesk_interface_web

### 3. Web Client RustDesk - Inclus

**Localisation** : `resources/web/` et `resources/web2/`

Client de bureau à distance fonctionnant directement dans le navigateur :

- **Technologie** : Flutter compilé en WebAssembly
- **Fonctionnalités** :
  - Connexion bureau à distance via navigateur
  - Transfert de fichiers
  - Audio bidirectionnel
  - Partage d'écran

---

## Architecture Technique

### Stack Technologique

```yaml
Backend:
  Language: Go 1.24
  Framework: Gin v1.11.0
  ORM: GORM v1.25.7
  Auth: JWT, bcrypt
  Documentation: Swagger/OpenAPI

Frontend Admin:
  Framework: React/Vue.js
  Build: Node.js 18, npm
  UI: Modern responsive design

Web Client:
  Framework: Flutter Web
  Compilation: WebAssembly (WASM)
  Codecs: OGV.js, Opus audio

Database:
  Primary: SQLite (embedded)
  Optional: MySQL, PostgreSQL

Deployment:
  Container: Docker, Alpine 3.19
  Orchestration: Docker Compose
  Security: Non-root, minimal capabilities
```

### Structure des Répertoires

```
rustdesk_interface/
├── cmd/                  # Point d'entrée Go
├── conf/                 # Fichiers de configuration
│   └── config.yaml      # Configuration principale
├── docs/                 # Documentation API Swagger
├── frontend/             # Sources frontend (à cloner)
├── http/                 # Couche HTTP
│   ├── controller/      # Contrôleurs
│   ├── middleware/      # Middlewares
│   ├── request/         # Validation requêtes
│   └── response/        # Structures réponses
├── lib/                  # Bibliothèques internes
│   ├── audit/           # Logging d'audit
│   ├── cache/           # Cache
│   ├── jwt/             # Gestion JWT
│   └── orm/             # Abstraction BDD
├── model/                # Modèles de données
├── resources/            # Ressources statiques
│   ├── admin/           # Frontend compilé (généré)
│   ├── web/             # Client web RustDesk
│   ├── web2/            # Client web alternatif
│   ├── templates/       # Templates HTML
│   └── i18n/            # Traductions
├── scripts/              # Scripts utilitaires
│   ├── prepare-offline.sh
│   └── install-offline.sh
├── service/              # Services métier
├── Dockerfile.dev        # Image Docker
├── docker-compose.yaml   # Orchestration
└── README.md             # Documentation
```

---

## Frontend Web Admin

### Le Frontend est Intégré

Le frontend d'administration Vue.js est **directement inclus** dans ce repository dans le dossier `frontend/`. Aucune action supplémentaire n'est nécessaire.

### Développement du Frontend

```bash
cd frontend

# Installer les dépendances
npm install

# Mode développement avec hot reload
npm run dev
# → http://localhost:5173

# Build production
npm run build
# → Fichiers dans dist/
```

### Structure du Frontend

```
frontend/
├── src/
│   ├── views/              # Pages Vue.js
│   │   ├── login/         # Authentification
│   │   ├── my/            # Espace utilisateur
│   │   ├── user/          # Gestion utilisateurs (admin)
│   │   ├── group/         # Gestion groupes
│   │   ├── peer/          # Gestion appareils
│   │   └── oauth/         # Configuration OAuth
│   ├── components/         # Composants réutilisables
│   ├── utils/
│   │   ├── i18n/          # Traductions (fr, en, es, ru, ko)
│   │   ├── auth.js        # Gestion authentification
│   │   └── request.js     # Appels API axios
│   └── router/             # Configuration routes
├── public/                 # Fichiers statiques
├── package.json            # Dépendances npm
└── vite.config.js          # Configuration Vite
```

### Processus de Build Docker

Le `Dockerfile.dev` compile automatiquement le frontend :

```dockerfile
# Stage 1: Build Backend Go
FROM crazymax/xgo:1.24 AS builder-backend
# ... compilation Go ...

# Stage 2: Build Frontend (automatique)
FROM node:18-alpine AS builder-admin-frontend
WORKDIR /frontend
COPY frontend/ .
RUN npm install && npm run build

# Stage 3: Image Finale
FROM alpine:3.19
# Copie backend + frontend compilé
COPY --from=builder-admin-frontend /frontend/dist/ /app/resources/admin/
```

---

## Flux de Données

### Authentification

```
Client RustDesk                    API Server                      Database
      │                                │                               │
      │ POST /api/login               │                               │
      │ {username, password}          │                               │
      │──────────────────────────────▶│                               │
      │                                │ SELECT user WHERE...         │
      │                                │──────────────────────────────▶│
      │                                │◀──────────────────────────────│
      │                                │ Verify bcrypt                 │
      │                                │ Generate JWT                  │
      │ {token, user_info}            │                               │
      │◀──────────────────────────────│                               │
```

### Gestion des Peers

```
Admin Panel                        API Server                    RustDesk Server
      │                                │                               │
      │ GET /_admin/peers             │                               │
      │──────────────────────────────▶│                               │
      │                                │ Query database               │
      │                                │ Query ID server              │
      │                                │──────────────────────────────▶│
      │                                │◀──────────────────────────────│
      │ {peers: [...]}                │                               │
      │◀──────────────────────────────│                               │
```

---

## Ports et Services

| Service | Port | Description |
|---------|------|-------------|
| **RustDesk Interface API** | 21114 | API REST + Admin + Web Client |
| RustDesk ID Server (hbbs) | 21115 | NAT traversal |
| RustDesk ID Server (hbbs) | 21116 | ID registration |
| RustDesk Relay (hbbr) | 21117 | Relay connections |
| RustDesk ID Server | 21118 | Web socket |
| RustDesk ID Server | 21119 | Web socket |

### URLs Importantes

| URL | Description |
|-----|-------------|
| `http://server:21114/api/` | API pour clients RustDesk |
| `http://server:21114/_admin/` | Interface d'administration |
| `http://server:21114/webclient/` | Client web RustDesk |
| `http://server:21114/swagger/` | Documentation API |
| `http://server:21114/_admin/swagger/` | Documentation Admin API |

---

## Sécurité

### Mesures Implémentées

1. **Authentification**
   - JWT tokens avec expiration configurable
   - bcrypt pour le hashage des mots de passe
   - OAuth2 / LDAP optionnels

2. **Protection des Endpoints**
   - Rate limiting (10 req/min sur login)
   - Bannissement automatique après tentatives échouées
   - Validation des entrées

3. **Headers HTTP**
   - `X-Content-Type-Options: nosniff`
   - `X-Frame-Options: DENY`
   - `Content-Security-Policy`
   - CORS configurable

4. **Docker**
   - Exécution non-root (UID 1000)
   - Filesystem read-only optionnel
   - Healthcheck intégré

5. **Audit**
   - Logs JSON structurés
   - Traçabilité des connexions
   - Alertes de sécurité

---

## Développement

### Prérequis

```bash
# Backend
go version  # >= 1.24

# Frontend
node --version  # >= 18
npm --version

# Docker
docker --version
docker-compose --version
```

### Lancer en Développement

```bash
# 1. Backend seul
go run cmd/apimain.go

# 2. Frontend seul
cd frontend
npm install
npm run dev

# 3. Avec Docker
docker-compose -f docker-compose-dev.yaml up
```

### Générer la Documentation API

```bash
# Installer swag
go install github.com/swaggo/swag/cmd/swag@latest

# Générer les docs
swag init -g cmd/apimain.go --output docs/api --instanceName api
swag init -g cmd/apimain.go --output docs/admin --instanceName admin
```

---

## Contribuer

1. Fork le repository
2. Créer une branche (`git checkout -b feature/ma-feature`)
3. Commiter les changements (`git commit -m 'Add feature'`)
4. Pousser (`git push origin feature/ma-feature`)
5. Ouvrir une Pull Request

---

## Licence

MIT License - Voir [LICENSE](LICENSE)
