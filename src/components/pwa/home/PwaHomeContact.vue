<script setup lang="ts">
import type { Agence } from '../../../types';
import { Share2, Phone, MapPin } from 'lucide-vue-next';

interface Props {
  schoolName: string;
  schoolPhone: string;
  schoolAddress: string;
  agences: Agence[];
  whatsappUrl: string;
  callUrl: string;
}

const props = defineProps<Props>();

const displayAgences = () => {
  if (props.agences && props.agences.length > 0) {
    return props.agences;
  }
  return [
    {
      id: 'siege',
      auto_ecole_id: '',
      code_agence: 'SIEGE',
      nom: 'Siège Principal',
      adresse: props.schoolAddress || 'Centre Ville',
      telephone: props.schoolPhone,
      ville: 'Douala',
      commune: 'Makepe',
      est_siege: true,
    },
  ];
};
</script>

<template>
  <section class="shrink-0 w-full py-10 max-w-4xl mx-auto px-4 sm:px-6">
    <div class="rounded-[4px] border border-border bg-card p-6 sm:p-8 shadow-sm">
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 pb-6 border-b border-border">
        <div>
          <h4 class="text-base sm:text-lg font-bold text-foreground">
            Besoin d'un renseignement ou d'une inscription ?
          </h4>
          <p class="text-xs text-muted-foreground mt-0.5">
            Contactez directement le guichet de l'auto-école {{ schoolName }}.
          </p>
        </div>

        <div class="flex items-center gap-2">
          <a
            :href="whatsappUrl"
            target="_blank"
            class="inline-flex items-center gap-1.5 px-4 py-2 rounded-[4px] bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-semibold transition shadow-xs"
          >
            <Share2 class="w-3.5 h-3.5" />
            <span>WhatsApp Direct</span>
          </a>

          <a
            v-if="callUrl !== '#'"
            :href="callUrl"
            class="inline-flex items-center gap-1.5 px-4 py-2 rounded-[4px] border border-border bg-background hover:bg-muted text-foreground text-xs font-semibold transition"
          >
            <Phone class="w-3.5 h-3.5 text-brand" />
            <span>Appeler</span>
          </a>
        </div>
      </div>

      <!-- Liste des agences physiques -->
      <div class="mt-6">
        <h5 class="text-xs font-bold uppercase tracking-wider text-muted-foreground mb-3">
          Succursales et Agences d'accueil
        </h5>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div
            v-for="agence in displayAgences()"
            :key="agence.id"
            class="p-3.5 rounded-[4px] border border-border bg-background/50 space-y-1 hover:border-brand/40 transition"
          >
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-foreground">{{ agence.nom }}</span>
              <span v-if="agence.est_siege" class="px-2 py-0.5 rounded text-[9px] font-bold bg-brand/10 text-brand border border-brand/20">
                Siège
              </span>
            </div>
            <p class="text-[11px] text-muted-foreground flex items-center gap-1">
              <MapPin class="w-3 h-3 text-muted-foreground shrink-0" />
              <span class="truncate">{{ agence.adresse || agence.commune || agence.ville }}</span>
            </p>
            <p v-if="agence.telephone" class="text-[11px] text-muted-foreground font-mono flex items-center gap-1">
              <Phone class="w-3 h-3 text-muted-foreground shrink-0" />
              <span>{{ agence.telephone }}</span>
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
