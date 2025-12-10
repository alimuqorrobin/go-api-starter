package jwt

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(claims map[string]interface{}, secret string) (string, error) {
    m := jwt.MapClaims{}
    for k, v := range claims {
        m[k] = v
    }
    if _, ok := m["exp"]; !ok {
        m["exp"] = time.Now().Add(24 * time.Hour).Unix()
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, m)
    return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr, secret string) (map[string]interface{}, error) {
    token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(secret), nil
    })
    if err != nil {
        return nil, err
    }
    if !token.Valid {
        return nil, errors.New("invalid token")
    }
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        out := map[string]interface{}{}
        for k, v := range claims {
            out[k] = v
        }
        return out, nil
    }
    return nil, errors.New("cannot parse claims")
}
