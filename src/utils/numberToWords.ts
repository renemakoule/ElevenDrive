// Convertisseur de montants numériques en lettres (Français / Devise FCFA)
// Conforme aux règles comptables pour les reçus et quittances légales

const UNITES = ['', 'un', 'deux', 'trois', 'quatre', 'cinq', 'six', 'sept', 'huit', 'neuf'];
const DIZAINES_SPECIALES = ['dix', 'onze', 'douze', 'treize', 'quatorze', 'quinze', 'seize', 'dix-sept', 'dix-huit', 'dix-neuf'];
const DIZAINES = ['', 'dix', 'vingt', 'trente', 'quarante', 'cinquante', 'soixante', 'soixante-dix', 'quatre-vingts', 'quatre-vingt-dix'];

function convertUnderHundred(n: number): string {
  if (n < 10) return UNITES[n];
  if (n >= 10 && n < 20) return DIZAINES_SPECIALES[n - 10];

  const diz = Math.floor(n / 10);
  const unit = n % 10;

  if (diz === 7) {
    if (unit === 1) return 'soixante-et-onze';
    return 'soixante-' + DIZAINES_SPECIALES[unit];
  }
  if (diz === 9) {
    return 'quatre-vingt-' + DIZAINES_SPECIALES[unit];
  }

  const dizText = DIZAINES[diz];
  if (unit === 0) return dizText;
  if (unit === 1 && diz !== 8) return dizText + '-et-un';
  return (diz === 8 ? 'quatre-vingt' : dizText) + '-' + UNITES[unit];
}

function convertUnderThousand(n: number): string {
  if (n === 0) return '';
  const cent = Math.floor(n / 100);
  const reste = n % 100;

  let result = '';
  if (cent === 1) {
    result = 'cent';
  } else if (cent > 1) {
    result = UNITES[cent] + (reste === 0 ? ' cents' : ' cent');
  }

  if (reste > 0) {
    const resteText = convertUnderHundred(reste);
    result = result ? `${result} ${resteText}` : resteText;
  }

  return result;
}

export function montantEnLettresFCFA(montant: number): string {
  if (montant === 0) return 'Zéro Franc CFA';

  const m = Math.floor(Math.abs(montant));
  const millions = Math.floor(m / 1000000);
  const milliers = Math.floor((m % 1000000) / 1000);
  const unites = m % 1000;

  const parts: string[] = [];

  if (millions > 0) {
    if (millions === 1) {
      parts.push('un million');
    } else {
      parts.push(`${convertUnderThousand(millions)} millions`);
    }
  }

  if (milliers > 0) {
    if (milliers === 1) {
      parts.push('mille');
    } else {
      parts.push(`${convertUnderThousand(milliers)} mille`);
    }
  }

  if (unites > 0) {
    parts.push(convertUnderThousand(unites));
  }

  const phrase = parts.join(' ').trim();
  // Mettre première lettre en majuscule
  const phraseCap = phrase.charAt(0).toUpperCase() + phrase.slice(1);

  return `${phraseCap} Francs CFA`;
}
