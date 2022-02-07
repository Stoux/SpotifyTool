<template>
  <div>
    <div v-if="!playlist">
      Loading...
    </div>
    <div v-else>
      <h1>Changelog of {{ playlist.name }}</h1>
      <ol class="list-group list-group-numbered">
        <li class="list-group-item" v-for="track of tracks">
            {{ track.Name }} - {{ track.Artists }} <br />
          Added: {{ track.AddedAt.Time }} {{ track.DeletedAt ? '- Deleted: ' + track.DeletedAt : '' }}
        </li>
      </ol>
    </div>
  </div>
</template>

<script>
import {mapState} from "vuex";
import {getApi} from "../../api/api";

export default {
  name: "Detail",
  props: {
    id: String,
  },
  data() {
    return {
      tracksForId: undefined,
      /** @type {SpotifyPlaylistTrack[]} */
      tracks: [],
    }
  },
  computed: {
    ...mapState([
        "playlists"
    ]),
    /**
     * @returns {SpotifyPlaylist}
     */
    playlist() {
      return this.playlists.find(p => p.id === this.id)
    },
  },
  watch: {
    playlist() {
      this.fetchTracks()
    },
  },
  methods: {
    fetchTracks() {
      if (!this.playlist || this.playlist.id === this.tracksForId) {
        return
      }

      this.tracksForId = this.playlist.id
      getApi().get(`/playlists/${this.playlist.id}/tracks`).then(result => {
        this.tracks = result.data
      })
    },
  },
  mounted() {
    this.fetchTracks()
  }


}
</script>

<style scoped>

</style>