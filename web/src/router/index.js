import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '@/store'
import Login from '@/views/Login'
import Layout from '@/views/Layout'
import User from '@/views/User'
import Users from '@/views/Users'
import Files from '@/views/Files'
import Error403 from '@/views/errors/403'
import Error404 from '@/views/errors/404'
import Error500 from '@/views/errors/500'
import {baseURL} from '@/utils/constants'
Vue.use(VueRouter)

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
        beforeEnter: (to, from, next) => {
            if (store.getters.isLogged) {
                return next({ path: '/files' })
            }

            document.title = 'Login'
            next()
        }
    },
    {
        path: '/*',
        component: Layout,
        children: [
            {
                path: '/files/*',
                name: 'Files',
                component: Files,
                meta: {
                    requiresAuth: true
                }
            },

            {
                path: '/403',
                name: 'Forbidden',
                component: Error403
            },
            {
                path: '/404',
                name: 'Not Found',
                component: Error404
            },
            {
                path: '/500',
                name: 'Internal Server Error',
                component: Error500
            },
            {
                path: '/files',
                redirect: {
                    path: '/files/'
                }
            },
            {
                path: '/users',
                name: 'Users',
                component: Users,
                meta: {
                    requiresAdmin: true
                }
            },
            {
                path: '/users/*',
                name: 'User',
                component: User,
                meta: {
                    requiresAdmin: true
                }
            },
            {
                path: '/*',
                redirect: to => `/files${to.path}`
            },

        ]
    }
]

// eslint-disable-next-line no-new
const router = new VueRouter({
  mode: 'history',
  base: baseURL,
  routes
});
router.beforeEach((to, from, next) => {
    document.title = to.name

    if (to.matched.some(record => record.meta.requiresAuth)) {
        if (!store.getters.isLogged) {
            next({
                path: '/login',
                query: { redirect: to.fullPath }
            })

            return
        }

        if (to.matched.some(record => record.meta.requiresAdmin)) {
            if (!store.state.user.perm.admin) {
                next({ path: '/403' })
                return
            }
        }
    }

    next()
})
export default router
