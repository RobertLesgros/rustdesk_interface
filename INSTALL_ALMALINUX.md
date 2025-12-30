# Installation sur AlmaLinux (sans Docker)

Guide d'installation native de RustDesk Interface sur AlmaLinux 8/9.

## Prérequis

```bash
# Mettre à jour le système
sudo dnf update -y

# Installer les outils de base
sudo dnf install -y git gcc make sqlite-devel
```

## Étape 1 : Installer Go 1.22+

```bash
# Télécharger Go (version 1.22.5)
cd /tmp
curl -LO https://go.dev/dl/go1.22.5.linux-amd64.tar.gz

# Installer Go
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

# Configurer le PATH
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc

# Vérifier
go version
# → go version go1.22.5 linux/amd64
```

## Étape 2 : Installer Node.js 18+ (pour le frontend)

```bash
# Installer Node.js via dnf module (AlmaLinux 8/9)
sudo dnf module enable nodejs:18 -y
sudo dnf install -y nodejs npm

# Vérifier
node --version  # → v18.x.x
npm --version   # → 9.x.x
```

## Étape 3 : Cloner et compiler le projet

```bash
# Créer le répertoire d'installation
sudo mkdir -p /opt/rustdesk-interface
sudo chown $USER:$USER /opt/rustdesk-interface
cd /opt/rustdesk-interface

# Cloner le projet
git clone https://github.com/RobertLesgros/rustdesk_interface.git .

# Installer swag (générateur de documentation API)
go install github.com/swaggo/swag/cmd/swag@latest

# Télécharger les dépendances Go
go mod tidy
go mod download

# Générer la documentation Swagger
swag init -g cmd/apimain.go --output docs/api --instanceName api --exclude http/controller/admin
swag init -g cmd/apimain.go --output docs/admin --instanceName admin --exclude http/controller/api

# Compiler le backend
CGO_ENABLED=1 go build -o apimain cmd/apimain.go

# Vérifier
./apimain --help
```

## Étape 4 : Compiler le frontend

```bash
cd /opt/rustdesk-interface/frontend

# Installer les dépendances npm
npm install

# Compiler pour production
npm run build

# Copier vers le dossier resources
cp -r dist/* ../resources/admin/
```

## Étape 5 : Configurer

```bash
cd /opt/rustdesk-interface

# Copier la configuration LAN-only
cp conf/config-lan-only.yaml.example conf/config.yaml

# Éditer la configuration
nano conf/config.yaml
```

**Modifications importantes dans `config.yaml` :**

```yaml
# Adaptez ces valeurs à votre réseau LAN
rustdesk:
  id-server: "VOTRE_IP:21116"
  relay-server: "VOTRE_IP:21117"
  api-server: "http://VOTRE_IP:21114"

# Configuration LDAP (si utilisé)
ldap:
  enable: true
  url: "ldap://VOTRE_SERVEUR_AD:389"
  # ... autres paramètres LDAP
```

## Étape 6 : Créer les répertoires de données

```bash
cd /opt/rustdesk-interface
mkdir -p data runtime
chmod 755 data runtime
```

## Étape 7 : Tester le lancement

```bash
cd /opt/rustdesk-interface
./apimain
```

Vous devriez voir :
```
[GIN] Listening and serving HTTP on 0.0.0.0:21114
```

Accédez à : http://VOTRE_IP:21114/_admin/

**Mot de passe admin initial** : affiché dans les logs au premier démarrage.

## Étape 8 : Créer un service systemd

```bash
sudo tee /etc/systemd/system/rustdesk-interface.service << 'EOF'
[Unit]
Description=RustDesk Interface API Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/rustdesk-interface
ExecStart=/opt/rustdesk-interface/apimain
Restart=always
RestartSec=5

# Sécurité
NoNewPrivileges=true
ProtectSystem=strict
ReadWritePaths=/opt/rustdesk-interface/data /opt/rustdesk-interface/runtime

# Environnement
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

# Recharger systemd
sudo systemctl daemon-reload

# Activer au démarrage
sudo systemctl enable rustdesk-interface

# Démarrer le service
sudo systemctl start rustdesk-interface

# Vérifier le statut
sudo systemctl status rustdesk-interface
```

## Étape 9 : Ouvrir le firewall

```bash
# Ouvrir le port API
sudo firewall-cmd --permanent --add-port=21114/tcp

# Si RustDesk Server est sur la même machine
sudo firewall-cmd --permanent --add-port=21115-21119/tcp
sudo firewall-cmd --permanent --add-port=21116/udp

# Appliquer
sudo firewall-cmd --reload
```

## Commandes utiles

```bash
# Voir les logs
sudo journalctl -u rustdesk-interface -f

# Redémarrer
sudo systemctl restart rustdesk-interface

# Arrêter
sudo systemctl stop rustdesk-interface

# Réinitialiser le mot de passe admin
cd /opt/rustdesk-interface
./apimain reset-admin-pwd NouveauMotDePasse
```

## Mise à jour

```bash
cd /opt/rustdesk-interface

# Arrêter le service
sudo systemctl stop rustdesk-interface

# Sauvegarder les données
cp -r data data.backup
cp conf/config.yaml conf/config.yaml.backup

# Mettre à jour le code
git pull

# Recompiler
go mod tidy
CGO_ENABLED=1 go build -o apimain cmd/apimain.go

# Recompiler le frontend
cd frontend
npm install
npm run build
cp -r dist/* ../resources/admin/

# Redémarrer
sudo systemctl start rustdesk-interface
```

## Dépannage

### Erreur "gcc not found"
```bash
sudo dnf install -y gcc
```

### Erreur "sqlite3.h not found"
```bash
sudo dnf install -y sqlite-devel
```

### Le port 21114 est déjà utilisé
```bash
# Vérifier quel processus utilise le port
sudo ss -tlnp | grep 21114

# Changer le port dans config.yaml
gin:
  api-addr: "0.0.0.0:8080"  # Autre port
```

### Problème de permissions
```bash
sudo chown -R $USER:$USER /opt/rustdesk-interface
chmod +x /opt/rustdesk-interface/apimain
```
