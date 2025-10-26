<template>
  <div class="rank-card">
    <div>Attack Rank</div>
    <div class="container">
      <a-table :data="props.data" :column="props.columns" style="margin-top: 10px" :bordered="false" :pagination="false">
      </a-table>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
    data: Data[];
    columns: Column[];
}
interface Data {
    key: string;
    attack_type: string;
    total: number;
}
interface Column {
    title: string;
    dataIndex: string;
}
const props = defineProps<Props>();
</script>

<style scoped>
.rank-card {
    width: 300px;
    border-radius: 15px;
    padding: 15px;
    /* background-color: darkslateblue; */
}

.rank-card .container :deep(.a-table) {
    text-align: center;
}
</style>
