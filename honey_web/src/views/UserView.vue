<script setup lang="ts">
import BasicLayout from '@/layouts/BasicLayout.vue';
import {
  House,
  MonitorCog,
} from "lucide-vue-next";

const userMenuList = [
  { title: "Home", icon: House, href: "/dashboard" },
  { title: "Settings", icon: MonitorCog, href: "/dashboard/settings" },
];
</script>

<template>
    <BasicLayout :menuList="userMenuList"/>
</template>