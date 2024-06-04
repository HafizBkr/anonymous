package validator

import (
    "regexp"
    "time"

    "github.com/google/uuid"
)

// IsUUID vérifie si une chaîne est un UUID valide.
func IsUUID(v string) bool {
    _, err := uuid.Parse(v)
    return err == nil
}

// IsEmptyString vérifie si une chaîne est vide.
func IsEmptyString(v string) bool {
    return v == ""
}

// IsEqual vérifie si deux valeurs sont égales.
func IsEqual(x interface{}, y interface{}) bool {
    return x == y
}

// IsNotEqual vérifie si deux valeurs sont différentes.
func IsNotEqual(x interface{}, y interface{}) bool {
    return x != y
}

// IsEmail vérifie si une chaîne est une adresse e-mail valide.
func IsEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    regex := regexp.MustCompile(pattern)
    return regex.MatchString(email)
}

// IsOneOf vérifie si une valeur est égale à l'une des valeurs spécifiées.
func IsOneOf(x interface{}, ys ...string) bool {
    ok := false
    for _, y := range ys {
        if x == y {
            ok = true
            break
        }
    }
    return ok
}

// IsValidDate vérifie si une chaîne est une date valide au format "YYYY-MM-DD".
func IsValidDate(v string) bool {
    _, err := time.Parse("2006-01-02", v)
    return err == nil
}

// ValidateStruct valide une structure en appliquant des règles de validation spécifiques.
func ValidateStruct(payload interface{}) error {
    // Implémentez votre logique de validation de structure ici
    // Utilisez les fonctions de validation existantes pour valider chaque champ de la structure
    // Retournez une erreur si la validation échoue
    return nil
}

// IsCorrectPhoneNumber vérifie si une chaîne est un numéro de téléphone valide.
func IsCorrectPhoneNumber(v string) bool {
    // Implémentez la validation du numéro de téléphone selon vos besoins
    return true
}
