<template>
  <div ref="root">
    <h2>{{ newSync ? 'Create backup / sync' : 'Editing ?' }}</h2>

    <div class="d-flex align-items-stretch">
      <div class="p-2 w-50">
        <PlaylistPicker title="From playlist"
                        :disabled="disabled"
                        v-model="source"
                        :playlists="playlists ? playlists : []"/>
      </div>
      <div class="p-2 w-50">
        <PlaylistPicker title="To playlist"
                        :disabled="disabled"
                        v-model="target"
                        :playlists="modifiablePlaylists"/>
      </div>
    </div>
    <div class="mt-3 p-2">
      <div class="d-flex flex-column">
        <label for="backup-comment" class="form-label">An optional (personal) comment about this sync</label>
        <textarea class="form-control" id="backup-comment" rows="3" maxlength="1000" :disabled="disabled" v-model="comment"></textarea>
        <div class="pt-3 d-flex justify-content-end">
          <button class="btn btn-success" :disabled="disabled || !canSave" @click="submit">
            Save
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {mapState} from "vuex";
import PlaylistPicker from "./Parts/PlaylistPicker";
import {getApi} from "../../api/api";

export default {
  name: "Edit",
  components: {PlaylistPicker},
  props: {
    newSync: {
      type: Boolean,
      default: false,
    },
    id: String,
  },
  data() {
    return {
      source: '',
      target: '',
      comment: '',
      disabled: false,
    }
  },
  computed: {
    ...mapState([
      'user',
      'playlists',
      'backupConfigs',
    ]),
    modifiablePlaylists() {
      if (this.playlists) {
        return this.playlists.filter(p => p.collaborative || (this.user && p.owner_id === this.user.spotify_id))
      } else {
        return []
      }
    },

    /**
     *
     * @returns {PlaylistBackupConfig}
     */
    original() {
      if (this.newSync) {
        return undefined
      } else {
        return this.backupConfigs && this.backupConfigs.find(p => p.ID === this.id)
      }
    },
    canSave() {
      if (!this.source || !this.target) {
        return false
      }

      if (!this.original) {
        return this.newSync
      }


      return this.source !== this.original.source.id
          || this.target !== this.original.target.id
          || this.comment !== this.original.comment
    },

  },
  watch: {
    original() {
      if (this.original) {
        this.source = this.original.source.id
        this.target = this.original.target.id
        this.comment = this.original.comment
      }
    },
  },
  methods: {
    async submit() {
      this.disabled = true

      const api = getApi()
      const method = this.newSync ? api.post : api.put;
      const path = this.newSync ? 'playlist-backups' : `playlist-backups/${this.id}`

      try {
        const result = await method(path, {
          source: this.source,
          target: this.target,
          comment: this.comment,
        })
        /** @type {PlaylistBackupConfig} */
        const config = result.data;

        await this.$store.dispatch('newToast', {
          title: 'Created backup / sync',
          text: 'New backup / sync config has been created with ID ' + config.ID,
          type: 'success',
        })

        await this.$store.dispatch('fetchBackupConfigs', {forceFetch: true})
        await this.$router.push("/backups")
      } catch (e) {
        await this.$store.dispatch('newToast', {
          title: 'Failed to save',
          text: e.response && e.response.data && e.response.data.error
              ? e.response.data.error
              : 'Something went wrong. Try again.',
          type: 'danger',
        })
      } finally {
        this.disabled = false
      }
    }
  },
  mounted() {
    this.$store.dispatch('fetchPlaylists')
  }
}
</script>

<style scoped>

</style>