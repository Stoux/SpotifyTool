<template>
  <div>
    <h1>Playlists</h1>
    <router-view></router-view>
  </div>
</template>

<script>
import {getApi} from "../../api/api";
import {mapState} from "vuex";

export default {
  name: "Overview",
  computed: {
    ...mapState([
        'playlists'
    ]),
  },
  mounted() {
    if (!this.playlists || !this.playlists.length) {
      getApi().get('playlists').then(result => {
        this.$store.commit('newPlaylists', result.data)
      })
    }
  }
}
</script>

<style scoped>

</style>