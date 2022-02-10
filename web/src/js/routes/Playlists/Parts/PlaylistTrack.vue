<template>
  <li class="list-group-item bg-dark text-white d-flex flex-row align-items-center">
    <div style="padding-right: 16px;" class="">
      <div style="width: 64px; height: 64px">
        <img class="img-thumbnail img-fluid" :src="track.image" v-if="track.image"/>
        <div style="width: 100%; height: 100%;" class="d-flex align-items-center justify-content-center" v-else>
          <i class="bi bi-file-earmark-image-fill fs-3"></i>
        </div>
      </div>
    </div>
    <div class="flex-grow-1 lh-lg">
      <a :href="'https://open.spotify.com/track/' + track.TrackId" rel="noreferrer noopener" target="_blank"
         class="text-white">
        {{ track.Name }}
      </a>
      <span class="text-white-50"> by </span>
      <span>
        <span v-for="(artist, index) of artists">
          <a href="#" class="text-white">{{ artist }}</a><span class="text-white-50"
                                                               v-if="index < artists.length - 1"
        >{{ index === (artists.length - 2) ? ' & ' : ', ' }}</span>
        </span>
      </span>
      <span class="text-white-50"> (from <a href="#" class="text-white-50">{{ track.Album }})</a></span>
      <br/>
      <span class="text-white-50">
        {{ formatDate(track.timeline) }}
      </span>
      <span class="text-white-50"
            v-if="track.type === 'added' && track.AddedBy && track.AddedBy.Valid && playlist && (displayPlaylist || playlist.collaborative || track.AddedBy.String !== playlist.owner_id)">
        by <a href="#" class="text-white-50">{{ track.AddedBy.String }}</a>
      </span>
    </div>
    <div class="p-2" v-if="displayPlaylist && playlist">
      <router-link :to="'/playlists/' + playlist.id" class="text-white-50">
        {{ playlist.name }}
      </router-link>
      <a class="text-white-50 m-2" :href="'https://open.spotify.com/playlist/' + playlist.id" target="_blank" >
        <i class="bi bi-box-arrow-up-right" role="img"></i>
      </a>
    </div>
    <div class="">
      <span class="badge p-2 text-capitalize" :class="track.type === 'added' ? 'bg-success' : 'bg-danger'">
        {{ track.type }}
      </span>
    </div>
    <div style="padding-left: 16px;" class="">
      <button class="btn btn-outline-primary" disabled>
        <i class="bi bi-pause-fill" role="img" v-if="isPlaying === track.ID"></i>
        <i class="bi bi-play-fill" role="img" v-else></i>
      </button>
    </div>
  </li>
</template>

<script>/**
 * @param {SpotifyPlaylistTrack} track
 */
import {mapState} from "vuex";

export default {
  name: "PlaylistTrack",
  props: {
    /** @property {SpotifyPlaylistTrack} */
    track: Object,
    isPlaying: {
      type: String,
      default: "",
    },
    displayPlaylist: {
      type: Boolean,
      default: false,
    }
  },
  computed: {
    ...mapState([
      'idToPlaylist',
    ]),
    playlist() {
      return this.idToPlaylist
          ? this.idToPlaylist[this.track.SpotifyPlaylistID]
          : undefined
    },
    artists() {
      return this.track.Artists.split(' | ')
    },
  },
  methods: {
    formatDate(date) {
      return new Date(Date.parse(date)).toLocaleString()
    },
  },
}
</script>

<style scoped>

</style>