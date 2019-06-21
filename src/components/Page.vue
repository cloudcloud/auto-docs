<template>
  <v-container fluid fill-height>
    <v-layout>
      <div v-html="page.content"></div>
    </v-layout>
  </v-container>
</template>

<script>
import { mapActions, mapMutations, mapGetters } from 'vuex';
export default {
  data: () => ({
    path: "",
    page: "",
  }),
  watch: {
    '$route': 'setPath',
  },
  mounted() {
    this.setPath();
  },
  created() {
    this.setPath();
  },
  methods: {
    setPath() {
      // first set the path from the url
      this.path = this.$route.path;
      // load up the content
      this.loadPage();
    },
    loadPage() {
      this.$store.dispatch('getPage', this.path).then(() => {
        this.page = this.$store.getters.currentPage;
      });
    },
    ...mapMutations(['resetPage']),
    ...mapActions(['getPage']),
  },
  computed: {
    ...mapGetters(['currentPage']),
  },
};
</script>

<style>
table {
  padding-top: 10px;
}

th {
  background-color: #454545;
  padding: 5px 0px;
}

td {
  padding: 5px 10px 0px 5px;
}

code {
  background-color: transparent;
  border-radius: 0px;
  padding: 2px;
  margin: 0px;
  box-shadow: 0px 0px 0px 0px;
}

code::before {
  content: "";
}

h2 {
  padding-top: 10px;
  margin-bottom: 10px;
  border-bottom: 1px dashed #404040;
}

pre {
  padding-top: 0px;
  margin-bottom: 5px;
  margin-top: 0px;
  color: #f4303a;
}

strong {
  color: #41bd47;
}
</style>
