<script setup lang="ts">
import BonneteauGame from "./components/BonneteauGame.vue";
</script>

<template>
  <div id="app">
    <BonneteauGame />
  </div>
</template>

<style>
#app {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Arial', sans-serif;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
</style>

