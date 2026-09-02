<template>
  <div id="home-view">
    <div class="stat-cards" :style="{ marginBottom: '20px'}">
      <StatCard :nums="nodes" name="Node">
        <Workflow :size="32" color="#64ffda" />
      </StatCard>
      <StatCard :nums="nodes" name="Network">
        <Workflow :size="32" color="#64ffda" />
      </StatCard>
      <StatCard :nums="nodes" name="Docker">
        <Workflow :size="32" color="#64ffda" />
      </StatCard>
      <StatCard :nums="nodes" name="Cloud">
        <Workflow :size="32" color="#64ffda" />
      </StatCard>
      <StatCard :nums="nodes" name="IP">
        <Workflow :size="32" color="#64ffda" />
      </StatCard>
    </div>
    <div class="rank-list" :style="{ marginBottom: '20px'}">
      <a-row :gutter="24">
        <a-col :span="6">
          <div class="rank-card">
            <div>攻击排行</div>
            <a-table :columns="columns1" :data="data1" style="margin-top: 10px" :bordered="false" :pagination="false" />
          </div>
        </a-col>
        <a-col :span="9">
          <div class="rank-card">
            <div>IP排行</div>
            <a-table :columns="columnsIp" :data="dataIp" style="margin-top: 10px" :bordered="false" :pagination="false" />
          </div>
        </a-col>
        <a-col :span="9">
          <div class="rank-card">
            <div>服务排行</div>
            <a-table :columns="columnService" :data="dataSerive" style="margin-top: 10px" :bordered="false" :pagination="false" />
          </div>
        </a-col>
      </a-row>
    </div>
    <div class="line-echart">
      <div>告警时间</div>
      <div :style="{ width: '100%', display: 'flex', justifyContent: 'center' }">
        <div :style="{ width: '1440px', margin: '10px' }">
          <DateLineECharts />
        </div>
      </div>
  
    </div>
  </div>
</template>

<script setup lang="ts">
import DateLineECharts from "@/components/DateLineECharts.vue";
import StatCard from "@/components/StatCard.vue";
import { Workflow } from "lucide-vue-next";
import { reactive, ref } from "vue";
const nodes = ref<number>(32);

const columns1 = [
  {
    title: 'Top',
    dataIndex: 'key',
  },
  {
    title: "type",
    dataIndex: "attack_type",
  },
  {
    title: "total",
    dataIndex: "total",
  },
];

const data1 = reactive([
  {
    key: "1",
    attack_type: "Jane Doe",
    total: 6,
  },
  {
    key: "2",
    attack_type: "Alisa Ross",
    total: 3,
  },
  {
    key: "3",
    attack_type: "Kevin Sandra",
    total: 2,
  },
  {
    key: "4",
    attack_type: "Kevin Sandra",
    total: 2,
  },
]);

const columnsIp = [
  {
    title: 'Top',
    dataIndex: 'key',
  },
  {
    title: "type",
    dataIndex: "ip_addr",
  },
  {
    title: "total",
    dataIndex: "total",
  },
  {
    title: "last time",
    dataIndex: "last_time",
  },
];

const dataIp = reactive([
  {
    key: "1",
    ip_addr: "Jane Doe",
    total: 6,
    last_time: "2025-06-30 12:00:00"
  },
]);

const columnService = [
  {
    title: 'Top',
    dataIndex: 'key',
  },
  {
    title: "service",
    dataIndex: "service_info",
  },
  {
    title: "total",
    dataIndex: "total",
  },
  {
    title: "last time",
    dataIndex: "last_time",
  },
];

const dataSerive = reactive([
  {
    key: "1",
    service_info: "Jane Doe",
    total: 6,
    last_time: "2025-06-30 12:00:00"
  },
]);

</script>

<style scoped>
#home-view .stat-cards {
  display: grid;
  width: 100vw;
  max-width: max-content;
  grid-template-columns: repeat(5, 25%);
}

#home-view .rank-list .rank-card {
    /* width: 300px; */
    border: 1px solid gray;
    border-radius: 15px;
    padding: 15px;
    /* background-color: var(--color-neutral-4); */
}

#home-view .line-echart {
  padding: 15px;
  border: 1px solid gray;
  border-radius: 10px;
}
</style>
