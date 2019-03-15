import { sync } from 'vuex-router-sync'
import App from './App.vue'
import router from './router'
import i18n from './i18n'
import store from './store/index'
import Vue from '@/utils/vue'
import { loginPage } from '@/utils/constants'
import { login, validateLogin } from '@/utils/auth'
Vue.config.productionTip = false;

sync(store, router)


async function start () {
    if (loginPage) {
        await validateLogin()
    } else {
        await login('', '', '')
    }
    new Vue({
        router,
        store,
        i18n,
        render: h => h(App)
    }).$mount("#app");

    /*new Vue({
        el: '#app',
        store,
        router,
        i18n,
        template: '<App/>',
        components: { App }
    })*/
}
start()
