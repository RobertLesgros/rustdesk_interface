# RustDesk API - Version Sécurisée

Serveur API pour RustDesk avec améliorations de sécurité et localisation française.

## Table des Matières

- [Fonctionnalités](#fonctionnalités)
- [Améliorations de Sécurité](#améliorations-de-sécurité)
- [Prérequis](#prérequis)
- [Installation](#installation)
- [Configuration](#configuration)
- [Déploiement Docker](#déploiement-docker)
- [Logs d'Audit](#logs-daudit)
- [Variables d'Environnement](#variables-denvironnement)
- [Fonctionnement Hors-Ligne](#fonctionnement-hors-ligne)

---

## Fonctionnalités

- **API complète** pour la gestion des utilisateurs RustDesk
- **Interface d'administration** web intégrée
- **Authentification** : mot de passe, OAuth2 (GitHub, Google, OIDC), LDAP/Active Directory
- **Gestion des groupes** et des permissions
- **Carnet d'adresses** partagé
- **Logs d'audit** structurés en JSON
- **Localisation française** complète

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
# Construire l'image
docker build -f Dockerfile.dev -t rustdesk-interface:latest .

# Lancer avec docker-compose
docker-compose up -d
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
  cors-allowed-origins: "https://votre-domaine.fr"

gorm:
  type: "sqlite"

rustdesk:
  id-server: "votre-serveur:21116"
  relay-server: "votre-serveur:21117"
  api-server: "http://127.0.0.1:21114"
  key-file: "/app/data/id_ed25519.pub"

audit:
  enabled: true
  file-path: "./runtime/audit.log"
```

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
