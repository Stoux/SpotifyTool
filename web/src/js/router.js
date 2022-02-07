import {createRouter, createWebHistory} from "vue-router";
import Home from "./routes/Home";
import Login from "./routes/Login";
import Account from "./routes/Account";
import store from "./store";


const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: '/', component: Home },
        { path: '/login', component: Login },
        { path: '/account', component: Account },
        { path : '/authenticated', redirect: to => {
            store.commit('newAccessToken', to.query.token)
            return { path: '/account', query: {} }
        }}
    ],
})

// Set up login guard
router.beforeEach((to, from) => {
    if (to.path === "/login" || store.state.accessToken) {
        console.log('Can access');
        return true
    } else {
        console.log('/login')
        return '/login'
    }
})


export default router