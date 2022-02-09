<template>
  <div ref="root">
    <div v-if="playlist">
      <h2>Changelog of {{ playlist.name }}</h2>
      <ul class="btn-group p-0 m-1">
        <a :href="'https://open.spotify.com/playlist/' + playlist.id" class="btn btn-outline-primary" target="_blank">Open on Spotify</a>
        <a :href="'https://open.spotify.com/user/' + playlist.owner_id" class="btn btn-outline-secondary" target="_blank">
          View user '{{ playlist.owner_display_name }}'
        </a>
      </ul>

      <hr />

      <ol class="list-group">
        <li class="list-group-item bg-dark text-white d-flex flex-row align-items-center" v-for="track of tracks">
          <div style="padding-right: 16px;" class="">
            <div style="width: 64px; height: 64px">
              <img class="img-thumbnail img-fluid" :src="track.image" v-if="track.image" />
              <div style="width: 100%; height: 100%;" class="d-flex align-items-center justify-content-center" v-else>
                <i class="bi bi-file-earmark-image-fill fs-3"></i>
              </div>
            </div>
          </div>
          <div class="flex-grow-1 lh-lg">
            <a :href="'https://open.spotify.com/track/' + track.TrackId" rel="noreferrer noopener" target="_blank" class="text-white">
              {{ track.Name }}
            </a>
            <span class="text-white-50"> by </span>
            <span>
              <span v-for="(artist, index) of track.Artists.split(' | ')">
                <a href="#" class="text-white">{{ artist }}</a><span class="text-white-50"
                                                                     v-if="index < track.Artists.split(' | ').length - 1"
              >{{ index === (track.Artists.split(' | ').length - 2) ? ' & ' : ', ' }}</span>
              </span>
            </span>
            <span class="text-white-50"> (from <a href="#" class="text-white-50">{{ track.Album }})</a></span>
            <br />
            <span class="text-white-50">
              {{ formatDate(track.timeline) }}
            </span>
            <span class="text-white-50" v-if="track.type === 'added' && track.AddedBy && track.AddedBy.Valid && (playlist.collaborative || track.AddedBy.String !== playlist.owner_id)">
              by <a href="#" class="text-white-50">{{ track.AddedBy.String }}</a>
            </span>
          </div>
          <div class="">
            <span class="badge p-2 text-capitalize" :class="track.type === 'added' ? 'bg-success' : 'bg-danger'">
              {{  track.type }}
            </span>
          </div>
          <div style="padding-left: 16px;" class="">
            <button class="btn btn-outline-primary" disabled>
              <i class="bi-pause-fill" role="img" v-if="isPlaying === track.ID"></i>
              <i class="bi-play-fill" role="img" v-else></i>
            </button>
          </div>
        </li>
      </ol>
    </div>
    <hr ref="loadingTrigger" />
    <div class="d-flex justify-content-center" v-if="!playlist || fetchingForId">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>
    <div v-if="playlist && !hasNext && !fetchingForId" class="p-3 pb-5">
      End of changelog.
    </div>
  </div>
</template>

<script>
import {mapState} from "vuex";
import {getApi} from "../../api/api";

const perPage = 25

export default {
  name: "PlaylistDetail",
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
      /** @type {IntersectionObserver} */
      loadingObserver: undefined,
      isPlaying: undefined,
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
      return this.playlists && this.playlists.find(p => p.id === this.id)
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
    formatDate(date) {
      return new Date(Date.parse(date)).toLocaleString()
    }
  },
  mounted() {
    // Create an observer for the loading trigger at the bottom of the page
    this.loadingObserver = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && this.playlist && !this.fetchingForId && this.tracks.length > 0) {
        console.log('Fetching next page')
        this.nextPage()
      }
    }, {
      root: document.querySelector('.container'),
      rootMargin: '400px',
      threshold: 0,
    })
    this.loadingObserver.observe(this.$refs.loadingTrigger)

    // Fetch the first page of tracks
    this.fetchTracks()
  },
  beforeUnmount() {
    this.loadingObserver.disconnect();
  }


}
</script>

<style scoped lang="scss">



</style>