<template>
  <div>
    <input class="form-control form-control-lg" type="text" placeholder="Search..." v-model="search" autofocus>
    <hr />

    <div class="list-group">
      <router-link :to="'/playlists/' + playlist.id" tag="button" type="button" class="list-group-item list-group-item-action f"
                   v-for="playlist of displayPlaylists">
        <div style="display: flex;">
          <span>{{ playlist.name }}</span>
          <div style="flex-grow: 1"></div>
          <span>{{ playlist.owner_display_name }}</span>
        </div>
      </router-link>
    </div>

  </div>
</template>

<script>
import {mapState} from "vuex";

export default {
  name: "List",
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
      if (!this.search) {
        return this.playlists
      }

      const s = this.search.toLowerCase()
      return this.playlists.filter(playlist => {
        return playlist.name.toLowerCase().indexOf(s) >= 0
            || playlist.owner_display_name.toLowerCase().indexOf(s) >= 0;
      })

    }
  },
}
</script>

<style scoped>

</style>