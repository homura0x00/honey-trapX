<template>
  <template v-for="menu in menuList" :key="menu.key">
    <a-sub-menu v-if="menu.children && menu.children.length" :key="menu.key">
      <template #icon>
        <component :is="menu.icon as Component" :size="18" v-if="menu.icon" />
      </template>
      <template #title>{{ menu.title }}</template>

      <!-- 递归渲染子菜单 -->
      <MenuRecursive :menuList="menu.children" />
    </a-sub-menu>
    <a-menu-item v-else :key="menu.key" :path="menu.href">
      <template #icon>
        <component :is="menu.icon as Component" :size="18" v-if="menu.icon" />
      </template>
      {{ menu.title }}
    </a-menu-item>
  </template>
</template>

<script setup lang="ts">
import { type Component, defineProps } from "vue";

// 定义菜单接口
interface MenuItem {
  key: string; // 唯一标识
  title: string;
  icon?: Component;
  href?: string; // 只有叶子节点有href
  children?: MenuItem[]; // 有子菜单的父节点
}

defineProps<{
  menuList: MenuItem[];
}>();
</script>
