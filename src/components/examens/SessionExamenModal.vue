<script setup lang="ts">
import { ref, computed } from 'vue';
import { X, Calendar, MapPin, Award, UserCheck, FileText, CheckCircle2, Trash2 } from 'lucide-vue-next';
import { api } from '../../db/api';
import { useGuichetStore } from '../../stores/guichetStore';
import type { SessionExamen, CreerOuModifierSessionExamenInput } from '../../types';

const props = defineProps<{
  session?: SessionExamen | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'saved', session: SessionExamen): void;
  (e: 'deleted'): void;
}>();

const store = useGuichetStore();

const isEditing = computed(() => !!props.session && !!props.session.id);
const isSubmitting = ref(false);
const isDeleting = ref(false);
const errorMessage = ref('');

// Formulaire : aucune valeur par défaut arbitraire
const form = ref<CreerOuModifierSessionExamenInput>({
  id: props.session?.id || undefined,
  code_session: props.session?.code_session || '',
  date_examen: props.session?.date_examen || '',
  centre_examen: props.session?.centre_examen || '',
  categorie_permis: props.session?.categorie_permis || 'B',
  inspecteur_nom: props.session?.inspecteur_nom || '',
  observations: '',
  statut_session: props.session?.statut_session || 'EN_PREPARATION',
});

const categoriesPermis = [
  { code: 'B', label: 'Permis B (Véhicules Légers / Tourisme)' },
  { code: 'A', label: 'Permis A (Grosses Cylindrées / Motocyclettes)' },
  { code: 'A1', label: 'Permis A1 (Motos légères)' },
  { code: 'C', label: 'Permis C (Poids Lourds)' },
  { code: 'D', label: 'Permis D (Transport en commun)' },
  { code: 'E', label: 'Permis E (Remorques & Semi-remorques)' },
];

const statutsSession = [
  { code: 'EN_PREPARATION', label: 'En préparation (Inscription des candidats)' },
  { code: 'BORDEREAU_DEPOSE', label: 'Bordereau déposé au Ministère / Centre' },
  { code: 'RESULTATS_PUBLIES', label: 'Résultats publiés (PV validé)' },
  { code: 'CLOTUREE', label: 'Session clôturée & archivée' },
];

const enregistrer = async () => {
  errorMessage.value = '';

  if (!form.value.code_session.trim()) {
    errorMessage.value = 'Veuillez saisir un intitulé ou code officiel de session (ex: EXAM-NOVEMBRE-2026).';
    return;
  }

  if (!form.value.date_examen) {
    errorMessage.value = 'Veuillez renseigner la date exacte de déroulement de l\'examen.';
    return;
  }

  if (!form.value.centre_examen.trim()) {
    errorMessage.value = 'Veuillez indiquer le centre officiel où se déroulera l\'examen.';
    return;
  }

  isSubmitting.value = true;
  try {
    const payload: CreerOuModifierSessionExamenInput = {
      id: form.value.id,
      code_session: form.value.code_session.trim(),
      date_examen: form.value.date_examen,
      centre_examen: form.value.centre_examen.trim(),
      categorie_permis: form.value.categorie_permis,
      inspecteur_nom: form.value.inspecteur_nom?.trim() || undefined,
      observations: form.value.observations?.trim() || undefined,
      statut_session: form.value.statut_session,
    };

    const saved = await api.enregistrerSessionExamen(payload);
    store.showToast(
      isEditing.value ? 'Session d\'examen mise à jour avec succès !' : 'Nouvelle session d\'examen planifiée avec succès !',
      'success'
    );
    emit('saved', saved);
    emit('close');
  } catch (err: any) {
    errorMessage.value = err.message || 'Erreur lors de l\'enregistrement de la session.';
    store.showToast(errorMessage.value, 'error');
  } finally {
    isSubmitting.value = false;
  }
};

