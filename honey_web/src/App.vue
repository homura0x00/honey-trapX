<template>
  <div id="app">
    <router-view />
  </div>
  <!-- <RouterView /> -->
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'

</script>

<style scoped> 
#app {
  height: 100vh;
  margin: 0;
}
</style>
