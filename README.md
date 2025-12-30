# RustDesk API - Version Sécurisée

Serveur API complet pour RustDesk avec interface d'administration web, client de bureau à distance via navigateur, et améliorations de sécurité.

<div align=center>
<img src="https://img.shields.io/badge/golang-1.24-blue"/>
<img src="https://img.shields.io/badge/gin-v1.11.0-lightBlue"/>
<img src="https://img.shields.io/badge/gorm-v1.25.7-green"/>
<img src="https://img.shields.io/badge/swag-v1.16.3-yellow"/>
</div>

## Qu'est-ce que RustDesk Interface ?

**RustDesk Interface** est une solution de gestion complète pour [RustDesk](https://rustdesk.com/), le logiciel de bureau à distance open source. Ce projet fournit :

- **Un backend API** (Go) pour gérer les utilisateurs, appareils et permissions
- **Une interface d'administration web** pour les administrateurs IT
- **Un client web RustDesk** pour accéder aux bureaux à distance depuis un navigateur
- **Une infrastructure sécurisée** conçue pour les environnements d'entreprise hors-ligne

### Architecture du Projet

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          RUSTDESK INTERFACE                                  │
│                                                                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────────┐ │
│  │   Backend API   │  │  Admin Panel    │  │     Web Client RustDesk     │ │
│  │      (Go)       │  │  (React/Vue)    │  │       (Flutter/WASM)        │ │
│  │   Ce repo       │  │  Repo séparé    │  │      Inclus ici             │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────────────────┘ │
│                               │                                             │
│                        Port 21114                                           │
└─────────────────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         RUSTDESK SERVER                                      │
│           (hbbs ID Server port 21116 + hbbr Relay port 21117)               │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Documentation

| Document | Description |
|----------|-------------|
| [README.md](README.md) | Ce fichier - guide principal en français |
| [README_EN.md](README_EN.md) | Documentation en anglais |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Architecture technique détaillée |
| [INSTALL_OFFLINE.md](INSTALL_OFFLINE.md) | Installation hors-ligne |

## Table des Matières

- [Fonctionnalités](#fonctionnalités)
- [Composants du Projet](#composants-du-projet)
- [Améliorations de Sécurité](#améliorations-de-sécurité)
- [Prérequis](#prérequis)
- [Installation](#installation)
- [Intégration du Frontend Web](#intégration-du-frontend-web)
- [Configuration](#configuration)
- [Déploiement Docker](#déploiement-docker)
- [Logs d'Audit](#logs-daudit)
- [Variables d'Environnement](#variables-denvironnement)
- [Fonctionnement Hors-Ligne](#fonctionnement-hors-ligne)

---

## Fonctionnalités

- **API complète** pour la gestion des utilisateurs RustDesk
- **Interface d'administration** web intégrée
- **Client web RustDesk** pour accès bureau à distance via navigateur
- **Authentification** : mot de passe, OAuth2 (GitHub, Google, OIDC), LDAP/Active Directory
- **Gestion des groupes** et des permissions
- **Carnet d'adresses** partagé
- **Logs d'audit** structurés en JSON
- **Localisation française** complète
- **Installation hors-ligne** pour environnements sécurisés

---

## Composants du Projet

Ce projet est composé de **trois parties** :

### 1. Backend API (Go) - Ce Repository

Le cœur du système, fournissant l'API REST et les services :

```
rustdesk_interface/
├── cmd/apimain.go        # Point d'entrée
├── http/controller/      # Handlers API
├── service/              # Logique métier
├── model/                # Modèles de données
├── resources/            # Ressources statiques
│   ├── admin/           # Frontend compilé (voir ci-dessous)
│   ├── web/             # Client web RustDesk
│   └── i18n/            # Traductions
└── conf/                 # Configuration
```

### 2. Frontend Admin Panel - Repository Séparé

**L'interface d'administration web est dans un repository séparé :**

> **Repository** : https://github.com/RobertLesgros/rustdesk_interface_web

Cette séparation permet :
- Un développement indépendant du frontend et backend
- La personnalisation de l'interface sans modifier le backend
- Des cycles de release séparés

**Important** : Pour construire l'image Docker complète, vous devez d'abord cloner le frontend. Voir [Intégration du Frontend Web](#intégration-du-frontend-web).

### 3. Client Web RustDesk - Inclus

Le client de bureau à distance web (Flutter/WASM) est **déjà inclus** dans ce repository dans `resources/web/`. Il permet d'accéder aux bureaux à distance directement depuis un navigateur.

---

## Améliorations de Sécurité

Cette version inclut plusieurs corrections de sécurité critiques :

### 1. Protection contre l'Injection LDAP
- **Fichier** : `service/ldap.go`
- Les entrées utilisateur sont maintenant échappées avec `ldap.EscapeFilter()` pour prévenir les attaques par injection LDAP.

### 2. Protection contre la Traversée de Chemin (Path Traversal)
- **Fichier** : `http/controller/admin/file.go`
- Les noms de fichiers sont sanitisés pour empêcher l'accès à des fichiers en dehors du répertoire prévu.

### 3. En-têtes de Sécurité HTTP
- **Fichier** : `http/middleware/cors.go`
- Ajout des en-têtes de sécurité :
  - `X-Content-Type-Options: nosniff`
  - `X-Frame-Options: DENY`
  - `X-XSS-Protection: 1; mode=block`
  - `Referrer-Policy: strict-origin-when-cross-origin`
  - `Permissions-Policy` restrictive

### 4. Configuration CORS Sécurisée
- Les origines autorisées sont maintenant configurables via `cors-allowed-origins` dans `config.yaml`
- Par défaut, seules les origines spécifiées sont autorisées

### 5. Suppression du Fallback MD5
- **Fichier** : `utils/password.go`
- Le support des anciens hashes MD5 a été supprimé
- Seul bcrypt est maintenant accepté pour les mots de passe

### 6. Limitation de Débit (Rate Limiting)
- **Fichier** : `http/middleware/limiter.go`
- Protection contre les attaques par force brute :
  - 10 requêtes par minute sur les opérations sensibles
  - Bannissement automatique après trop de tentatives de connexion échouées

### 7. Aucune Connexion Externe
- Firebase Analytics a été supprimé du client web
- L'application fonctionne entièrement en mode hors-ligne
- Les services OAuth/LDAP sont optionnels et désactivés par défaut

### 8. Logs d'Audit Structurés
- **Fichier** : `lib/audit/audit.go`
- Journalisation JSON des événements de sécurité :
  - Connexions réussies/échouées
  - Modifications de mots de passe
  - Création/suppression d'utilisateurs
  - Tentatives d'accès refusées
  - Dépassements de limite de débit

### 9. Docker Sécurisé
- Exécution en tant qu'utilisateur non-root
- Système de fichiers en lecture seule
- Capacités Linux minimales
- Healthcheck intégré

---

## Prérequis

- **Go 1.24+** (pour la compilation)
- **Docker** et **Docker Compose** (pour le déploiement)
- **SQLite** (par défaut) ou **MySQL/PostgreSQL**

---

## Installation

### Option 1 : Compilation depuis les Sources

```bash
# Cloner le dépôt
git clone https://github.com/votre-repo/rustdesk-interface.git
cd rustdesk-interface

# Compiler
go mod tidy
go build -o apimain cmd/apimain.go

# Lancer
./apimain
```

### Option 2 : Docker (Recommandé)

```bash
# Construire l'image (nécessite le frontend, voir section suivante)
docker build -f Dockerfile.dev -t rustdesk-interface:latest .

# Lancer avec docker-compose
docker-compose up -d
```

---

## Intégration du Frontend Web

L'interface d'administration web est développée dans un **repository séparé**. Vous devez l'intégrer avant de construire l'image Docker.

### Pourquoi est-ce nécessaire ?

Le `Dockerfile.dev` s'attend à trouver le frontend dans le dossier `frontend/` :

```dockerfile
# Stage 2: Frontend Build
FROM node:18-alpine AS builder-admin-frontend
COPY frontend/ .  # ← Requiert le dossier frontend/
RUN npm install && npm run build
```

Sans ce dossier, la construction Docker **échouera**.

### Méthode 1 : Script Automatique (Recommandé)

```bash
# Ce script clone le frontend et construit tout
./scripts/prepare-offline.sh
```

### Méthode 2 : Clonage Manuel

```bash
# Cloner le frontend
git clone https://github.com/RobertLesgros/rustdesk_interface_web.git frontend

# Puis construire l'image Docker
docker build -f Dockerfile.dev -t rustdesk-interface:latest .
```

### Méthode 3 : Git Submodule (Pour Développeurs)

```bash
# Ajouter comme submodule
git submodule add https://github.com/RobertLesgros/rustdesk_interface_web.git frontend

# Pour cloner le projet avec les submodules
git clone --recursive https://github.com/RobertLesgros/rustdesk_interface.git
```

### Personnalisation du Frontend

Vous pouvez modifier le frontend avant la construction :

```bash
cd frontend
# Modifier les fichiers...
npm install
npm run build
# Les fichiers compilés seront dans dist/
```

---

## Configuration

Le fichier de configuration principal est `conf/config.yaml`.

### Configuration Minimale

```yaml
# Langue française
lang: "fr"

app:
  web-client: 1
  register: false
  captcha-threshold: 3
  ban-threshold: 10

gin:
  api-addr: "0.0.0.0:21114"
  mode: "release"
  cors-allowed-origins: "*"  # En LAN, peut rester ouvert

gorm:
  type: "sqlite"

rustdesk:
  id-server: "192.168.1.10:21116"    # Adaptez à votre LAN
  relay-server: "192.168.1.10:21117"
  api-server: "http://192.168.1.10:21114"
  key-file: "/app/data/id_ed25519.pub"

audit:
  enabled: true
  file-path: "./runtime/audit.log"
```

### Configuration LAN-Only avec LDAP/Active Directory

Pour un déploiement en réseau local sans connexion Internet, utilisez le fichier de configuration dédié :

```bash
# Copier la configuration LAN-only
cp conf/config-lan-only.yaml.example conf/config.yaml

# Éditer et adapter à votre environnement
nano conf/config.yaml
```

**Caractéristiques du mode LAN-only :**

| Fonctionnalité | État | Notes |
|----------------|------|-------|
| OAuth (GitHub, Google, OIDC) | **Désactivé** | Pas de connexion Internet requise |
| LDAP / Active Directory | **Activé** | Authentification via votre annuaire |
| Inscription publique | **Désactivée** | Sécurité renforcée |
| Client web RustDesk | **Activé** | Bureau à distance via navigateur |
| Logs d'audit | **Activé** | Traçabilité pour conformité |

### Configuration LDAP / Active Directory

Pour authentifier vos utilisateurs via Active Directory :

```yaml
ldap:
  enable: true
  url: "ldap://192.168.1.5:389"  # Ou ldaps:// pour TLS
  base-dn: "DC=entreprise,DC=local"
  bind-dn: "CN=rustdesk-svc,OU=Service Accounts,DC=entreprise,DC=local"
  bind-password: "VotreMotDePasse"

  user:
    base-dn: "OU=Utilisateurs,DC=entreprise,DC=local"
    username: "sAMAccountName"  # Pour Active Directory
    filter: "(&(objectClass=user)(objectCategory=person))"
    admin-group: "CN=RustDesk-Admins,OU=Groupes,DC=entreprise,DC=local"
    allow-group: "CN=RustDesk-Users,OU=Groupes,DC=entreprise,DC=local"
```

> **Note** : Créez les groupes `RustDesk-Admins` et `RustDesk-Users` dans votre Active Directory et ajoutez-y les utilisateurs autorisés.

### Désactivation de OAuth (GitHub, Google, OIDC)

OAuth est configuré via l'interface d'administration, pas dans `config.yaml`. Pour s'assurer qu'il reste désactivé :

1. **Ne configurez rien** dans l'interface admin (`/_admin/`) → Paramètres OAuth
2. Définissez dans `config.yaml` :
   ```yaml
   app:
     web-sso: false  # Désactive le SSO via OAuth
   ```
3. Si OAuth était configuré, supprimez les entrées via l'interface admin

---

## Déploiement Docker

### 1. Préparer les Secrets

```bash
# Copier le fichier d'exemple
cp .env.production.example .env.production

# Éditer et remplir les valeurs
nano .env.production
```

**Contenu de `.env.production` :**

```env
# Clé JWT (OBLIGATOIRE - générez avec: openssl rand -hex 32)
RUSTDESK_API_JWT_KEY=votre_cle_secrete_de_32_caracteres

# Configuration
TZ=Europe/Paris
RUSTDESK_API_LANG=fr
RUSTDESK_API_GIN_MODE=release

# RustDesk Server
RUSTDESK_API_RUSTDESK_ID_SERVER=192.168.1.66:21116
RUSTDESK_API_RUSTDESK_RELAY_SERVER=192.168.1.66:21117
```

### 2. Lancer le Service

```bash
# Démarrage
docker-compose up -d

# Vérifier les logs
docker-compose logs -f

# Arrêter
docker-compose down
```

### 3. Accéder à l'Interface

- **API** : http://votre-serveur:21114/api
- **Administration** : http://votre-serveur:21114/admin

**Identifiants par défaut** :
- Utilisateur : `admin`
- Mot de passe : affiché dans les logs au premier démarrage

---

## Logs d'Audit

Les logs d'audit sont enregistrés en format JSON dans `runtime/audit.log`.

### Format des Événements

```json
{
  "timestamp": "2024-01-15T14:30:22Z",
  "event_type": "LOGIN_SUCCESS",
  "severity": "INFO",
  "user_id": 1,
  "username": "admin",
  "client_ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "method": "POST",
  "path": "/admin/login",
  "message": "User logged in successfully",
  "success": true
}
```

### Types d'Événements

| Type | Description |
|------|-------------|
| `LOGIN_SUCCESS` | Connexion réussie |
| `LOGIN_FAILED` | Tentative de connexion échouée |
| `LOGOUT` | Déconnexion |
| `PASSWORD_CHANGED` | Changement de mot de passe |
| `USER_CREATED` | Création d'utilisateur |
| `USER_DELETED` | Suppression d'utilisateur |
| `ACCESS_DENIED` | Accès refusé |
| `RATE_LIMITED` | Limite de débit atteinte |
| `IP_BANNED` | Adresse IP bannie |
| `SECURITY_ALERT` | Alerte de sécurité |

---

## Variables d'Environnement

Toutes les valeurs de `config.yaml` peuvent être remplacées par des variables d'environnement.

| Variable | Description | Valeur par défaut |
|----------|-------------|-------------------|
| `RUSTDESK_API_LANG` | Langue | `fr` |
| `RUSTDESK_API_GIN_MODE` | Mode Gin | `release` |
| `RUSTDESK_API_GIN_API_ADDR` | Adresse d'écoute | `0.0.0.0:21114` |
| `RUSTDESK_API_JWT_KEY` | Clé secrète JWT | (requis) |
| `RUSTDESK_API_GORM_TYPE` | Type de BDD | `sqlite` |
| `RUSTDESK_API_AUDIT_ENABLED` | Activer les logs d'audit | `true` |
| `RUSTDESK_API_AUDIT_FILE_PATH` | Chemin des logs d'audit | `./runtime/audit.log` |

---

## Fonctionnement Hors-Ligne

Cette version est conçue pour fonctionner **sans connexion Internet**.

### Ce qui a été supprimé :
- ✅ Firebase Analytics (tracking Google)
- ✅ Dépendances CDN externes

### Ce qui reste optionnel (désactivé par défaut) :
- OAuth2 (GitHub, Google, OIDC)
- Stockage Alibaba Cloud OSS
- Proxy externe

### Architecture Réseau Recommandée

```
┌─────────────────────────────────────────────────────────────┐
│                     Réseau Interne                          │
│                                                             │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐   │
│  │  Clients    │────▶│  RustDesk   │────▶│  RustDesk   │   │
│  │  RustDesk   │     │  Server     │     │  API        │   │
│  └─────────────┘     │  (hbbs/hbbr)│     │  (ce projet)│   │
│                      └─────────────┘     └─────────────┘   │
│                                                             │
│  Ports : 21115, 21116, 21117, 21118, 21119 (RustDesk)      │
│  Port  : 21114 (API)                                        │
└─────────────────────────────────────────────────────────────┘
         │
         X (Pas de connexion Internet requise)
```

---

## Sécurité - Bonnes Pratiques

1. **Changez immédiatement le mot de passe admin** après le premier démarrage
2. **Générez une clé JWT forte** : `openssl rand -hex 32`
3. **Utilisez HTTPS** en production (via reverse proxy nginx/traefik)
4. **Sauvegardez régulièrement** le volume `/app/data`
5. **Surveillez les logs d'audit** pour détecter les anomalies
6. **Mettez à jour régulièrement** avec `docker-compose pull`

---

## Support

Pour signaler un problème ou demander une fonctionnalité :
- Ouvrez une issue sur le dépôt GitHub

---

## Licence

MIT License - Voir le fichier LICENSE pour plus de détails.
