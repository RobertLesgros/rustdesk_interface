#!/bin/bash
# =============================================================================
# Script de préparation pour installation hors ligne de RustDesk Interface
# =============================================================================
# Ce script doit être exécuté sur une machine AVEC connexion Internet.
# Il prépare tous les fichiers nécessaires pour une installation hors ligne.
#
# Prérequis:
#   git clone https://github.com/RobertLesgros/rustdesk_interface.git
#   cd rustdesk_interface
#   ./scripts/prepare-offline.sh
#
# Usage:
#   ./scripts/prepare-offline.sh [--export-only] [--skip-build]
#
# Options:
#   --export-only    : Exporte uniquement l'image existante sans reconstruire
#   --skip-build     : Ne pas construire l'image Docker
#
# Sorties:
#   - offline-bundle/        : Bundle complet pour transfert hors ligne
#   - rustdesk-interface-offline.tar : Image Docker exportée
# =============================================================================

set -e

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
IMAGE_NAME="rustdesk-interface"
IMAGE_TAG="offline"
BUNDLE_DIR="$PROJECT_DIR/offline-bundle"

# Options
EXPORT_ONLY=false
SKIP_BUILD=false

# Parser les arguments
for arg in "$@"; do
    case $arg in
        --export-only)
            EXPORT_ONLY=true
            ;;
        --skip-build)
            SKIP_BUILD=true
            ;;
        --help|-h)
            echo "Usage: $0 [--export-only] [--skip-build]"
            echo ""
            echo "Options:"
            echo "  --export-only    Exporte uniquement l'image existante sans reconstruire"
            echo "  --skip-build     Ne pas construire l'image Docker"
            echo "  --help, -h       Affiche cette aide"
            exit 0
            ;;
        *)
            echo -e "${RED}Option inconnue: $arg${NC}"
            exit 1
            ;;
    esac
done

# Fonctions utilitaires
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[ATTENTION]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERREUR]${NC} $1"
}

# Vérifier les prérequis
check_prerequisites() {
    log_info "Vérification des prérequis..."

    if ! command -v docker &> /dev/null; then
        log_error "Docker n'est pas installé. Veuillez l'installer avant de continuer."
        exit 1
    fi

    if ! docker info &> /dev/null; then
        log_error "Docker n'est pas en cours d'exécution. Veuillez le démarrer."
        exit 1
    fi

    log_success "Tous les prérequis sont satisfaits"
}

# Construire l'image Docker
# Utilise Dockerfile.dev qui gère automatiquement le téléchargement du frontend
build_docker_image() {
    if [ "$SKIP_BUILD" = true ] || [ "$EXPORT_ONLY" = true ]; then
        log_warning "Construction de l'image ignorée"
        return
    fi

    log_info "Construction de l'image Docker (cela peut prendre plusieurs minutes)..."
    log_info "Le Dockerfile.dev va automatiquement télécharger et compiler le frontend."

    cd "$PROJECT_DIR"

    # Construire l'image avec Dockerfile.dev
    # Ce Dockerfile gère tout : backend Go + frontend Node.js
    docker build \
        -f Dockerfile.dev \
        -t ${IMAGE_NAME}:${IMAGE_TAG} \
        --build-arg BUILDARCH=$(dpkg --print-architecture 2>/dev/null || echo "amd64") \
        .

    log_success "Image Docker construite: ${IMAGE_NAME}:${IMAGE_TAG}"
}

# Exporter l'image Docker
export_docker_image() {
    log_info "Export de l'image Docker..."

    cd "$PROJECT_DIR"

    # Vérifier que l'image existe
    if ! docker image inspect ${IMAGE_NAME}:${IMAGE_TAG} &> /dev/null; then
        log_error "L'image ${IMAGE_NAME}:${IMAGE_TAG} n'existe pas. Construisez-la d'abord."
        exit 1
    fi

    # Exporter l'image
    docker save ${IMAGE_NAME}:${IMAGE_TAG} -o rustdesk-interface-offline.tar

    log_success "Image exportée: rustdesk-interface-offline.tar"

    # Afficher la taille
    local size=$(du -h rustdesk-interface-offline.tar | cut -f1)
    log_info "Taille de l'image: $size"
}

