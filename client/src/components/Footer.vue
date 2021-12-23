<script lang="ts" setup>
import {isDark, toggleDark} from '~/composables'

const router = useRouter()
</script>

<template>
  <nav gap-6 inline-flex mt-4 text-xl>
    <button class="icon-btn !outline-none" @click="router.back()">
      <div i-carbon-undo/>
    </button>
    <button class="icon-btn !outline-none" @click="router.push('/')">
      <div i-carbon-home/>
    </button>
    <button class="icon-btn !outline-none" @click="toggleDark()">
      <div v-if="isDark" i-carbon-moon/>
      <div v-else i-carbon-sun/>
    </button>
  </nav>
</template>
