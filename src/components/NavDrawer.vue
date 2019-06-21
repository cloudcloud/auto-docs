<template>
  <v-navigation-drawer fixed app v-model="drawer">
    <v-list class="pa-0" dense expand>
      <template v-for="item in items">

        <NavGroup v-if="item.children" :key="item.text" :item="item" />
        <NavItem v-else :key="item.text" :item="item" />

      </template>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
import NavGroup from './NavGroup.vue';
import NavItem from './NavItem.vue';

export default {
  data: () => ({}),
  props: {
    items: {
      type: Object
    },
    drawer: {
      type: Boolean
    },
  },
  name: 'NavDrawer',
  components: {
    NavGroup,
    NavItem,
  },
};
</script>
