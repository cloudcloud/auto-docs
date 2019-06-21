import Vue from 'vue';
import Vuex from 'vuex';

import apiClient from './api';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    pages: {},
    page: "",
  },
  mutations: {
    resetPages (state, pages) {
      state.pages = pages;
    },
    resetPage (state, page) {
      state.page = page;
    },
  },
  getters: {
    allPages: state => {
      return state.pages;
    },
    currentPage: state => {
      return state.page;
    },
  },
  actions: {
    getPages({commit}) {
      return new Promise((resolve) => {
        apiClient.getPages().then((data) => {
          commit('resetPages', data.pages);
          resolve();
        });
      });
    },
    getPage({commit}, page) {
      return new Promise((resolve) => {
        apiClient.getPage(page).then((data) => {
          commit('resetPage', data);
          resolve();
        });
      });
    },
  },
});
