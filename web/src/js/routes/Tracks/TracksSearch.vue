<template>
  <div>
    <input class="form-control form-control-lg" type="text" placeholder="Search..." v-model="query" autofocus>
    <hr/>

    <p v-if="searching">Searching...</p>
    <p v-if="!searching && foundTracks === undefined">Search for a artist, title or spotify ID.</p>

    <div class="list-group" v-if="foundTracks !== undefined">
      <div class="list-group-item list-group-item-action bg-dark text-white" v-if="foundTracks.length === 0">
        No items found.
      </div>

      <div class="list-group-item list-group-item-action bg-dark text-white d-flex flex-row align-items-center"
           v-for="track of foundTracks">
        <a class="d-flex flex-column text-white gap-1" target="_blank"
           :href="'https://open.spotify.com/track/' + track.TrackId">
          <div>{{ track.TrackName }}</div>
          <small>{{ track.Artists }}</small>
        </a>
        <div class="flex-grow-1"><!-- spacer --></div>
        <div class="d-flex flex-column align-items-end gap-1">
          <a :href="'https://open.spotify.com/playlist/' + track.PlaylistId" target="_blank" class="text-white-50">
            {{ track.PlaylistName }}
          </a>

          <span class="badge bg-success rounded-pill m-1">
            + {{ formatDate( track.AddedAt ) }}
          </span>
          <span class="badge bg-danger rounded-pill m-1" v-if="track.DeletedAt">
            - {{ formatDate( track.DeletedAt ) }}
          </span>


        </div>
      </div>
    </div>

  </div>
</template>

<script>
import {getApi} from "../../api/api";

export default {
  name: "TracksSearch",
  data() {
    return {
      query: '',
      searching: false,
      /** @type {SearchFoundTrack[]|undefined} */
      foundTracks: undefined,
      debounce: undefined,
    }
  },
  watch: {
    query() {
      clearTimeout(this.debounce);
      const query = this.query;
      if (query === '') {
        this.searching = false;
        this.foundTracks = undefined;
        return;
      }

      this.searching = true;
      this.debounce = setTimeout(() => {
        getApi().get(`/tracks/search`, {params: {query}}).then(result => {
          if (query === this.query) {
            this.foundTracks = result.data;
          }
        }).finally(() => {
          if (query === this.query) {
            this.searching = false;
          }
        })
      }, 250);
    },
  },
  methods: {
    formatDate(date) {
      return new Date(Date.parse(date)).toLocaleString()
    },
  },
}

/**
 * @typedef {object} SearchFoundTrack
 * @property {string} PlaylistId
 * @property {string} PlaylistName
 * @property {string} TrackId
 * @property {string} TrackName
 * @property {string} Artists
 * @property {string} Album
 * @property {string} AddedAt
 * @property {string|null} [DeletedAt]
 */

</script>

<style scoped>

</style>