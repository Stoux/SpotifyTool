import axios from "axios";
import {createStore} from "vuex";
import {get, remove, set} from "./util/storage";

const KEY_ACCESS_TOKEN = 'access_token';


const store = createStore({
    state() {
        return {
            accessToken: get(KEY_ACCESS_TOKEN),
            user: undefined,
            /** @type {SpotifyPlaylist[]} */
            playlists: [],
        }
    },
    getters: {
        loggedIn: state => !!state.accessToken,
    },
    mutations: {
        newAccessToken( state, payload ) {
            state.accessToken = payload
            set(KEY_ACCESS_TOKEN, payload)
        },
        invalidateAccessToken( state, payload ) {
            state.accessToken = ''
            remove(KEY_ACCESS_TOKEN)
        },
        newPlaylists( state, playlists ) {
            state.playlists = playlists
        },
    }
})

export default store

/**
 * @typedef {object} SpotifyPlaylist
 *
 * @property {string} id
 * @property {string} name
 * @property {string} owner_display_name
 * @property {string} owner_id
 * @property {boolean} public
 * @property {boolean} collaborative
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
 * @property {string|'added'|'removed'} type
 * @property {string} timeline
 */