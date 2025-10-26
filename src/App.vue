<script setup lang="ts">
import MainPage from "./components/MainPage.vue";
</script>

<template>
    <div id="app">
        <MainPage />
    </div>
</template>

<style scoped>
#app {
    min-height: 100vh;
}
</style>

<style>
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family:
        -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen",
        "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
        sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    background-color: #f5f5f5;
}
</style>
