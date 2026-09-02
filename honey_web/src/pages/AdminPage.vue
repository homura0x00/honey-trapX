<script setup lang="ts">
import BasicLayout from '@/layouts/BasicLayout.vue';
import {
  House,
  CloudCog,
  GlobeLock,
  MonitorCog,
  ShieldAlert,
  Workflow,
  User2,
} from "lucide-vue-next";

const menuList = [
  { key: "/admin", title: "Home", icon: House, href: "/admin" },
  { key: "/admin/alert", title: "Alert", icon: ShieldAlert, href: "/admin/alert" },
  { key: "/admin/node", title: "Node", icon: Workflow, href: "/admin/node" },
  { key: "4", title: "Network", icon: CloudCog, href: "/admin/network" },
  { key: "5", title: "Virtual", icon: GlobeLock, href: "/admin/virtual" },
  { key: "6", title: "System", icon: MonitorCog, children: [
      { key: "6_1", title: 'User', icon: User2, href: '/admin/user' }
    ] 
  },
];
</script>

<template>
    <BasicLayout :menuList="menuList"/>
</template>