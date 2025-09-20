<template>
  <a-layout class="layout-demo">
    <a-layout-sider collapsible breakpoint="xl" style="margin: 10px; border-radius: 5px;">
      <div class="logo" />
      <a-menu
        :default-open-keys="['1']"
        :default-selected-keys="['0_3']"
        :style="{ width: '100%' }"
        @menu-item-click="onClickMenuItem"
      >
      <template v-for="menu in menuList">
        <a-menu-item>
          <template #icon>
            <component :is="menu.icon as Component">
            </component>
          </template>
          {{ menu.title }}
        </a-menu-item>
      </template>
      </a-menu>
      <!-- trigger -->
      <template #trigger="{ collapsed }">
        <IconCaretRight v-if="collapsed"></IconCaretRight>
        <IconCaretLeft v-else></IconCaretLeft>
      </template>
    </a-layout-sider>
    <a-layout class="contrainer">
      <a-layout-header>
        <a-breadcrumb :style="{ margin: '16px 0'}">
          <a-breadcrumb-item>
            <House :size="15" />
          </a-breadcrumb-item>
          <a-breadcrumb-item>List</a-breadcrumb-item>
          <a-breadcrumb-item>App</a-breadcrumb-item>
        </a-breadcrumb>
        <a-button type="text" style="color: var(--color-text-2);">Text</a-button>
      </a-layout-header>
      <a-layout style="margin: 10px; overflow: auto;">
        <a-layout-content style="background: inherit;">
          <RouterView />
        </a-layout-content>
      </a-layout>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { Message } from "@arco-design/web-vue";
import {
  IconCaretRight,
  IconCaretLeft,
} from "@arco-design/web-vue/es/icon";
import { House, CloudCog, GlobeLock, MonitorCog, ShieldAlert, Workflow } from "lucide-vue-next";
import type { Component } from "vue";

const onClickMenuItem = (key: any) => {
  Message.info({ content: `You select ${key}`, showIcon: true });
};

const menuList = [
  {title: 'Home', icon: House, href: '/'},
  {title: 'Alert', icon: ShieldAlert, href: '/alert'},
  {title: 'Node', icon: Workflow, href: '/node'},
  {title: 'Network', icon: CloudCog, href: '/net'},
  {title: 'Virtual', icon: GlobeLock, href: '/virtual'},
  {title: 'System', icon: MonitorCog, href: '/sys'},
]
</script>

<style scoped>
.layout-demo {
  height: calc(100% - 2px);
  background: var(--color-bg-1);
  border: 1px solid var(--color-border);
}
.layout-demo :deep(.arco-layout-sider) .logo {
  height: 32px;
  margin: 12px 8px;
  background: rgba(255, 255, 255, 0.2);
}
.layout-demo :deep(.arco-layout-sider-light) .logo {
  background: var(--color-fill-2);
}
.layout-demo :deep(.arco-layout-header) {
  height: 44px;
  line-height: 64px;
  margin: 10px;
  padding: 0 10;
  border-radius: 5px;
  background: var(--color-bg-2);
  display: flex;
  align-items: center;
  place-content: space-between;
}
.layout-demo :deep(.arco-layout-content) {
  color: var(--color-text-2);
  font-weight: 400;
  font-size: 14px;
  background: var(--color-bg-3);
}
.layout-demo :deep(.arco-layout-content) {
  display: flex;
  flex-direction: column;
  /* justify-content: center; */
  color: var(--color-white);
  font-size: 16px;
  font-stretch: condensed;
  /* text-align: center; */
}
</style>
