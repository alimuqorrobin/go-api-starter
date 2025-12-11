package validator

import (
    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
}

// ValidateStruct validates a struct
func ValidateStruct(s interface{}) error {
    return validate.Struct(s)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
    return validate
}

type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// Gin auto-validate dengan ShouldBindJSON
if err := c.ShouldBindJSON(&req); err != nil {
    response.ValidationError(c, err.Error())
    return
}
```

**Validation Tags:**
```
required     → Field wajib diisi
min=3        → Minimum 3 characters
max=50       → Maximum 50 characters
email        → Valid email format
omitempty    → Optional field
gt=0         → Greater than 0
gte=18       → Greater than or equal 18
```

---

### ✅ Part 4 Selesai!

Sudah selesai semua file di Part 4:
- ✅ pkg/response/response.go (Generic JSON response)
- ✅ pkg/logger/logger.go (Logger dengan rotation)
- ✅ pkg/jwt/jwt.go (JWT token service)
- ✅ pkg/validator/validator.go (Input validation)

**Struktur sekarang:**
```
// golang-api-starter/
// ├── cmd/server/main.go
// ├── config/config.go
// ├── internal/
// │   ├── database/ (4 files)
// │   ├── domain/ (3 files)
// │   ├── repository/ (2 files)
// │   └── service/ (4 files)
// ├── pkg/
// │   ├── response/
// │   │   └── response.go
// │   ├── logger/
// │   │   └── logger.go
// │   ├── jwt/
// │   │   └── jwt.go
// │   └── validator/
// │       └── validator.go
// ├── go.mod
// └── .env.example