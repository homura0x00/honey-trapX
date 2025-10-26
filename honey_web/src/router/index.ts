import {createRouter, createWebHistory} from 'vue-router'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/admin',
            name: 'Admin',
            component: () => import('../views/AdminView.vue'),
            props: true,
            children: [
                {
                    path: '',
                    name: 'Home',
                    component: () => import('../views/admin/HomeView.vue'),
                },
                {
                    path: 'alert',
                    name: 'Alert',
                    component: () => import('../views/admin/AlertView.vue'),
                },
                {
                    path: 'node',
                    name: 'Node',
                    component: () => import('../views/admin/NodeView.vue'),
                },
                {
                    path: 'network',
                    name: 'Network',
                    component: () => import('../views/admin/NetworkView.vue'),
                },
                {
                    path: 'virtual',
                    name: 'Virtual',
                    component: () => import('../views/admin/VirtualView.vue'),
                },
                {
                    path: 'user',
                    name: 'User',
                    component: () => import('../views/admin/UserManagement.vue'),
                }
            ]
        },
        {
            path: '/dashboard',
            name: 'dashboard',
            component: () => import('../views/UserView.vue'),
            children: [
                {
                    path: '',
                    name: 'Home',
                    component: () => import('../views/user/HomeView.vue'),
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
