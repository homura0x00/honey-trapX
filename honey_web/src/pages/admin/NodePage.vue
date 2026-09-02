<template>
    <p>Node Management</p>
    <a-input-search :style="{width:'320px'}" placeholder="Please enter something" />
    <div :style="{marginBottom: '10px'}"></div>
    <a-table row-key="name" :columns="columns" :data="dataList" :row-selection="rowSelection" v-model:selectedKeys="selectedKeys">
      <template #ip="{record }">
        <router-link :to="`/admin/network`">{{ record.ip }}</router-link>
      </template>
      <template #optional="{ record }">
        <div class="optional-btn">
          <a-button @click="handleEdit" type="text">操作</a-button>
          <a-button @click="$modal.info({ title:'Name', content:record.name })" type="text">详情</a-button>
          <a-button @click="handleDelete" type="text" status="danger">删除</a-button>
        </div>
      </template>
    </a-table>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
const columns = [
  { title: '#', dataIndex: 'key' },
  { title: 'Name', dataIndex: 'name' },
  { title: 'IP', slotName: 'ip' },
  { title: 'Port', dataIndex: 'port' },
  { title: 'Status', dataIndex: 'status'},
  { title: '网络数', dataIndex: 'network_count' },
  { title: '诱捕IP数', dataIndex: 'capture_ip_count' },
  { title: 'CPU/使用率', dataIndex: 'cpu_usage' },
  { title: '内存/空闲率', dataIndex: 'memory_usage' },
  { title: '磁盘/使用率', dataIndex: 'disk_usage' },
  { title: '创建时间', dataIndex: 'create_time' },
  { title: '操作', slotName: 'optional' },
]

const dataList = reactive([
  {
    key: '1',
    name: 'node1',
    ip: '127.0.1.1',
    port: '8080',
    status: 'running',
    network_count: '3',
    capture_ip_count: '3%',
    cpu_usage: '3%',
    memory_usage: '30%',
    disk_usage: '40%',
    create_time: '20025-05-01',
  }
]);

const selectedKeys = ref<string[]>([]);
const rowSelection = {
  type: 'checkbox',
  showCheckedAll: true,
}

const handleEdit = () => {}

const handleDelete = () => {}
</script>

<style scoped>
.optional-btn {
  width: 100%;
  display: flex;
  justify-content: space-around;
}

</style>