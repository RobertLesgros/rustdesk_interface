#!/bin/bash
# =============================================================================
# Script d'installation hors ligne pour RustDesk Interface
# =============================================================================
# Ce script installe RustDesk Interface sur une machine SANS connexion Internet.
#
# Prérequis:
#   - Docker installé et en cours d'exécution
#   - L'image Docker exportée (rustdesk-interface-offline.tar)
#
# Usage:
#   ./install-offline.sh [OPTIONS]
#
# Options:
#   --image-file FILE  : Spécifie le fichier image Docker (défaut: rustdesk-interface-offline.tar)
#   --install-dir DIR  : Répertoire d'installation (défaut: /opt/rustdesk-interface)
#   --skip-load        : Ne pas charger l'image (déjà chargée)
#   --start            : Démarrer le service après installation
#   --help             : Affiche l'aide
# =============================================================================

set -e

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Variables par défaut
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGE_FILE="$SCRIPT_DIR/rustdesk-interface-offline.tar"
INSTALL_DIR="/opt/rustdesk-interface"
IMAGE_NAME="rustdesk-interface:offline"
SKIP_LOAD=false
START_SERVICE=false

# Parser les arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --image-file)
            IMAGE_FILE="$2"
            shift 2
            ;;
        --install-dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --skip-load)
            SKIP_LOAD=true
            shift
            ;;
        --start)
            START_SERVICE=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --image-file FILE  Spécifie le fichier image Docker"
            echo "  --install-dir DIR  Répertoire d'installation (défaut: /opt/rustdesk-interface)"
            echo "  --skip-load        Ne pas charger l'image Docker"
            echo "  --start            Démarrer le service après installation"
            echo "  --help             Affiche cette aide"
            exit 0
            ;;
        *)
            echo -e "${RED}Option inconnue: $1${NC}"
            exit 1
            ;;
    esac
done

# Fonctions
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[OK]${NC} $1"; }
log_warning() { echo -e "${YELLOW}[ATTENTION]${NC} $1"; }
log_error() { echo -e "${RED}[ERREUR]${NC} $1"; }

# Vérifier si root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        log_warning "Ce script doit être exécuté en tant que root pour certaines opérations."
        log_warning "Continuez sans droits root ? Certaines fonctionnalités peuvent échouer."
        read -p "Continuer ? (o/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Oo]$ ]]; then
            exit 1
        fi
    fi
}

# Vérifier Docker
check_docker() {
    log_info "Vérification de Docker..."

    if ! command -v docker &> /dev/null; then
        log_error "Docker n'est pas installé."
        log_info "Installez Docker avec: apt install docker.io docker-compose"
        exit 1
    fi

    if ! docker info &> /dev/null; then
        log_error "Docker n'est pas en cours d'exécution."
        log_info "Démarrez Docker avec: systemctl start docker"
        exit 1
    fi

    log_success "Docker est opérationnel"
}

# Charger l'image Docker
load_docker_image() {
    if [ "$SKIP_LOAD" = true ]; then
        log_warning "Chargement de l'image ignoré (--skip-load)"
        return
    fi

    log_info "Chargement de l'image Docker..."

    if [ ! -f "$IMAGE_FILE" ]; then
        log_error "Fichier image introuvable: $IMAGE_FILE"
        log_info "Assurez-vous que le fichier rustdesk-interface-offline.tar est présent."
        exit 1
    fi

    docker load -i "$IMAGE_FILE"

    log_success "Image chargée: $IMAGE_NAME"
}

# Créer la structure de répertoires
create_directories() {
    log_info "Création des répertoires..."

    mkdir -p "$INSTALL_DIR"/{conf,data,runtime}

    log_success "Répertoires créés dans: $INSTALL_DIR"
}

# Copier les fichiers de configuration
copy_config_files() {
    log_info "Copie des fichiers de configuration..."

    # Copier config.yaml si présent
    if [ -d "$SCRIPT_DIR/conf" ]; then
        cp -r "$SCRIPT_DIR/conf/"* "$INSTALL_DIR/conf/" 2>/dev/null || true
        log_success "Configuration copiée"
    fi

    # Copier docker-compose si présent
    if [ -f "$SCRIPT_DIR/docker-compose.yml" ]; then
        cp "$SCRIPT_DIR/docker-compose.yml" "$INSTALL_DIR/"
        log_success "docker-compose.yml copié"
    fi

    # Copier ou créer .env.production
    if [ -f "$SCRIPT_DIR/.env.production" ]; then
        cp "$SCRIPT_DIR/.env.production" "$INSTALL_DIR/"
        log_success ".env.production copié"
    elif [ -f "$SCRIPT_DIR/.env.production.example" ]; then
        cp "$SCRIPT_DIR/.env.production.example" "$INSTALL_DIR/.env.production"
        log_warning ".env.production créé depuis l'exemple. PENSEZ À LE MODIFIER !"
    fi
}

