import axios from "axios";
import {createStore} from "vuex";
import {get, remove, set} from "./util/storage";

const KEY_ACCESS_TOKEN = 'access_token';


const store = createStore({
    state() {
        return {
            accessToken: get(KEY_ACCESS_TOKEN),
            user: undefined,
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
        }
    }
})

export default store