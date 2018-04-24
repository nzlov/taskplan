<template>
  <v-container fluid>
    <gantt class="left-container" :tasks="tasks"></gantt>
  </v-container>
</template>

<script>
import DateUtil from '../utils/date';

export default {
  data() {
    return {
      tasks: {
        data: [
          { id: 1, text: 'Task #1', start_date: '15-04-2017', duration: 3, progress: 0.6 },
          { id: 2, text: 'Task #2', start_date: '18-04-2017', duration: 3, progress: 0.4 },
        ],
        links: [
          { id: 1, source: 1, target: 2, type: '0' },
        ]
      },
    };
  },
  created() {
    // this.load();
  },
  methods: {
    load() {
      new Promise((resolve) => {
        const items = [];
        this.$http.LGet(this.$store.state, '/task?all=t')
          .then((resp) => {
            switch (resp.data.code) {
              case 0: {
                resp.data.data.data.forEach((element) => {
                  items.push({
                    id: element.ID,
                    text: element.Name,
                    start_date: this.formatDate1(element.Start),
                    open: true,
                    parent: element.ParentTaskID,
                    duration: (element.End - element.Start) / 3600,
                  });
                });
                break;
              }
              default:
            }
            resolve({
              items,
            });
          }).catch((e) => {
            console.dir(e);
            if (e.response.data.code === 101) {
              this.$store.commit('logout');
              this.$router.replace('/login');
            }
            this.nodata = '服务器报错';
            resolve({
              items,
            });
          });
      }).then((data) => {
        this.tasks.data = data.items;
        console.dir(this.tasks);
      });
    },
    formatDate1(d) {
      return DateUtil.formatDate1(new Date(d));
    },
  },
};
</script>