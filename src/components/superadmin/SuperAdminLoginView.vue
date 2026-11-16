<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, computed } from "vue";
import { useAuthStore } from "../../stores/authStore";
import {
  ShieldCheck,
  Lock,
  User,
  ArrowRight,
  ArrowLeft,
  Loader2,
  AlertCircle,
  Terminal,
  Server,
  Database,
  Cpu,
  Activity,
  Sun,
  Moon,
  Eye,
  EyeOff,
  CheckCircle2,
  Building2,
  Monitor,
  CalendarCheck,
  KeyRound,
  Layers,
  ScrollText,
  WifiOff,
  Car,
  BarChart3,
} from "lucide-vue-next";

const emit = defineEmits<{
  (e: "login-success"): void;
}>();

const authStore = useAuthStore();

const form = reactive({
  identifiant: "",
  motDePasse: "",
});

const isLoading = ref(false);
const showPassword = ref(false);
const errorMessage = ref<string | null>(null);

/* ───────── Vague de points mathématique (Canvas identique à LandingHero) ───────── */
const GRID = 14;
const WAVELENGTH = 260;
const PERIOD = 9;
const SHARPNESS = 8;
const BASE_ALPHA = 0.07;
const PEAK_ALPHA = 0.28;
const BASE_RADIUS = 0.6;
const PEAK_GROW = 0.5;
const WAVE_COLOR = "rgb(16, 185, 129)"; // Émeraude / Brand
const FRAME_MS = 33;

const hostRef = ref<HTMLElement | null>(null);
const canvasRef = ref<HTMLCanvasElement | null>(null);
let raf = 0;
let running = false;
let visible = true;
let last = 0;
let ro: ResizeObserver | null = null;
let io: IntersectionObserver | null = null;
let reduced = false;
let w = 0;
let h = 0;
let dpr = 1;
let baseColor = "rgb(0, 0, 0)";
let lastColorRead = -Infinity;

function resize() {
  const c = canvasRef.value;
  if (!c) return;
  const rect = c.getBoundingClientRect();
  dpr = Math.min(window.devicePixelRatio || 1, 2);
  w = rect.width;
  h = rect.height;
  c.width = Math.round(w * dpr);
  c.height = Math.round(h * dpr);
}

function draw(now: number) {
  const c = canvasRef.value;
  const ctx = c?.getContext("2d");
  if (!c || !ctx || w === 0) return;

  if (now - lastColorRead > 1000) {
    baseColor = getComputedStyle(c).color;
    lastColorRead = now;
  }

  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  ctx.clearRect(0, 0, w, h);

  const t = now / 1000;
  const cx = w / 2;
  const cy = h / 2;
  const rx = 0.6 * w;
  const ry = 0.7 * h;
  const cols = Math.ceil(w / GRID) + 1;
  const rows = Math.ceil(h / GRID) + 1;

  let currentFill = "";

  for (let i = 0; i < cols; i++) {
    for (let j = 0; j < rows; j++) {
      const x = i * GRID + GRID / 2;
      const y = j * GRID + GRID / 2;

      const r = Math.hypot((x - cx) / rx, (y - cy) / ry);
      const mask = Math.min(Math.max((r - 0.22) / 0.65, 0), 1);
      if (mask < 0.03) continue;

      const d = Math.hypot(x - cx, y - cy);
      const wave =
        0.5 + 0.5 * Math.cos(2 * Math.PI * (d / WAVELENGTH - t / PERIOD));
      const k = Math.pow(wave, SHARPNESS);

      const fill = k < 0.05 ? baseColor : WAVE_COLOR;
      if (fill !== currentFill) {
        ctx.fillStyle = fill;
        currentFill = fill;
      }
      ctx.globalAlpha = (BASE_ALPHA + (PEAK_ALPHA - BASE_ALPHA) * k) * mask;
      const s = (BASE_RADIUS + PEAK_GROW * k) * 2;
      ctx.fillRect(x - s / 2, y - s / 2, s, s);
    }
  }
  ctx.globalAlpha = 1;
}

