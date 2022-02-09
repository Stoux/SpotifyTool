import axios from "axios";
import store from "../store";
import {API_ROOT} from "../constants";

/**
 *
 * @param {GetApiOptions} [options]
 *
 * @returns {AxiosInstance}
 */
export function getApi(options) {
    const headers = {};
    if ( options && options.withToken ) {
        headers.Authorization = options.withToken
    } else if (store.state.accessToken) {
        headers.Authorization = store.state.accessToken
    }

    let instance = axios.create({
        baseURL: API_ROOT,
        headers,
    });

    // Middleware for handling 401 unauthorized responses
    instance.interceptors.response.use(
        undefined,
        error => {
            if (error && error.response && error.response.status === 401) {
                store.dispatch('_onInvalidAccessToken')
            }
            throw error
        }
    )

    return instance
}

/**
 * @typedef {interface} GetApiOptions
 * @property {string} [withToken]
 */