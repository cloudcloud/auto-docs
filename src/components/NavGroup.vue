<template>
  <v-list-group v-model="item.model" :key="item.text" append-icon="">
    <template slot="activator">
      <NavItem :item="item" />
    </template>

    <template v-for="child in item.children">
      <NavGroup v-if="child.children" :key="child.text" :item="child" />
      <NavItem v-else :key="child.text" :item="child" />
    </template>
  </v-list-group>
</template>

<script>
import NavGroup from './NavGroup.vue';
import NavItem from './NavItem.vue';

export default {
  props: {
    item: {
      type: Object
    },
  },
  components: {
    NavGroup,
    NavItem,
  },
  name: 'NavGroup',
};
</script>
