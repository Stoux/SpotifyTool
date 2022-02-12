<template>
  <div class="card bg-dark h-100">
    <div class="card-body">
      <h5 class="card-title">{{ title }}</h5>
      <select class="form-select mb-4" :value="modelValue" @change="event => $emit('update:modelValue', event.target.value)" :disabled="disabled">
        <option disabled value="">Select playlist</option>
        <optgroup v-for="group of groupedPlaylists" :label="group.name" v-if="groupByUser">
          <option v-for="playlist of group.playlists" :value="playlist.id">
            {{ playlist.name }}
          </option>
        </optgroup>
        <option v-for="playlist of playlists" :value="playlist.id" v-else>
          {{ playlist.name }} ({{ playlist.owner_display_name }})
        </option>
      </select>
      <p class="card-text" v-if="!resolvedPlaylist">
        Selected playlist data not available.
      </p>
      <p class="card-text" v-else>
        <strong>{{ resolvedPlaylist.name }}</strong>
        <span class="text-white-50"> by </span>
        {{ resolvedPlaylist.owner_display_name }}
        <span v-if="resolvedPlaylist.collaborative" class="text-white-50">(Collab)</span>
        <br />
        <a :href="'https://open.spotify.com/playlist/' + resolvedPlaylist.id" target="_blank">
          <small>ID: {{ resolvedPlaylist.id }}</small>
        </a>
      </p>
    </div>
  </div>
</template>

<script>
import {mapState} from "vuex";

export default {
  name: "PlaylistPicker",
  props: {
    title: String,
    modelValue: String,
    /** @type {SpotifyPlaylist[]} */
    playlists: Array,
    disabled: Boolean,
    groupByUser: {
      type: Boolean,
      default: false,
    }
  },
  computed: {
    ...mapState([
        'user',
    ]),
    groupedPlaylists() {
      if (!this.groupByUser || !this.user) {
        return undefined;
      }

      const userToPlaylists = {}
      userToPlaylists[this.user.spotify_id] = {
        name: this.user.display_name,
        playlists: [],
      }
      /** @type {SpotifyPlaylist} */
      for (const playlist of this.playlists) {
        if (!userToPlaylists.hasOwnProperty(playlist.owner_id)) {
          userToPlaylists[playlist.owner_id] = {
            name: playlist.owner_display_name,
            playlists: [],
          }
        }

        userToPlaylists[playlist.owner_id].playlists.push(playlist)
      }

      return userToPlaylists
    },
    resolvedPlaylist() {
      const modelValue = parseInt(this.modelValue, 10);
      return this.playlists && this.playlists.find(p => p.id === modelValue)
    },
  }

}
</script>

<style scoped>

</style>