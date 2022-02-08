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

    return axios.create({
        baseURL: API_ROOT,
        headers,
    })
}

/**
 * @typedef {interface} GetApiOptions
 * @property {string} [withToken]
 */