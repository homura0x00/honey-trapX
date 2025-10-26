<template>
  <a-card
    hoverable
    :style="{
      width: '180px',
      height: 'auto',
      borderRadius: '15px',
    }"
  >
    <div
      :style="{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-around',
      }"
    >
      <span
        :style="{ display: 'flex', alignItems: 'center', color: '#1D2129' }"
      >
        <a-avatar
          :style="{ marginRight: '8px', backgroundColor: 'transparent' }"
          :size="50"
        >
          <slot></slot>
        </a-avatar>
      </span>
      <div
        class="status"
        :style="{
          display: 'flex',
          flexDirection: 'column',
          textAlign: 'end',
        }"
      >
        <a-typography-text :style="{fontSize: '1.2rem'}"><strong>{{ props.nums }}</strong></a-typography-text>
        <a-typography-text>{{ props.name }}</a-typography-text>
      </div>
    </div>
  </a-card>
</template>

<script setup lang="ts">
interface Props {
  nums: number,
  name: string,
}
const props = defineProps<Props>()
</script>
