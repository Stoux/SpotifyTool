import axios from "axios";
import store from "../store";



export function getApi() {
    const headers = {};
    if (store.state.accessToken) {
        headers.Authorization = store.state.accessToken
    }

    return axios.create({
        baseURL: 'http://localhost:8080/',
        headers,
    })
}

