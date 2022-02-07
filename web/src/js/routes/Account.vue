<template>
  <div>
    <div v-if="!account">
      Loading...
    </div>
    <div v-else>
      <h1>Hi {{ account.display_name }}</h1>
    </div>
  </div>
</template>

<script>
import {getApi} from "../api/api";

export default {
  name: "Account",
  data() {
    return {
      account: undefined,
    }
  },
  mounted() {
    getApi().get('/auth/me').then(result => {
      this.account = result.data
    }).catch(error => {

    })

  }
}
</script>

<style scoped>

</style>