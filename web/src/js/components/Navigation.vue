<template>
  <div class="d-flex flex-column flex-shrink-0 p-3 text-white bg-dark vh-100" style="width: 280px;">
    <a href="/" class="d-flex align-items-center mb-3 mb-md-0 me-md-auto text-white text-decoration-none">
      <span class="fs-4">Spotify Tools</span>
    </a>
    <hr>
    <ul class="nav nav-pills flex-column mb-auto">
      <li class="nav-item">
        <router-link to="/" active-class="active" class="nav-link text-white">
          Home
        </router-link>
      </li>
      <li class="nav-item">
        <router-link to="/playlists" active-class="active" class="nav-link text-white">
          Playlists
        </router-link>
      </li>
      <li class="nav-item">
        <router-link to="/playlists/combined-changelog" active-class="active" class="nav-link text-white">
          Combined changelog <small>(all playlists combined)</small>
        </router-link>
      </li>
      <li class="nav-item">
        <router-link to="/tracks" active-class="active" class="nav-link text-white">
          Find tracks
        </router-link>
      </li>
      <li class="nav-item">
        <router-link to="/backups" active-class="active" class="nav-link text-white">
          Backups / Sync
        </router-link>
      </li>
    </ul>
    <hr>
    <ul class="nav nav-pills flex-column">
      <li class="nav-item" v-if="loggedIn">
        <router-link to="/account" active-class="active" class="nav-link text-white">
          Account: {{ user ? user.display_name : '...' }}
        </router-link>
      </li>
      <li class="nav-item" v-if="!loggedIn">
        <router-link to="/login" active-class="active" class="nav-link text-white">
          Login
        </router-link>
      </li>
    </ul>
  </div>
</template>

<script>
import {mapGetters, mapState} from "vuex";

export default {
  name: "Navigation",
  computed: {
    ...mapState([
        'user',
    ]),
    ...mapGetters([
      'loggedIn',
    ]),

  }
}
</script>

<style scoped>

</style>