<template>
  <div ref="root">
    <h2>{{ newSync ? 'Create backup / sync' : 'Editing ' + this.id }}</h2>

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
        <textarea class="form-control" id="backup-comment" rows="3" maxlength="1000" :disabled="disabled"
                  v-model="comment"></textarea>
        <div class="d-flex justify-content-end pt-3">
          <div class="p-1">
            <button class="btn btn-danger" :disabled="disabled" @click="showModal = true" v-if="!newSync">
              Delete
            </button>
          </div>
          <div class="p-1">
            <button class="btn btn-success" :disabled="disabled || !canSave" @click="submit">
              Save
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="modal fade" id="exampleModal" :class="{ 'show d-block': showModal }" tabindex="-1"
         aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content bg-dark">
          <div class="modal-header">
            <h5 class="modal-title">Are you sure?</h5>
            <button type="button" class="btn-close" @click="showModal = false" :disabled="disabled" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            Are you sure you want to delete this configuration?
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="showModal = false" :disabled="disabled">Cancel</button>
            <button type="button" class="btn btn-danger" @click="deleteConfig" :disabled="disabled">Delete</button>
          </div>
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
      showModal: false,
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
        const id = parseInt(this.id, 10);
        return this.backupConfigs && this.backupConfigs.find(p => p.ID === id)
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
      this.copyFromOriginal()
    },
  },
  methods: {
    copyFromOriginal() {
      if (this.original) {
        this.source = this.original.source.id
        this.target = this.original.target.id
        this.comment = this.original.comment
      }
    },

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
    },
    async deleteConfig() {
      if (this.newSync) {
        return;
      }


      try {
        this.disabled = true
        await getApi().delete(`playlist-backups/${this.id}`)

        await this.$store.dispatch('newToast', {
          title: 'Deleted config',
          text: `The backup with config ID ${this.id} has been deleted.`,
          type: 'success',
        })

        await this.$store.dispatch('fetchBackupConfigs', {forceFetch: true})
        await this.$router.push("/backups")

      } catch (e) {
        await this.$store.dispatch('newToast', {
          title: 'Failed to delete config',
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
    if (!this.newSync) {
      this.copyFromOriginal()
    }

    this.$store.dispatch('fetchPlaylists')
  }
}
</script>

<style scoped>

</style>