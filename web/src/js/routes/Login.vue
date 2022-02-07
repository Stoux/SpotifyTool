<template>
  <h1>Login</h1>
  <p>
    Note: this module only works if I have whitelisted your account to use this app. <br />
    Otherwise, you'll just hit a 401 or something.
  </p>
  <button @click="login" :disabled="!buttonEnabled" class="btn btn-success">
    Login with Spotify
  </button>
</template>

<script>
import {getApi} from "../api/api";

export default {
  name: "Login",
  data() {
    return {
      buttonEnabled: true,
    }
  },
  methods: {
    async login() {
      this.buttonEnabled = false;

      const result = await getApi().get("/auth/start")
      window.location.href = result.data.url
    }
  }
}

</script>

<style scoped>

</style>