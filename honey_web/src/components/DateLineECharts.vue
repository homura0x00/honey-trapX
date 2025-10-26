<script setup lang="ts">
import * as echarts from "echarts/core";
import { GridComponent, type GridComponentOption } from "echarts/components";
import { LineChart, type LineSeriesOption } from "echarts/charts";
import { UniversalTransition } from "echarts/features";
import { CanvasRenderer } from "echarts/renderers";
import { onMounted } from "vue";

echarts.use([GridComponent, LineChart, CanvasRenderer, UniversalTransition]);

type EChartsOption = echarts.ComposeOption<
  GridComponentOption | LineSeriesOption
>;

onMounted(() => {
  var chartDom = document.getElementById("date-line-echart")!;
  var myChart = echarts.init(chartDom, "dark");
  var option: EChartsOption;

  
  option = {
    color: "#64ffda",
    backgroundColor: "transparent",
    grid: {
      left: "7%",
      right: "7%",
      top: "12%",
      bottom: "12%",
    },
    xAxis: {
      type: "category",
      data: ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"],
      axisLine: {
        lineStyle: {
          color: "white",
        },
      },
      axisLabel: {
        color: "white",
      },
    },
    yAxis: {
      type: "value",
      splitLine: {
        show: false,
      },
    },
    series: [
      {
        data: [820, 932, 901, 934, 1290, 1330, 1320],
        type: "line",
        smooth: true,
      },
    ],
  };

  option && myChart.setOption(option);
});
</script>

<template>
  <div id="date-line-echart" :style="{ width: '100%', height: '300px' }"></div>
</template>

<style scoped></style>
