<template>
  <a-layout class="layout-demo">
    <a-layout-sider
      collapsible
      breakpoint="xl"
      style="margin: 10px; border-radius: 5px"
    >
      <div class="logo" />
      <a-menu
        :selected-keys="[route.path]"
        :style="{ width: '100%' }"
        @menuItemClick="onClickMenuItem"
      >
        <template v-for="menu in menuList" :key="menu.key">
          <template v-if="menu.children && menu.children.length">
            <a-sub-menu>
              <template #icon>
                <component :is="menu.icon as Component" :size="18" />
              </template>
              <template #title>{{ menu.title }}</template>

              <template v-for="subMenu in menu.children" :key="subMenu.key">
                <a-menu-item :path="subMenu.href" >
                  <template #icon>
                    <component
                      :is="subMenu.icon as Component"
                      :size="18"
                      v-if="subMenu.icon"
                    />
                  </template>
                  {{ subMenu.title }}
                </a-menu-item>
              </template>
            </a-sub-menu>
          </template>
          <template v-else>
            <a-menu-item :key="menu.href" :path="menu.href">
              <template #icon>
                <component :is="menu.icon as Component" :size="18" />
              </template>
              {{ menu.title }}
            </a-menu-item>
          </template>
        </template>
      </a-menu>
      <!-- trigger -->
      <template #trigger="{ collapsed }">
        <IconCaretRight v-if="collapsed"></IconCaretRight>
        <IconCaretLeft v-else></IconCaretLeft>
      </template>
    </a-layout-sider>
    <a-layout class="container">
      <a-layout-header>
        <a-breadcrumb :style="{ margin: '16px 12px' }">
          <a-breadcrumb-item :style="{ display: 'flex', alignItems: 'center' }">
            <House :size="15" />
          </a-breadcrumb-item>
          <a-breadcrumb-item>{{ $route.name }}</a-breadcrumb-item>
        </a-breadcrumb>
        <div class="user-center">
          <div :style="{ color: 'var(--color-text-2)', marginRight: '10px' }">
            Welcome, {{ user.name }}
          </div>
          <div :style="{ color: 'var(--color-text-2)' }">|</div>
          <a-button type="text" @click="handleLogout" status="danger">
            <template #icon>
              <icon-export />
            </template>
            <template #default>Logout</template>
          </a-button>
          <a-modal
            :visible="visible"
            @ok="handleOk"
            @cancel="handleCancel"
            okText="确定"
            unmountOnClose
          >
            <template #title> Logout </template>
            <div>该操作不可逆, 你确定要退出吗?</div>
          </a-modal>
        </div>
      </a-layout-header>
      <a-layout style="margin: 10px; overflow: auto">
        <a-layout-content style="background: inherit">
          <RouterView />
        </a-layout-content>
      </a-layout>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import router from "@/router";
import { Message } from "@arco-design/web-vue";
import {
  IconCaretRight,
  IconCaretLeft,
  IconExport,
} from "@arco-design/web-vue/es/icon";
import { House } from "lucide-vue-next";
import { ref, type Component } from "vue";
import {useRoute} from "vue-router";

const route = useRoute();

// user center
const visible = ref(false);
const user = ref({
  name: "admin",
});
const handleOk = () => {
  visible.value = false;
  router.push({ path: "/" });
};
const handleCancel = () => {
  visible.value = false;
};
const handleLogout = () => {
  visible.value = true;
};

// 处理菜单项点击 - 只处理叶子节点
const onClickMenuItem = (key) => {
  console.log(key)
  // console.log(route.path)
  // 找到对应的菜单配置
  const findHref = (items: MenuItem[]): string | undefined => {
    for (const item of items) {
      if (item.title == key) return item.href;
      if (item.children) {
        const found = findHref(item.children);
        if (found) return found;
      }
    }
    return undefined;
  };

  const href = findHref(props.menuList);
  if (href) {
    router.push({ path: href });
  }
};

interface MenuItem {
  key: string;
  title: string;
  icon?: Component;
  href?: string;
  children?: MenuItem[];
}

const props = defineProps<{
  menuList: MenuItem[];
}>();
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
.user-center {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-right: 10px;
}
</style>
