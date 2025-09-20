<template>
  <div id="home-view">
    <StatCard :nums="nodes" name="node">
      <Workflow :size="32" color="#64ffda" />
    </StatCard>
    <StatCard :nums="nodes" name="node">
      <Workflow :size="32" color="#64ffda" />
    </StatCard>
    <StatCard :nums="nodes" name="node">
      <Workflow :size="32" color="#64ffda" />
    </StatCard>
    <StatCard :nums="nodes" name="node">
      <Workflow :size="32" color="#64ffda" />
    </StatCard>
  </div>
</template>

<script setup lang="ts">
import StatCard from '@/components/StatCard.vue';
import { Workflow } from 'lucide-vue-next';
import { ref } from 'vue';
const nodes = ref<number>(32)
</script>

<style scoped>
#home-view {
  display: grid;
  width: 100vw;
  max-width: max-content;
  grid-template-columns: repeat(4, 30.3%);
}
</style>