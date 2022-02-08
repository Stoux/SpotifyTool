<template>
  <div>
    <div v-if="!user">
      Loading...
    </div>
    <div class="p-3" v-else>
      <h1>Hi {{ user.display_name }}</h1>
      <button class="btn btn-danger" @click="logout">
        Logout
      </button>
    </div>

  </div>
</template>

<script>
import {mapState} from "vuex";

export default {
  name: "Account",
  computed: {
    ...mapState([
      'user',
    ]),
  },
  methods: {
    logout() {
      this.$store.commit('invalidateAccessToken')
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>

</style>