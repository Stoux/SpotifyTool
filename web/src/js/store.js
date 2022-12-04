import {createStore} from "vuex";
import {get, remove, set} from "./util/storage";
import {getApi} from "./api/api";
import router from "./router";

const KEY_ACCESS_TOKEN = 'access_token';
const ACCESS_TOKEN_TOAST_ID = 'AccessTokenToast';


const store = createStore({
    state() {
        return {
            accessToken: get(KEY_ACCESS_TOKEN),

            /** @type {ToolUser} */
            user: undefined,
            /** @type {SpotifyPlaylist[]} */
            playlists: undefined,
            /** @type {Object.<string, SpotifyPlaylist>} */
            idToPlaylist: undefined,
            /** @type {PlaylistBackupConfig[]} */
            backupConfigs: undefined,

            /** @type {Toast[]} */
            toasts: [],
        }
    },
    getters: {
        loggedIn: state => !!state.accessToken,
    },
    mutations: {
        _setAccessToken(state, {token, user}) {
            state.accessToken = token
            state.user = user
            set(KEY_ACCESS_TOKEN, token)
        },
        invalidateAccessToken(state, payload) {
            state.accessToken = ''
            state.user = undefined
            remove(KEY_ACCESS_TOKEN)
        },


        /**
         *
         * @param state
         * @param {ToolUserPlaylist[]} userPlaylists
         */
        newPlaylists(state, userPlaylists) {
            const spotifyPlaylists = [];
            const idToPlaylist = [];

            // Map the userPlaylists to state values
            userPlaylists.forEach(userPlaylist => {
                const spotifyPlaylist = userPlaylist.spotify_playlist;
                spotifyPlaylist.is_tracked = userPlaylist.is_tracked;
                spotifyPlaylists.push(spotifyPlaylist);
                idToPlaylist[spotifyPlaylist.id] = spotifyPlaylist;
            });

            // Update the state
            state.playlists = spotifyPlaylists.sort((a, b) => a.name.localeCompare(b.name))
            state.idToPlaylist = idToPlaylist
        },
        newBackupConfigs(state, configs) {
            state.backupConfigs = configs
        },

        /**
         * @private shouldn't be used directly
         * @see newAccessToken
         */
        _updateToasts(state, toasts) {
            state.toasts = toasts
        }
    },
    actions: {
        /**
         * Create a new toast
         * @param {string} [id] Optional ID (otherwise one is generated
         * @param {string} title
         * @param {string} [text]
         * @param {number} [closeInSeconds] if the toast should auto close & in how many seconds
         * @param {string|'default'|'success'|'danger'|'warning'|'primary'|'secondary'} [type]
         */
        newToast(context, {id, title, text, closeInSeconds, type}) {
            id = id ? id : (new Date()).getTime() + '-' + Math.random() * 10000
            const toasts = [...context.state.toasts];
            toasts.push({
                id,
                title,
                text,
                type,
                autoClose: closeInSeconds && closeInSeconds > 0,
                _timeLeft: closeInSeconds,
            })
            context.commit('_updateToasts', toasts)

            return id
        },
        /**
         * Close the given toast
         * @param {string} id
         */
        closeToast(context, id) {
            context.commit('_updateToasts', context.state.toasts.filter(t => t.id !== id))
        },

        /**
         * Check & validate a new supplied token.
         * @param context
         * @param {string} token
         * @param {function()} [onSuccess] Callback on success
         * @param {bool|false} [noSuccessToast] Hide the success toast
         * @returns {Promise<void>}
         */
        async newAccessToken(context, {token, onSuccess, noSuccessToast}) {
            // Notify the user of what's happening
            await context.dispatch('newToast', {
                id: ACCESS_TOKEN_TOAST_ID,
                title: 'Getting user',
                text: 'We\'re getting your user data, one moment please.',
            })

            // Fetch the user
            try {
                const user = (await getApi({withToken: token}).get('/auth/me')).data
                context.commit('_setAccessToken', {user, token})
                await context.dispatch('closeToast', ACCESS_TOKEN_TOAST_ID)
                await context.dispatch('newToast', {
                    title: 'Welcome back',
                    text: 'You\'re logged in as ' + user.display_name,
                    closeInSeconds: 5,
                    type: 'success',
                })

                if (onSuccess) {
                    onSuccess()
                }
            } catch {
                await context.dispatch('_onInvalidAccessToken')
            }
        },

        async _onInvalidAccessToken(context) {
            context.commit('invalidateAccessToken')
            await context.dispatch('closeToast', ACCESS_TOKEN_TOAST_ID)
            await context.dispatch('newToast', {
                title: 'Invalid token',
                text: 'Unable to fetch user details. Try to logging in again.',
                type: 'danger',
            })
            await router.push('/login')
        },

        /**
         * @param context
         * @param {{forceFetch: bool}} payload
         * @returns {Promise<void>}
         */
        async fetchPlaylists(context, payload ) {
            if (!(payload && payload.forceFetch) && context.state.playlists !== undefined) {
                // Already fetching or no results.
                return
            }
            context.commit('newPlaylists', [])
            try {
                const result = await getApi().get('playlists')
                context.commit('newPlaylists', result.data)
            } catch {
                context.commit('newPlaylists', undefined)
                await context.dispatch('newToast', {
                    title: 'Error',
                    text: 'Something went wrong trying to fetch the playlists, try again.',
                    type: 'danger',
                    closeInSeconds: 10,
                })
            }
        },
        /**
         * @param context
         * @param {{forceFetch: bool}} payload
         * @returns {Promise<void>}
         */
        async fetchBackupConfigs(context, payload ) {
            if (!(payload && payload.forceFetch) && context.state.backupConfigs !== undefined) {
                // Already fetching or no results.
                return
            }
            context.commit('newBackupConfigs', [])
            try {
                const result = await getApi().get('playlist-backups')
                context.commit('newBackupConfigs', result.data)
            } catch {
                context.commit('newBackupConfigs', undefined)
                await context.dispatch('newToast', {
                    title: 'Error',
                    text: 'Something went wrong trying to fetch the backup configs, try again.',
                    type: 'danger',
                    closeInSeconds: 10,
                })
            }
        },

    },
})

