<script setup lang="ts">
import type { FormulePermis } from "../../../types";

defineProps<{
  telephoneWhatsapp: string;
  typePiece: string;
  numeroPiece: string;
  visiteMedicale: "VALIDE_SANS_RESERVE" | "VALIDE_AVEC_VERRES" | "EN_ATTENTE";
  acompteInitial: number;
  modeReglement: string;
  acompteTropBas: boolean;
  soldeApresAcompte: number;
  selectedFormule?: FormulePermis;
}>();

const emit = defineEmits<{
  (e: "update:telephoneWhatsapp", val: string): void;
  (e: "update:typePiece", val: string): void;
  (e: "update:numeroPiece", val: string): void;
  (
    e: "update:visiteMedicale",
    val: "VALIDE_SANS_RESERVE" | "VALIDE_AVEC_VERRES" | "EN_ATTENTE",
  ): void;
  (e: "update:acompteInitial", val: number): void;
  (e: "update:modeReglement", val: string): void;
}>();

const fcfa = (n: number) => n.toLocaleString("fr-FR");

const medicalOptions = [
  {
    value: "VALIDE_SANS_RESERVE",
    label: "Apte sans réserve",
    on: "border-emerald-500 bg-emerald-50 text-emerald-800 dark:bg-emerald-950/40 dark:text-emerald-300",
  },
  {
    value: "VALIDE_AVEC_VERRES",
    label: "Port de verres",
    on: "border-sky-500 bg-sky-50 text-sky-800 dark:bg-sky-950/40 dark:text-sky-300",
  },
  {
    value: "EN_ATTENTE",
    label: "Certificat en attente",
    on: "border-amber-500 bg-amber-50 text-amber-800 dark:bg-amber-950/40 dark:text-amber-300",
  },
] as const;

const field =
  "h-8 w-full rounded-[5px] border border-border bg-card px-2.5 text-xs text-foreground transition placeholder:text-muted-foreground focus:border-brand focus:outline-none focus:ring-2 focus:ring-brand/20";
const label = "mb-1 block text-xs text-muted-foreground";
</script>

<template>
  <div>
    <!-- Contact et pièce -->
    <section class="space-y-3 px-4 py-4">
      <h3 class="text-xs font-semibold text-foreground">
        Contact et pièce d'identité
      </h3>

      <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
        <div>
          <label for="tel" :class="label"
            >WhatsApp <span class="text-rose-500">*</span></label
          >
          <input
            id="tel"
            type="tel"
            :value="telephoneWhatsapp"
            @input="
              emit(
                'update:telephoneWhatsapp',
                ($event.target as HTMLInputElement).value,
              )
            "
            required
            placeholder="+237 651 727 932"
            :class="[field, 'tabular-nums']"
          />
        </div>
        <div>
          <label for="type-piece" :class="label">Type de pièce</label>
          <select
            id="type-piece"
            :value="typePiece"
            @change="
              emit(
                'update:typePiece',
                ($event.target as HTMLSelectElement).value,
              )
            "
            :class="field"
          >
            <option value="CNI">CNI nationale</option>
            <option value="PASSEPORT">Passeport</option>
            <option value="ATTESTATION">Attestation d'identité</option>
            <option value="CONSULAIRE">Carte consulaire</option>
          </select>
        </div>
        <div>
          <label for="num-piece" :class="label"
            >N° de la pièce <span class="text-rose-500">*</span></label
          >
          <input
            id="num-piece"
            type="text"
            :value="numeroPiece"
            @input="
              emit(
                'update:numeroPiece',
                ($event.target as HTMLInputElement).value,
              )
            "
            required
            placeholder="CMR00389218"
            :class="[field, 'font-mono']"
          />
        </div>
      </div>

      <fieldset>
        <legend :class="label">Visite médicale</legend>
        <div class="grid grid-cols-1 gap-2 sm:grid-cols-3">
          <label
            v-for="o in medicalOptions"
            :key="o.value"
            class="flex h-8 cursor-pointer items-center justify-center rounded-lg border px-2 text-xs font-medium transition focus-within:ring-2 focus-within:ring-brand/40"
            :class="
              visiteMedicale === o.value
                ? o.on
                : 'border-border bg-card text-muted-foreground hover:bg-muted'
            "
          >
            <input
              type="radio"
              :value="o.value"
              :checked="visiteMedicale === o.value"
              @change="emit('update:visiteMedicale', o.value)"
              class="sr-only"
            />
            {{ o.label }}
          </label>
        </div>
      </fieldset>
    </section>

    <!-- Acompte initial -->
    <section class="space-y-3 border-t border-border px-4 py-4">
      <h3 class="text-xs font-semibold text-foreground">Acompte initial</h3>

      <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <div>
          <label for="acompte" :class="label">Montant versé</label>
          <div class="relative">
            <input
              id="acompte"
              type="number"
              :value="acompteInitial"
              @input="
                emit(
                  'update:acompteInitial',
                  Number(($event.target as HTMLInputElement).value) || 0,
                )
              "
              step="1000"
              min="0"
              class="h-8 w-full rounded-[4px] border bg-card pl-2.5 pr-12 text-right text-xs font-semibold tabular-nums text-foreground transition focus:outline-none focus:ring-2"
              :class="
                acompteTropBas
                  ? 'border-destructive/40 focus:border-destructive focus:ring-destructive/20'
                  : 'border-border focus:border-brand focus:ring-brand/20'
              "
            />
            <span
              class="pointer-events-none absolute right-2.5 top-1/2 -translate-y-1/2 text-[11px] text-muted-foreground"
              >FCFA</span
            >
          </div>
          <p
            class="mt-1 text-[11px] tabular-nums"
            :class="
              acompteTropBas ? 'text-destructive' : 'text-muted-foreground'
            "
          >
            <template v-if="acompteTropBas && selectedFormule">
              Minimum requis :
              {{ fcfa(selectedFormule.acompte_minimum) }} FCFA
            </template>
            <template v-else
              >Reste à payer après acompte :
              {{ fcfa(soldeApresAcompte) }} FCFA</template
            >
          </p>
        </div>

        <div>
          <label for="mode" :class="label">Mode de règlement</label>
          <select
            id="mode"
            :value="modeReglement"
            @change="
              emit(
                'update:modeReglement',
                ($event.target as HTMLSelectElement).value,
              )
            "
            :class="field"
          >
            <option value="ESPECES">Espèces (tiroir-caisse)</option>
            <option value="WAVE">Wave</option>
            <option value="ORANGE_MONEY">Orange Money</option>
            <option value="MTN_MOMO">MTN MoMo</option>
            <option value="MOOV">Moov Money</option>
          </select>
        </div>
      </div>
    </section>
  </div>
</template>
