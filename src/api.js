import axios from 'axios';

const client = axios.create({
  baseURL: '',
  json: true,
});

const apiClient = {
  getPages() {
    return this.perform('get', '/_api/pages');
  },

  getPage(page) {
    return this.perform('get', `/_api/page${page}`);
  },

  async perform(method, resource, data) {
    return client({
      method,
      url: resource,
      data,
      headers: {
        'X-Client': 'Ahoy-hoy'
      }
    }).then(req => {
      return req.data;
    });
  }
};

export default apiClient;