if (store.state.accessToken) {
    store.dispatch('newAccessToken', {token: store.state.accessToken})
}

export default store

/**
 * @typedef {object} ToolUser
 * @property {number} ID
 * @property {string} display_name
 * @property {string} spotify_id
 * @property {{String: string}} [plan]
 */

/**
 * @typedef {object} ToolUserPlaylist
 * @property {number} tool_user_id
 * @property {ToolUser} [tool_user]
 * @property {string} spotify_playlist_id
 * @property {SpotifyPlaylist} spotify_playlist
 * @property {boolean} is_tracked Whether the user is following this playlist
 */

/**
 * @typedef {object} SpotifyPlaylist
 *
 * @property {string} id
 * @property {string} name
 * @property {string} owner_display_name
 * @property {string} owner_id
 * @property {boolean} public
 * @property {boolean} collaborative
 * @property {boolean} [is_tracked] Whether the user is following this playlist
 */


/**
 * @typedef {object} SpotifyPlaylistTrack
 *
 * @property {string} ID
 * @property {string} Name
 * @property {string} Artists
 * @property {string} Album
 * @property {string} CreatedAt
 * @property {{Time: string}} AddedAt
 * @property {{String: string, Valid: boolean}} AddedBy
 * @property {string} DeletedAt
 * @property {string} TrackId
 * @property {string} SpotifyPlaylistID
 * @property {string|'added'|'removed'} type
 * @property {string} timeline
 */

/**
 * @typedef {object} PlaylistBackupConfig
 *
 * @property {number} ID
 * @property {SpotifyPlaylist} source
 * @property {SpotifyPlaylist} target
 * @property {string} last_sync
 * @property {string} [comment] Optional comment by the user
 */