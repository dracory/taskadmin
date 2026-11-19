function initTaskDefinitionManager() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        urls: window.taskDefinitionUrls || {},
        definitions: [],
        loading: false,
        error: null,
        showFilterModal: false,
        showEditModal: false,
        saving: false,
        editForm: {
          id: '',
          title: '',
          alias: '',
          status: 'active',
          description: ''
        },
        currentPage: 0,
        perPage: 10,
        totalDefinitions: 0,
        conditions: [],
        modalConditions: [],
        sortByColumn: 'created_at',
        sortOrder: 'desc',
        conditionOptions: [
          { value: 'status', label: 'Status', operators: ['equals'], inputType: 'select', placeholder: '', options: [
            { value: 'active', label: 'Active' },
            { value: 'inactive', label: 'Inactive' },
            { value: 'canceled', label: 'Canceled' },
          ] },
          { value: 'id', label: 'ID', operators: ['equals'], inputType: 'text', placeholder: 'definition id' },
          { value: 'alias', label: 'Alias', operators: ['equals'], inputType: 'text', placeholder: 'task-alias' },
          { value: 'date_from', label: 'Created after', operators: ['equals'], inputType: 'date', placeholder: '' },
          { value: 'date_to', label: 'Created before', operators: ['equals'], inputType: 'date', placeholder: '' },
        ],
        opLabels: { equals: '=', contains: 'contains' }
      };
    },
    computed: {
      totalPages() {
        return Math.ceil(this.totalDefinitions / this.perPage);
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
        if (this.conditions.length === 0) return 'Showing all task definitions';
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
      this.loadDefinitions();
    },
    methods: {
      async loadDefinitions() {
        this.loading = true;
        try {
          const response = await fetch(this.urls.loadDefinitions, {
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
            this.definitions = data.data.definitions || [];
            this.totalDefinitions = data.data.total || 0;
          } else {
            this.notifyError(data.message || 'Failed to load task definitions');
          }
        } catch (e) {
          this.notifyError('Failed to load task definitions: ' + e.message);
        } finally {
          this.loading = false;
        }
      },
      openCreateModal() {
        this.editForm = { id: '', title: '', alias: '', status: 'active', description: '' };
        this.showEditModal = true;
      },
      async openEditModal(definition) {
        try {
          const response = await fetch(this.urls.definitionDetails, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: definition.id })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.editForm = {
              id: data.data.id,
              title: data.data.title,
              alias: data.data.alias,
              status: data.data.status,
              description: data.data.description || ''
            };
            this.showEditModal = true;
          } else {
            this.notifyError(data.message || 'Failed to load task definition');
          }
        } catch (e) {
          this.notifyError('Failed to load task definition: ' + e.message);
        }
      },
      async saveDefinition() {
        if (!this.editForm.title) {
          this.notifyError('Title is required');
          return;
        }
        if (!this.editForm.alias) {
          this.notifyError('Alias is required');
          return;
        }
        this.saving = true;
        try {
          const url = this.editForm.id ? this.urls.definitionUpdate : this.urls.definitionCreate;
          const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(this.editForm)
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.notifySuccess(this.editForm.id ? 'Task definition updated' : 'Task definition created');
            this.showEditModal = false;
            this.loadDefinitions();
          } else {
            this.notifyError(data.message || 'Failed to save task definition');
          }
        } catch (e) {
          this.notifyError('Failed to save task definition: ' + e.message);
        } finally {
          this.saving = false;
        }
      },
      async deleteDefinition(definition) {
        const confirmed = await this.confirmDialog('Delete task definition "' + definition.title + '"?');
        if (!confirmed) return;
        try {
          const response = await fetch(this.urls.definitionDelete, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: definition.id })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.notifySuccess('Task definition deleted');
            this.loadDefinitions();
          } else {
            this.notifyError(data.message || 'Failed to delete task definition');
          }
        } catch (e) {
          this.notifyError('Failed to delete task definition: ' + e.message);
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
        this.loadDefinitions();
      },
      clearFilters() {
        this.conditions = [];
        this.currentPage = 0;
        this.updateUrl();
        this.loadDefinitions();
      },
      sortBy(column) {
        if (this.sortByColumn === column) {
          this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
        } else {
          this.sortByColumn = column;
          this.sortOrder = 'asc';
        }
        this.updateUrl();
        this.loadDefinitions();
      },
      goToPage(page) {
        if (page < 0 || page >= this.totalPages) return;
        this.currentPage = page;
        this.updateUrl();
        this.loadDefinitions();
      },
      changePerPage() {
        this.currentPage = 0;
        this.updateUrl();
        this.loadDefinitions();
      },
      updateUrl() {
        const params = new URLSearchParams(window.location.search);
        params.set('controller', 'task-definition-manager');
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
          active: 'bg-success',
          inactive: 'bg-secondary',
          canceled: 'bg-dark'
        };
        return map[status] || 'bg-secondary';
      },
      formatDate(value) {
        if (!value) return '';
        const d = new Date(value);
        if (isNaN(d.getTime())) return value;
        return d.toLocaleString();
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
  document.addEventListener('DOMContentLoaded', initTaskDefinitionManager);
} else {
  initTaskDefinitionManager();
}
