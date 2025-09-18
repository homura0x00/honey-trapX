<template>
  <div id="app">
    <BasicLayout />
  </div>
  <!-- <RouterView /> -->
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import BasicLayout from './layouts/BasicLayout.vue';

</script>

<style scoped> 
#app {
  height: 100vh;
}
</style>
