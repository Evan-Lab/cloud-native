#!/bin/bash

# Script de déploiement de la Cloud Function Discord Auth
# Usage: ./deploy.sh

set -e

echo "🚀 Déploiement de la Cloud Function discord-auth..."

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Vérifier que les variables sont définies
if [ -z "$DISCORD_CLIENT_ID" ] || [ -z "$DISCORD_CLIENT_SECRET" ] || [ -z "$DISCORD_REDIRECT_URI" ] || [ -z "$FRONTEND_URL" ]; then
    echo -e "${RED}❌ Erreur: Variables d'environnement manquantes${NC}"
    echo ""
    echo "Veuillez définir les variables suivantes:"
    echo "  export DISCORD_CLIENT_ID='votre_client_id'"
    echo "  export DISCORD_CLIENT_SECRET='votre_client_secret'"
    echo "  export DISCORD_REDIRECT_URI='https://votre-frontend.com/auth/callback'"
    echo "  export FRONTEND_URL='https://votre-frontend.com'"
    echo ""
    echo "Ou créez un fichier .env et sourcez-le: source .env"
    exit 1
fi

# Configuration
FUNCTION_NAME="discord-auth"
REGION="europe-west1"  # Changez selon votre région
RUNTIME="nodejs20"

echo -e "${YELLOW}📋 Configuration:${NC}"
echo "  - Fonction: ${FUNCTION_NAME}"
echo "  - Région: ${REGION}"
echo "  - Runtime: ${RUNTIME}"
echo "  - Frontend URL: ${FRONTEND_URL}"
echo ""

# Déploiement
echo -e "${YELLOW}🔨 Déploiement en cours...${NC}"
gcloud functions deploy ${FUNCTION_NAME} \
  --gen2 \
  --runtime=${RUNTIME} \
  --region=${REGION} \
  --trigger-http \
  --allow-unauthenticated \
  --entry-point=discordAuth \
  --set-env-vars DISCORD_CLIENT_ID="${DISCORD_CLIENT_ID}",DISCORD_CLIENT_SECRET="${DISCORD_CLIENT_SECRET}",DISCORD_REDIRECT_URI="${DISCORD_REDIRECT_URI}",FRONTEND_URL="${FRONTEND_URL}"

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Déploiement réussi!${NC}"
    echo ""
    echo -e "${YELLOW}📝 Récupération de l'URL de la fonction...${NC}"
    
    FUNCTION_URL=$(gcloud functions describe ${FUNCTION_NAME} \
      --gen2 \
      --region=${REGION} \
      --format='value(serviceConfig.uri)')
    
    echo ""
    echo -e "${GREEN}🎉 URL de la Cloud Function:${NC}"
    echo "  ${FUNCTION_URL}"
    echo ""
    echo -e "${YELLOW}📋 Prochaines étapes:${NC}"
    echo "  1. Ajoutez cette URL dans votre fichier .env frontend:"
    echo "     VITE_BACKEND_URL=${FUNCTION_URL}"
    echo ""
    echo "  2. Configurez cette URL comme 'Redirect URL' dans Discord Developer Portal:"
    echo "     ${DISCORD_REDIRECT_URI}"
    echo ""
else
    echo ""
    echo -e "${RED}❌ Échec du déploiement${NC}"
    exit 1
fi

