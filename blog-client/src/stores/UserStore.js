import { defineStore } from "pinia"; // 引入pinia

export const UserStore = defineStore("admin", {
  state: () => {
    return {
      token: "",
      phoneNumber: "",
      user_id: 0,
    };
  },
  persist: {
    enabled: true,
    // 自定义持久化参数
    strategies: [
      {
        // 自定义key
        key: 'token',
        // 自定义存储方式，默认sessionStorage
        storage: sessionStorage,
        // 指定要持久化的数据，默认所有 state 都会进行缓存，可以通过 paths 指定要持久化的字段，其他的则不会进行持久化。
        paths: ['token'],
      },
      {
        // 自定义key
        key: 'user_id',
        // 自定义存储方式，默认sessionStorage
        storage: sessionStorage,
        // 指定要持久化的数据，默认所有 state 都会进行缓存，可以通过 paths 指定要持久化的字段，其他的则不会进行持久化。
        paths: ['user_id'],
      },
    ]
  },
  getters: {},
});