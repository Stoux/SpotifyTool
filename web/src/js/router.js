import {createRouter, createWebHistory} from "vue-router";
import Home from "./routes/Home";
import Login from "./routes/Login";
import Account from "./routes/Account";
import store from "./store";
import PlaylistOverview from "./routes/Playlists/Overview";
import PlaylistList from "./routes/Playlists/List";
import PlaylistDetail from "./routes/Playlists/Detail";


const router = createRouter({
    history: createWebHistory(),
    routes: [
        // Account / general routes
        { path: '/', component: Home },
        { path: '/login', component: Login },
        { path: '/account', component: Account },
        { path : '/authenticated', redirect: to => {
            store.commit('newAccessToken', to.query.token)
            return { path: '/account', query: {} }
        }},

        // Playlist routes
        {
            path: '/playlists',
            component: PlaylistOverview,
            children: [
                {
                    path: '',
                    component: PlaylistList
                },
                {
                    path: ':id',
                    component: PlaylistDetail,
                    props: true,
                }
            ]
        },


    ],
})

// Set up login guard
router.beforeEach((to, from) => {
    if (to.path === "/login" || store.state.accessToken) {
        return true
    } else {
        console.log('/login')
        return '/login'
    }
})


export default router