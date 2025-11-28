import { PubSub } from '@google-cloud/pubsub';

const pubsub = new PubSub();
const TOPIC_PIXEL_EVENTS = 'drawing-pixel';
const TOPIC_SESSION_EVENTS = 'session-events';
const DISCORD_API_URL = 'https://discord.com/api/users/@me';

export const webProxyRouter = async (req, res) => {
    try {
        const allowedOrigins = [
            'http://localhost:5173',
            'https://serverless-epitech-dev-476110.ew.r.appspot.com'
        ];

        const origin = req.headers.origin;
        if (allowedOrigins.includes(origin)) {
            res.setHeader('Access-Control-Allow-Origin', origin);
        }

        res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
        res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-API-KEY, X-Discord-Token');
        res.setHeader('Access-Control-Allow-Credentials', 'true');

        if (req.method === 'OPTIONS') {
            return res.status(204).send('');
        }

        const path = req.path;
        const method = req.method;
        console.log(`Routing request for path: ${path}, method: ${method}`);

        // --- ROUTING POUR LES ENDPOINTS DISCORD (PROXY) ---
        // Ces endpoints ne nécessitent pas de validation préalable, juste le token Discord
        if ((path === '/web/api/discord/oauth2/@me' || path === '/web/api/discord/users/@me') && method === 'GET') {
            const authHeader = req.headers['x-discord-token'] || req.headers['authorization'];

            if (!authHeader) {
                return res.status(401).json({ error: 'Unauthorized', message: 'Discord token missing' });
            }

            const accessToken = authHeader.startsWith('Bearer ')
                ? authHeader.split(' ')[1]
                : authHeader;

            return await handleDiscordUserInfo(req, res, accessToken);
        }

        // --- ROUTING POUR LES AUTRES ENDPOINTS (NÉCESSITENT VALIDATION DISCORD) ---
        // Le header X-API-KEY a été vérifié par la Gateway (Couche 1)
        const authHeader = req.headers['x-discord-token'];

        if (!authHeader) {
            return res.status(401).send('Unauthorized: X-Discord-Token missing.');
        }

        const accessToken = authHeader?.startsWith('Bearer ') ? authHeader.split(' ')[1] : null;

        if (!accessToken) {
            return res.status(401).send('Unauthorized: Bearer token missing.');
        }

        // Validation du jeton Discord (Couche 2)
        const discordResponse = await fetch(DISCORD_API_URL, {
            headers: { authorization: `Bearer ${accessToken}` },
        });

        if (!discordResponse.ok) {
            const errorBody = await discordResponse.text();
            console.error(
                "Discord API call failed. Status:",
                discordResponse.status,
                "Body:",
                errorBody
            );
            return res.status(401).send('Unauthorized: Invalid Discord Token.');
        }

        const userData = await discordResponse.json();
        const userId = userData.id;

        // Routing interne
        if (path === '/web/api/draw-pixel' && method === 'POST') {
            await handleDrawPixel(req, res, userId);
        }
        else if (path === '/web/api/retrieve-canvas' && method === 'GET') {
            handleRetrieveCanvas(req, res);
        }
        else if (path.startsWith('/web/api/session/') && method === 'POST') {
            await handleSessionAction(req, res, userId, path.split('/').pop()); // Extrait 'start', 'pause', 'reset'
        }
        else if (path === '/web/api/session/snapshot' && method === 'GET') {
            handleRetrieveSnapshot(req, res);
        }
        else {
            res.status(404).send('Not Found');
        }

    } catch (error) {
        console.error("Critical error in web proxy router:", error);
        res.status(500).send('Internal Server Error.');
    }
};

// --------------------------------------------------------
// FONCTIONS DE GESTION (PUBLICATION À PUB/SUB)
// --------------------------------------------------------

/** Gère l'événement /draw-pixel. */
const handleDrawPixel = async (req, res, userId) => {
    const { x, y, color, canvasId } = req.body;
    if (x === undefined || y === undefined || !color || !canvasId) {
        return res.status(400).send('Missing pixel data.');
    }

    // Publication de l'événement dans Pub/Sub
    await pubsub.topic(TOPIC_PIXEL_EVENTS).publishMessage({
        json: { x, y, color, canvasId, authorId: userId, updatedAt: Date.now() }
    });

    res.status(200).json({ ok: true, message: 'Pixel event queued.' });
};

/** Gère les événements /session/action. */
const handleSessionAction = async (req, res, userId, action) => {
    if (!['start', 'pause', 'reset'].includes(action)) {
        return res.status(400).send('Invalid session action.');
    }

    let eventData = {};

    if (action === 'start') {
        const { name, width, height, endDate } = req.body;

        if (!name || !width || !height || !endDate) {
            return res.status(400).send('Missing session start parameters.');
        }
        eventData = { adminId: userId, name, width, height, status: action.toUpperCase(), startDate: Date.now(), endDate };
    } else {
        const { canvasId } = req.body;
        if (!canvasId) {
            return res.status(400).send('Missing canvasId for session action.');
        }
        eventData = { adminId: userId, canvasId, status: action.toUpperCase() };
    }

    // Publication de l'événement dans le topic des sessions
    await pubsub.topic(TOPIC_SESSION_EVENTS).publishMessage({
        json: eventData
    });

    res.status(200).json({ ok: true, message: `Session action ${action} queued.` });
};

/** Gère la récupération de données (READ) qui ne va pas à Pub/Sub. */
const handleRetrieveCanvas = (req, res) => {
    // Les requêtes GET ne sont généralement pas asynchrones (pas de Pub/Sub).
    // Ici, le proxy doit demander l'état actuel de la toile au système de stockage (e.g., Firestore ou Cloud Storage).

    // NOTE : Pour ce projet, le proxy ne doit pas contenir de logique de stockage/lecture complexe.
    // Pour l'instant, on simule une réponse pour valider le routing :
    res.status(200).json({
        ok: true,
        status: 'live',
        message: 'Canvas retrieval route is active.'
    });
};

const handleRetrieveSnapshot = (req, res) => {
    // Peut rediriger l'utilisateur vers un lien Cloud Storage
    res.status(200).json({
        ok: true,
        snapshot_url: 'https://storage.googleapis.com/votre-bucket/latest-snapshot.png'
    });
};

/** Gère les appels proxy vers l'API Discord pour récupérer les infos utilisateur. */
const handleDiscordUserInfo = async (req, res, accessToken) => {
    try {
        const discordResponse = await fetch(DISCORD_API_URL, {
            headers: {
                authorization: `Bearer ${accessToken}`,
                'User-Agent': 'PixelPlace/1.0'
            },
        });

        if (!discordResponse.ok) {
            const errorBody = await discordResponse.text();
            console.error(
                "Discord API call failed. Status:",
                discordResponse.status,
                "Body:",
                errorBody
            );
            return res.status(discordResponse.status).json({
                error: 'Discord API Error',
                message: errorBody || 'Failed to fetch user info'
            });
        }

        const userData = await discordResponse.json();

        // Retourne les données dans le format attendu par le frontend
        // Le format peut varier selon l'endpoint appelé
        if (req.path === '/web/api/discord/oauth2/@me') {
            // Format pour AuthCallback.vue qui attend { user: {...} }
            res.status(200).json({ user: userData });
        } else {
            // Format direct pour useDiscordAuth.ts
            res.status(200).json(userData);
        }
    } catch (error) {
        console.error("Error in handleDiscordUserInfo:", error);
        res.status(500).json({ error: 'Internal Server Error', message: error.message });
    }
};