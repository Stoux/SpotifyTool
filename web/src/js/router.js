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
        {path: '/', component: Home},
        {path: '/login', component: Login},
        {path: '/account', component: Account},
        {
            path: '/authenticated', redirect: to => {
                store.dispatch('newAccessToken', {
                    token: to.query.token, onSuccess: () => {
                        // Redirect to home
                        router.push('/')

                        // Check if first time login
                        store.dispatch('newToast', {
                            title: 'First time login',
                            text: 'It might take a while before all your data has been indexed if this is your first login...',
                            type: 'primary',
                            closeInSeconds: 10,
                        })
                    }
                })
                return {path: '/login', query: {}}
            }
        },

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