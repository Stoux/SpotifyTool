<template>
  <div>
    <main class="d-flex flex-row bg-black text-white">
      <Navigation/>
      <div class="container-fluid vh-100 overflow-auto">
        <router-view></router-view>
      </div>
    </main>
    <Toasts/>
  </div>
</template>
<script>

import Navigation from "./components/Navigation";
import Toasts from "./components/Toasts";

export default {
  name: "App",
  components: {Toasts, Navigation},
  data() {
    return {}
  },
  computed: {}
}

</script>
