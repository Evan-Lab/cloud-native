import fetch from 'node-fetch'

const API_URL = 'https://rplace-gateway-5uir24en.ew.gateway.dev/web/api/draw-pixel'
const DISCORD_TOKEN = 'MTQ0MTM1OTEzNTQyNDY0NzMyMA.CD18tQQrFcR2m1mr9x7VdMxY2sHRrM'
const API_KEY = import.meta.env.VITE_API_KEY

const pixel = {
  x: 60,
  y: 99,
  color: '#FF1111',
}

async function drawPixel() {
  try {
    const body = {
      x: pixel.x,
      y: pixel.y,
      color: pixel.color,
      canvasId: 'zGYJpT1GTkWY95li4q0q',
    }

    const res = await fetch(API_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-KEY': API_KEY,
        'X-Discord-Token': `Bearer ${DISCORD_TOKEN}`,
      },
      body: JSON.stringify(body),
    })

    if (!res.ok) {
      const errText = await res.text()
      throw new Error(`Erreur serveur (${res.status}): ${errText}`)
    }

    console.log(`Pixel placé en (${pixel.x}, ${pixel.y}) avec couleur ${pixel.color}`)
  } catch (err) {
    console.error('Erreur placement pixel:', err.message)
  }
}

drawPixel()
