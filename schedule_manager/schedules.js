function initScheduleManager() {
  const { createApp } = Vue;

  const MAX_DATETIME = '9999-12-31';
  const NULL_DATETIME = '0002-01-01';

  function emptyForm() {
    return {
      id: '',
      name: '',
      description: '',
      status: 'draft',
      queue_name: 'default',
      task_definition_id: '',
      start_at: '',
      end_at: '',
      max_execution_count: 0,
      frequency: 'none',
      interval: 1,
      days_of_week: [],
      days_of_month: [],
      months_of_year: []
    };
  }

  // "2006-01-02 15:04:05" -> "2006-01-02T15:04"; sentinel dates -> ""
  function toDatetimeLocal(value) {
    if (!value) return '';
    if (value.indexOf(MAX_DATETIME) === 0 || value.indexOf(NULL_DATETIME) === 0) return '';
    return value.replace(' ', 'T').slice(0, 16);
  }

  createApp({
    data() {
      return {
        urls: window.scheduleUrls || {},
        schedules: [],
        taskDefinitions: [],
        loading: false,
        saving: false,
        showEditModal: false,
        editForm: emptyForm(),
        taskParametersInput: '',
        daysOfMonthInput: '',
        filterName: '',
        filterStatus: '',
        filterDefinitionId: '',
        filterQueueName: '',
        currentPage: 0,
        perPage: 10,
        totalSchedules: 0,
        daysOfWeekOptions: [
          { value: 'monday', label: 'Mon' },
          { value: 'tuesday', label: 'Tue' },
          { value: 'wednesday', label: 'Wed' },
          { value: 'thursday', label: 'Thu' },
          { value: 'friday', label: 'Fri' },
          { value: 'saturday', label: 'Sat' },
          { value: 'sunday', label: 'Sun' }
        ],
        monthsOfYearOptions: [
          { value: 'JANUARY', label: 'January' },
          { value: 'FEBRUARY', label: 'February' },
          { value: 'MARCH', label: 'March' },
          { value: 'APRIL', label: 'April' },
          { value: 'MAY', label: 'May' },
          { value: 'JUNE', label: 'June' },
          { value: 'JULY', label: 'July' },
          { value: 'AUGUST', label: 'August' },
          { value: 'SEPTEMBER', label: 'September' },
          { value: 'OCTOBER', label: 'October' },
          { value: 'NOVEMBER', label: 'November' },
          { value: 'DECEMBER', label: 'December' }
        ]
      };
    },
    computed: {
      totalPages() {
        return Math.ceil(this.totalSchedules / this.perPage);
      },
      pageNumbers() {
        const pages = [];
        const start = Math.max(0, this.currentPage - 2);
        const end = Math.min(this.totalPages, start + 5);
        for (let i = start; i < end; i++) {
          pages.push(i);
        }
        return pages;
      },
      hasActiveFilters() {
        return this.filterName !== '' || this.filterStatus !== '' ||
          this.filterDefinitionId !== '' || this.filterQueueName !== '';
      }
    },
    mounted() {
      this.loadDefinitions();
      this.loadSchedules();
    },
    methods: {
      async loadSchedules() {
        this.loading = true;
        try {
          const response = await fetch(this.urls.loadSchedules, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              page: this.currentPage,
              per_page: this.perPage,
              status: this.filterStatus,
              name: this.filterName,
              queue_name: this.filterQueueName,
              task_definition_id: this.filterDefinitionId
            })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.schedules = data.data.schedules || [];
            this.totalSchedules = data.data.total || 0;
          } else {
            this.notifyError(data.message || 'Failed to load schedules');
          }
        } catch (e) {
          this.notifyError('Failed to load schedules: ' + e.message);
        } finally {
          this.loading = false;
        }
      },
      async loadDefinitions() {
        try {
          const response = await fetch(this.urls.loadDefinitions, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: '{}'
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.taskDefinitions = data.data.definitions || [];
          }
        } catch (e) {
          console.error('Failed to load task definitions:', e);
        }
      },
      applyFilters() {
        this.currentPage = 0;
        this.loadSchedules();
      },
      clearFilters() {
        this.filterName = '';
        this.filterStatus = '';
        this.filterDefinitionId = '';
        this.filterQueueName = '';
        this.currentPage = 0;
        this.loadSchedules();
      },
      openCreateModal() {
        this.editForm = emptyForm();
        this.taskParametersInput = '';
        this.daysOfMonthInput = '';
        this.showEditModal = true;
      },
      async openEditModal(schedule) {
        try {
          const response = await fetch(this.urls.scheduleDetails, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: schedule.id })
          });
          const data = await response.json();
          if (data.status === 'success') {
            const s = data.data;
            const rule = s.recurrence_rule || {};
            this.editForm = {
              id: s.id,
              name: s.name || '',
              description: s.description || '',
              status: s.status || 'draft',
              queue_name: s.queue_name || 'default',
              task_definition_id: s.task_definition_id || '',
              start_at: toDatetimeLocal(s.start_at || rule.startsAt),
              end_at: toDatetimeLocal(s.end_at || rule.endsAt),
              max_execution_count: s.max_execution_count || 0,
              frequency: rule.frequency || 'none',
              interval: rule.interval || 1,
              days_of_week: rule.daysOfWeek || [],
              days_of_month: rule.daysOfMonth || [],
              months_of_year: rule.monthsOfYear || []
            };
            this.taskParametersInput = s.task_parameters && Object.keys(s.task_parameters).length > 0
              ? JSON.stringify(s.task_parameters, null, 2)
              : '';
            this.daysOfMonthInput = (this.editForm.days_of_month || []).join(', ');
            this.showEditModal = true;
          } else {
            this.notifyError(data.message || 'Failed to load schedule');
          }
        } catch (e) {
          this.notifyError('Failed to load schedule: ' + e.message);
        }
      },
      async saveSchedule() {
        if (!this.editForm.name) {
          this.notifyError('Name is required');
          return;
        }
        if (!this.editForm.task_definition_id) {
          this.notifyError('Task definition is required');
          return;
        }

        let taskParameters = {};
        if (this.taskParametersInput.trim() !== '') {
          try {
            taskParameters = JSON.parse(this.taskParametersInput);
          } catch (e) {
            this.notifyError('Task parameters must be valid JSON');
            return;
          }
          if (typeof taskParameters !== 'object' || taskParameters === null || Array.isArray(taskParameters)) {
            this.notifyError('Task parameters must be a JSON object');
            return;
          }
        }

        const daysOfMonth = this.daysOfMonthInput
          .split(',')
          .map(v => parseInt(v.trim(), 10))
          .filter(v => !isNaN(v) && v >= 1 && v <= 31);

        this.saving = true;
        try {
          const payload = {
            id: this.editForm.id,
            name: this.editForm.name,
            description: this.editForm.description,
            status: this.editForm.status,
            queue_name: this.editForm.queue_name,
            task_definition_id: this.editForm.task_definition_id,
            task_parameters: taskParameters,
            start_at: this.editForm.start_at,
            end_at: this.editForm.end_at,
            max_execution_count: this.editForm.max_execution_count || 0,
            frequency: this.editForm.frequency,
            interval: this.editForm.interval || 1,
            days_of_week: this.editForm.days_of_week,
            days_of_month: daysOfMonth,
            months_of_year: this.editForm.months_of_year
          };
          const url = this.editForm.id ? this.urls.scheduleUpdate : this.urls.scheduleCreate;
          const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.notifySuccess(this.editForm.id ? 'Schedule updated' : 'Schedule created');
            this.showEditModal = false;
            this.loadSchedules();
          } else {
            this.notifyError(data.message || 'Failed to save schedule');
          }
        } catch (e) {
          this.notifyError('Failed to save schedule: ' + e.message);
        } finally {
          this.saving = false;
        }
      },
      async deleteSchedule(schedule) {
        const confirmed = await this.confirmDialog('Delete schedule "' + schedule.name + '"?');
        if (!confirmed) return;
        try {
          const response = await fetch(this.urls.scheduleDelete, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: schedule.id })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.notifySuccess('Schedule deleted');
            this.loadSchedules();
          } else {
            this.notifyError(data.message || 'Failed to delete schedule');
          }
        } catch (e) {
          this.notifyError('Failed to delete schedule: ' + e.message);
        }
      },
      goToPage(page) {
        if (page < 0 || page >= this.totalPages) return;
        this.currentPage = page;
        this.loadSchedules();
      },
      changePerPage() {
        this.currentPage = 0;
        this.loadSchedules();
      },
      recurrenceLabel(schedule) {
        const rule = schedule.recurrence_rule || {};
        const freq = rule.frequency || 'none';
        if (freq === 'none') return 'once';
        const interval = rule.interval > 1 ? 'every ' + rule.interval + ' ' : '';
        const unit = rule.interval > 1 ? freq.replace('ly', 's') : freq.replace('ly', '');
        return interval + unit;
      },
      statusBadgeClass(status) {
        const map = {
          active: 'bg-success',
          draft: 'bg-warning text-dark',
          inactive: 'bg-secondary',
          completed: 'bg-dark'
        };
        return map[status] || 'bg-secondary';
      },
      formatDate(value) {
        if (!value) return '';
        if (value.indexOf(MAX_DATETIME) === 0) return 'never';
        if (value.indexOf(NULL_DATETIME) === 0) return '—';
        const d = new Date(value.replace(' ', 'T') + 'Z');
        if (isNaN(d.getTime())) return value;
        return d.toLocaleString();
      },
      truncate(value, length) {
        if (!value) return '';
        return value.length > length ? value.slice(0, length) + '...' : value;
      },
      notifySuccess(message) {
        if (window.Notiflix && Notiflix.Notify) {
          Notiflix.Notify.success(message);
        } else {
          alert(message);
        }
      },
      notifyError(message) {
        if (window.Notiflix && Notiflix.Notify) {
          Notiflix.Notify.failure(message);
        } else {
          alert(message);
        }
      },
      confirmDialog(message) {
        if (window.Notiflix && Notiflix.Confirm) {
          return new Promise(resolve => {
            Notiflix.Confirm.show('Confirm', message, 'Yes', 'No',
              () => resolve(true),
              () => resolve(false)
            );
          });
        }
        return Promise.resolve(confirm(message));
      }
    }
  }).mount('#app');
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initScheduleManager);
} else {
  initScheduleManager();
}
