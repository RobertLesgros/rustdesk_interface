# Installation Hors Ligne de RustDesk Interface

Ce guide explique comment installer RustDesk Interface sur un serveur **sans connexion Internet**.

## Table des Matières

- [Vue d'ensemble](#vue-densemble)
- [Phase 1 : Préparation (machine connectée)](#phase-1--préparation-machine-connectée)
- [Phase 2 : Installation (machine hors ligne)](#phase-2--installation-machine-hors-ligne)
- [Configuration](#configuration)
- [Démarrage et vérification](#démarrage-et-vérification)
- [Dépannage](#dépannage)
- [Mises à jour hors ligne](#mises-à-jour-hors-ligne)

---

## Vue d'ensemble

L'installation hors ligne se déroule en deux phases :

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           PHASE 1 : PRÉPARATION                              │
│                        (Machine avec Internet)                               │
│                                                                              │
│  1. Cloner le dépôt                                                          │
│  2. Exécuter le script de préparation                                        │
│  3. Transférer le bundle vers la machine cible                               │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
                    ┌────────────────────────────────┐
                    │      Transfert USB/SCP/etc.    │
                    │  rustdesk-interface-offline-   │
                    │       bundle.tar.gz            │
                    └────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           PHASE 2 : INSTALLATION                             │
│                        (Machine SANS Internet)                               │
│                                                                              │
│  1. Extraire le bundle                                                       │
│  2. Charger l'image Docker                                                   │
│  3. Configurer et démarrer                                                   │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Phase 1 : Préparation (machine connectée)

### Prérequis sur la machine de préparation

- **Docker** (version 20.10 ou supérieure)
- **Git** (pour cloner le dépôt)
- **Connexion Internet**
- Environ **2 Go** d'espace disque libre

### Étape 1.1 : Cloner le dépôt

```bash
git clone https://github.com/RobertLesgros/rustdesk_interface.git
cd rustdesk_interface
```

### Étape 1.2 : Exécuter le script de préparation

```bash
./scripts/prepare-offline.sh
```

Ce script effectue automatiquement :
1. Construit l'image Docker complète (backend Go + frontend)
2. Exporte l'image dans un fichier tar
3. Crée un bundle avec tous les fichiers nécessaires

**Note** : Le `Dockerfile.dev` gère automatiquement le téléchargement et la compilation du frontend. Tout est inclus dans une seule commande.

### Étape 1.3 : Vérifier les fichiers créés

```bash
ls -la
```

Vous devriez voir :
- `rustdesk-interface-offline.tar` - Image Docker (~200-300 Mo)
- `offline-bundle/` - Dossier avec tous les fichiers
- `rustdesk-interface-offline-bundle.tar.gz` - Archive complète pour transfert

### Étape 1.4 : Transférer le bundle

Copiez l'archive vers la machine cible par le moyen de votre choix :

**Via USB :**
```bash
cp rustdesk-interface-offline-bundle.tar.gz /media/usb/
```

**Via SCP :**
```bash
scp rustdesk-interface-offline-bundle.tar.gz user@machine-cible:/tmp/
```

**Via partage réseau :**
```bash
cp rustdesk-interface-offline-bundle.tar.gz /mnt/partage/
```

---

## Phase 2 : Installation (machine hors ligne)

### Prérequis sur la machine cible

- **Docker** installé (peut être installé hors ligne si nécessaire)
- **Docker Compose** (inclus avec Docker sur les versions récentes)
- Environ **500 Mo** d'espace disque libre

### Étape 2.1 : Extraire le bundle

```bash
cd /tmp  # ou là où vous avez copié l'archive
tar -xzvf rustdesk-interface-offline-bundle.tar.gz
cd offline-bundle
```

### Étape 2.2 : Exécuter l'installateur automatique

```bash
sudo ./install-offline.sh --start
```

Options disponibles :
- `--start` : Démarre le service immédiatement après installation
- `--install-dir /chemin` : Choisir un répertoire d'installation personnalisé
- `--skip-load` : Ne pas recharger l'image Docker (si déjà chargée)

### Étape 2.3 : Installation manuelle (alternative)

Si vous préférez une installation manuelle :

```bash
# 1. Charger l'image Docker
docker load -i rustdesk-interface-offline.tar

# 2. Créer le répertoire d'installation
sudo mkdir -p /opt/rustdesk-interface/{conf,data,runtime}
cd /opt/rustdesk-interface

# 3. Copier les fichiers
cp /tmp/offline-bundle/docker-compose.yml .
cp /tmp/offline-bundle/.env.production.example .env.production
cp -r /tmp/offline-bundle/conf/* ./conf/

# 4. Configurer les permissions
sudo chown -R 1000:1000 data runtime
sudo chmod 600 .env.production

# 5. Éditer la configuration
nano .env.production

# 6. Démarrer
docker-compose up -d
```

---

## Configuration

### Fichier .env.production

Éditez le fichier `.env.production` avec vos paramètres :

```env
# Timezone
TZ=Europe/Paris

# Langue
RUSTDESK_API_LANG=fr

# Mode API
RUSTDESK_API_GIN_MODE=release

# Clé JWT (IMPORTANT: générez une clé unique)
# Si laissée à la valeur par défaut, le script en génère une automatiquement
RUSTDESK_API_JWT_KEY=votre_cle_secrete_64_caracteres_hex

# Configuration RustDesk Server
# Adaptez ces valeurs à votre réseau
RUSTDESK_API_RUSTDESK_ID_SERVER=192.168.1.10:21116
RUSTDESK_API_RUSTDESK_RELAY_SERVER=192.168.1.10:21117
```

### Fichier conf/config.yaml

Pour une configuration plus avancée, éditez `conf/config.yaml` :

```yaml
lang: "fr"

app:
  web-client: 1
  register: false
  captcha-threshold: 3
  ban-threshold: 10

gin:
  api-addr: "0.0.0.0:21114"
  mode: "release"
  cors-allowed-origins: "*"

gorm:
  type: "sqlite"

rustdesk:
  id-server: "192.168.1.10:21116"
  relay-server: "192.168.1.10:21117"
  api-server: "http://192.168.1.10:21114"

audit:
  enabled: true
  file-path: "./runtime/audit.log"
```

### Génération d'une clé JWT

Si vous devez générer une nouvelle clé JWT :

```bash
# Méthode 1: OpenSSL
openssl rand -hex 32

# Méthode 2: /dev/urandom
head -c 32 /dev/urandom | xxd -p
```

---

## Démarrage et vérification

### Démarrer le service

```bash
cd /opt/rustdesk-interface  # ou votre répertoire d'installation
docker-compose up -d
```

### Vérifier l'état

```bash
# État du conteneur
docker-compose ps

# Logs en temps réel
docker-compose logs -f

# Santé du service
docker inspect --format='{{.State.Health.Status}}' rustdesk-interface
```

### Accéder à l'interface

- **API** : http://VOTRE_IP:21114/api
- **Administration** : http://VOTRE_IP:21114/admin
- **Documentation API** : http://VOTRE_IP:21114/api/swagger/index.html

### Identifiants par défaut

- **Utilisateur** : `admin`
- **Mot de passe** : Affiché dans les logs au premier démarrage

Pour voir le mot de passe initial :
```bash
docker-compose logs | grep -i password
```

### Changer le mot de passe admin

```bash
docker exec -it rustdesk-interface ./apimain reset-admin-pwd NouveauMotDePasse
```

---

## Dépannage

### L'image ne se charge pas

```bash
# Vérifier l'intégrité du fichier
file rustdesk-interface-offline.tar

# Vérifier l'espace disque
df -h

# Recharger avec verbose
docker load -i rustdesk-interface-offline.tar --quiet=false
```

### Le conteneur ne démarre pas

```bash
# Voir les logs détaillés
docker-compose logs --tail=100

# Vérifier les permissions
ls -la data/ runtime/ conf/

# Corriger les permissions
sudo chown -R 1000:1000 data/ runtime/
```

### Erreur "permission denied"

```bash
# Le conteneur tourne en tant qu'utilisateur 1000
sudo chown -R 1000:1000 /opt/rustdesk-interface/data
sudo chown -R 1000:1000 /opt/rustdesk-interface/runtime
```

### Le health check échoue

```bash
# Tester manuellement
curl http://localhost:21114/api/health

# Vérifier que le port est ouvert
netstat -tlnp | grep 21114
```

### Réinitialiser complètement

```bash
cd /opt/rustdesk-interface
docker-compose down -v
rm -rf data/* runtime/*
docker-compose up -d
```

---

## Mises à jour hors ligne

Pour mettre à jour RustDesk Interface :

### Sur la machine connectée

```bash
cd rustdesk_interface
git pull
./scripts/prepare-offline.sh
```

### Sur la machine hors ligne

```bash
# Arrêter le service
cd /opt/rustdesk-interface
docker-compose down

# Sauvegarder les données
cp -r data data.backup
cp .env.production .env.production.backup

# Charger la nouvelle image
docker load -i rustdesk-interface-offline.tar

# Redémarrer
docker-compose up -d

# Vérifier
docker-compose logs -f
```

---

## Sauvegarde et restauration

### Sauvegarder

```bash
cd /opt/rustdesk-interface
tar -czvf backup-$(date +%Y%m%d).tar.gz data/ conf/ .env.production
```

### Restaurer

```bash
cd /opt/rustdesk-interface
docker-compose down
tar -xzvf backup-YYYYMMDD.tar.gz
docker-compose up -d
```

---

## Service systemd (optionnel)

Le script d'installation crée automatiquement un service systemd si exécuté en root.

### Activer le démarrage automatique

```bash
sudo systemctl enable rustdesk-interface
```

### Commandes de gestion

```bash
sudo systemctl start rustdesk-interface    # Démarrer
sudo systemctl stop rustdesk-interface     # Arrêter
sudo systemctl restart rustdesk-interface  # Redémarrer
sudo systemctl status rustdesk-interface   # État
```

---

## Architecture réseau

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Réseau Interne Sécurisé                            │
│                                                                              │
│   ┌─────────────┐        ┌─────────────────────┐        ┌─────────────┐    │
│   │   Clients   │───────▶│   RustDesk Server   │◀───────│  RustDesk   │    │
│   │  RustDesk   │        │   (hbbs + hbbr)     │        │  Interface  │    │
│   └─────────────┘        │                      │        │    (API)    │    │
│                          │  Ports: 21115-21119  │        │ Port: 21114 │    │
│                          └─────────────────────┘        └─────────────┘    │
│                                                                              │
│   Aucune connexion Internet requise                                         │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Support

Pour toute question ou problème :
- Consultez les logs : `docker-compose logs`
- Vérifiez les logs d'audit : `cat runtime/audit.log`
- Ouvrez une issue sur le dépôt GitHub

---

## Licence

MIT License - Voir le fichier LICENSE pour les détails.
