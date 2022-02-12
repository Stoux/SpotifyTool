<template>
  <div class="card bg-dark h-100">
    <div class="card-body">
      <h5 class="card-title">{{ title }}</h5>
      <select class="form-select mb-4" :value="modelValue" @change="event => $emit('update:modelValue', event.target.value)" :disabled="disabled">
        <option disabled value="">Select playlist</option>
        <option v-for="playlist of playlists" :value="playlist.id">
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
export default {
  name: "PlaylistPicker",
  props: {
    title: String,
    modelValue: String,
    playlists: Array,
    disabled: Boolean,
  },
  computed: {
    resolvedPlaylist() {
      const modelValue = parseInt(this.modelValue, 10);
      return this.playlists && this.playlists.find(p => p.id === modelValue)
    },
  }

}
</script>

<style scoped>

</style>