function loop(now: number) {
  if (!running) return;
  raf = requestAnimationFrame(loop);
  if (now - last < FRAME_MS) return;
  last = now;
  draw(now);
}

function start() {
  if (running || reduced || !visible) return;
  running = true;
  raf = requestAnimationFrame(loop);
}

function stop() {
  running = false;
  cancelAnimationFrame(raf);
}

onMounted(() => {
  reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  resize();

  ro = new ResizeObserver(() => {
    resize();
    if (!running) draw(performance.now());
  });
  if (canvasRef.value) ro.observe(canvasRef.value);

  io = new IntersectionObserver(([entry]) => {
    visible = entry.isIntersecting;
    visible ? start() : stop();
  });
  if (hostRef.value) io.observe(hostRef.value);

  if (reduced) draw(0);
  else start();
});

onBeforeUnmount(() => {
  stop();
  ro?.disconnect();
  io?.disconnect();
});

const handleSubmit = async () => {
  errorMessage.value = null;
  if (!form.identifiant.trim() || !form.motDePasse) {
    errorMessage.value = "Identifiant et clé secrète requis.";
    return;
  }

  isLoading.value = true;
  try {
    await authStore.loginSuperAdmin(form.identifiant.trim(), form.motDePasse);
    if (typeof window !== "undefined") {
      window.history.pushState({}, "", "/super-admin/dashboard");
    }
    emit("login-success");
  } catch (err: any) {
    errorMessage.value =
      err?.message || "Échec d'authentification Super-Admin.";
  } finally {
    isLoading.value = false;
  }
};

const handleBack = () => {
  if (
    typeof window !== "undefined" &&
    window.location.pathname.startsWith("/super-admin")
  ) {
    window.history.pushState({}, "", "/");
  }
  authStore.webPortalView =
    authStore.platform === "desktop" ? "desktop_auth" : "landing";
};

// Marquee des services système & capacités de la plateforme (12 items comme LandingHero)
const marqueeItems = computed(() => [
  { icon: Building2, label: "Flotte Auto-Écoles" },
  { icon: Monitor, label: "Quotas Postes Desktop" },
  { icon: CalendarCheck, label: "Licences & Expirations" },
  { icon: KeyRound, label: "Matrice des Permissions" },
  { icon: Database, label: "Cluster Neon PostgreSQL" },
  { icon: Layers, label: "Isolation RLS Multi-Tenant" },
  { icon: ScrollText, label: "Journal d'Audit Système" },
  { icon: Activity, label: "Télémétrie & Latence Live" },
  { icon: ShieldCheck, label: "Sécurité & Chiffrement" },
  { icon: WifiOff, label: "Mode Résilience Hors Ligne" },
  { icon: Car, label: "Parc Véhicules & Moniteurs" },
  { icon: BarChart3, label: "Supervision des Flux" },
]);
</script>

