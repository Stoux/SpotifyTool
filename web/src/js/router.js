import {createRouter, createWebHistory} from "vue-router";
import Home from "./routes/Home";
import Login from "./routes/Login";
import Account from "./routes/Account";
import store from "./store";
import PlaylistOverview from "./routes/Playlists/Overview";
import PlaylistList from "./routes/Playlists/List";
import PlaylistDetail from "./routes/Playlists/Detail";
import TracksOverview from "./routes/Tracks/Overview";
import BackupOverview from "./routes/Backup/Overview";
import BackupList from "./routes/Backup/List";
import BackupEdit from "./routes/Backup/Edit";
import CombinedPlaylistChangelog from "./routes/Playlists/CombinedChangelog";
import TracksSearch from "./routes/Tracks/TracksSearch.vue";


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
                    path: 'combined-changelog',
                    component: CombinedPlaylistChangelog,
                },
                {
                    path: ':id',
                    component: PlaylistDetail,
                    props: true,
                }
            ]
        },

        // Track/songs routes
        {
            path: '/tracks',
            component: TracksOverview,
            children: [
                {
                    path: '',
                    component: TracksSearch,
                },
            ]
        },

        // Backup / Sync configs
        {
            path: '/backups',
            component: BackupOverview,
            children: [
                {
                    path: '',
                    component: BackupList,
                },
                {
                    path: 'new',
                    component: BackupEdit,
                    props: {
                        newSync: true,
                    }
                },
                {
                    path: ':id',
                    component: BackupEdit,
                    props: true,
                }
            ],
        }


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