<script setup lang="ts">
import { Search } from 'lucide-vue-next';

const searchQuery = defineModel<string>('searchQuery', { default: '' });
const selectedCategory = defineModel<string>('selectedCategory', { default: 'ALL' });
const onlyWithCorrection = defineModel<boolean>('onlyWithCorrection', { default: false });

const availableCategories = ['B', 'A', 'A1', 'C', 'D', 'E'];
</script>

<template>
  <div class="p-3 rounded-[4px] bg-muted/30 border border-border flex flex-col md:flex-row md:items-center justify-between gap-3">
    <!-- Recherche textuelle -->
    <div class="relative flex-1 max-w-md">
      <Search class="w-3.5 h-3.5 text-muted-foreground absolute left-3 top-1/2 -translate-y-1/2" />
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Rechercher par titre, mot-clé ou session..."
        class="w-full h-8 pl-8 pr-3 text-xs bg-card border border-border rounded-[4px] text-foreground focus:outline-none focus:border-brand-dark"
      />
    </div>

    <!-- Filtres Catégories de Permis (Pills) -->
    <div class="flex items-center space-x-1.5 overflow-x-auto text-xs pb-1 md:pb-0">
      <button
        type="button"
        @click="selectedCategory = 'ALL'"
        class="px-2.5 py-1 rounded-[4px] border font-semibold transition cursor-pointer"
        :class="selectedCategory === 'ALL'
          ? 'bg-foreground text-background border-foreground'
          : 'bg-card text-muted-foreground border-border hover:text-foreground'"
      >
        Toutes
      </button>

      <button
        v-for="cat in availableCategories"
        :key="cat"
        type="button"
        @click="selectedCategory = cat"
        class="px-2.5 py-1 rounded-[4px] border font-semibold transition cursor-pointer"
        :class="selectedCategory === cat
          ? 'bg-brand text-brand-dark border-brand font-bold'
          : 'bg-card text-muted-foreground border-border hover:text-foreground'"
      >
        Permis {{ cat }}
      </button>

      <label class="flex items-center space-x-1.5 text-xs text-muted-foreground ml-2 cursor-pointer select-none">
        <input
          v-model="onlyWithCorrection"
          type="checkbox"
          class="rounded border-border text-brand focus:ring-0"
        />
        <span class="text-[11px]">Corrigé inclus</span>
      </label>
    </div>
  </div>
</template>