# Créer le docker-compose s'il n'existe pas
create_docker_compose() {
    if [ -f "$INSTALL_DIR/docker-compose.yml" ]; then
        return
    fi

    log_info "Création du fichier docker-compose.yml..."

    cat > "$INSTALL_DIR/docker-compose.yml" << 'EOF'
# Docker Compose pour RustDesk Interface (Installation Hors Ligne)
version: '3.8'

services:
  rustdesk-interface:
    image: rustdesk-interface:offline
    container_name: rustdesk-interface
    restart: unless-stopped

    # Sécurité: Utilisateur non-root
    user: "1000:1000"

    # Sécurité: Système de fichiers en lecture seule
    read_only: true

    # Sécurité: Supprimer toutes les capacités
    cap_drop:
      - ALL

    # Sécurité: Pas d'escalade de privilèges
    security_opt:
      - no-new-privileges:true

    # Ports
    ports:
      - "21114:21114"

    # Volumes
    volumes:
      - ./data:/app/data
      - ./runtime:/app/runtime
      - ./conf:/app/conf:ro

    # Tmpfs pour fichiers temporaires
    tmpfs:
      - /tmp:mode=1777,size=64M

    # Variables d'environnement
    env_file:
      - .env.production

    # Health check
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:21114/api/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s

    # Limites de ressources
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 128M

    # Logging
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "5"

# Pas de volumes nommés - on utilise des chemins locaux pour faciliter les sauvegardes
EOF

    log_success "docker-compose.yml créé"
}

# Générer une clé JWT
generate_jwt_key() {
    if [ -f "$INSTALL_DIR/.env.production" ]; then
        if grep -q "CHANGEZ_MOI" "$INSTALL_DIR/.env.production" 2>/dev/null; then
            log_warning "La clé JWT n'est pas configurée !"

            # Générer une nouvelle clé
            local new_key=$(openssl rand -hex 32 2>/dev/null || head -c 32 /dev/urandom | xxd -p)

            log_info "Génération d'une nouvelle clé JWT..."
            sed -i "s/CHANGEZ_MOI_AVEC_UNE_CLE_SECRETE/$new_key/" "$INSTALL_DIR/.env.production"
            log_success "Clé JWT générée et configurée"
        fi
    fi
}

# Définir les permissions
set_permissions() {
    log_info "Configuration des permissions..."

    # L'utilisateur 1000:1000 doit pouvoir écrire dans data et runtime
    chown -R 1000:1000 "$INSTALL_DIR/data" "$INSTALL_DIR/runtime" 2>/dev/null || true
    chmod -R 755 "$INSTALL_DIR/data" "$INSTALL_DIR/runtime" 2>/dev/null || true

    # Configuration en lecture seule
    chmod 644 "$INSTALL_DIR/conf/"* 2>/dev/null || true
    chmod 600 "$INSTALL_DIR/.env.production" 2>/dev/null || true

    log_success "Permissions configurées"
}

# Démarrer le service
start_service() {
    if [ "$START_SERVICE" = true ]; then
        log_info "Démarrage du service..."

        cd "$INSTALL_DIR"
        docker-compose up -d

        log_success "Service démarré"

        # Attendre et vérifier
        sleep 5
        if docker-compose ps | grep -q "Up"; then
            log_success "RustDesk Interface est opérationnel"
        else
            log_warning "Le service pourrait avoir des problèmes. Vérifiez les logs:"
            log_info "docker-compose logs"
        fi
    fi
}

# Créer le service systemd
create_systemd_service() {
    if [ "$EUID" -ne 0 ]; then
        return
    fi

    log_info "Création du service systemd..."

    cat > /etc/systemd/system/rustdesk-interface.service << EOF
[Unit]
Description=RustDesk Interface API
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=$INSTALL_DIR
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    log_success "Service systemd créé: rustdesk-interface"
    log_info "Pour activer au démarrage: systemctl enable rustdesk-interface"
}

# Afficher les instructions finales
print_instructions() {
    echo ""
    echo "============================================================================="
    echo -e "${GREEN}Installation terminée !${NC}"
    echo "============================================================================="
    echo ""
    echo "Répertoire d'installation: $INSTALL_DIR"
    echo ""
    echo "IMPORTANT - Avant le premier démarrage:"
    echo "  1. Éditez la configuration: nano $INSTALL_DIR/.env.production"
    echo "  2. Configurez votre serveur RustDesk (ID Server, Relay Server)"
    echo "  3. Vérifiez la clé JWT (générée automatiquement si nécessaire)"
    echo ""
    echo "Pour démarrer le service:"
    echo "  cd $INSTALL_DIR"
    echo "  docker-compose up -d"
    echo ""
    echo "Pour voir les logs:"
    echo "  docker-compose logs -f"
    echo ""
    echo "Accès à l'interface:"
    echo "  - API:           http://<IP>:21114/api"
    echo "  - Administration: http://<IP>:21114/admin"
    echo ""
    echo "Identifiants par défaut:"
    echo "  - Utilisateur: admin"
    echo "  - Mot de passe: Affiché dans les logs au premier démarrage"
    echo ""
    if [ "$EUID" -eq 0 ]; then
        echo "Service systemd disponible:"
        echo "  systemctl enable rustdesk-interface  # Activer au démarrage"
        echo "  systemctl start rustdesk-interface   # Démarrer"
        echo ""
    fi
    echo "============================================================================="
}

# Main
main() {
    echo "============================================================================="
    echo "  Installation hors ligne de RustDesk Interface"
    echo "============================================================================="
    echo ""

    check_root
    check_docker
    load_docker_image
    create_directories
    copy_config_files
    create_docker_compose
    generate_jwt_key
    set_permissions
    create_systemd_service
    start_service
    print_instructions
}

main "$@"
