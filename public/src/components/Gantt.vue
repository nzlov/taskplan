<template>
  <v-container>
    <remote-js src="//export.dhtmlx.com/gantt/api.js"></remote-js>
    <div style="height:90%;" ref="gantt"></div>
  </v-container>
</template>

<script>
import 'dhtmlx-gantt';
import 'dhtmlx-gantt/codebase/locale/locale_cn';
import 'dhtmlx-gantt/codebase/ext/dhtmlxgantt_marker';
import 'dhtmlx-gantt/codebase/ext/dhtmlxgantt_tooltip';

export default {
  components: {
    'remote-js': {
      render(createElement) {
        return createElement('script', {
          attrs: {
            type: 'text/javascript',
            src: this.src,
          },
        });
      },
      props: {
        src: { type: String, required: true },
      },
    },
  },
  mounted() {
    gantt.config.scale_unit = 'month';
    gantt.config.date_scale = '%F, %Y';

    gantt.config.scale_height = 50;

    gantt.config.subscales = [
      { unit: 'day', step: 1, date: '%j, %D' },
    ];
    gantt.config.add_column = false;

    gantt.templates.grid_row_class = ((startdate, enddate, item) => {
      return this.getColorStyle(item);
    });
    gantt.templates.task_class = ((startdate, enddate, item) => {
      return this.getColorStyle(item);
    });
    gantt.templates.tooltip_text = function(start,end,task){
      return `<b>任务:</b>${task.text}<br/>
              <b>组:</b>${task.usergroup}<br/>
              <b>资源:</b>${task.user}<br/>
              <b>开始时间:</b>${task.start_date2}<br/>
              <b>结束时间:</b>${task.end_date2}<br/>
              <b>真正结束时间:</b>${task.real_date}<br/>
              <b>状态:</b>${task.statuss}<br/>
              `
    };
    gantt.config.columns = [
      {
        name: 'text',
        label: '任务',
        tree: true,
        width: '*',
      },
      {
        name: 'statuss',
        label: '状态',
        align: 'center',
      },
      {
        name: 'assigned',
        label: '执行人',
        align: 'center',
        width: 100,
        template: (item) => {
          if (!item.user) return '';
          return item.user;
        },
      },
    ];
    gantt.config.readonly = true;
    gantt.init(this.$refs.gantt);
  },
  methods: {
    load(data) {
      const ids = [];
      data.data.forEach((element) => {
        ids.push(element.id);
      });
      const newdata = [];
      data.data.forEach((element) => {
       if (ids.indexOf(element.parent) === -1) {
          element.parent = 0;
        }
        newdata.push(element);
      });
      data.data = newdata;
      gantt.clearAll();
      gantt.parse(data);
      const today = new Date();
      gantt.addMarker({
        start_date: today,
        css: 'today',
        text: '当前时间',
      });
      gantt.render();
    },
    exportdata() {
      gantt.exportToPNG({
        header: '<link rel="stylesheet" href="//dhtmlx.com/docs/products/dhtmlxGantt/common/customstyles.css" type="text/css">',
      });
    },
    getColorStyle(item) {
      switch (item.statuss) {
        case '计划中': {
          return 'plan';
        }
        case '已超期': {
          return 'expired';
        }
        case '临近过期': {
          return 'nearexpire';
        }
        case '进行中': {
          return 'work';
        }
        case '超期完成': {
          return 'doneexpired';
        }
        case '提前完成': {
          return 'donebefore';
        }
        case '按时完成': {
          return 'done';
        }
        case '重新打开并超期': {
          return 'reopenexpired';
        }
        case '重新打开': {
          return 'reopen';
        }
        default: {
          return '';
        }
      }
    },
  },
};
</script>

<style>
  @import "../../node_modules/dhtmlx-gantt/codebase/skins/dhtmlxgantt_material.css";
</style>

<style>
  .plan, .plan .gantt_cell, .odd.plan .gantt_cell,
  .plan .gantt_task_cell, .odd.plan .gantt_task_cell {
    background: #99CCFF;
  }
  .expired, .expired .gantt_cell, .odd.expired .gantt_cell,
  .expired .gantt_task_cell, .odd.expired .gantt_task_cell {
    background: #FF0033;
  }
  .nearexpire, .nearexpire .gantt_cell, .odd.nearexpire .gantt_cell,
  .nearexpire .gantt_task_cell, .odd.nearexpire .gantt_task_cell {
    background: #FF6600;
  }
  .work, .work .gantt_cell, .odd.work .gantt_cell,
  .work .gantt_task_cell, .odd.work .gantt_task_cell {
    background: #66CC99;
  }
  .doneexpired, .doneexpired .gantt_cell, .odd.doneexpired .gantt_cell,
  .doneexpired .gantt_task_cell, .odd.doneexpired .gantt_task_cell {
    background: #FF6666;
  }
  .donebefore, .donebefore .gantt_cell, .odd.donebefore .gantt_cell,
  .donebefore .gantt_task_cell, .odd.donebefore .gantt_task_cell {
    background: #33CC99;
  }
  .done, .done .gantt_cell, .odd.done .gantt_cell,
  .done .gantt_task_cell, .odd.done .gantt_task_cell {
    background: #009933;
  }
  .reopenexpired, .reopenexpired .gantt_cell, .odd.reopenexpired .gantt_cell,
  .reopenexpired .gantt_task_cell, .odd.reopenexpired .gantt_task_cell {
    background: #CC0066;
  }
  .reopen, .reopen .gantt_cell, .odd.reopen .gantt_cell,
  .reopen .gantt_task_cell, .odd.reopen .gantt_task_cell {
    background: #009999;
  }
</style>