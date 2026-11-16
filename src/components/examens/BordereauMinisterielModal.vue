<script setup lang="ts">
import { Printer, X, FileText } from 'lucide-vue-next';
import type { SessionExamen, CandidatExamenItem } from '../../types';
import { useAuthStore } from '../../stores/authStore';

const authStore = useAuthStore();

defineProps<{
  session: SessionExamen;
  candidats: CandidatExamenItem[];
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const dateCourte = (iso: string) => {
  const d = new Date(`${iso}T12:00:00`);
  return Number.isNaN(d.getTime()) ? iso : d.toLocaleDateString('fr-FR');
};

const aujourdhui = new Date().toLocaleDateString('fr-FR');

const resultat = (valeur: string) => (valeur !== 'EN_ATTENTE' ? valeur : '..........');

const imprimer = () => window.print();
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 select-none backdrop-blur-sm">
    <div class="flex max-h-[95vh] w-full max-w-4xl flex-col overflow-hidden rounded-xl border border-border bg-card shadow-2xl">
      <!-- En-tête -->
      <header class="flex items-center justify-between border-b border-border px-4 py-3">
        <div class="flex items-center gap-2 text-foreground">
          <FileText class="h-4 w-4 text-brand-dark" />
          <span class="text-xs font-semibold">Bordereau de présentation à l'examen</span>
          <span class="rounded bg-muted px-1.5 py-0.5 text-[10px] text-muted-foreground">Format A4</span>
        </div>
        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="imprimer"
            class="inline-flex items-center gap-1.5 rounded-lg bg-brand-dark px-3 py-1.5 text-xs font-semibold text-white hover:bg-brand transition"
          >
            <Printer class="h-3.5 w-3.5" />
            Imprimer
          </button>
          <button
            @click="emit('close')"
            class="rounded-lg p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition"
            aria-label="Fermer"
          >
            <X class="h-4 w-4" />
          </button>
        </div>
      </header>

      <!-- Feuille A4 -->
      <div class="flex justify-center overflow-y-auto bg-muted/40 p-4">
        <div
          id="bordereau-a4-container"
          class="w-full select-text space-y-3 border border-gray-300 bg-white p-6 font-sans text-[9px] leading-normal text-black"
        >
          <!-- Cartouche -->
          <div class="flex items-start justify-between border-b border-black pb-2">
            <div class="w-1/3">
              <div class="text-[9px] font-bold leading-tight">RÉPUBLIQUE DE CÔTE D'IVOIRE</div>
              <div class="text-[8px] italic text-gray-600">Union – Discipline – Travail</div>
              <div class="mt-1 text-[8px] font-semibold text-gray-800">Ministère des Transports</div>
              <div class="text-[7px] text-gray-600">Direction générale des transports terrestres et de la circulation</div>
            </div>

            <div class="w-1/3 text-center">
              <div class="text-[10px] font-bold">Bordereau d'examen</div>
              <div class="text-[8px] text-gray-700">
                Session <span class="font-mono">{{ session.code_session }}</span>
              </div>
            </div>

            <div class="w-1/3 text-right">
              <div class="text-[10px] font-bold">{{ authStore.schoolName.toUpperCase() }}</div>
              <div class="text-[8px] text-gray-700">Agrément {{ authStore.schoolAgrement }}</div>
              <div class="text-[7px] text-gray-600">{{ authStore.schoolAddress }}</div>
            </div>
          </div>

          <!-- Session -->
          <div class="grid grid-cols-4 gap-2 rounded border border-gray-300 bg-gray-50 p-2 text-[8px]">
            <div>
              <span class="block text-gray-500">Date de l'épreuve</span>
              <strong class="font-mono text-gray-900">{{ dateCourte(session.date_examen) }}</strong>
            </div>
            <div>
              <span class="block text-gray-500">Centre d'examen</span>
              <strong class="text-gray-900">{{ session.centre_examen }}</strong>
            </div>
            <div>
              <span class="block text-gray-500">Catégorie de permis</span>
              <strong class="text-gray-900">{{ session.categorie_permis }}</strong>
            </div>
            <div>
              <span class="block text-gray-500">Candidats présentés</span>
              <strong class="font-mono text-gray-900">{{ candidats.length }}</strong>
            </div>
          </div>

          <!-- Candidats -->
          <div class="overflow-x-auto border border-black">
            <table class="w-full border-collapse text-left text-[8px]">
              <thead>
                <tr class="border-b border-black bg-gray-100 font-semibold text-gray-900">
                  <th class="w-6 border-r border-black p-1 text-center">N°</th>
                  <th class="w-10 border-r border-black p-1 text-center">Photo</th>
                  <th class="border-r border-black p-1">Nom et prénoms</th>
                  <th class="border-r border-black p-1">Né(e) le, à</th>
                  <th class="border-r border-black p-1">Pièce d'identité</th>
                  <th class="border-r border-black p-1 text-center">Dossier</th>
                  <th class="w-14 border-r border-black p-1 text-center">Table</th>
                  <th class="w-14 border-r border-black p-1 text-center">Code</th>
                  <th class="w-14 p-1 text-center">Conduite</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(cand, idx) in candidats" :key="cand.id" class="border-b border-gray-300 last:border-b-0">
                  <td class="border-r border-black p-1 text-center font-mono">{{ idx + 1 }}</td>
                  <td class="border-r border-black p-1 text-center">
                    <div class="mx-auto flex h-9 w-7 items-center justify-center overflow-hidden border border-gray-400 bg-gray-100 text-[6px] text-gray-500">
                      <img v-if="cand.photo_url" :src="cand.photo_url" class="h-full w-full object-cover" alt="" />
                      <span v-else>Photo</span>
                    </div>
                  </td>
                  <td class="border-r border-black p-1">
                    <strong class="uppercase">{{ cand.nom }}</strong><br />
                    {{ cand.prenoms }}
                  </td>
                  <td class="border-r border-black p-1">
                    {{ dateCourte(cand.date_naissance) }}<br />
                    <span class="text-gray-600">à {{ cand.lieu_naissance }}</span>
                  </td>
                  <td class="border-r border-black p-1 font-mono">
                    {{ cand.type_piece_id }}<br />
                    {{ cand.numero_piece_id }}
                  </td>
                  <td class="border-r border-black p-1 text-center font-mono">{{ cand.numero_dossier }}</td>
                  <td class="border-r border-black p-1 text-center font-mono font-semibold">{{ cand.numero_table }}</td>
                  <td class="border-r border-black p-1 text-center font-mono">{{ resultat(cand.resultat_code) }}</td>
                  <td class="p-1 text-center font-mono">{{ resultat(cand.resultat_conduite) }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Attestation -->
          <p class="text-[8px] italic leading-snug text-gray-700">
            Je soussigné, directeur de l'établissement d'enseignement de la conduite désigné ci-dessus, atteste sur l'honneur que les candidats inscrits sur ce bordereau ont suivi les volumes horaires théoriques et pratiques obligatoires, sont médicalement aptes et sont à jour de leurs obligations financières.
          </p>

          <!-- Signatures -->
          <div class="grid grid-cols-2 gap-6 pt-1">
            <div class="flex h-24 flex-col justify-between border border-black p-2 text-center">
              <div class="text-[8px] font-semibold">
                Pour l'auto-école {{ authStore.schoolName }}<br />
                <span class="font-normal text-gray-600">Le directeur pédagogique (cachet et signature)</span>
              </div>
              <div class="text-[7px] text-gray-500">Fait à {{ authStore.currentUser.agence || 'Abidjan' }}, le {{ aujourdhui }}</div>
            </div>

            <div class="flex h-24 flex-col justify-between border border-black p-2 text-center">
              <div class="text-[8px] font-semibold">
                Pour le Ministère des Transports<br />
                <span class="font-normal text-gray-600">L'inspecteur, président du centre d'examen (visa)</span>
              </div>
              <div class="text-[7px] text-gray-500">Date et signature</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>