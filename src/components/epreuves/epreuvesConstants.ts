export interface CategoriePermisOption {
  value: string;
  label: string;
  color: string;
}

export const CATEGORIES_PERMIS: CategoriePermisOption[] = [
  { value: 'B', label: 'Permis B (Véhicule Léger)', color: 'text-emerald-600 bg-emerald-500/10 border-emerald-500/20' },
  { value: 'A', label: 'Permis A (Motos & Gros Cubes)', color: 'text-amber-600 bg-amber-500/10 border-amber-500/20' },
  { value: 'A1', label: 'Permis A1 (Léger 2 Roues)', color: 'text-orange-600 bg-orange-500/10 border-orange-500/20' },
  { value: 'C', label: 'Permis C (Poids Lourd Marchandises)', color: 'text-blue-600 bg-blue-500/10 border-blue-500/20' },
  { value: 'D', label: 'Permis D (Transport en Commun)', color: 'text-purple-600 bg-purple-500/10 border-purple-500/20' },
  { value: 'E', label: 'Permis E (Ensembles & Remorques)', color: 'text-rose-600 bg-rose-500/10 border-rose-500/20' },
  { value: 'TOUTES', label: 'Générale (Toutes catégories)', color: 'text-zinc-600 bg-zinc-500/10 border-zinc-500/20' },
];

export const getCategoryBadgeClass = (cat: string): string => {
  const item = CATEGORIES_PERMIS.find((c) => c.value.toUpperCase() === cat.toUpperCase());
  return item ? item.color : 'text-brand bg-brand/10 border-brand/20';
};
