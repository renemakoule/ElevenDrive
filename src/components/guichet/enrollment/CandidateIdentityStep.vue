<script setup lang="ts">
import { Camera, User } from "lucide-vue-next";

defineProps<{
  nom: string;
  prenoms: string;
  genre: "M" | "F";
  dateNaissance: string;
  lieuNaissance: string;
  photoUrl: string | null;
  ageCalcule: number;
  ageMinimum: number;
  isAgeLegal: boolean;
  categoriePermis?: string;
}>();

const emit = defineEmits<{
  (e: "update:nom", val: string): void;
  (e: "update:prenoms", val: string): void;
  (e: "update:genre", val: "M" | "F"): void;
  (e: "update:dateNaissance", val: string): void;
  (e: "update:lieuNaissance", val: string): void;
  (e: "open-webcam"): void;
  (e: "remove-photo"): void;
}>();

const field =
  "h-8 w-full rounded-[5px] border border-border bg-card px-2.5 text-xs text-foreground transition placeholder:text-muted-foreground focus:border-brand focus:outline-none focus:ring-2 focus:ring-brand/20";
const label = "mb-1 block text-xs text-muted-foreground";
</script>

<template>
  <section class="space-y-3 px-4 py-4">
    <div class="flex items-center justify-between gap-3">
      <h3 class="text-xs font-semibold text-foreground">Identité</h3>
      <span
        class="flex items-center gap-1.5 text-[11px] font-medium tabular-nums"
        :class="isAgeLegal ? 'text-brand-dark' : 'text-rose-700'"
        aria-live="polite"
      >
        <span
          class="h-1.5 w-1.5 rounded-full"
          :class="isAgeLegal ? 'bg-brand' : 'bg-rose-500'"
        />
        {{ ageCalcule }} ans ·
        {{
          isAgeLegal
            ? `éligible au permis ${categoriePermis ?? ""}`
            : `minimum ${ageMinimum} ans`
        }}
      </span>
    </div>

    <div class="flex items-start gap-4">
      <div class="w-20 shrink-0">
        <div
          class="flex h-24 w-20 items-center justify-center overflow-hidden rounded-lg border border-border bg-muted/40"
        >
          <img
            v-if="photoUrl"
            :src="photoUrl"
            alt="Photo du candidat"
            class="h-full w-full object-cover"
          />
          <span
            v-else
            class="flex flex-col items-center gap-1 text-muted-foreground"
          >
            <User class="h-6 w-6" />
            <span class="text-[10px]">Photo 4x4</span>
          </span>
        </div>
        <button
          type="button"
          @click="emit('open-webcam')"
          class="mt-2 inline-flex h-7 w-full items-center justify-center gap-1 rounded-[4px] border border-border bg-card text-[11px] font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
        >
          <Camera class="h-3.5 w-3.5 text-muted-foreground" />
          {{ photoUrl ? "Reprendre" : "Webcam" }}
        </button>
        <button
          v-if="photoUrl"
          type="button"
          @click="emit('remove-photo')"
          class="mt-1 w-full text-center text-[11px] text-muted-foreground transition hover:text-destructive"
        >
          Retirer
        </button>
      </div>

      <div class="min-w-0 flex-1 space-y-3">
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label for="nom" :class="label"
              >Nom <span class="text-rose-500">*</span></label
            >
            <input
              id="nom"
              type="text"
              :value="nom"
              @input="emit('update:nom', ($event.target as HTMLInputElement).value)"
              required
              placeholder="MAKOULE"
              :class="[field, 'uppercase']"
            />
          </div>
          <div>
            <label for="prenoms" :class="label"
              >Prénoms <span class="text-rose-500">*</span></label
            >
            <input
              id="prenoms"
              type="text"
              :value="prenoms"
              @input="
                emit('update:prenoms', ($event.target as HTMLInputElement).value)
              "
              required
              placeholder="Rene Boris"
              :class="field"
            />
          </div>
        </div>

        <div class="grid grid-cols-3 gap-3">
          <div>
            <label for="genre" :class="label">Genre</label>
            <select
              id="genre"
              :value="genre"
              @change="
                emit(
                  'update:genre',
                  ($event.target as HTMLSelectElement).value as 'M' | 'F',
                )
              "
              :class="field"
            >
              <option value="M">Masculin</option>
              <option value="F">Féminin</option>
            </select>
          </div>
          <div>
            <label for="naissance" :class="label">Né(e) le</label>
            <input
              id="naissance"
              type="date"
              :value="dateNaissance"
              @input="
                emit(
                  'update:dateNaissance',
                  ($event.target as HTMLInputElement).value,
                )
              "
              required
              :class="field"
            />
          </div>
          <div>
            <label for="lieu" :class="label">Lieu de naissance</label>
            <input
              id="lieu"
              type="text"
              :value="lieuNaissance"
              @input="
                emit(
                  'update:lieuNaissance',
                  ($event.target as HTMLInputElement).value,
                )
              "
              required
              placeholder="Douala"
              :class="field"
            />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

