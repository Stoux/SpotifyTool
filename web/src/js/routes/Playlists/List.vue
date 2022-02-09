<template>
  <div>
    <input class="form-control form-control-lg" type="text" placeholder="Search..." v-model="search" autofocus>
    <hr />

    <div class="list-group">
      <router-link :to="'/playlists/' + playlist.id" tag="button" type="button" class="list-group-item list-group-item-action bg-dark text-white"
                   v-for="playlist of displayPlaylists">
        <div style="display: flex;">
          <span>{{ playlist.name }}</span>
          <div style="flex-grow: 1"></div>
          <span class="badge bg-secondary rounded-pill m-1" v-if="playlist.collaborative">
            Collab
          </span>
          <span>{{ playlist.owner_display_name }}</span>
        </div>
      </router-link>
    </div>

  </div>
</template>

<script>
import {mapState} from "vuex";

export default {
  name: "PlaylistList",
  data() {
    return {
      search: '',
    }
  },
  computed: {
    ...mapState([
      'playlists',
    ]),
    displayPlaylists() {
      let playlists = this.playlists ? this.playlists : [];

      // Filter
      if (this.search) {
        const s = this.search.toLowerCase()
        playlists = playlists.filter(playlist => {
          return playlist.name.toLowerCase().indexOf(s) >= 0
              || playlist.owner_display_name.toLowerCase().indexOf(s) >= 0;
        })
      }

      // Order
      playlists.sort((a,b) => {
        return a.name.localeCompare(b.name)
      })

      return playlists
    }
  },
}
</script>

<style scoped>

</style>