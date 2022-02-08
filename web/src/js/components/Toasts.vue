<template>
  <div class="position-fixed bottom-0 end-0 p-3" style="z-index: 11" v-if="toasts.length">
    <div class="toast show" :class="{ 'mt-2': index > 0 }" role="alert" aria-live="assertive" aria-atomic="true" v-for="(toast, index) of toasts">
      <div class="toast-header" :class="getToastClasses(toast)">
        <strong class="me-auto">{{ toast.title }}</strong>
        <small v-if="toast.autoClose">Closes in {{  toast._timeLeft }}..</small>
        <button type="button" class="btn-close" aria-label="Close" @click.prevent="closeToast(toast.id)"></button>
      </div>
      <div class="toast-body">
        {{ toast.text }}
      </div>
    </div>
  </div>
</template>

<script>
import {mapActions, mapState} from "vuex";

export default {
  name: "Toasts",
  computed: {
    ...mapState([
        'toasts',
    ])
  },
  methods: {
    ...mapActions([
        'closeToast'
    ]),
    getToastClasses(toast) {
      const type = toast.type
      if (!type || type === 'default') {
        return 'text-black';
      } else if (type === 'warning') {
        return 'bg-warning text-black'
      } else {
        return `bg-${type} text-white`
      }
    }
  },
  mounted() {
    setInterval(() => {
      if (this.toasts.length && this.toasts.find(t => t.autoClose)) {
        this.$store.commit('_updateToasts', this.toasts.map(t => {
          if (t.autoClose) {
            t._timeLeft--
            if (t._timeLeft <= 0) {
              return null
            }
          }
          return t
        }).filter(t => t !== null))
      }
    }, 1000);
  }
}

/**
 * @typedef {object} Toast
 *
 * @property {string} id
 * @property {string} title
 * @property {string} text
 * @property {string|'default'|'success'|'danger'|'warning'|'primary'|'secondary'} [type]
 * @property {boolean|false} autoClose
 *
 * @property {number} _timeLeft in seconds (internal, only if autoClose is true)
 */

</script>

<style scoped>

</style>