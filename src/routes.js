import Vue from 'vue';
import VueRouter from 'vue-router';

// import individual components in here
import Page from './components/Page';

Vue.use(VueRouter);

export default new VueRouter({
  mode: 'history',
  routes: [
    {path: '*', name: 'Page', component: Page}
  ]
});