<template>
  <div
    class="relative h-full w-full bg-background text-foreground flex flex-col overflow-hidden"
  >
    <!-- Zone supérieure & centrale scrollable sur écrans ultra-compacts -->
    <div class="flex-1 overflow-y-auto flex flex-col min-h-0">
      <!-- Conteneur principal aux traits verticaux identiques à LandingHero -->
      <div
        class="relative mx-auto w-full max-w-6xl flex-1 border-x border-border flex flex-col min-h-0"
      >
        <!-- Barre d'actions supérieure discrète (Retour & Toggle thème) -->
        <header
          class="shrink-0 flex items-center justify-between border-b border-border px-5 py-2.5 sm:px-6"
        >
          <button
            type="button"
            @click="handleBack"
            class="inline-flex items-center gap-1.5 text-[12px] font-medium text-muted-foreground transition hover:text-foreground cursor-pointer"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>Portail Auto-École</span>
          </button>

          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="authStore.toggleDarkMode"
              class="inline-flex items-center justify-center h-8 w-8 rounded-[4px] border border-border bg-card text-muted-foreground transition hover:text-foreground hover:bg-muted cursor-pointer"
              :title="authStore.isDark ? 'Mode clair' : 'Mode sombre'"
            >
              <Sun v-if="authStore.isDark" class="w-4 h-4 text-amber-400" />
              <Moon v-else class="w-4 h-4 text-slate-600" />
            </button>
          </div>
        </header>

        <!-- Zone centrale avec la trame animée et le formulaire -->
        <div
          ref="hostRef"
          class="relative flex-1 flex flex-col items-center justify-center overflow-hidden px-4 py-3 sm:py-5 my-auto"
        >
          <!-- Canvas d'animation mathématique identique à LandingHero -->
          <canvas
            ref="canvasRef"
            class="pointer-events-none absolute inset-0 h-full w-full text-foreground"
            aria-hidden="true"
          />

          <!-- Halos lumineux doux aux couleurs du thème -->
          <div
            class="pointer-events-none absolute -left-16 -top-16 h-80 w-80 rounded-full bg-brand/15 blur-3xl"
            aria-hidden="true"
          />
          <div
            class="pointer-events-none absolute -bottom-20 -right-16 h-80 w-80 rounded-full bg-brand/10 blur-3xl"
            aria-hidden="true"
          />

          <!-- Contenu au-dessus de la trame -->
          <div class="relative w-full max-w-md mx-auto text-center">
            <!-- Badge Pill supérieur -->
            <span
              class="inline-flex items-center gap-2 rounded-[4px] border border-border bg-card px-2.5 py-1 text-[11px] font-medium text-muted-foreground mb-2.5 shadow-xs"
            >
              <span>Root Infrastructure · Super-Admin</span>
            </span>

            <!-- Titre principal -->
            <h1
              class="text-[24px] sm:text-[30px] font-semibold leading-[1.15] tracking-[-0.02em] text-foreground mb-1"
            >
              ElevenDrive
              <span class="text-brand">Core Console</span>
            </h1>

            <p class="text-[12px] sm:text-[13px] leading-relaxed text-muted-foreground mb-4">
              Accès sécurisé réservé à l'administration centrale du système et de
              la flotte.
            </p>

            <!-- Carte de Connexion épurée et moderne -->
            <div
              class="text-left bg-card border border-border rounded-[8px] p-4 sm:p-5 shadow-sm backdrop-blur-sm"
            >
              <form @submit.prevent="handleSubmit" class="space-y-3">
                <!-- Message d'erreur élégant -->
                <div
                  v-if="errorMessage"
                  class="p-2.5 rounded-[4px] bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-900/50 text-rose-800 dark:text-rose-300 text-[12px] flex items-start gap-2 animate-fadeIn"
                >
                  <AlertCircle
                    class="w-4 h-4 shrink-0 mt-0.5 text-rose-600 dark:text-rose-400"
                  />
                  <span class="leading-relaxed">{{ errorMessage }}</span>
                </div>

                <!-- Champ Identifiant -->
                <div>
                  <label
                    class="block text-[12px] font-medium text-foreground mb-1"
                  >
                    Identifiant Super-Administrateur
                  </label>
                  <div class="relative">
                    <User
                      class="w-4 h-4 text-muted-foreground absolute left-3 top-2.5 pointer-events-none"
                    />
                    <input
                      v-model="form.identifiant"
                      type="text"
                      required
                      autocomplete="username"
                      placeholder="admin@elevendrive.com ou superadmin"
                      class="w-full h-9 pl-9 pr-3 rounded-[4px] bg-background border border-input text-[13px] text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring transition-colors font-mono"
                    />
                  </div>
                </div>

                <!-- Champ Mot de passe -->
                <div>
                  <div class="flex items-center justify-between mb-1">
                    <label class="block text-[12px] font-medium text-foreground">
                      Clé secrète maître / Mot de passe
                    </label>
                    <span class="text-[10px] font-mono text-muted-foreground">
                      Protocole SHA-256
                    </span>
                  </div>
                  <div class="relative">
                    <Lock
                      class="w-4 h-4 text-muted-foreground absolute left-3 top-2.5 pointer-events-none"
                    />
                    <input
                      v-model="form.motDePasse"
                      :type="showPassword ? 'text' : 'password'"
                      required
                      autocomplete="current-password"
                      placeholder="****************"
                      class="w-full h-9 pl-9 pr-10 rounded-[4px] bg-background border border-input text-[13px] text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring transition-colors font-mono"
                    />
                    <button
                      type="button"
                      @click="showPassword = !showPassword"
                      class="absolute right-3 top-2.5 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                    >
                      <EyeOff v-if="showPassword" class="w-4 h-4" />
                      <Eye v-else class="w-4 h-4" />
                    </button>
                  </div>
                </div>

                <!-- Bouton de Connexion conforme au style LandingHero -->
                <button
                  type="submit"
                  :disabled="isLoading"
                  class="mt-1 inline-flex w-full items-center justify-center gap-1.5 rounded-[4px] bg-brand-600 px-4 py-2 text-[13px] font-medium text-white transition hover:bg-brand-700 active:scale-[0.98] disabled:opacity-50 disabled:pointer-events-none cursor-pointer shadow-xs"
                >
                  <Loader2
                    v-if="isLoading"
                    class="w-4 h-4 animate-spin text-white"
                  />
                  <template v-else>
                    <span>Accéder à la console</span>
                    <ArrowRight class="w-3.5 h-3.5" />
                  </template>
                </button>
              </form>

              <!-- Télémétrie en pied de formulaire -->
              <div
                class="mt-3.5 pt-2.5 border-t border-border flex items-center justify-between text-[10px] text-muted-foreground font-mono"
              >
                <div class="flex items-center gap-1.5">
                  <span class="h-1.5 w-1.5 rounded-full bg-brand-500 animate-pulse" />
                  <span>PostgreSQL Cluster</span>
                </div>
                <span>Télémétrie 24/7</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Trait horizontal pleine largeur & Marquee (Ancré au bas du viewport, TOUJOURS visible) -->
    <div class="shrink-0 relative border-t border-border bg-card/60 backdrop-blur-xs z-10 w-full">
      <div class="mx-auto max-w-6xl border-x border-border">
        <div
          class="marquee relative overflow-hidden py-3 sm:py-3.5"
          aria-label="Services système et capacités de la plateforme"
        >
          <div class="marquee-track flex w-max">
            <!-- Deux groupes identiques : le second est masqué aux lecteurs d'écran -->
            <div
              v-for="n in 2"
              :key="n"
              class="flex shrink-0 items-center gap-10 pr-10"
              :aria-hidden="n === 2"
            >
              <div
                v-for="item in marqueeItems"
                :key="`${n}-${item.label}`"
                class="flex shrink-0 items-center gap-2 text-muted-foreground"
              >
                <component
                  :is="item.icon"
                  class="h-4 w-4"
                  :stroke-width="1.75"
                />
                <span class="whitespace-nowrap text-[13px] font-medium">
                  {{ item.label }}
                </span>
              </div>
            </div>
          </div>

          <div
            class="pointer-events-none absolute inset-y-0 left-0 w-24 bg-gradient-to-r from-background to-transparent"
          />
          <div
            class="pointer-events-none absolute inset-y-0 right-0 w-24 bg-gradient-to-l from-background to-transparent"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Couche isolée : le défilement ne dépend plus du reste de la page */
.marquee {
  contain: layout paint;
  transform: translateZ(0);
}
.marquee-track {
  animation: marquee 60s linear infinite;
  will-change: transform;
  backface-visibility: hidden;
}
.marquee:hover .marquee-track {
  animation-play-state: paused;
}
/* Chaque groupe = moitié exacte de la piste : boucle sans saut */
@keyframes marquee {
  from {
    transform: translate3d(0, 0, 0);
  }
  to {
    transform: translate3d(-50%, 0, 0);
  }
}
/* Animations réduites : on ralentit au lieu d'arrêter */
@media (prefers-reduced-motion: reduce) {
  .marquee-track {
    animation-duration: 120s;
  }
}
</style>
