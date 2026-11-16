<script setup lang="ts">
import type { DossierArchiveItem } from "../../types";
import { ShieldCheck, X, FileCheck, Download, CreditCard } from "lucide-vue-next";

interface Props {
  isOpen: boolean;
  archive: DossierArchiveItem | null;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const fcfa = (n: number) => (n || 0).toLocaleString("fr-FR");
</script>

<template>
  <div
    v-if="isOpen && archive"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="w-full max-w-2xl rounded-[4px] border border-border bg-card p-6 space-y-4 max-h-[90vh] overflow-y-auto"
    >
      <div
        class="flex items-center justify-between border-b border-border pb-3"
      >
        <div class="flex items-center gap-2">
          <ShieldCheck class="h-5 w-5 text-emerald-500" />
          <div>
            <h2 class="text-base font-bold text-foreground">
              Archive Scellée : {{ archive.candidat_nom }} {{ archive.candidat_prenoms }}
            </h2>
            <p class="text-xs text-muted-foreground">
              Dossier N° {{ archive.numero_dossier }} • Boîte :
              {{ archive.reference_boite_physique || "Non classé" }}
            </p>
          </div>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground"
        >
          <X class="h-4 w-4" />
        </button>
      </div>

      <!-- Bilan scellé -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 text-xs">
        <div class="p-2.5 rounded-[4px] border border-border bg-muted/20">
          <span class="text-muted-foreground text-[10px]">Catégorie Permis</span>
          <p class="font-bold text-foreground mt-0.5">
            Permis {{ archive.categorie_permis }}
          </p>
        </div>

        <div class="p-2.5 rounded-[4px] border border-border bg-muted/20">
          <span class="text-muted-foreground text-[10px]">Total Règlement</span>
          <p class="font-bold text-emerald-600 dark:text-emerald-400 mt-0.5">
            {{ fcfa(archive.total_paye) }} FCFA
          </p>
        </div>

        <div class="p-2.5 rounded-[4px] border border-border bg-muted/20">
          <span class="text-muted-foreground text-[10px]">Code Théorique</span>
          <p class="font-bold text-foreground mt-0.5">
            {{ archive.heures_code_effectuees }} h validées
          </p>
        </div>

        <div class="p-2.5 rounded-[4px] border border-border bg-muted/20">
          <span class="text-muted-foreground text-[10px]">Conduite Pratique</span>
          <p class="font-bold text-foreground mt-0.5">
            {{ archive.heures_conduite_effectuees }} h validées
          </p>
        </div>
      </div>

      <!-- Documents scellés dans l'archive -->
      <div class="space-y-2">
        <h4
          class="text-xs font-bold uppercase tracking-wider text-muted-foreground"
        >
          Documents Scellés ({{ (archive.documents || []).length }})
        </h4>

        <div
          v-if="!archive.documents || archive.documents.length === 0"
          class="p-4 rounded-[4px] border border-border bg-muted/20 text-center text-xs text-muted-foreground"
        >
          Aucun document numérisé n'était joint à ce dossier lors de l'archivage.
        </div>

        <div
          v-else
          class="divide-y divide-border border border-border rounded-[4px] overflow-hidden"
        >
          <div
            v-for="d in archive.documents"
            :key="d.id"
            class="p-3 flex items-center justify-between text-xs hover:bg-muted/30"
          >
            <div class="flex items-center gap-2">
              <FileCheck class="h-4 w-4 text-emerald-500" />
              <div>
                <span class="font-semibold text-foreground">{{ d.titre }}</span>
                <span class="text-muted-foreground ml-2">({{ d.nom_fichier }})</span>
              </div>
            </div>

            <a
              v-if="d.fichier_url_ou_base64"
              :href="d.fichier_url_ou_base64"
              :download="d.nom_fichier"
              class="inline-flex items-center gap-1 text-brand font-medium hover:underline text-[11px]"
            >
              <Download class="h-3 w-3" />
              Télécharger
            </a>
          </div>
        </div>
      </div>

      <!-- Paiements scellés -->
      <div class="space-y-2">
        <h4
          class="text-xs font-bold uppercase tracking-wider text-muted-foreground"
        >
          Reçus & Règlements Thermiques ({{ (archive.paiements || []).length }})
        </h4>

        <div
          class="divide-y divide-border border border-border rounded-[4px] overflow-hidden max-h-40 overflow-y-auto"
        >
          <div
            v-for="p in archive.paiements"
            :key="p.id"
            class="p-2.5 flex items-center justify-between text-xs hover:bg-muted/30"
          >
            <div class="flex items-center gap-2">
              <CreditCard class="h-3.5 w-3.5 text-sky-500" />
              <span>Reçu N° {{ p.reference_recu || p.id.slice(0, 8) }} • {{ p.mode_reglement }}</span>
            </div>
            <span class="font-bold text-foreground">
              {{ fcfa(p.montant) }} FCFA
            </span>
          </div>
        </div>
      </div>

      <div class="flex items-center justify-end border-t border-border pt-3">
        <button
          type="button"
          @click="emit('close')"
          class="px-4 py-1.5 rounded-[4px] bg-brand text-brand-dark font-bold text-xs"
        >
          Fermer la consultation
        </button>
      </div>
    </div>
  </div>
</template>

