# Discord Auth Cloud Function

Cloud Function GCP pour gérer l'authentification Discord OAuth2.

## 🚀 Déploiement Rapide

### 1. Configurer les variables d'environnement

```bash
export DISCORD_CLIENT_ID="votre_client_id"
export DISCORD_CLIENT_SECRET="votre_client_secret"
export DISCORD_REDIRECT_URI="http://localhost:5173/auth/callback"
export FRONTEND_URL="http://localhost:5173"
```

Ou créez un fichier `.env` :

```bash
cp .env.local.example .env.local
# Éditez .env.local.local avec vos valeurs
source .env.local
```

### 2. Déployer

```bash
./deploy.sh
```

### 3. Récupérer l'URL

L'URL de la fonction sera affichée après le déploiement :
```
https://REGION-PROJECT.cloudfunctions.net/discord-auth
```

Utilisez cette URL dans votre frontend (variable `VITE_BACKEND_URL`).

## 🧪 Test Local

```bash
npm install
npm start
```

La fonction sera disponible sur `http://localhost:8080`

### Test avec curl

```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"code":"test_code"}'
```

## 📋 Variables d'Environnement

| Variable | Description | Exemple |
|----------|-------------|---------|
| `DISCORD_CLIENT_ID` | Client ID de votre application Discord | `123456789012345678` |
| `DISCORD_CLIENT_SECRET` | Client Secret de votre application Discord | `abcdef123456...` |
| `DISCORD_REDIRECT_URI` | URL de redirection après auth Discord | `http://localhost:5173/auth/callback` |
| `FRONTEND_URL` | URL de votre frontend (pour CORS) | `http://localhost:5173` |

## 🔍 Vérification

### Voir les logs

```bash
gcloud functions logs read discord-auth \
  --gen2 \
  --region=europe-west1 \
  --limit=50
```

### Vérifier les variables d'environnement

```bash
gcloud functions describe discord-auth \
  --gen2 \
  --region=europe-west1
```

## 🔧 Structure

```
discord-auth/
├── index.js          # Code de la fonction
├── package.json      # Dépendances
├── deploy.sh         # Script de déploiement
├── .env.example      # Exemple de variables
└── README.md         # Ce fichier
```

## 📖 API

### POST /

Échange un code d'autorisation Discord contre un token d'accès.

**Body:**
```json
{
  "code": "authorization_code_from_discord"
}
```

**Response Success (200):**
```json
{
  "access_token": "...",
  "token_type": "Bearer",
  "expires_in": 604800,
  "refresh_token": "...",
  "user": {
    "id": "123456789",
    "username": "username",
    "discriminator": "0",
    "avatar": "...",
    "email": "user@example.com",
    "verified": true,
    "global_name": "Display Name"
  }
}
```

**Response Error (400/500):**
```json
{
  "error": "Message d'erreur",
  "details": "Détails supplémentaires"
}
```

## 🛡️ CORS

La fonction gère automatiquement CORS pour :
- `http://localhost:5173` (Vite dev server)
- `http://localhost:3000` 
- `http://localhost:4173` (Vite preview)
- Votre `FRONTEND_URL` configurée

## 📚 Ressources

- [Discord OAuth2 Documentation](https://discord.com/developers/docs/topics/oauth2)
- [Google Cloud Functions](https://cloud.google.com/functions/docs)
