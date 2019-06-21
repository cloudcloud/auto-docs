<template>
  <v-list-tile :to="item.path">
    <v-list-tile-action v-if="item.icon" :prepend-icon="item.model ? item.icon : item['icon-alt']" append-icon="">
      <v-icon v-text="item.model ? item.icon : item['icon-alt']" />
    </v-list-tile-action>

    <v-list-tile-content>
      <v-list-tile-title>{{ item.text }}</v-list-tile-title>
    </v-list-tile-content>
    <slot />
  </v-list-tile>
</template>

<script>
export default {
  props: {
    item: {
      type: Object
    },
  },
  name: 'NavItem',
};
</script>

<style>
.v-list--dense .v-list__tile:not(.v-list__tile--avatar) {
  height: 30px;
}
</style>
