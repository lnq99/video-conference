<template>
  <main font-sans p="x-4 y-4" text="center gray-700 dark:gray-200">
    <router-view/>
    <Footer/>
  </main>
</template>

<style scoped>
main {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style>
