import { createRouter, createWebHistory } from 'vue-router'
import { authGuard } from '@auth0/auth0-vue'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: LoginView,
      meta: { guest: true }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      // Lazy load the dashboard
      component: () => import('../views/DashboardView.vue'),
      beforeEnter: authGuard
    },
    // Catch-all redirect to dashboard (or login if unauth)
    {
        path: '/:pathMatch(.*)*',
        redirect: '/dashboard'
    }
  ]
})

export default router
