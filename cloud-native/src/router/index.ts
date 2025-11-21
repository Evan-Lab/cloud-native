import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// Import des vues
import LoginPage from '@/views/LoginPage.vue'
import HomePage from '@/views/HomePage.vue'
import AuthCallback from '@/views/AuthCallback.vue'

// Définition des routes
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: HomePage,
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/login',
    name: 'login',
    component: LoginPage,
    meta: {
      requiresAuth: false
    }
  },
  {
    path: '/auth/callback',
    name: 'auth-callback',
    component: AuthCallback,
    meta: {
      requiresAuth: false
    }
  },
  // Route 404
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    redirect: { name: 'login' }
  }
]

// Création du router
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Guard de navigation pour vérifier l'authentification
router.beforeEach((to, from, next) => {
  // Récupérer le token depuis localStorage
  const token = localStorage.getItem('discord_token')
  const userInfo = localStorage.getItem('discord_user')
  const isAuthenticated = !!(token && userInfo)

  // Si la route nécessite une authentification et que l'utilisateur n'est pas connecté
  if (to.meta.requiresAuth && !isAuthenticated) {
    next({ name: 'login' })
  }
  // Si l'utilisateur est connecté et tente d'accéder à la page de login
  else if (to.name === 'login' && isAuthenticated) {
    next({ name: 'home' })
  }
  // Sinon, laisser passer
  else {
    next()
  }
})

export default router

