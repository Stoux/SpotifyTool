<template>
  <div>
    <Changelog :tracks="tracks ? tracks : []"
               :fetching="tracks === undefined || fetchingOffset !== undefined"
               :reached-end="!hasNext && fetchingOffset === undefined"
               :display-playlist="true"
               @next-page="nextPage">
      <h2>Combined changelog of all playlists</h2>
    </Changelog>
  </div>
</template>

<script>
import Changelog from "./Parts/Changelog";
import PlaylistTrack from "./Parts/PlaylistTrack";
import {getApi} from "../../api/api";

const perPage = 25

export default {
  name: "CombinedChangelog",
  components: {Changelog, PlaylistTrack},
  data() {
    return {
      /** @type {SpotifyPlaylistTrack[]} */
      tracks: undefined,
      offset: -perPage,
      hasNext: true,
      fetchingOffset: undefined,
      isPlaying: undefined,
    }
  },
  methods: {
    nextPage() {
      if (this.fetchingOffset !== undefined || !this.hasNext) {
        // Already fetching
        return;
      }

      this.offset += perPage
      this.fetchingOffset = this.offset
      this.hasNext = false
      getApi().get(`/playlists/combined-changelog?offset=${this.offset}&limit=${perPage}`).then(result => {
        // Get the tracks from the response
        /** @type {SpotifyPlaylistTrack[]} */
        const tracks = result.data
        if (this.tracks === undefined) {
          this.tracks = tracks
        } else {
          this.tracks.push(...tracks)
        }

        // Check if there's a next page
        if (tracks.length === perPage) {
          this.hasNext = true
        }

        this.fetchingOffset = undefined
      })
    },
  },
  mounted() {
    this.nextPage()
  },
}
</script>

<style scoped>

</style>