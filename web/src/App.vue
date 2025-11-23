<script setup lang="ts">

import { RouterView } from 'vue-router'

</script>

<template>
  <RouterView />
</template>

<style>
/* We rely on the global styles imported in main.ts (assets/main.css)
  or the specific styles within LoginView/DashboardView.

  Resetting margin here ensures the views take up the full screen.
*/
body {
  margin: 0;
  padding: 0;
}
</style>
