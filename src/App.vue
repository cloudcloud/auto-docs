<template>
  <v-app id="ad" dark>
    <v-toolbar fixed dark app clipped>
      <v-toolbar-side-icon @click.stop="drawer = !drawer" />
      <v-toolbar-title>auto-docs</v-toolbar-title>
    </v-toolbar>

    <NavDrawer :items="pages" :drawer="drawer" />

    <v-content>
      <router-view />
    </v-content>

    <v-footer app fixed>
      <span>&copy; auto-docs 2019.</span>
    </v-footer>
  </v-app>
</template>

<script>
import { mapActions, mapMutations, mapGetters } from 'vuex';
import NavDrawer from './components/NavDrawer.vue';

export default {
  data: () => ({
    drawer: true,
    pages: [],
  }),
  created() {
    this.$store.dispatch('getPages').then(() => {
      this.loadPages();
    });
  },
  methods: {
    loadPages() {
      this.pages = this.$store.getters.allPages;
    },
    ...mapMutations(['resetPages']),
    ...mapActions(['getPages']),
  },
  computed: {
    ...mapGetters(['allPages']),
  },
  components: {
    NavDrawer,
  },
};
</script>

<style></style>