# Créer le bundle complet
create_bundle() {
    log_info "Création du bundle hors ligne..."

    cd "$PROJECT_DIR"

    # Créer le dossier bundle
    rm -rf "$BUNDLE_DIR"
    mkdir -p "$BUNDLE_DIR"

    # Copier les fichiers nécessaires
    cp -r conf "$BUNDLE_DIR/"
    cp docker-compose-offline.yml "$BUNDLE_DIR/docker-compose.yml" 2>/dev/null || true
    cp rustdesk-interface-offline.tar "$BUNDLE_DIR/" 2>/dev/null || true
    cp INSTALL_OFFLINE.md "$BUNDLE_DIR/" 2>/dev/null || true
    cp scripts/install-offline.sh "$BUNDLE_DIR/" 2>/dev/null || true

    # Créer le fichier .env.production exemple
    cat > "$BUNDLE_DIR/.env.production.example" << 'EOF'
# =============================================================================
# Configuration RustDesk Interface - Production
# =============================================================================
# Copiez ce fichier vers .env.production et modifiez les valeurs

# Timezone
TZ=Europe/Paris

# Langue
RUSTDESK_API_LANG=fr

# Mode de l'API (release ou debug)
RUSTDESK_API_GIN_MODE=release

# Clé JWT (OBLIGATOIRE - générez avec: openssl rand -hex 32)
RUSTDESK_API_JWT_KEY=CHANGEZ_MOI_AVEC_UNE_CLE_SECRETE

# Configuration RustDesk Server (adaptez à votre réseau)
RUSTDESK_API_RUSTDESK_ID_SERVER=192.168.1.66:21116
RUSTDESK_API_RUSTDESK_RELAY_SERVER=192.168.1.66:21117

# Base de données (sqlite par défaut)
# RUSTDESK_API_GORM_TYPE=sqlite

# Pour MySQL:
# RUSTDESK_API_GORM_TYPE=mysql
# RUSTDESK_API_MYSQL_HOST=localhost
# RUSTDESK_API_MYSQL_PORT=3306
# RUSTDESK_API_MYSQL_USERNAME=rustdesk
# RUSTDESK_API_MYSQL_PASSWORD=CHANGEZ_MOI
# RUSTDESK_API_MYSQL_DBNAME=rustdesk
EOF

    log_success "Bundle créé dans: $BUNDLE_DIR/"

    # Lister le contenu
    log_info "Contenu du bundle:"
    ls -la "$BUNDLE_DIR/"

    # Calculer la taille totale
    local total_size=$(du -sh "$BUNDLE_DIR" | cut -f1)
    log_info "Taille totale du bundle: $total_size"
}

# Créer une archive du bundle
create_archive() {
    log_info "Création de l'archive..."

    cd "$PROJECT_DIR"

    # Créer l'archive tar.gz
    tar -czvf rustdesk-interface-offline-bundle.tar.gz offline-bundle/

    local archive_size=$(du -h rustdesk-interface-offline-bundle.tar.gz | cut -f1)
    log_success "Archive créée: rustdesk-interface-offline-bundle.tar.gz ($archive_size)"
}

# Afficher les instructions finales
print_instructions() {
    echo ""
    echo "============================================================================="
    echo -e "${GREEN}Préparation terminée !${NC}"
    echo "============================================================================="
    echo ""
    echo "Fichiers créés:"
    echo "  - rustdesk-interface-offline.tar      : Image Docker"
    echo "  - offline-bundle/                     : Bundle complet"
    echo "  - rustdesk-interface-offline-bundle.tar.gz : Archive pour transfert"
    echo ""
    echo "Pour installer sur la machine hors ligne:"
    echo ""
    echo "  1. Transférez l'archive vers la machine cible:"
    echo "     scp rustdesk-interface-offline-bundle.tar.gz user@machine:/tmp/"
    echo ""
    echo "  2. Sur la machine cible, extrayez et installez:"
    echo "     cd /tmp"
    echo "     tar -xzvf rustdesk-interface-offline-bundle.tar.gz"
    echo "     cd offline-bundle"
    echo "     ./install-offline.sh"
    echo ""
    echo "  Ou manuellement:"
    echo "     docker load -i rustdesk-interface-offline.tar"
    echo "     cp .env.production.example .env.production"
    echo "     # Éditez .env.production avec vos valeurs"
    echo "     docker-compose up -d"
    echo ""
    echo "Consultez INSTALL_OFFLINE.md pour la documentation complète."
    echo "============================================================================="
}

# Main
main() {
    echo "============================================================================="
    echo "  Préparation RustDesk Interface pour installation hors ligne"
    echo "============================================================================="
    echo ""

    check_prerequisites

    if [ "$EXPORT_ONLY" = false ]; then
        build_docker_image
    fi

    export_docker_image
    create_bundle
    create_archive
    print_instructions
}

main "$@"
