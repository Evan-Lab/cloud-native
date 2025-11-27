import { PubSub } from '@google-cloud/pubsub';

const pubsub = new PubSub();
const TOPIC_PIXEL_EVENTS = 'pixel-events';
const TOPIC_SESSION_EVENTS = 'session-events';
const DISCORD_API_URL = 'https://discord.com/api/users/@me';

export const webProxyRouter = async (req, res) => {
    try {
        // --- 1. SÉCURITÉ UTILISATEUR (AUTHENTIFICATION DISCORD) ---
        // Le header X-API-KEY a été vérifié par la Gateway (Couche 1)

        const authHeader = req.headers['x-discord-token'];

        if (!authHeader) {
            return res.status(401).send('Unauthorized: X-Discord-Token missing.');
        }

        const accessToken = authHeader?.startsWith('Bearer ') ? authHeader.split(' ')[1] : null;

        if (!accessToken) {
            return res.status(401).send('Unauthorized: Bearer token missing.');
        }
        const authHeaderToSend = `Bearer ${accessToken}`;

        console.log("DEBUG: Header sent to Discord:", authHeaderToSend);
        // 1.1. Validation du jeton Discord (Couche 2)
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

        // --- 2. ROUTING INTERNE (Basé sur le chemin) ---
        const path = req.path;
        console.log(`Routing request for path: ${path}, method: ${req.method}`);
        const method = req.method;

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
    const { x, y, color } = req.body;
    if (x === undefined || y === undefined || !color) {
        return res.status(400).send('Missing pixel data.');
    }

    // Publication de l'événement dans Pub/Sub
    await pubsub.topic(TOPIC_PIXEL_EVENTS).publishMessage({
        json: { userId, x, y, color, action: 'DRAW', timestamp: Date.now() }
    });

    res.status(200).json({ ok: true, message: 'Pixel event queued.' });
};

/** Gère les événements /session/action. */
const handleSessionAction = async (req, res, userId, action) => {
    if (!['start', 'pause', 'reset'].includes(action)) {
        return res.status(400).send('Invalid session action.');
    }

    // Publication de l'événement dans le topic des sessions
    await pubsub.topic(TOPIC_SESSION_EVENTS).publishMessage({
        json: { userId, action: action.toUpperCase(), timestamp: Date.now() }
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