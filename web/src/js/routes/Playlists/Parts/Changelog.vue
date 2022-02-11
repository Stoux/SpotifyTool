<template>
  <div>
    <slot></slot>

    <hr/>

    <ol class="list-group">
      <PlaylistTrack v-for="track of tracks" :track="track" :display-playlist="displayPlaylist"/>
    </ol>

    <hr ref="loadingTrigger"/>
    <div class="d-flex justify-content-center" v-if="fetching">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <div v-if="reachedEnd" class="p-3 pb-5">
      End of changelog.
    </div>

  </div>
</template>

<script>
import PlaylistTrack from "./PlaylistTrack";

export default {
  name: "Changelog",
  components: {PlaylistTrack},
  emits: [
      'next-page',
  ],
  props: {
    tracks: Array,
    fetching: Boolean,
    reachedEnd: Boolean,
    displayPlaylist: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      loadingObserver: undefined,
    }
  },
  mounted() {
    // Create an observer for the loading trigger at the bottom of the page
    this.loadingObserver = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && !this.fetching && this.tracks.length > 0) {
        console.log('Fetching next page')
        this.$emit('next-page')
      }
    }, {
      root: document.querySelector('.js-observe-container'),
      rootMargin: '400px',
      threshold: 0,
    })
    this.loadingObserver.observe(this.$refs.loadingTrigger)
  },
  beforeUnmount() {
    this.loadingObserver.disconnect();
  }
}
</script>

<style scoped>

</style>