const supprimer = async () => {
  if (!props.session?.id) return;
  if (!confirm(`Êtes-vous sûr de vouloir supprimer définitivement la session « ${props.session.code_session} » ? Cette action est irréversible.`)) {
    return;
  }

  isDeleting.value = true;
  try {
    await api.supprimerSessionExamen(props.session.id);
    store.showToast('Session d\'examen supprimée avec succès.', 'success');
    emit('deleted');
    emit('close');
  } catch (err: any) {
    errorMessage.value = err.message || 'Impossible de supprimer la session.';
    store.showToast(errorMessage.value, 'error');
  } finally {
    isDeleting.value = false;
  }
};
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 select-none backdrop-blur-xs">
    <div class="flex max-h-[92vh] w-full max-w-lg flex-col overflow-hidden rounded-xl border border-border bg-card shadow-2xl">
      <!-- En-tête -->
      <header class="flex items-center justify-between border-b border-border px-5 py-3.5 bg-muted/30">
        <div class="flex items-center gap-2.5">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-brand-light text-brand">
            <Calendar class="h-4 w-4" />
          </div>
          <div>
            <h3 class="text-sm font-semibold text-foreground">
              {{ isEditing ? 'Modifier la session d\'examen' : 'Planifier une session d\'examen d\'État' }}
            </h3>
            <p class="text-[10px] text-muted-foreground">
              Définissez la date réelle, le centre et la catégorie d'examen officiel.
            </p>
          </div>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="rounded-lg p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <!-- Corps du formulaire -->
      <form @submit.prevent="enregistrer" class="flex flex-1 flex-col overflow-y-auto p-5 space-y-4 text-xs">
        <!-- Message d'erreur -->
        <div
          v-if="errorMessage"
          class="rounded-md border border-rose-200 bg-rose-50 p-2.5 text-[11px] text-rose-700 dark:border-rose-900/50 dark:bg-rose-950/40 dark:text-rose-300"
        >
          {{ errorMessage }}
        </div>

        <!-- Code / Intitulé de la session -->
        <div>
          <label class="mb-1 block font-medium text-foreground">
            Intitulé / Code officiel de la session <span class="text-rose-500">*</span>
          </label>
          <div class="relative">
            <input
              type="text"
              v-model="form.code_session"
              placeholder="Ex: SESSION-NOVEMBRE-2026 ou EXAM-B-01"
              required
              class="w-full rounded-md border border-border bg-background px-3 py-2 text-xs text-foreground placeholder:text-muted-foreground/60 focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand font-mono"
            />
          </div>
          <p class="mt-1 text-[10px] text-muted-foreground">
            Saisissez le code ou la désignation officielle délivrée par le Ministère ou l'Auto-École.
          </p>
        </div>

        <!-- Date de l'examen -->
        <div>
          <label class="mb-1 block font-medium text-foreground">
            Date officielle de déroulement <span class="text-rose-500">*</span>
          </label>
          <div class="relative">
            <input
              type="date"
              v-model="form.date_examen"
              required
              class="w-full rounded-md border border-border bg-background px-3 py-2 text-xs text-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand"
            />
          </div>
          <p class="mt-1 text-[10px] text-muted-foreground">
            Aucune date n'est imposée par défaut : choisissez le jour exact du passage des épreuves.
          </p>
        </div>

        <!-- Centre d'examen & Catégorie -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="mb-1 block font-medium text-foreground">
              Centre officiel d'examen <span class="text-rose-500">*</span>
            </label>
            <div class="relative">
              <input
                type="text"
                v-model="form.centre_examen"
                placeholder="Ex: Centre DGTTC / Ministère"
                required
                class="w-full rounded-md border border-border bg-background px-3 py-2 text-xs text-foreground placeholder:text-muted-foreground/60 focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand"
              />
            </div>
          </div>

          <div>
            <label class="mb-1 block font-medium text-foreground">
              Catégorie de permis <span class="text-rose-500">*</span>
            </label>
            <select
              v-model="form.categorie_permis"
              class="w-full rounded-md border border-border bg-background px-3 py-2 text-xs text-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand"
            >
              <option v-for="c in categoriesPermis" :key="c.code" :value="c.code">
                {{ c.code }} - {{ c.label }}
              </option>
            </select>
          </div>
        </div>

        <!-- Inspecteur / Examinateur d'État -->
        <div>
          <label class="mb-1 block font-medium text-foreground">
            Inspecteur / Examinateur d'État (Optionnel)
          </label>
          <div class="relative">
            <input
              type="text"
              v-model="form.inspecteur_nom"
              placeholder="Ex: Inspecteur DGTTC / Ministère des Transports"
              class="w-full rounded-md border border-border bg-background px-3 py-2 text-xs text-foreground placeholder:text-muted-foreground/60 focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand"
            />
          </div>
        </div>

        <!-- Statut de la session (si modification) -->
        <div v-if="isEditing">
          <label class="mb-1 block font-medium text-foreground">
            Statut de la session d'examen
          </label>
          <select
            v-model="form.statut_session"
            class="w-full rounded-md border border-border bg-background px-3 py-2 text-xs text-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand"
          >
            <option v-for="s in statutsSession" :key="s.code" :value="s.code">
              {{ s.label }}
            </option>
          </select>
        </div>

        <!-- Observations / Consignes particulières -->
        <div>
          <label class="mb-1 block font-medium text-foreground">
            Observations / Consignes ministérielles (Optionnel)
          </label>
          <textarea
            v-model="form.observations"
            rows="2"
            placeholder="Consignes particulières, heure de convocation ou mentions spéciales..."
            class="w-full rounded-md border border-border bg-background px-3 py-2 text-xs text-foreground placeholder:text-muted-foreground/60 focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand"
          ></textarea>
        </div>

        <!-- Pied de page & Boutons d'action -->
        <div class="flex items-center justify-between pt-3 border-t border-border mt-2">
          <div>
            <button
              v-if="isEditing"
              type="button"
              @click="supprimer"
              :disabled="isDeleting || isSubmitting"
              class="inline-flex items-center gap-1.5 rounded-md px-2.5 py-1.5 text-xs font-medium text-rose-600 hover:bg-rose-50 hover:text-rose-700 dark:hover:bg-rose-950/40 transition disabled:opacity-50"
            >
              <Trash2 class="h-3.5 w-3.5" />
              Supprimer la session
            </button>
          </div>

          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="emit('close')"
              class="rounded-md border border-border bg-background px-3 py-1.5 text-xs font-medium text-muted-foreground hover:bg-muted hover:text-foreground transition"
            >
              Annuler
            </button>
            <button
              type="submit"
              :disabled="isSubmitting || isDeleting"
              class="inline-flex items-center gap-1.5 rounded-md bg-brand-dark px-3.5 py-1.5 text-xs font-semibold text-white hover:bg-brand transition disabled:opacity-50 shadow-sm"
            >
              <CheckCircle2 class="h-3.5 w-3.5" />
              {{ isEditing ? 'Mettre à jour' : 'Enregistrer et planifier' }}
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>
