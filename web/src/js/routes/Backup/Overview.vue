<template>
  <div>
    <h1>Backups / Sync</h1>
    <p>
      Here you can configure backups / syncs for playlists. <br />
      The app will automatically sync any new tracks added to playlist A to the destination playlist B. It will do this ~daily.
    </p>
    <router-view></router-view>
  </div>
</template>

<script>
import {getApi} from "../../api/api";
import {mapState} from "vuex";

export default {
  name: "BackupOverview",
  computed: {
    ...mapState([
        'backupConfigs',
    ]),
  },
  mounted() {
    if (this.backupConfigs === undefined) {
      this.$store.dispatch('fetchBackupConfigs')
    }
  }
}
</script>

<style scoped>

</style>