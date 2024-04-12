import { createApp } from "vue";
import "./style.css";
import App from "./App.vue";
import naive from "naive-ui"; // 引入ui框架
import { createDiscreteApi } from "naive-ui"; // 引入createDiscreteApi
import { createPinia } from "pinia"; // 引入pinia
import piniaPluginPersist from 'pinia-plugin-persist'
import { router } from "./common/router"; // 引入路由
import axios from "axios"; // 引入axios
import { UserStore } from "./stores/UserStore" // 引入UserStore

axios.defaults.baseURL = "http://127.0.0.1:8080"; // 服务端地址全局配置
const { message, notification, dialog } = createDiscreteApi(["message", "notification", "dialog"])

const app = createApp(App);

app.provide("axios", axios); // 将axios全局放入
app.provide("message", message)
app.provide("notification", notification)
app.provide("dialog", dialog)
app.provide("serverUrl", axios.defaults.baseURL)

app.use(naive); // 引入ui框架
const store = createPinia()
store.use(piniaPluginPersist)
app.use(store); // 引入pinia

const userStore = UserStore()
// 拦截器传token
axios.interceptors.request.use((config) => {
    config.headers.Authorization = `Bearer ${userStore.token}`
    return config
})

app.use(router); // 引入路由
app.mount("#app");
