<template>
  <main class="d-flex flex-row">
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
      </ul>
      <hr>
      <ul class="nav nav-pills flex-column">
        <li class="nav-item" v-if="loggedIn">
          <router-link to="/account" active-class="active" class="nav-link text-white">
            Account
          </router-link>
        </li>
        <li class="nav-item" v-if="!loggedIn">
          <router-link to="/login" active-class="active" class="nav-link text-white">
            Login
          </router-link>
        </li>
      </ul>
    </div>
    <div class="container vh-100 overflow-auto">
      <router-view></router-view>
    </div>
  </main>
</template>
<script>

import {mapGetters} from "vuex";

export default {
  name: "App",
  data() {
    return {

    }
  },
  computed: {
    ...mapGetters([
        'loggedIn',
    ]),
  }
}

</script>
