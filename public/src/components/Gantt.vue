<template>
  <div style="width:100%;height:100%" ref="gantt"></div>
</template>

<script>
import 'dhtmlx-gantt';

export default {
  name: 'gantt',
  props: {
    tasks: {
      type: Object,
      default() {
        return { data: [], links: [] };
      },
    },
  },

  mounted() {
    gantt.init(this.$refs.gantt);
    gantt.parse(this.$props.tasks);
  },
};
</script>

<style>
    @import "dhtmlx-gantt/codebase/dhtmlxgantt.css";
</style>