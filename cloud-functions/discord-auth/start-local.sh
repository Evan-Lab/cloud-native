
#!/bin/bash

# Script pour démarrer la Cloud Function localement
# Usage: ./start-local.sh

set -e

echo "🔧 Démarrage de la Cloud Function en local..."

# Couleurs
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m'

# Vérifier si .env.local existe
if [ -f ".env.local" ]; then
    echo -e "${GREEN}✅ Chargement des variables depuis .env.local${NC}"
    set -a  # Exporter automatiquement les variables
    . .env.local
    set +a
else
    echo -e "${YELLOW}⚠️  Fichier .env.local non trouvé${NC}"
    echo "Créez un fichier .env.local avec vos variables d'environnement:"
    echo ""
    cat .env.example
    echo ""
    exit 1
fi

# Vérifier les variables
if [ -z "$DISCORD_CLIENT_ID" ] || [ -z "$DISCORD_CLIENT_SECRET" ]; then
    echo -e "${YELLOW}⚠️  Variables d'environnement manquantes${NC}"
    echo "Assurez-vous que .env.local contient:"
    echo "  DISCORD_CLIENT_ID"
    echo "  DISCORD_CLIENT_SECRET"
    echo "  DISCORD_REDIRECT_URI"
    echo "  FRONTEND_URL"
    exit 1
fi

# Installer les dépendances si nécessaire
if [ ! -d "node_modules" ]; then
    echo "📦 Installation des dépendances..."
    npm install
fi

echo ""
echo -e "${GREEN}🚀 Démarrage du serveur local...${NC}"
echo ""
echo "📍 URL: http://localhost:8080"
echo "📋 Utilisez cette URL dans VITE_BACKEND_URL de votre frontend"
echo ""
echo "Configuration chargée:"
echo "  DISCORD_CLIENT_ID: ${DISCORD_CLIENT_ID:0:10}..."
echo "  REDIRECT_URI: $DISCORD_REDIRECT_URI"
echo ""
echo "Pour tester:"
echo "  curl -X POST http://localhost:8080 -H 'Content-Type: application/json' -d '{\"code\":\"test\"}'"
echo ""

npm start