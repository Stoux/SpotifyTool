<template>
  <div>
    <router-link to="/backups/new" tag="button" class="btn btn-primary">
      Add new sync
    </router-link>

    <div class="list-group" v-if="backupConfigs">
      <router-link :to="'/backups/' + config.ID" tag="button" type="button" class="list-group-item list-group-item-action bg-dark text-white"
                   v-for="config of backupConfigs">
        <div class="d-flex flex-column">
          <div class="d-flex align-items-center">
            <span class="badge bg-secondary">
              {{ config.source.name }}
            </span>
            <span>-></span>
            <span class="badge bg-primary">
              {{ config.target.name }}
            </span>
          </div>
          <div>
            <span>Last synced</span>
            <span>{{ config.last_sync }}</span>
          </div>
          <small v-if="config.comment">
            {{ config.comment }}
          </small>
        </div>
      </router-link>
    </div>

  </div>
</template>

<script>
import {mapState} from "vuex";

export default {
  name: "BackupList",
  computed: {
    ...mapState([
      'backupConfigs',
    ]),
  },
}
</script>

<style scoped>

</style>