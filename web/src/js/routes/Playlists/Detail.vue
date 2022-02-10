<template>
  <div ref="root">
    <Changelog v-if="playlist" :tracks="tracks" :fetching="!playlist || fetchingForId" :reached-end="playlist && !hasNext && !fetchingForId" @next-page="nextPage">
      <h2>Changelog of {{ playlist.name }}</h2>
      <ul class="btn-group p-0 m-1">
        <a :href="'https://open.spotify.com/playlist/' + playlist.id" class="btn btn-outline-primary" target="_blank">Open on Spotify</a>
        <a :href="'https://open.spotify.com/user/' + playlist.owner_id" class="btn btn-outline-secondary" target="_blank">
          View user '{{ playlist.owner_display_name }}'
        </a>
      </ul>
    </Changelog>
  </div>
</template>

<script>
import {mapState} from "vuex";
import {getApi} from "../../api/api";
import PlaylistTrack from "./Parts/PlaylistTrack";
import Changelog from "./Parts/Changelog";

const perPage = 25

export default {
  name: "PlaylistDetail",
  components: {Changelog, PlaylistTrack},
  props: {
    id: String,
  },
  data() {
    return {
      tracksForId: undefined,
      fetchingForId: undefined,
      /** @type {SpotifyPlaylistTrack[]} */
      tracks: [],
      offset: 0,
      hasNext: true,
      isPlaying: undefined,
    }
  },
  computed: {
    ...mapState([
        "idToPlaylist"
    ]),
    /**
     * @returns {SpotifyPlaylist}
     */
    playlist() {
      return this.idToPlaylist && this.idToPlaylist[this.id]
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
      this.fetchingForId = undefined
      this.tracks = []
      this.hasNext = true
      this.offset = -perPage
      this.nextPage()
    },
    nextPage() {
      if (this.fetchingForId || !this.hasNext) {
        // Already fetching
        return;
      }

      const playlistId = this.fetchingForId = this.tracksForId

      this.offset += perPage
      this.hasNext = false
      getApi().get(`/playlists/${this.playlist.id}/tracks?offset=${this.offset}&limit=${perPage}`).then(result => {
        // Check if we fetching for the right page
        if (this.fetchingForId !== playlistId) {
          return
        }

        // Get the tracks from the response
        /** @type {SpotifyPlaylistTrack[]} */
        const tracks = result.data
        this.tracks.push(...tracks)

        // Check if there's a next page
        if (tracks.length === perPage) {
          this.hasNext = true
        }

        this.fetchingForId = undefined
      })
    },
  },
  mounted() {
    // Fetch the first page of tracks
    this.fetchTracks()
  },
}
</script>

<style scoped lang="scss">



</style>