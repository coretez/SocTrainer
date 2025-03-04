<template>
  <div>
    <MenuBar v-if="showMenu" />
    <router-view />
  </div>
</template>

<script>
import MenuBar from './components/MenuBar.vue';
import { useRoute } from 'vue-router';
import { computed } from 'vue';

export default {
  components: {
    MenuBar
  },
  setup() {
    const route = useRoute();
    const showMenu = computed(() => route.path !== '/login');

    return { showMenu };
  }
};
</script>