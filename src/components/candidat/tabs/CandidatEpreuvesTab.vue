<script setup lang="ts">
import type { EpreuvePdf } from "../../../types";
import { useEpreuvesStore } from "../../../stores/epreuvesStore";
import { useOfflineDownloadsStore } from "../../../stores/offlineDownloadsStore";
import {
  BookOpen,
  CheckCircle2,
  Lock,
  FileText,
  Clock,
  HelpCircle,
  Eye,
  DownloadCloud,
} from "lucide-vue-next";

defineProps<{
  studentCategory: string;
  epreuves: EpreuvePdf[];
}>();

const epreuvesStore = useEpreuvesStore();
const offlineStore = useOfflineDownloadsStore();
</script>

<template>
  <div class="space-y-4">
    <!-- Bannière explicative & Catégorie active -->
    <div
      class="p-4 rounded-[4px] bg-brand-light/30 border border-brand-border flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3"
    >
      <div class="flex items-center space-x-3">
        <div
          class="w-10 h-10 rounded-[4px] bg-brand-light border border-brand-border text-brand-dark flex items-center justify-center font-black"
        >
          <BookOpen class="w-5 h-5 text-brand" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-sm font-bold text-foreground">
              Annales & Épreuves Officielles · Permis {{ studentCategory }}
            </h2>
            <span
              class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-[4px] bg-emerald-500/10 text-emerald-600 border border-emerald-500/20"
            >
              <CheckCircle2 class="w-3 h-3" /> Catégorie active
            </span>
          </div>
          <p class="text-xs text-muted-foreground mt-0.5">
            Consultez et téléchargez les sujets d'examens théoriques mis à
            disposition par votre auto-école pour votre catégorie.
          </p>
        </div>
      </div>

      <div
        class="px-2.5 py-1 rounded-[4px] bg-card border border-border text-[11px] text-muted-foreground flex items-center gap-1.5 shrink-0"
      >
        <Lock class="w-3.5 h-3.5 text-brand" />
        <span>Accès réservé au <strong>Permis {{ studentCategory }}</strong></span>
      </div>
    </div>

    <!-- Aucune épreuve -->
    <div
      v-if="epreuves.length === 0"
      class="p-12 text-center rounded-[4px] border border-dashed border-border bg-card space-y-2"
    >
      <FileText class="w-10 h-10 text-muted-foreground/40 mx-auto" />
      <h3 class="text-sm font-semibold text-foreground">
        Aucune épreuve disponible pour le Permis {{ studentCategory }}
      </h3>
      <p class="text-xs text-muted-foreground max-w-md mx-auto">
        Votre auto-école n'a pas encore téléversé d'épreuve pour la catégorie
        Permis {{ studentCategory }}. Vos moniteurs et votre secrétariat en
        ajouteront prochainement.
      </p>
    </div>

    <!-- Grille des épreuves PDF -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="ep in epreuves"
        :key="ep.id"
        class="p-4 rounded-[4px] bg-card border border-border hover:border-brand-border/60 transition shadow-xs flex flex-col justify-between space-y-3"
      >
        <div class="space-y-2">
          <div class="flex items-center justify-between gap-2">
            <span
              class="inline-flex items-center text-[10px] font-bold px-2 py-0.5 rounded-[4px] bg-brand-light text-brand-dark border border-brand-border uppercase tracking-wider"
            >
              Permis {{ ep.categorie_permis }}
            </span>

            <span
              v-if="ep.avec_corrige"
              class="inline-flex items-center gap-1 text-[10px] font-medium text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 px-1.5 py-0.5 rounded border border-emerald-500/20"
            >
              <CheckCircle2 class="w-3 h-3" /> Corrigé inclus
            </span>
            <span
              v-else
              class="text-[10px] text-muted-foreground bg-muted px-1.5 py-0.5 rounded"
            >
              Sujet d'examen blanc
            </span>
          </div>

          <h3 class="text-xs sm:text-sm font-bold text-foreground leading-snug">
            {{ ep.titre }}
          </h3>

          <p v-if="ep.description" class="text-xs text-muted-foreground line-clamp-2">
            {{ ep.description }}
          </p>
        </div>

        <div
          class="space-y-2 pt-2 border-t border-border text-[11px] text-muted-foreground"
        >
          <div class="flex items-center justify-between">
            <span class="flex items-center gap-1">
              <Clock class="w-3.5 h-3.5" /> {{ ep.duree_minutes }} minutes
            </span>
            <span class="flex items-center gap-1">
              <HelpCircle class="w-3.5 h-3.5" />
              {{ ep.nombre_questions }} questions
            </span>
            <span class="font-mono text-[10px]">
              {{ ep.fichier_taille || "PDF" }}
            </span>
          </div>

          <div class="text-[10px]">
            Session :
            <strong class="text-foreground">{{ ep.session || ep.annee }}</strong>
          </div>
        </div>

        <div class="pt-2 flex items-center gap-2 border-t border-border">
          <button
            type="button"
            @click="epreuvesStore.previewPdf(ep)"
            class="flex-1 flex items-center justify-center space-x-1.5 px-3 py-1.5 rounded-[4px] bg-brand hover:bg-brand-dark text-brand-dark font-bold text-xs transition"
          >
            <Eye class="w-3.5 h-3.5" />
            <span>Consulter le Sujet</span>
          </button>

          <button
            type="button"
            @click="offlineStore.downloadEpreuve(ep)"
            class="flex items-center justify-center p-1.5 rounded-[4px] border transition cursor-pointer shrink-0"
            :class="
              offlineStore.isDownloaded(ep.id)
                ? 'border-emerald-500/40 bg-emerald-500/10 text-emerald-600'
                : 'border-border hover:bg-muted text-muted-foreground hover:text-foreground'
            "
            :title="
              offlineStore.isDownloaded(ep.id)
                ? 'Disponible hors-ligne'
                : 'Conserver pour révision hors-ligne sans connexion'
            "
          >
            <component
              :is="offlineStore.isDownloaded(ep.id) ? CheckCircle2 : DownloadCloud"
              class="w-3.5 h-3.5"
            />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

