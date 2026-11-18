export interface ColorThemeOption {
  id: string;
  label: string;
  primary: string;
  secondary: string;
  accent: string;
  background: string;
}

export const CUSTOM_COLOR_THEMES: ColorThemeOption[] = [
  {
    id: "default",
    label: "ElevenDrive (Défaut)",
    primary: "oklch(0.461 0.07 245.64)",
    secondary: "oklch(0.929 0.003 264.541)",
    accent: "oklch(0.854 0.012 264.504)",
    background: "oklch(1 0 180)",
  },
  {
    id: "minato",
    label: "Minato",
    primary: "oklch(0.842 0.044 223.823)",
    secondary: "oklch(0.823 0.047 42.719)",
    accent: "oklch(0.823 0.047 42.719)",
    background: "oklch(1 0 180)",
  },
  {
    id: "naruto",
    label: "Naruto",
    primary: "oklch(0.9 0.058 162.52)",
    secondary: "oklch(0.845 0.061 346.042)",
    accent: "oklch(0.845 0.061 346.042)",
    background: "oklch(0.988 0.007 164.116)",
  },
  {
    id: "itachi",
    label: "Itachi",
    primary: "oklch(0.386 0.134 330.604)",
    secondary: "oklch(0.503 0.095 113.009)",
    accent: "oklch(0.465 0.067 201.891)",
    background: "oklch(0.982 0.009 327.635)",
  },
  {
    id: "sasuke",
    label: "Sasuke",
    primary: "oklch(0.834 0.103 113.302)",
    secondary: "oklch(0.775 0.073 77.436)",
    accent: "oklch(0.805 0.118 133.643)",
    background: "oklch(0.996 0.004 111.086)",
  },
  {
    id: "ben10",
    label: "Ben 10",
    primary: "oklch(0.798 0.062 182.021)",
    secondary: "oklch(0.884 0.011 183.982)",
    accent: "oklch(0.884 0.011 183.982)",
    background: "oklch(0.994 0.002 184.255)",
  },
  {
    id: "sakamoto",
    label: "Sakamoto",
    primary: "oklch(0.649 0.107 191.186)",
    secondary: "oklch(0.849 0.006 194.667)",
    accent: "oklch(0.849 0.006 194.667)",
    background: "oklch(0.989 0.006 194.684)",
  },
  {
    id: "kaiju8",
    label: "Kaiju 8",
    primary: "oklch(0.719 0.053 324.541)",
    secondary: "oklch(0.846 0.027 324.311)",
    accent: "oklch(0.784 0.051 324.496)",
    background: "oklch(1 0 180)",
  },
  {
    id: "nagato",
    label: "Nagato",
    primary: "oklch(0.745 0.141 170.13)",
    secondary: "oklch(0.905 0.007 179.629)",
    accent: "oklch(0.905 0.007 179.629)",
    background: "oklch(1 0 180)",
  },
  {
    id: "samsam",
    label: "Samsam",
    primary: "oklch(0.818 0.052 318.267)",
    secondary: "oklch(0.804 0.046 297.566)",
    accent: "oklch(0.822 0.048 335.521)",
    background: "oklch(0.991 0.004 318.326)",
  },
  {
    id: "makapaka",
    label: "Makapaka",
    primary: "oklch(0.883 0.062 6.126)",
    secondary: "oklch(0.895 0.008 4.553)",
    accent: "oklch(0.895 0.008 4.553)",
    background: "oklch(0.982 0.005 4.441)",
  },
  {
    id: "toupi-binou",
    label: "Toupi & Binou",
    primary: "oklch(0.481 0.08 198.818)",
    secondary: "oklch(0.476 0.116 154.911)",
    accent: "oklch(0.342 0.108 255.769)",
    background: "oklch(0.988 0.005 199.536)",
  },
  {
    id: "lou",
    label: "Lou",
    primary: "oklch(0.485 0.079 198.612)",
    secondary: "oklch(0.387 0.163 329.73)",
    accent: "oklch(0.524 0.11 111.699)",
    background: "oklch(0.994 0.002 199.578)",
  },
  {
    id: "oui-oui",
    label: "Oui-Oui",
    primary: "oklch(0.492 0.089 160.792)",
    secondary: "oklch(0.395 0.108 355.195)",
    accent: "oklch(0.395 0.108 355.195)",
    background: "oklch(0.987 0.004 167.709)",
  },
  {
    id: "petite-princesse",
    label: "Petite-Princesse",
    primary: "oklch(0.389 0.052 276.995)",
    secondary: "oklch(0.47 0.046 100.493)",
    accent: "oklch(0.47 0.046 100.493)",
    background: "oklch(1 0 180)",
  }
];

export type BackgroundStyleId = "default" | "dots" | "grid" | "gradient";

export interface BackgroundStyleOption {
  id: BackgroundStyleId;
  label: string;
  description: string;
}

export const CUSTOM_BACKGROUND_STYLES: BackgroundStyleOption[] = [
  { id: "default", label: "Uni (Défaut)", description: "Fond épuré standard" },
  { id: "dots", label: "Points subtils", description: "Trame de points discrets" },
  { id: "grid", label: "Grille technique", description: "Quadrillage léger" },
  { id: "gradient", label: "Halo doux", description: "Dégradé ambiant du thème" },
];
