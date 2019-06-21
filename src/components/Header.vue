<template>
  <NavDrawer :drawer="drawer" :items="pages" />
</template>

<script>
import { mapActions, mapMutations, mapGetters } from 'vuex';
import NavDrawer from './NavDrawer.vue';

export default {
  data: () => ({
    drawer: null,
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
