<script setup lang="ts">
import {
  Archive,
  FolderArchive,
  UploadCloud,
  Printer,
  ShieldCheck,
  Cloud,
  FileCheck,
  CreditCard,
  AlertCircle,
} from "lucide-vue-next";

interface Props {
  totalArchives: number;
  totalDocs: number;
  totalFonds: number;
  eligibleCount: number;
  bucketName: string;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: "open-archive"): void;
  (e: "open-upload"): void;
  (e: "print-inventory"): void;
}>();

const fcfa = (n: number) => (n || 0).toLocaleString("fr-FR");
</script>

<template>
  <div class="space-y-4">
    <!-- EN-TÊTE PRINCIPAL & ACTIONS -->
    <div
      class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 border-b border-border pb-4"
    >
      <div>
        <div class="flex items-center gap-2 flex-wrap">
          <div
            class="h-9 w-9 rounded-[4px] bg-brand/10 text-brand flex items-center justify-center font-bold"
          >
            <Archive class="h-5 w-5" />
          </div>
          <h1 class="text-lg font-bold text-foreground">
            Centre d'Archivage & Rétention Documentaire
          </h1>
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-emerald-500/10 text-emerald-700 border border-emerald-500/30 dark:text-emerald-400"
          >
            <ShieldCheck class="h-3.5 w-3.5 text-emerald-500" />
            Coffre-fort Scellé
          </span>
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-sky-500/10 text-sky-700 border border-sky-500/30 dark:text-sky-400"
            :title="'Stockage Cloud : Cloudflare R2 (' + bucketName + ')'"
          >
            <Cloud class="h-3.5 w-3.5 text-sky-500" />
            Cloudflare R2 ({{ bucketName }})
          </span>
        </div>
        <p class="text-xs text-muted-foreground mt-1">
          Gestion sécurisée des dossiers clos, conservation des pièces d'identité scannées, reçus thermiques et certificats légaux.
        </p>
      </div>

      <div class="flex items-center gap-2 shrink-0 flex-wrap">
        <button
          type="button"
          @click="emit('open-archive')"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-[4px] bg-brand text-brand-dark hover:bg-brand-dark hover:text-white text-xs font-semibold transition active:scale-95"
        >
          <FolderArchive class="h-4 w-4" />
          Archiver un Dossier
        </button>

        <button
          type="button"
          @click="emit('open-upload')"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-[4px] border border-border bg-card text-xs font-medium text-foreground hover:bg-muted transition active:scale-95"
        >
          <UploadCloud class="h-4 w-4 text-brand" />
          Téléverser Document
        </button>

        <button
          type="button"
          @click="emit('print-inventory')"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-[4px] border border-border bg-card text-xs font-medium text-foreground hover:bg-muted transition active:scale-95"
          title="Imprimer le bordereau officiel de versement"
        >
          <Printer class="h-3.5 w-3.5" />
          Bordereau A4
        </button>
      </div>
    </div>

    <!-- CARTES STATISTIQUES RÉCAPITULATIVES -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
      <div class="rounded-[4px] border border-border bg-card p-4 flex flex-col justify-between">
        <div class="flex items-center justify-between text-muted-foreground">
          <span class="text-xs font-medium">Dossiers Archivés</span>
          <Archive class="h-4 w-4 text-brand" />
        </div>
        <div class="mt-2">
          <span class="text-2xl font-bold text-foreground tabular-nums">
            {{ totalArchives }}
          </span>
          <p class="text-[11px] text-muted-foreground mt-0.5">
            Dossiers complets scellés
          </p>
        </div>
      </div>

      <div class="rounded-[4px] border border-border bg-card p-4 flex flex-col justify-between">
        <div class="flex items-center justify-between text-muted-foreground">
          <span class="text-xs font-medium">Pièces Numérisées</span>
          <FileCheck class="h-4 w-4 text-emerald-500" />
        </div>
        <div class="mt-2">
          <span class="text-2xl font-bold text-foreground tabular-nums">
            {{ totalDocs }}
          </span>
          <p class="text-[11px] text-emerald-600 dark:text-emerald-400 mt-0.5">
            CNI, Visites médicales, Contrats
          </p>
        </div>
      </div>

      <div class="rounded-[4px] border border-border bg-card p-4 flex flex-col justify-between">
        <div class="flex items-center justify-between text-muted-foreground">
          <span class="text-xs font-medium">Fonds Clôturés</span>
          <CreditCard class="h-4 w-4 text-sky-500" />
        </div>
        <div class="mt-2">
          <span class="text-xl font-bold text-foreground tabular-nums">
            {{ fcfa(totalFonds) }} FCFA
          </span>
          <p class="text-[11px] text-sky-600 dark:text-sky-400 mt-0.5">
            Total des recettes archivées
          </p>
        </div>
      </div>

      <div class="rounded-[4px] border border-border bg-card p-4 flex flex-col justify-between">
        <div class="flex items-center justify-between text-muted-foreground">
          <span class="text-xs font-medium">En Attente d'Archivage</span>
          <AlertCircle class="h-4 w-4 text-amber-500" />
        </div>
        <div class="mt-2">
          <span class="text-2xl font-bold text-foreground tabular-nums">
            {{ eligibleCount }}
          </span>
          <p class="text-[11px] text-amber-600 dark:text-amber-400 mt-0.5">
            Dossiers actifs éligibles
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

