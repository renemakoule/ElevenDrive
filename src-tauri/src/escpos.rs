/// Module générateur de séquences ESC/POS pour imprimantes thermiques 80mm
/// Supporte le tiroir-caisse, la coupure papier, le formatage gras/double hauteur et les codes de sécurité

pub struct EscposTicketBuilder {
    buffer: Vec<u8>,
}

impl EscposTicketBuilder {
    pub fn new() -> Self {
        let mut builder = Self {
            buffer: Vec::with_capacity(1024),
        };
        // Initialisation imprimante (ESC @)
        builder.buffer.extend_from_slice(&[0x1B, 0x40]);
        builder
    }

    /// Ouvre le tiroir-caisse connecté via le port RJ11 de l'imprimante (ESC p 0 25 250)
    pub fn open_cash_drawer(&mut self) -> &mut Self {
        self.buffer.extend_from_slice(&[0x1B, 0x70, 0x00, 0x19, 0xFA]);
        self
    }

    /// Aligner le texte : 0 = Gauche, 1 = Centre, 2 = Droite
    pub fn align(&mut self, align: u8) -> &mut Self {
        self.buffer.extend_from_slice(&[0x1B, 0x61, align]);
        self
    }

    /// Activer / désactiver le mode gras
    pub fn bold(&mut self, enable: bool) -> &mut Self {
        let val = if enable { 1 } else { 0 };
        self.buffer.extend_from_slice(&[0x1B, 0x45, val]);
        self
    }

    /// Activer double hauteur / largeur pour le titre
    pub fn double_size(&mut self, enable: bool) -> &mut Self {
        let val = if enable { 0x11 } else { 0x00 };
        self.buffer.extend_from_slice(&[0x1D, 0x21, val]);
        self
    }

    /// Ajouter du texte brut avec saut de ligne
    pub fn println(&mut self, text: &str) -> &mut Self {
        self.buffer.extend_from_slice(text.as_bytes());
        self.buffer.push(0x0A);
        self
    }

    /// Ajouter une ligne pointillée de séparation 80mm (environ 44-48 caractères)
    pub fn separator(&mut self) -> &mut Self {
        self.println("------------------------------------------------");
        self
    }

    pub fn double_separator(&mut self) -> &mut Self {
        self.println("================================================");
        self
    }

    /// Imprimer deux colonnes alignées (Gauche et Droite)
    pub fn print_row(&mut self, left: &str, right: &str) -> &mut Self {
        let max_width = 44;
        let left_len = left.chars().count();
        let right_len = right.chars().count();

        if left_len + right_len >= max_width {
            self.println(&format!("{} {}", left, right));
        } else {
            let spaces = max_width - (left_len + right_len);
            let space_str = " ".repeat(spaces);
            self.println(&format!("{}{}{}", left, space_str, right));
        }
        self
    }

    /// Avance le papier et effectue la coupure automatique (GS V 66 0)
    pub fn cut(&mut self) -> &mut Self {
        // Avancer de 4 lignes pour dégager le ticket du massicot
        self.buffer.extend_from_slice(&[0x1B, 0x64, 0x04]);
        // Coupure partielle/complète
        self.buffer.extend_from_slice(&[0x1D, 0x56, 0x42, 0x00]);
        self
    }

    /// Récupérer les octets bruts pour l'envoi vers le port USB / Réseau / COM
    pub fn build(self) -> Vec<u8> {
        self.buffer
    }
}
