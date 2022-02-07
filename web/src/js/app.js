import {createApp} from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

// Create the app
const app = createApp({
    components: { App },
})
app.use(router)
app.use(store)
app.mount('#app')