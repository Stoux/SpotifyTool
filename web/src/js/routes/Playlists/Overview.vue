<template>
  <div>
    <h1>Playlists</h1>
    <router-view></router-view>
  </div>
</template>

<script>
import {mapState} from "vuex";

export default {
  name: "PlaylistOverview",
  computed: {
    ...mapState([
      'playlists'
    ]),
  },
  mounted() {
    if (this.playlists === undefined) {
      this.$store.dispatch('fetchPlaylists')
    }
  }
}
</script>

<style scoped>

</style>