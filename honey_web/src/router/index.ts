import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/admin',
      name: 'Admin',
      component: () => import('../views/AdminView.vue'),
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('../views/web/HomeView.vue'),
        }
      ]
    },
    {
      path: '/',
      name: 'Login',
      component: () => import('../views/LoginView.vue'),
    }
  ],
})

router.afterEach(async (to, from, next) => {

})

export default router
