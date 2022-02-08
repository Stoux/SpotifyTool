import axios from "axios";
import store from "../store";
import {API_ROOT} from "../constants";



export function getApi() {
    const headers = {};
    if (store.state.accessToken) {
        headers.Authorization = store.state.accessToken
    }

    return axios.create({
        baseURL: API_ROOT,
        headers,
    })
}

