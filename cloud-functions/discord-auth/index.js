/**
 * Google Cloud Function pour Discord OAuth2
 * 
 * Déploiement :
 * gcloud functions deploy discord-auth \
 *   --runtime nodejs20 \
 *   --trigger-http \
 *   --allow-unauthenticated \
 *   --entry-point discordAuth \
 *   --set-env-vars DISCORD_CLIENT_ID=xxx,DISCORD_CLIENT_SECRET=xxx,DISCORD_REDIRECT_URI=xxx
 */

const functions = require('@google-cloud/functions-framework');

/**
 * Fonction Cloud Function pour gérer l'authentification Discord
 */
functions.http('discordAuth', async (req, res) => {
  // Configuration CORS améliorée
  const allowedOrigins = [
    'http://localhost:5173',
    'http://localhost:3000',
    'http://localhost:4173',
    process.env.FRONTEND_URL
  ].filter(Boolean);

  const origin = req.headers.origin;
  if (allowedOrigins.includes(origin) || process.env.FRONTEND_URL === '*') {
    res.set('Access-Control-Allow-Origin', origin || process.env.FRONTEND_URL || '*');
  } else {
    res.set('Access-Control-Allow-Origin', '*');
  }
  
  res.set('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
  res.set('Access-Control-Allow-Headers', 'Content-Type, Authorization');
  res.set('Access-Control-Max-Age', '3600');

  // Gestion preflight CORS
  if (req.method === 'OPTIONS') {
    return res.status(204).send('');
  }

  // Vérifier la méthode HTTP
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Méthode non autorisée' });
  }

  try {
    const { code } = req.body;

    if (!code) {
      return res.status(400).json({ error: 'Code d\'autorisation manquant' });
    }

    // Configuration depuis les variables d'environnement
    const DISCORD_CLIENT_ID = process.env.DISCORD_CLIENT_ID;
    const DISCORD_CLIENT_SECRET = process.env.DISCORD_CLIENT_SECRET;
    const REDIRECT_URI = process.env.DISCORD_REDIRECT_URI;

    if (!DISCORD_CLIENT_ID || !DISCORD_CLIENT_SECRET || !REDIRECT_URI) {
      console.error('Variables d\'environnement manquantes');
      return res.status(500).json({ 
        error: 'Configuration serveur incorrecte' 
      });
    }

    console.log('🔄 Échange du code OAuth...');

    // Échanger le code contre un token Discord
    const tokenResponse = await fetch('https://discord.com/api/oauth2/token', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: new URLSearchParams({
        client_id: DISCORD_CLIENT_ID,
        client_secret: DISCORD_CLIENT_SECRET,
        grant_type: 'authorization_code',
        code: code,
        redirect_uri: REDIRECT_URI,
      }),
    });

    const tokenData = await tokenResponse.json();

    if (!tokenResponse.ok) {
      console.error('❌ Erreur Discord token:', tokenData);
      return res.status(400).json({ 
        error: 'Échec de l\'obtention du token Discord',
        details: tokenData.error_description || tokenData.error
      });
    }

    console.log('✅ Token obtenu, récupération des données utilisateur...');

    // Récupérer les informations de l'utilisateur
    const userResponse = await fetch('https://discord.com/api/users/@me', {
      headers: {
        Authorization: `Bearer ${tokenData.access_token}`,
      },
    });

    const userData = await userResponse.json();

    if (!userResponse.ok) {
      console.error('❌ Erreur Discord user:', userData);
      return res.status(400).json({ 
        error: 'Échec de la récupération des données utilisateur',
        details: userData.message
      });
    }

    console.log('✅ Authentification réussie pour:', userData.username);

    // Retourner les données au frontend
    return res.status(200).json({
      access_token: tokenData.access_token,
      token_type: tokenData.token_type,
      expires_in: tokenData.expires_in,
      refresh_token: tokenData.refresh_token,
      user: {
        id: userData.id,
        username: userData.username,
        discriminator: userData.discriminator,
        avatar: userData.avatar,
        email: userData.email,
        verified: userData.verified,
        global_name: userData.global_name,
      }
    });

  } catch (error) {
    console.error('❌ Erreur serveur:', error);
    return res.status(500).json({ 
      error: 'Erreur interne du serveur',
      message: error.message 
    });
  }
});

