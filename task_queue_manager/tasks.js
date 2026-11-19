function initTaskQueueManager() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        urls: window.taskQueueUrls || {},
        definitions: window.taskDefinitions || [],
        tasks: [],
        selectedTasks: [],
        selectAll: false,
        loading: false,
        error: null,
        showFilterModal: false,
        showEnqueueModal: false,
        showRequeueModal: false,
        showParametersModal: false,
        showDetailsModal: false,
        enqueuing: false,
        requeuing: false,
        savingParameters: false,
        enqueueForm: {
          task_definition_id: '',
          parameters: '{}'
        },
        requeueForm: {
          id: '',
          parameters: '{}'
        },
        parametersForm: {
          id: '',
          parameters: '{}'
        },
        details: null,
        currentPage: 0,
        perPage: 10,
        totalTasks: 0,
        conditions: [],
        modalConditions: [],
        sortByColumn: 'created_at',
        sortOrder: 'desc',
        conditionOptions: [
          { value: 'status', label: 'Status', operators: ['equals'], inputType: 'select', placeholder: '', options: [
            { value: 'queued', label: 'Queued' },
            { value: 'running', label: 'Running' },
            { value: 'success', label: 'Success' },
            { value: 'failed', label: 'Failed' },
            { value: 'paused', label: 'Paused' },
            { value: 'canceled', label: 'Canceled' },
            { value: 'deleted', label: 'Deleted' },
          ] },
          { value: 'id', label: 'ID', operators: ['equals'], inputType: 'text', placeholder: 'queue id' },
          { value: 'task_id', label: 'Task ID', operators: ['equals'], inputType: 'text', placeholder: 'task definition id' },
          { value: 'task_name', label: 'Task Name', operators: ['contains'], inputType: 'text', placeholder: 'task name or alias' },
          { value: 'date_from', label: 'Created after', operators: ['equals'], inputType: 'date', placeholder: '' },
          { value: 'date_to', label: 'Created before', operators: ['equals'], inputType: 'date', placeholder: '' },
        ],
        opLabels: { equals: '=', contains: 'contains' }
      };
    },
    computed: {
      totalPages() {
        return Math.ceil(this.totalTasks / this.perPage);
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
      filterStatus() {
        if (this.conditions.length === 0) return 'Showing all tasks';
        return this.conditions.length + ' filter(s) active';
      },
      hasActiveFilters() {
        return this.conditions.length > 0;
      }
    },
    mounted() {
      const urlParams = new URLSearchParams(window.location.search);

      const filtersRaw = urlParams.get('filters');
      if (filtersRaw) {
        try {
          const parsed = JSON.parse(filtersRaw);
          if (Array.isArray(parsed)) {
            this.conditions = parsed.filter(c => c && c.field && c.value);
          }
        } catch (e) { /* ignore malformed */ }
      } else {
        if (urlParams.get('status')) this.conditions.push({ field: 'status', operator: 'equals', value: urlParams.get('status') });
        if (urlParams.get('created_from')) this.conditions.push({ field: 'date_from', operator: 'equals', value: urlParams.get('created_from') });
        if (urlParams.get('created_to')) this.conditions.push({ field: 'date_to', operator: 'equals', value: urlParams.get('created_to') });
      }

      this.sortByColumn = urlParams.get('sort_by') || 'created_at';
      this.sortOrder = urlParams.get('sort_order') || 'desc';
      this.currentPage = parseInt(urlParams.get('page') || '0', 10);
      this.perPage = parseInt(urlParams.get('per_page') || '10', 10);
      this.loadTasks();
    },
    methods: {
      async loadTasks() {
        this.loading = true;
        try {
          const response = await fetch(this.urls.loadTasks, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              page: this.currentPage,
              per_page: this.perPage,
              sort_by: this.sortByColumn,
              sort: this.sortOrder,
              conditions: this.conditions
            })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.tasks = data.data.tasks || [];
            this.totalTasks = data.data.total || 0;
          } else {
            this.notifyError(data.message || 'Failed to load tasks');
          }
        } catch (e) {
          this.notifyError('Failed to load tasks: ' + e.message);
        } finally {
          this.loading = false;
        }
      },
      async enqueueTask() {
        if (!this.enqueueForm.task_definition_id) {
          this.notifyError('Please select a task');
          return;
        }
        if (!this.isJSON(this.enqueueForm.parameters)) {
          this.notifyError('Parameters must be valid JSON');
          return;
        }
        this.enqueuing = true;
        try {
          const response = await fetch(this.urls.taskEnqueue, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(this.enqueueForm)
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.notifySuccess('Task enqueued');
            this.showEnqueueModal = false;
            this.enqueueForm = { task_definition_id: '', parameters: '{}' };
            this.loadTasks();
          } else {
            this.notifyError(data.message || 'Failed to enqueue task');
          }
        } catch (e) {
          this.notifyError('Failed to enqueue task: ' + e.message);
        } finally {
          this.enqueuing = false;
        }
      },
      async requeueTask() {
        if (!this.isJSON(this.requeueForm.parameters)) {
          this.notifyError('Parameters must be valid JSON');
          return;
        }
        this.requeuing = true;
        try {
          const response = await fetch(this.urls.taskRequeue, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(this.requeueForm)
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.notifySuccess('Task requeued');
            this.showRequeueModal = false;
            this.loadTasks();
          } else {
            this.notifyError(data.message || 'Failed to requeue task');
          }
        } catch (e) {
          this.notifyError('Failed to requeue task: ' + e.message);
        } finally {
          this.requeuing = false;
        }
      },
      async restartTask(task) {
        const confirmed = await this.confirmDialog('Restart this task?');
        if (!confirmed) return;
        try {
          const response = await fetch(this.urls.taskRestart, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: task.id })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.notifySuccess('Task restarted');
            this.loadTasks();
          } else {
            this.notifyError(data.message || 'Failed to restart task');
          }
        } catch (e) {
          this.notifyError('Failed to restart task: ' + e.message);
        }
      },
      async deleteTask(task) {
        const confirmed = await this.confirmDialog('Delete this queued task?');
        if (!confirmed) return;
        try {
          const response = await fetch(this.urls.taskDelete, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: task.id })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.notifySuccess('Task deleted');
            this.loadTasks();
          } else {
            this.notifyError(data.message || 'Failed to delete task');
          }
        } catch (e) {
          this.notifyError('Failed to delete task: ' + e.message);
        }
      },
      async deleteSelectedTasks() {
        const confirmed = await this.confirmDialog('Delete ' + this.selectedTasks.length + ' selected task(s)?');
        if (!confirmed) return;
        try {
          const response = await fetch(this.urls.taskDeleteSelected, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ids: this.selectedTasks })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.notifySuccess('Tasks deleted');
            this.selectedTasks = [];
            this.selectAll = false;
            this.loadTasks();
          } else {
            this.notifyError(data.message || 'Failed to delete tasks');
          }
        } catch (e) {
          this.notifyError('Failed to delete tasks: ' + e.message);
        }
      },
      async updateParameters() {
        if (!this.isJSON(this.parametersForm.parameters)) {
          this.notifyError('Parameters must be valid JSON');
          return;
        }
        this.savingParameters = true;
        try {
          const response = await fetch(this.urls.taskUpdateParams, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(this.parametersForm)
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.notifySuccess('Parameters updated');
            this.showParametersModal = false;
            this.loadTasks();
          } else {
            this.notifyError(data.message || 'Failed to update parameters');
          }
        } catch (e) {
          this.notifyError('Failed to update parameters: ' + e.message);
        } finally {
          this.savingParameters = false;
        }
      },
      async openDetailsModal(task) {
        try {
          const response = await fetch(this.urls.taskDetails, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: task.id })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.details = data.data;
            this.showDetailsModal = true;
          } else {
            this.notifyError(data.message || 'Failed to load task details');
          }
        } catch (e) {
          this.notifyError('Failed to load task details: ' + e.message);
        }
      },
      async openParametersModal(task) {
        try {
          const response = await fetch(this.urls.taskDetails, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: task.id })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.parametersForm = { id: task.id, parameters: data.data.parameters || '{}' };
            this.showParametersModal = true;
          } else {
            this.notifyError(data.message || 'Failed to load parameters');
          }
        } catch (e) {
          this.notifyError('Failed to load parameters: ' + e.message);
        }
      },
      async openRequeueModal(task) {
        try {
          const response = await fetch(this.urls.taskDetails, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: task.id })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.requeueForm = { id: task.id, parameters: data.data.parameters || '{}' };
            this.showRequeueModal = true;
          } else {
            this.notifyError(data.message || 'Failed to load parameters');
          }
        } catch (e) {
          this.notifyError('Failed to load parameters: ' + e.message);
        }
      },
      openFilterModal() {
        this.modalConditions = JSON.parse(JSON.stringify(this.conditions));
        if (this.modalConditions.length === 0) {
          this.addModalCondition();
        }
        this.showFilterModal = true;
      },
      closeFilterModal() {
        this.showFilterModal = false;
      },
      addModalCondition() {
        this.modalConditions.push({ field: 'status', operator: 'equals', value: '' });
      },
      removeModalCondition(index) {
        this.modalConditions.splice(index, 1);
      },
      clearModalConditions() {
        this.modalConditions = [];
      },
      applyFilters() {
        this.conditions = this.modalConditions.filter(c => c.field && c.value);
        this.showFilterModal = false;
        this.currentPage = 0;
        this.updateUrl();
        this.loadTasks();
      },
      clearFilters() {
        this.conditions = [];
        this.currentPage = 0;
        this.updateUrl();
        this.loadTasks();
      },
      sortBy(column) {
        if (this.sortByColumn === column) {
          this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
        } else {
          this.sortByColumn = column;
          this.sortOrder = 'asc';
        }
        this.updateUrl();
        this.loadTasks();
      },
      goToPage(page) {
        if (page < 0 || page >= this.totalPages) return;
        this.currentPage = page;
        this.updateUrl();
        this.loadTasks();
      },
      changePerPage() {
        this.currentPage = 0;
        this.updateUrl();
        this.loadTasks();
      },
      toggleSelectAll() {
        if (this.selectAll) {
          this.selectedTasks = this.tasks.map(t => t.id);
        } else {
          this.selectedTasks = [];
        }
      },
      updateUrl() {
        const params = new URLSearchParams(window.location.search);
        params.set('controller', 'task-queue-manager');
        if (this.conditions.length > 0) {
          params.set('filters', JSON.stringify(this.conditions));
        } else {
          params.delete('filters');
        }
        params.set('sort_by', this.sortByColumn);
        params.set('sort_order', this.sortOrder);
        params.set('page', this.currentPage);
        params.set('per_page', this.perPage);
        const newUrl = window.location.pathname + '?' + params.toString();
        window.history.replaceState({}, '', newUrl);
      },
      operatorsFor(field) {
        const opt = this.conditionOptions.find(o => o.value === field);
        return opt ? opt.operators : ['equals'];
      },
      inputTypeFor(field) {
        const opt = this.conditionOptions.find(o => o.value === field);
        return opt ? opt.inputType : 'text';
      },
      placeholderFor(field) {
        const opt = this.conditionOptions.find(o => o.value === field);
        return opt ? opt.placeholder : '';
      },
      optionsFor(field) {
        const opt = this.conditionOptions.find(o => o.value === field);
        return opt && opt.options ? opt.options : [];
      },
      statusBadgeClass(status) {
        const map = {
          queued: 'bg-primary',
          running: 'bg-info text-dark',
          success: 'bg-success',
          failed: 'bg-danger',
          paused: 'bg-warning text-dark',
          canceled: 'bg-secondary',
          deleted: 'bg-dark'
        };
        return map[status] || 'bg-secondary';
      },
      formatDate(value) {
        if (!value) return '';
        const d = new Date(value);
        if (isNaN(d.getTime())) return value;
        return d.toLocaleString();
      },
      prettyJSON(value) {
        try {
          return JSON.stringify(JSON.parse(value), null, 2);
        } catch (e) {
          return value;
        }
      },
      isJSON(value) {
        try {
          JSON.parse(value);
          return true;
        } catch (e) {
          return false;
        }
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
  document.addEventListener('DOMContentLoaded', initTaskQueueManager);
} else {
  initTaskQueueManager();
}
