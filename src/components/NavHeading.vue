<template>
  <v-layout row align-center :key="item.heading">
    <v-flex xs6>
      <v-subheader v-if="item.heading">{{ item.heading }}</v-subheader>
    </v-flex>
    <v-flex xs6 class="text-xs-center">
      <a href="#!" class="body-2 white--text">EDIT</a>
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  props: {
    item: {
      type: Object,
    },
  },
  name: 'NavHeading',
};
</script>
