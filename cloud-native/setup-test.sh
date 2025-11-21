#!/bin/bash

# Script de setup pour tester l'authentification Discord en local
# Usage: ./setup-test.sh

set -e

echo "🎯 Setup de l'environnement de test Discord OAuth2"
echo ""

# Couleurs
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Vérifier si .env.local existe (prioritaire)
if [ -f ".env.local" ]; then
    echo -e "${BLUE}ℹ️  Fichier .env.local détecté${NC}"
    
    # Vérifier si les variables nécessaires sont présentes
    if grep -q "VITE_DISCORD_CLIENT_ID" .env.local && grep -q "VITE_REDIRECT_URI" .env.local; then
        echo -e "${GREEN}✅ Configuration trouvée dans .env.local${NC}"
        
        # Afficher les valeurs (masquer le Client ID sauf les 4 derniers caractères)
        CLIENT_ID=$(grep "VITE_DISCORD_CLIENT_ID" .env.local | cut -d '=' -f2)
        REDIRECT_URI=$(grep "VITE_REDIRECT_URI" .env.local | cut -d '=' -f2)
        
        CLIENT_ID_MASKED="${CLIENT_ID:0:4}...${CLIENT_ID: -4}"
        
        echo ""
        echo "  📋 VITE_DISCORD_CLIENT_ID: $CLIENT_ID_MASKED"
        echo "  📋 VITE_REDIRECT_URI: $REDIRECT_URI"
        echo ""
        
        # Vérifier si l'URI de redirection est correcte
        if [[ "$REDIRECT_URI" != "http://localhost:5173/auth/callback" ]]; then
            echo -e "${YELLOW}⚠️  L'URI de redirection ne correspond pas à l'environnement de dev${NC}"
            echo -e "${YELLOW}   Attendu: http://localhost:5173/auth/callback${NC}"
            echo -e "${YELLOW}   Trouvé:  $REDIRECT_URI${NC}"
            echo ""
            read -p "Voulez-vous continuer quand même ? (y/N) " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                echo "💡 Mettez à jour votre .env.local avec la bonne URI"
                exit 1
            fi
        fi
        
        echo -e "${GREEN}✅ Configuration validée !${NC}"
        ENV_FILE=".env.local"
    else
        echo -e "${YELLOW}⚠️  .env.local existe mais les variables nécessaires sont manquantes${NC}"
        echo ""
        read -p "Voulez-vous créer un nouveau .env.local ? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "❌ Opération annulée"
            exit 1
        fi
        
        # Créer .env.local depuis .env.example
        if [ -f ".env.example" ]; then
            cp .env.example .env.local
            echo -e "${GREEN}✅ Fichier .env.local créé${NC}"
            ENV_FILE=".env.local"
            NEED_CONFIG=true
        else
            echo -e "${RED}❌ .env.example introuvable${NC}"
            exit 1
        fi
    fi
else
    # Aucun fichier .env.local, en créer un
    echo -e "${YELLOW}ℹ️  Aucun fichier .env.local trouvé${NC}"
    
    if [ -f ".env.example" ]; then
        cp .env.example .env.local
        echo -e "${GREEN}✅ Fichier .env.local créé depuis .env.example${NC}"
        ENV_FILE=".env.local"
        NEED_CONFIG=true
    else
        echo -e "${RED}❌ .env.example introuvable${NC}"
        exit 1
    fi
fi

# Si besoin de configuration
if [ "$NEED_CONFIG" = true ]; then
    echo ""
    echo "📝 Configuration Discord"
    echo ""
    echo "Pour obtenir votre Client ID :"
    echo "  1. Allez sur https://discord.com/developers/applications"
    echo "  2. Créez une application (ou sélectionnez-en une)"
    echo "  3. Dans OAuth2 → General, copiez le CLIENT ID"
    echo "  4. Ajoutez http://localhost:5173/auth/callback dans les Redirects"
    echo ""

    read -p "Entrez votre Discord Client ID: " CLIENT_ID

    if [ -z "$CLIENT_ID" ]; then
        echo -e "${RED}❌ Client ID vide${NC}"
        exit 1
    fi

    # Remplacer dans .env.local
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s/YOUR_DISCORD_CLIENT_ID/$CLIENT_ID/g" .env.local
    else
        # Linux
        sed -i "s/YOUR_DISCORD_CLIENT_ID/$CLIENT_ID/g" .env.local
    fi

    echo -e "${GREEN}✅ Client ID configuré dans .env.local${NC}"
    echo ""
fi

# Vérifier si node_modules existe
if [ ! -d "node_modules" ]; then
    echo "📦 Installation des dépendances..."
    npm install
    echo -e "${GREEN}✅ Dépendances installées${NC}"
    echo ""
else
    echo -e "${GREEN}✅ Dépendances déjà installées${NC}"
    echo ""
fi

echo -e "${GREEN}🎉 Configuration terminée !${NC}"
echo ""
echo "Pour lancer le serveur de développement :"
echo -e "${YELLOW}  npm run dev${NC}"
echo ""
echo "Puis ouvrez votre navigateur sur :"
echo -e "${YELLOW}  http://localhost:5173/login${NC}"
echo ""
echo "Pour debugger l'authentification :"
echo -e "${YELLOW}  http://localhost:5173/debug${NC}"
echo ""
echo "📖 Pour plus d'informations, consultez TEST_AUTH.md"

