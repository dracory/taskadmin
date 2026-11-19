const { createApp } = Vue;

createApp({
  data() {
    return {
      stats: {
        queued_count: 0,
        running_count: 0,
        failed_count: 0,
        definition_count: 0,
        schedule_count: 0
      },
      urlTaskQueue: urlTaskQueue,
      urlTaskDefinitions: urlTaskDefinitions,
      urlSchedules: urlSchedules
    };
  },
  mounted() {
    this.loadStats();
  },
  methods: {
    async loadStats() {
      try {
        const response = await fetch(urlLoadStats);
        const result = await response.json();
        if (result.status === 'success') {
          this.stats = result.data;
        }
      } catch (error) {
        console.error('Failed to load stats:', error);
      }
    }
  }
}).mount('#app');
